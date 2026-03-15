package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"wacrm-api/config"
	"wacrm-api/models"

	"github.com/gorilla/websocket"
)

// WhatsAppService 负责与WhatsApp Web建立WebSocket连接
type WhatsAppService struct {
	conn      *websocket.Conn
	accountID uint
	phone     string
	closed    bool
	mutex     sync.Mutex
}

// 消息处理函数类型
type MessageHandler func(accountID uint, message WhatsAppMessage)

// 全局服务实例
var (
	waServices   = make(map[uint]*WhatsAppService)
	waMutex      sync.RWMutex
	msgHandler   MessageHandler
)

// WhatsAppMessage WhatsApp消息结构
type WhatsAppMessage struct {
	ID           string `json:"id"`
	From         string `json:"from"`
	FromMe       bool   `json:"fromMe"`
	Body         string `json:"body"`
	Type         string `json:"type"`
	Timestamp    int64  `json:"timestamp"`
	MediaURL     string `json:"mediaUrl,omitempty"`
	Participant  string `json:"participant,omitempty"`
}

// QRCodeResponse QR码响应
type QRCodeResponse struct {
	Code     string `json:"code"`
	Ref      string `json:"ref"`
	Labels   []int  `json:"labels"`
	Version  int    `json:"version"`
	Type     string `json:"type"`
}

// NewWhatsAppService 创建新的WhatsApp服务
func NewWhatsAppService(accountID uint, phone string) *WhatsAppService {
	return &WhatsAppService{
		accountID: accountID,
		phone:     phone,
		closed:    false,
	}
}

// SetMessageHandler 设置消息处理函数
func SetMessageHandler(handler MessageHandler) {
	msgHandler = handler
}

// Connect 连接到WhatsApp Web
func (s *WhatsAppService) Connect() error {
	// WhatsApp Web连接地址
	url := "wss://web.whatsapp.com/ws"
	
	// 创建WebSocket连接
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WhatsApp Web: %v", err)
	}
	
	s.conn = conn
	
	// 开始读取消息
	go s.readLoop()
	
	// 更新账号状态
	s.updateStatus("online")
	
	log.Printf("WhatsApp service connected for account %d", s.accountID)
	return nil
}

// readLoop 持续读取WebSocket消息
func (s *WhatsAppService) readLoop() {
	defer func() {
		s.Close()
	}()
	
	for {
		if s.closed {
			break
		}
		
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
		
		s.handleMessage(message)
	}
}

// handleMessage 处理接收到的消息
func (s *WhatsAppService) handleMessage(data []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}
	
	// 检查消息类型
	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}
	
	switch msgType {
	case "message":
		s.handleIncomingMessage(msg)
	case "connection":
		s.handleConnectionUpdate(msg)
	case "qrcode":
		s.handleQRCode(msg)
	}
}

// handleIncomingMessage 处理收到的聊天消息
func (s *WhatsAppService) handleIncomingMessage(data map[string]interface{}) {
	// 解析消息内容
	body, _ := data["body"].(string)
	from, _ := data["from"].(string)
	id, _ := data["id"].(string)
	timestamp, _ := data["timestamp"].(float64)
	
	msg := WhatsAppMessage{
		ID:        id,
		From:      from,
		FromMe:    false,
		Body:      body,
		Type:      "text",
		Timestamp: int64(timestamp),
	}
	
	// 调用消息处理函数
	if msgHandler != nil {
		msgHandler(s.accountID, msg)
	}
	
	// 保存到数据库
	s.saveMessage(msg)
}

// handleConnectionUpdate 处理连接状态更新
func (s *WhatsAppService) handleConnectionUpdate(data map[string]interface{}) {
	status, _ := data["status"].(string)
	log.Printf("WhatsApp connection status: %s", status)
	
	if status == "close" {
		s.updateStatus("offline")
	}
}

// handleQRCode 处理QR码
func (s *WhatsAppService) handleQRCode(data map[string]interface{}) {
	log.Printf("Received QR code update")
}

// SendMessage 发送消息
func (s *WhatsAppService) SendMessage(to string, content string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if s.closed || s.conn == nil {
		return fmt.Errorf("connection is closed")
	}
	
	msg := map[string]interface{}{
		"type": "message",
		"to":   to,
		"body": content,
	}
	
	return s.conn.WriteJSON(msg)
}

// Close 关闭连接
func (s *WhatsAppService) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if s.closed {
		return
	}
	
	s.closed = true
	
	if s.conn != nil {
		s.conn.Close()
	}
	
	// 从全局map中移除
	waMutex.Lock()
	delete(waServices, s.accountID)
	waMutex.Unlock()
	
	s.updateStatus("offline")
	log.Printf("WhatsApp service closed for account %d", s.accountID)
}

// updateStatus 更新账号状态
func (s *WhatsAppService) updateStatus(status string) {
	config.DB.Model(&models.WhatsAppAccount{}).
		Where("id = ?", s.accountID).
		Update("status", status)
}

// saveMessage 保存消息到数据库
func (s *WhatsAppService) saveMessage(msg WhatsAppMessage) {
	// 查找或创建客户
	phone := extractPhone(msg.From)
	
	var customer models.Customer
	result := config.DB.Where("account_id = ? AND phone = ?", s.accountID, phone).First(&customer)
	
	if result.Error != nil {
		// 创建新客户
		customer = models.Customer{
			AccountID: s.accountID,
			Phone:     phone,
			Source:    "whatsapp",
		}
		config.DB.Create(&customer)
	}
	
	// 保存消息
	message := models.Message{
		AccountID:   s.accountID,
		CustomerID:  customer.ID,
		Direction:   "inbound",
		Content:     msg.Body,
		MessageType: msg.Type,
		Status:      "received",
	}
	
	sentAt := time.Unix(msg.Timestamp, 0)
	message.SentAt = &sentAt
	
	config.DB.Create(&message)
	
	// 更新客户最后消息时间
	config.DB.Model(&customer).Update("last_msg_at", sentAt)
}

// extractPhone 从WhatsApp ID提取电话号码
func extractPhone(waID string) string {
	// WhatsApp ID格式: phone@c.us
	for i := len(waID) - 1; i >= 0; i-- {
		if waID[i] == '@' {
			return waID[:i]
		}
	}
	return waID
}

// GetService 获取账号对应的服务实例
func GetService(accountID uint) *WhatsAppService {
	waMutex.RLock()
	defer waMutex.RUnlock()
	return waServices[accountID]
}

// StartService 启动账号的WhatsApp服务
func StartService(accountID uint, phone string) (*WhatsAppService, error) {
	// 检查是否已存在
	if service := GetService(accountID); service != nil {
		return service, nil
	}
	
	// 创建新服务
	service := NewWhatsAppService(accountID, phone)
	
	// 启动连接
	if err := service.Connect(); err != nil {
		return nil, err
	}
	
	// 保存到全局map
	waMutex.Lock()
	waServices[accountID] = service
	waMutex.Unlock()
	
	return service, nil
}

// StopService 停止账号的WhatsApp服务
func StopService(accountID uint) {
	if service := GetService(accountID); service != nil {
		service.Close()
	}
}

// StopAllServices 停止所有服务
func StopAllServices() {
	waMutex.Lock()
	defer waMutex.Unlock()
	
	for _, service := range waServices {
		service.Close()
	}
}
