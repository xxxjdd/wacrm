package handlers

import (
	"net/http"
	"strconv"
	"time"

	"wacrm-api/config"
	"wacrm-api/models"
	"wacrm-api/services"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	type RegisterInput struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, map[string]string{"message": "User created"})
}

func Login(c *gin.Context) {
	resp := map[string]interface{}{
		"token": "demo-token",
		"user":  map[string]interface{}{"id": 1, "username": "demo"},
	}
	c.JSON(http.StatusOK, resp)
}

func GetMe(c *gin.Context) {
	resp := map[string]interface{}{"id": 1, "username": "demo"}
	c.JSON(http.StatusOK, resp)
}

func ListAccounts(c *gin.Context) {
	resp := map[string]interface{}{"accounts": []interface{}{}}
	c.JSON(http.StatusOK, resp)
}

func CreateAccount(c *gin.Context) {
	resp := map[string]string{"message": "Account created"}
	c.JSON(http.StatusCreated, resp)
}

func DeleteAccount(c *gin.Context) {
	resp := map[string]string{"message": "Account deleted"}
	c.JSON(http.StatusOK, resp)
}

func LogoutAccount(c *gin.Context) {
	resp := map[string]string{"message": "Account logged out"}
	c.JSON(http.StatusOK, resp)
}

func GetQRCode(c *gin.Context) {
	resp := map[string]string{"qr_code": "demo", "message": "Scan QR code"}
	c.JSON(http.StatusOK, resp)
}

func VerifyAccount(c *gin.Context) {
	resp := map[string]string{"message": "Account verified"}
	c.JSON(http.StatusOK, resp)
}

func ConnectAccount(c *gin.Context) {
	accountID := c.Param("id")
	var account models.WhatsAppAccount
	if err := config.DB.First(&account, accountID).Error; err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
		return
	}

	// 启动WhatsApp连接
	_, err := services.StartService(account.ID, account.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := map[string]string{"message": "Connected"}
	c.JSON(http.StatusOK, resp)
}

func DisconnectAccount(c *gin.Context) {
	accountID := c.Param("id")
	id, _ := strconv.ParseUint(accountID, 10, 32)

	services.StopService(uint(id))

	resp := map[string]string{"message": "Disconnected"}
	c.JSON(http.StatusOK, resp)
}

func ListCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	resp := map[string]interface{}{
		"customers": []interface{}{},
		"total":    0,
		"page":     page,
		"limit":    limit,
	}
	c.JSON(http.StatusOK, resp)
}

func CreateCustomer(c *gin.Context) {
	resp := map[string]string{"message": "Customer created"}
	c.JSON(http.StatusCreated, resp)
}

func UpdateCustomer(c *gin.Context) {
	resp := map[string]string{"message": "Customer updated"}
	c.JSON(http.StatusOK, resp)
}

func DeleteCustomer(c *gin.Context) {
	resp := map[string]string{"message": "Customer deleted"}
	c.JSON(http.StatusOK, resp)
}

func ImportCustomers(c *gin.Context) {
	resp := map[string]interface{}{"message": "Import completed", "imported": 0, "skipped": 0}
	c.JSON(http.StatusOK, resp)
}

func ListMessages(c *gin.Context) {
	resp := map[string]interface{}{"messages": []interface{}{}}
	c.JSON(http.StatusOK, resp)
}

func SendMessage(c *gin.Context) {
	now := time.Now()
	resp := map[string]interface{}{
		"message": "Message sent",
		"data":    map[string]interface{}{"id": 1, "sent_at": now},
	}
	c.JSON(http.StatusOK, resp)
}

func ListConversations(c *gin.Context) {
	resp := map[string]interface{}{"conversations": []interface{}{}}
	c.JSON(http.StatusOK, resp)
}

func ListTemplates(c *gin.Context) {
	resp := map[string]interface{}{"templates": []interface{}{}}
	c.JSON(http.StatusOK, resp)
}

func CreateTemplate(c *gin.Context) {
	resp := map[string]string{"message": "Template created"}
	c.JSON(http.StatusCreated, resp)
}

func UpdateTemplate(c *gin.Context) {
	resp := map[string]string{"message": "Template updated"}
	c.JSON(http.StatusOK, resp)
}

func DeleteTemplate(c *gin.Context) {
	resp := map[string]string{"message": "Template deleted"}
	c.JSON(http.StatusOK, resp)
}

func ListTasks(c *gin.Context) {
	resp := map[string]interface{}{"tasks": []interface{}{}}
	c.JSON(http.StatusOK, resp)
}

func CreateTask(c *gin.Context) {
	resp := map[string]string{"message": "Task created"}
	c.JSON(http.StatusCreated, resp)
}

func UpdateTask(c *gin.Context) {
	resp := map[string]string{"message": "Task updated"}
	c.JSON(http.StatusOK, resp)
}

func DeleteTask(c *gin.Context) {
	resp := map[string]string{"message": "Task deleted"}
	c.JSON(http.StatusOK, resp)
}

func RunTaskNow(c *gin.Context) {
	resp := map[string]string{"message": "Task executed"}
	c.JSON(http.StatusOK, resp)
}

func GetStatsOverview(c *gin.Context) {
	resp := map[string]interface{}{
		"total_accounts":   0,
		"online_accounts": 0,
		"total_customers": 0,
		"total_messages":  0,
		"today_messages":  0,
	}
	c.JSON(http.StatusOK, resp)
}

func GetMessageStats(c *gin.Context) {
	resp := map[string]interface{}{"stats": []interface{}{}}
	c.JSON(http.StatusOK, resp)
}

func UpdateUser(c *gin.Context) {
	resp := map[string]string{"message": "User updated"}
	c.JSON(http.StatusOK, resp)
}
