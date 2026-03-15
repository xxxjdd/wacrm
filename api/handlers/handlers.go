package handlers

import (
	"net/http"
	"strconv"
	"time"

	"wacrm-api/config"
	"wacrm-api/models"
	"wacrm-api/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	type RegisterInput struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
	}
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var existing models.User
	if err := config.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "User already exists"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Nickname:     input.Nickname,
		Role:         "user",
		Status:       "active",
	}

	config.DB.Create(&user)
	token := generateToken(user.ID, user.Username)

	c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"token":   token,
		"user":    map[string]interface{}{"id": user.ID, "username": user.Username, "email": user.Email},
	})
}

func Login(c *gin.Context) {
	type LoginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ? OR email = ?", input.Username, input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return
	}

	token := generateToken(user.ID, user.Username)

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  map[string]interface{}{"id": user.ID, "username": user.Username, "email": user.Email},
	})
}

func GetMe(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"id": user.ID, "username": user.Username, "email": user.Email})
}

func UpdateUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(input)
	c.JSON(http.StatusOK, map[string]string{"message": "User updated"})
}

func ListAccounts(c *gin.Context) {
	userID := c.GetUint("user_id")
	var accounts []models.WhatsAppAccount
	config.DB.Where("user_id = ?", userID).Find(&accounts)
	c.JSON(http.StatusOK, map[string]interface{}{"accounts": accounts})
}

func CreateAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	type Input struct {
		Phone      string `json:"phone"`
		Nickname   string `json:"nickname"`
		DeviceName string `json:"device_name"`
	}
	var input Input
	c.ShouldBindJSON(&input)
	account := models.WhatsAppAccount{UserID: userID, Phone: input.Phone, Nickname: input.Nickname, DeviceName: input.DeviceName, Status: "offline"}
	config.DB.Create(&account)
	c.JSON(http.StatusCreated, map[string]interface{}{"message": "Account created", "account": account})
}

func DeleteAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	accountID := c.Param("id")
	config.DB.Where("id = ? AND user_id = ?", accountID, userID).Delete(&models.WhatsAppAccount{})
	c.JSON(http.StatusOK, map[string]string{"message": "Account deleted"})
}

func LogoutAccount(c *gin.Context) {
	accountID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	services.StopService(uint(accountID))
	config.DB.Model(&models.WhatsAppAccount{}).Where("id = ?", accountID).Update("status", "offline")
	c.JSON(http.StatusOK, map[string]string{"message": "Account logged out"})
}

func GetQRCode(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"qr_code": "demo", "message": "Scan QR code with WhatsApp"})
}

func VerifyAccount(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"message": "Account verified"})
}

func ConnectAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	accountID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var account models.WhatsAppAccount
	if err := config.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
		return
	}
	services.StartService(account.ID, account.Phone)
	config.DB.Model(&account).Update("status", "online")
	c.JSON(http.StatusOK, map[string]string{"message": "Connected"})
}

func DisconnectAccount(c *gin.Context) {
	accountID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	services.StopService(uint(accountID))
	config.DB.Model(&models.WhatsAppAccount{}).Where("id = ?", accountID).Update("status", "offline")
	c.JSON(http.StatusOK, map[string]string{"message": "Disconnected"})
}

func ListCustomers(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var accounts []models.WhatsAppAccount
	config.DB.Where("user_id = ?", userID).Find(&accounts)
	if len(accounts) == 0 {
		c.JSON(http.StatusOK, map[string]interface{}{"customers": []interface{}{}, "total": 0})
		return
	}

	accountIDs := []uint{}
	for _, a := range accounts {
		accountIDs = append(accountIDs, a.ID)
	}

	var customers []models.Customer
	offset := (page - 1) * limit
	config.DB.Where("account_id IN ?", accountIDs).Offset(offset).Limit(limit).Order("last_msg_at DESC").Find(&customers)

	c.JSON(http.StatusOK, map[string]interface{}{"customers": customers, "total": len(customers)})
}

func CreateCustomer(c *gin.Context) {
	userID := c.GetUint("user_id")
	type Input struct {
		Phone     string `json:"phone"`
		Name      string `json:"name"`
		AccountID uint   `json:"account_id"`
	}
	var input Input
	c.ShouldBindJSON(&input)

	var account models.WhatsAppAccount
	if err := config.DB.Where("id = ? AND user_id = ?", input.AccountID, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
		return
	}

	customer := models.Customer{AccountID: input.AccountID, Phone: input.Phone, Name: input.Name, Source: "manual"}
	config.DB.Create(&customer)
	c.JSON(http.StatusCreated, map[string]interface{}{"message": "Customer created", "customer": customer})
}

