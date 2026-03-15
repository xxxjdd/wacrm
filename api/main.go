package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"wacrm-api/config"
	"wacrm-api/handlers"
	"wacrm-api/middleware"
	"wacrm-api/models"
	"wacrm-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 解析命令行参数
	dsn := flag.String("dsn", "", "Database DSN (e.g., user:pass@tcp(host:port)/dbname)")
	port := flag.String("port", "8080", "Server port")
	flag.Parse()
	
	// 如果传入dsn参数，设置环境变量
	if *dsn != "" {
		os.Setenv("DB_DSN", *dsn)
	}
	if *port != "" {
		os.Setenv("PORT", *port)
	}

	// 初始化数据库
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	db.AutoMigrate(
		&models.User{},
		&models.WhatsAppAccount{},
		&models.Customer{},
		&models.Message{},
		&models.MessageTemplate{},
		&models.ScheduledTask{},
	)

	// 初始化配置
	config.InitRedis()
	config.StartSessionCleanup()

	// 设置WhatsApp消息处理函数
	services.SetMessageHandler(func(accountID uint, msg services.WhatsAppMessage) {
		log.Printf("Received message from %s: %s", msg.From, msg.Body)
	})

	// 启动任务调度器
	services.StartScheduler()
	defer services.StopScheduler()

	// 启动时恢复已有的WhatsApp连接
	go restoreWhatsAppConnections()

	// 设置路由
	r := gin.Default()
	r.Use(middleware.CORS())

	// 公开路由
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// 受保护的路由
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		api.GET("/me", handlers.GetMe)
		api.PUT("/users/:id", handlers.UpdateUser)

		api.GET("/accounts", handlers.ListAccounts)
		api.POST("/accounts", handlers.CreateAccount)
		api.DELETE("/accounts/:id", handlers.DeleteAccount)
		api.POST("/accounts/:id/logout", handlers.LogoutAccount)
		api.GET("/accounts/:id/qr", handlers.GetQRCode)
		api.POST("/accounts/:id/verify", handlers.VerifyAccount)
		api.POST("/accounts/:id/connect", handlers.ConnectAccount)
		api.POST("/accounts/:id/disconnect", handlers.DisconnectAccount)

		api.GET("/customers", handlers.ListCustomers)
		api.POST("/customers", handlers.CreateCustomer)
		api.PUT("/customers/:id", handlers.UpdateCustomer)
		api.DELETE("/customers/:id", handlers.DeleteCustomer)
		api.POST("/customers/import", handlers.ImportCustomers)

		api.GET("/messages", handlers.ListMessages)
		api.POST("/messages/send", handlers.SendMessage)
		api.GET("/messages/conversations", handlers.ListConversations)

		api.GET("/templates", handlers.ListTemplates)
		api.POST("/templates", handlers.CreateTemplate)
		api.PUT("/templates/:id", handlers.UpdateTemplate)
		api.DELETE("/templates/:id", handlers.DeleteTemplate)

		api.GET("/tasks", handlers.ListTasks)
		api.POST("/tasks", handlers.CreateTask)
		api.PUT("/tasks/:id", handlers.UpdateTask)
		api.DELETE("/tasks/:id", handlers.DeleteTask)
		api.POST("/tasks/:id/run", handlers.RunTaskNow)

		api.GET("/stats/overview", handlers.GetStatsOverview)
		api.GET("/stats/messages", handlers.GetMessageStats)
	}

	// 端口 (已通过命令行参数设置)
	// port := os.Getenv("PORT")

	// 优雅关闭
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")
		services.StopAllServices()
		services.StopScheduler()
		os.Exit(0)
	}()

	log.Printf("Server starting on port %s", *port)
	r.Run(":" + *port)
}

// restoreWhatsAppConnections 恢复之前的WhatsApp连接
func restoreWhatsAppConnections() {
	// 查找所有状态为online的账号
	var accounts []models.WhatsAppAccount
	config.DB.Where("status = ?", "online").Find(&accounts)

	for _, account := range accounts {
		log.Printf("Restoring WhatsApp connection for account %d", account.ID)
		_, err := services.StartService(account.ID, account.Phone)
		if err != nil {
			log.Printf("Failed to restore connection for account %d: %v", account.ID, err)
		}
	}
}
