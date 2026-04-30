package handlers

import (
	"net/http"

	"gatorpay-backend/models"
	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// BudgetHandler handles budget API requests
type BudgetHandler struct {
	service *services.BudgetService
}

// NewBudgetHandler creates a new BudgetHandler
func NewBudgetHandler(service *services.BudgetService) *BudgetHandler {
	return &BudgetHandler{service: service}
}

// CreateGoal creates a new savings goal
func (h *BudgetHandler) CreateGoal(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	goal, err := h.service.CreateGoal(userID.(string), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create goal")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Goal created successfully", goal)
}

// GetGoals returns all budget goals for the current user
func (h *BudgetHandler) GetGoals(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	goals, err := h.service.GetGoals(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch goals")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Goals retrieved", goals)
}

// CreateAutoSaveRule creates a new auto-save rule
func (h *BudgetHandler) CreateAutoSaveRule(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreateAutoSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	rule, err := h.service.CreateAutoSaveRule(userID.(string), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create rule")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Auto-save rule created", rule)
}

// ExecuteRoundup executes round-up savings
func (h *BudgetHandler) ExecuteRoundup(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	result, err := h.service.ExecuteRoundup(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to execute roundup")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Roundup executed", result)
}
