package directory

import (
	"context"
	"net/http"
	"strconv"

	"github.com/evencycu/TeamsNotifyGoV3/libs/response"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	directoryService services.DirectoryService
}

func NewHandler(directoryService services.DirectoryService) *Handler {
	return &Handler{directoryService: directoryService}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	dir := rg.Group("/directory")
	{
		dir.GET("/users/search", h.SearchUsers)
		dir.GET("/users", h.ListUsers)
		dir.GET("/groups", h.ListGroups)
		dir.GET("/groups/teams", h.ListTeamsGroups)
		dir.GET("/groups/teams/with-channels", h.ListTeamsGroupsWithChannels)
		dir.GET("/groups/:groupId/channels", h.ListGroupChannels)
		dir.GET("/channels", h.ListChannels)
		dir.POST("/sync", h.SyncDirectory)
		dir.GET("/sync/status", h.GetSyncStatus)
	}

	// Chat Group management (per project)
	projects := rg.Group("/projects/:id/chat-groups")
	{
		projects.GET("", h.ListChatGroups)
		projects.POST("", h.RegisterChatGroup)
		projects.DELETE("/:chatGroupId", h.RemoveChatGroup)
	}

	// All Chat Groups (admin view)
	dir.GET("/chat-groups", h.ListAllChatGroups)
	// Group chats from bot_installations
	dir.GET("/group-chats/from-bot-installations", h.GetGroupChatsFromBotInstallations)
}

func (h *Handler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "Query parameter 'q' is required", nil, nil)
		return
	}

	users, err := h.directoryService.SearchUsers(c.Request.Context(), query)
	if err != nil {
		response.InternalServerError(c, "Failed to search users", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Users found", users)
}

func (h *Handler) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1000"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := h.directoryService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to list users", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Users listed successfully", users)
}

func (h *Handler) ListGroups(c *gin.Context) {
	groups, err := h.directoryService.ListGroups(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to list groups", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Groups listed successfully", groups)
}

func (h *Handler) ListTeamsGroups(c *gin.Context) {
	groups, err := h.directoryService.ListTeamsGroups(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to list Teams groups", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Teams groups listed successfully", groups)
}

func (h *Handler) ListTeamsGroupsWithChannels(c *gin.Context) {
	groupsWithChannels, err := h.directoryService.ListTeamsGroupsWithChannels(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to list Teams groups with channels", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Teams groups with channels listed successfully", groupsWithChannels)
}

func (h *Handler) ListGroupChannels(c *gin.Context) {
	groupID := c.Param("groupId")
	if groupID == "" {
		response.BadRequest(c, "Path parameter 'groupId' is required", nil, nil)
		return
	}

	channels, err := h.directoryService.ListGroupChannels(c.Request.Context(), groupID)
	if err != nil {
		response.InternalServerError(c, "Failed to list group channels", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Channels listed successfully", channels)
}

func (h *Handler) ListChannels(c *gin.Context) {
	teamID := c.Query("team_id") // Optional: filter by team_id

	channels, err := h.directoryService.ListChannels(c.Request.Context(), teamID)
	if err != nil {
		response.InternalServerError(c, "Failed to list channels", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Channels listed successfully", channels)
}

func (h *Handler) ListChatGroups(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID", err, nil)
		return
	}

	chatGroups, err := h.directoryService.ListChatGroups(c.Request.Context(), projectID)
	if err != nil {
		response.InternalServerError(c, "Failed to list chat groups", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Chat groups listed successfully", chatGroups)
}

type RegisterChatGroupRequest struct {
	Name   string `json:"name" binding:"required"`
	ChatID string `json:"chatId" binding:"required"`
}

func (h *Handler) RegisterChatGroup(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID", err, nil)
		return
	}

	var req RegisterChatGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err, nil)
		return
	}

	chatGroup, err := h.directoryService.RegisterChatGroup(c.Request.Context(), projectID, req.Name, req.ChatID)
	if err != nil {
		response.InternalServerError(c, "Failed to register chat group", err, nil)
		return
	}

	response.Success(c, http.StatusCreated, "Chat group registered successfully", chatGroup)
}

func (h *Handler) RemoveChatGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("chatGroupId"))
	if err != nil {
		response.BadRequest(c, "Invalid chat group ID", err, nil)
		return
	}

	if err := h.directoryService.RemoveChatGroup(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, "Failed to remove chat group", err, nil)
		return
	}

	response.Success(c, http.StatusNoContent, "Chat group removed successfully", nil)
}

func (h *Handler) SyncDirectory(c *gin.Context) {
	// Launch sync in background
	go func() {
		_ = h.directoryService.SyncDirectory(context.Background())
	}()

	response.Success(c, http.StatusAccepted, "Directory sync started", nil)
}

func (h *Handler) GetSyncStatus(c *gin.Context) {
	status, err := h.directoryService.GetSyncStatus(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to get sync status", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Sync status retrieved", status)
}

func (h *Handler) ListAllChatGroups(c *gin.Context) {
	chatGroups, err := h.directoryService.ListAllChatGroups(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to list chat groups", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Chat groups listed successfully", chatGroups)
}

func (h *Handler) GetGroupChatsFromBotInstallations(c *gin.Context) {
	groupChats, err := h.directoryService.GetGroupChatsFromBotInstallations(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to get group chats from bot installations", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Group chats from bot installations retrieved successfully", groupChats)
}
