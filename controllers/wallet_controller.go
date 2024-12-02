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
// @Description Создает новый кошелек с начальным балансом.
// @Tags Wallet
// @Accept json
// @Produce json
// @Param initialBalance body int64 true "Начальный баланс"
// @Success 201 {object} map[string]interface{} "Созданный кошелек"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Ошибка на стороне сервера"
// @Router /wallets [post]
func CreateWallet(c *gin.Context) {
	var req struct {
		InitialBalance int64 `json:"initialBalance" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	walletId, err := db.CreateWallet(req.InitialBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
