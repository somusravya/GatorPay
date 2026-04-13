package handlers

import (
	"log"
	"net/http"

	"flowpay-trading/services"
	"flowpay-trading/utils"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	stockService *services.StockService
}

func NewStockHandler(ss *services.StockService) *StockHandler {
	return &StockHandler{stockService: ss}
}

func (h *StockHandler) GetQuote(c *gin.Context) {
	symbol := c.Param("symbol")
	data, err := h.stockService.GetQuote(symbol)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get quote")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, data)
}

func (h *StockHandler) GetDetails(c *gin.Context) {
	symbol := c.Param("symbol")
	data, err := h.stockService.GetDetails(symbol)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get stock details")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, data)
}

func (h *StockHandler) GetChart(c *gin.Context) {
	symbol := c.Param("symbol")
	data, err := h.stockService.GetChart(symbol)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get chart data")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, data)
}

func (h *StockHandler) Search(c *gin.Context) {
	query := c.Query("q")
	log.Printf("Searching for %s", query)
	data, err := h.stockService.Search(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Search failed")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, data)
}

func (h *StockHandler) MarketSummary(c *gin.Context) {
	data, err := h.stockService.MarketSummary()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get market summary")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, data)
}
