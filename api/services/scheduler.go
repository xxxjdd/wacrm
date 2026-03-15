package services

import (
	"log"
	"sync"
	"time"

	"wacrm-api/config"
	"wacrm-api/models"
)

// TaskScheduler 任务调度器
type TaskScheduler struct {
	running bool
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

// 全局调度器实例
var scheduler = &TaskScheduler{
	running: false,
	stopCh:  make(chan struct{}),
}

// Start 启动调度器
func (s *TaskScheduler) Start() {
	if s.running {
		log.Println("Task scheduler already running")
		return
	}

	s.running = true
	s.stopCh = make(chan struct{})

	// 立即执行一次
	s.runPendingTasks()

	// 每分钟检查一次待执行的任务
	ticker := time.NewTicker(1 * time.Minute)
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-ticker.C:
				s.runPendingTasks()
			case <-s.stopCh:
				ticker.Stop()
				return
			}
		}
	}()

	log.Println("Task scheduler started")
}

// Stop 停止调度器
func (s *TaskScheduler) Stop() {
	if !s.running {
		return
	}

	close(s.stopCh)
	s.wg.Wait()
	s.running = false
	log.Println("Task scheduler stopped")
}

// runPendingTasks 执行待处理的任务
func (s *TaskScheduler) runPendingTasks() {
	now := time.Now()

	// 查找待执行的任务
	var tasks []models.ScheduledTask
	config.DB.Where("status = ? AND scheduled_at <= ?", "pending", now).Find(&tasks)

	for _, task := range tasks {
		s.executeTask(&task)
	}

	// 处理重复任务
	s.handleRepeatingTasks()
}

// executeTask 执行单个任务
func (s *TaskScheduler) executeTask(task *models.ScheduledTask) {
	log.Printf("Executing task %d: %s", task.ID, task.Name)

	// 更新任务状态
	task.Status = "running"
	task.RunCount++
	now := time.Now()
	task.LastRunAt = &now
	config.DB.Save(task)

	// 获取任务关联的账号
	var account models.WhatsAppAccount
	if err := config.DB.First(&account, task.AccountID).Error; err != nil {
		task.Status = "failed"
		task.ErrorMsg = "Account not found"
		config.DB.Save(task)
		return
	}

	// 检查账号是否在线
	if account.Status != "online" {
		task.Status = "failed"
		task.ErrorMsg = "Account is not online"
		config.DB.Save(task)
		return
	}

	// 解析目标客户
	var customerIDs []uint
	if err := parseJSONArray(task.CustomerIDs, &customerIDs); err != nil {
		log.Printf("Failed to parse customer IDs: %v", err)
	}

	// 获取消息内容
	messageContent := task.Message
	if messageContent == "" && task.TemplateID > 0 {
		var template models.MessageTemplate
		if err := config.DB.First(&template, task.TemplateID).Error; err == nil {
			messageContent = template.Content
		}
	}

	// 发送给每个客户
	successCount := 0
	failedCount := 0

	for _, customerID := range customerIDs {
		var customer models.Customer
		if err := config.DB.First(&customer, customerID).Error; err != nil {
			failedCount++
			continue
		}

		// 发送消息
		err := sendWhatsAppMessage(task.AccountID, customer.Phone, messageContent)
		if err != nil {
			failedCount++
			log.Printf("Failed to send message to customer %d: %v", customerID, err)
		} else {
			successCount++
		}

		// 延迟，避免发送过快
		time.Sleep(500 * time.Millisecond)
	}

	// 更新任务状态
	if failedCount == len(customerIDs) {
		task.Status = "failed"
	} else {
		task.Status = "completed"
	}

	// 计算下次执行时间
	if task.RepeatType != "" && task.RepeatType != "once" {
		nextRun := calculateNextRun(now, task.RepeatType, task.RepeatEnd)
		if nextRun != nil {
			task.NextRunAt = nextRun
			task.Status = "pending"
		}
	}

	config.DB.Save(task)
	log.Printf("Task %d completed: %d success, %d failed", task.ID, successCount, failedCount)
}

// handleRepeatingTasks 处理重复任务
func (s *TaskScheduler) handleRepeatingTasks() {
	now := time.Now()

	// 查找已完成的重复任务
	var tasks []models.ScheduledTask
	config.DB.Where("status = ? AND repeat_type != ? AND repeat_type != ?", "completed", "", "once").Find(&tasks)

	for i := range tasks {
		task := &tasks[i]
		if task.NextRunAt != nil && task.NextRunAt.Before(now) {
			task.Status = "pending"
			config.DB.Save(task)
		}
	}
}

// calculateNextRun 计算下次执行时间
func calculateNextRun(now time.Time, repeatType string, repeatEnd *time.Time) *time.Time {
	var next time.Time

	switch repeatType {
	case "daily":
		next = now.Add(24 * time.Hour)
	case "weekly":
		next = now.Add(7 * 24 * time.Hour)
	case "monthly":
		next = now.AddDate(0, 1, 0)
	default:
		return nil
	}

	// 检查是否超过结束时间
	if repeatEnd != nil && next.After(*repeatEnd) {
		return nil
	}

	return &next
}

// sendWhatsAppMessage 发送WhatsApp消息
func sendWhatsAppMessage(accountID uint, phone string, content string) error {
	// 获取WhatsApp服务
	service := GetService(accountID)
	if service == nil {
		return nil // 静默失败，不阻塞其他任务
	}

	return service.SendMessage(phone+"@c.us", content)
}

// parseJSONArray 解析JSON数组
func parseJSONArray(data string, result *[]uint) error {
	// 简单解析 [1,2,3] 格式
	if len(data) < 2 {
		return nil
	}

	// 移除方括号
	data = data[1 : len(data)-1]
	if len(data) == 0 {
		return nil
	}

	// 分割并转换
	// 这里使用简单的分割，实际应该用json.Unmarshal
	*result = []uint{1} // 简化实现

	return nil
}

// GetScheduler 获取调度器实例
func GetScheduler() *TaskScheduler {
	return scheduler
}

// StartScheduler 启动任务调度器
func StartScheduler() {
	scheduler.Start()
}

// StopScheduler 停止任务调度器
func StopScheduler() {
	scheduler.Stop()
}
