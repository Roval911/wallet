package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"wallet/db"
)

type WalletRequest struct {
	WalletId      uuid.UUID `json:"walletId" binding:"required"`
	OperationType string    `json:"operationType" binding:"required"`
	Amount        int64     `json:"amount" binding:"required"`
}

// CreateWallet godoc
// @Summary Создать кошелек
// @Description Создает новый кошелек с начальным балансом. Баланс должен быть неотрицательным.
// @Tags Wallet
// @Accept json
// @Produce json
// @Param input body int64 true "Начальный баланс (неотрицательное число)"
// @Success 201 {object} map[string]interface{} "ID созданного кошелька"
// @Failure 400 {object} map[string]string "Неверный запрос, например, отрицательный баланс"
// @Failure 500 {object} map[string]string "Ошибка на стороне сервера, например, проблема с базой данных"
// @Router /wallets [post]
func CreateWallet(c *gin.Context) {
	var balance int64

	if err := c.ShouldBindJSON(&balance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if balance < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Balance must be non-negative"})
		return
	}

	walletId, err := db.CreateWallet(balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create wallet: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"walletId": walletId})
}

// GetWalletBalance godoc
// @Summary Получить баланс кошелька
// @Description Возвращает текущий баланс кошелька по его UUID.
// @Tags Wallet
// @Accept json
// @Produce json
// @Param walletId path string true "UUID кошелька"
// @Success 200 {object} map[string]interface{} "Баланс кошелька"
// @Failure 400 {object} map[string]string "Неверный UUID"
// @Failure 404 {object} map[string]string "Кошелек не найден"
// @Router /wallets/{walletId} [get]
func GetWalletBalance(c *gin.Context) {
	walletId, err := uuid.Parse(c.Param("walletId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный идентификатор кошелька"})
		return
	}

	balance, err := db.GetWalletBalance(walletId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"walletId": walletId, "balance": balance})
}

// UpdateWallet godoc
// @Summary Обновить баланс кошелька
// @Description Пополнение или списание средств с кошелька.
// @Tags Wallet
// @Accept json
// @Produce json
// @Param request body WalletRequest true "Запрос для обновления кошелька"
// @Success 200 {string} string "OK"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Ошибка на стороне сервера"
// @Router /wallet [post]
func UpdateWallet(c *gin.Context) {
	var req WalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.UpdateWallet(req.WalletId, req.OperationType, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
