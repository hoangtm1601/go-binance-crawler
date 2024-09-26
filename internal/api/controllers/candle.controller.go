package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hoangtm1601/go-binance-crawler/internal/api/services"
	"github.com/hoangtm1601/go-binance-crawler/internal/models"
)

type CandleController struct {
	service *services.CandleService
}

func NewCandleController(service *services.CandleService) *CandleController {
	return &CandleController{service: service}
}

func (ctrl *CandleController) CreateCandle(c *gin.Context) {
	var candle models.Candle
	if err := c.ShouldBindJSON(&candle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.CreateCandle(&candle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, candle)
}

func (ctrl *CandleController) GetCandleByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	candle, err := ctrl.service.GetCandleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Candle not found"})
		return
	}

	c.JSON(http.StatusOK, candle)
}

func (ctrl *CandleController) GetCandlesBySymbol(c *gin.Context) {
	symbol := c.Param("symbol")

	candles, err := ctrl.service.GetCandlesBySymbol(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, candles)
}