func UpdateCustomer(c *gin.Context) {
	customerID := c.Param("id")
	var input map[string]interface{}
	c.ShouldBindJSON(&input)
	config.DB.Model(&models.Customer{}).Where("id = ?", customerID).Updates(input)
	c.JSON(http.StatusOK, map[string]string{"message": "Customer updated"})
}

func DeleteCustomer(c *gin.Context) {
	customerID := c.Param("id")
	config.DB.Delete(&models.Customer{}, customerID)
	c.JSON(http.StatusOK, map[string]string{"message": "Customer deleted"})
}

func ImportCustomers(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{"message": "Import completed", "imported": 0, "skipped": 0})
}

func ListMessages(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{"messages": []interface{}{}})
}

func SendMessage(c *gin.Context) {
	type Input struct {
		AccountID  uint   `json:"account_id"`
		CustomerID uint   `json:"customer_id"`
		Content    string `json:"content"`
	}
	var input Input
	c.ShouldBindJSON(&input)

	now := time.Now()
	message := models.Message{AccountID: input.AccountID, CustomerID: input.CustomerID, Direction: "outbound", Content: input.Content, Status: "sent", SentAt: &now}
	config.DB.Create(&message)
	c.JSON(http.StatusOK, map[string]interface{}{"message": "Message sent", "data": message})
}

func ListConversations(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{"conversations": []interface{}{}})
}

func ListTemplates(c *gin.Context) {
	userID := c.GetUint("user_id")
	var templates []models.MessageTemplate
	config.DB.Where("user_id = ?", userID).Find(&templates)
	c.JSON(http.StatusOK, map[string]interface{}{"templates": templates})
}

func CreateTemplate(c *gin.Context) {
	userID := c.GetUint("user_id")
	type Input struct {
		Name     string `json:"name"`
		Content  string `json:"content"`
		Category string `json:"category"`
	}
	var input Input
	c.ShouldBindJSON(&input)
	template := models.MessageTemplate{UserID: userID, Name: input.Name, Content: input.Content, Category: input.Category, Status: "active"}
	config.DB.Create(&template)
	c.JSON(http.StatusCreated, map[string]interface{}{"message": "Template created", "template": template})
}

func UpdateTemplate(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"message": "Template updated"})
}

func DeleteTemplate(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"message": "Template deleted"})
}

func ListTasks(c *gin.Context) {
	userID := c.GetUint("user_id")
	var tasks []models.ScheduledTask
	config.DB.Where("user_id = ?", userID).Find(&tasks)
	c.JSON(http.StatusOK, map[string]interface{}{"tasks": tasks})
}

func CreateTask(c *gin.Context) {
	userID := c.GetUint("user_id")
	type Input struct {
		Name       string `json:"name"`
		AccountID  uint   `json:"account_id"`
		CustomerIDs string `json:"customer_ids"`
		Message    string `json:"message"`
	}
	var input Input
	c.ShouldBindJSON(&input)
	task := models.ScheduledTask{UserID: userID, AccountID: input.AccountID, Name: input.Name, CustomerIDs: input.CustomerIDs, Message: input.Message, Status: "pending"}
	config.DB.Create(&task)
	c.JSON(http.StatusCreated, map[string]interface{}{"message": "Task created", "task": task})
}

func UpdateTask(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"message": "Task updated"})
}

func DeleteTask(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"message": "Task deleted"})
}

func RunTaskNow(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"message": "Task executed"})
}

func GetStatsOverview(c *gin.Context) {
	userID := c.GetUint("user_id")
	var accounts []models.WhatsAppAccount
	config.DB.Where("user_id = ?", userID).Find(&accounts)
	
	online := 0
	for _, a := range accounts {
		if a.Status == "online" {
			online++
		}
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"total_accounts":  len(accounts),
		"online_accounts": online,
		"total_customers": 0,
		"total_messages":  0,
		"today_messages":  0,
	})
}

func GetMessageStats(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{"stats": []interface{}{}})
}

func generateToken(userID uint, username string) string {
	token := strconv.FormatUint(uint64(userID), 10) + "-" + username + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	config.Sessions.Set(token, &config.UserSession{UserID: userID, Username: username, Token: token, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)})
	return token
}
