package handlers

import (
	"fmt"
	"net/http"

	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

type StatementHandler struct {
	statementService *services.StatementService
}

func NewStatementHandler(ss *services.StatementService) *StatementHandler {
	return &StatementHandler{statementService: ss}
}

func (h *StatementHandler) GetStatement(c *gin.Context) {
	userID, _ := c.Get("userID")
	format := c.Query("format")
	start := c.Query("start")
	end := c.Query("end")

	if format != "csv" && format != "pdf" {
		format = "csv"
	}

	data, contentType, err := h.statementService.GenerateStatement(userID.(string), format, start, end)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate statement")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=statement.%s", format))
	c.Data(http.StatusOK, contentType, data)
}
