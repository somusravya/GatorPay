package handlers

import (
	"net/http"

	"gatorpay-backend/models"
	"gatorpay-backend/services"
	"gatorpay-backend/utils"

	"github.com/gin-gonic/gin"
)

// SocialHandler handles social feed API requests
type SocialHandler struct {
	service *services.SocialService
}

// NewSocialHandler creates a new SocialHandler
func NewSocialHandler(service *services.SocialService) *SocialHandler {
	return &SocialHandler{service: service}
}

// GetFeed returns the social payment feed
func (h *SocialHandler) GetFeed(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	feed, err := h.service.GetFeed(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch feed")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Feed retrieved", feed)
}

// CreatePost creates a social feed post
func (h *SocialHandler) CreatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreateSocialPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	post, err := h.service.CreatePost(userID.(string), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Post created", post)
}

// ReactToPost adds a reaction to a feed item
func (h *SocialHandler) ReactToPost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.ReactToPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	reaction, err := h.service.ReactToPost(userID.(string), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to react")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Reaction added", reaction)
}

// GetFriends returns the friend list
func (h *SocialHandler) GetFriends(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	friends, err := h.service.GetFriends(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch friends")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Friends retrieved", friends)
}

// AddFriend sends a friend request
func (h *SocialHandler) AddFriend(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.AddFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	friendship, err := h.service.AddFriend(userID.(string), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to add friend")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Friend added", friendship)
}
