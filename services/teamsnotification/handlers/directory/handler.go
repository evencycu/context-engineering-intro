package directory

import (
	"context"
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV3/libs/response"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/services"
	"github.com/gin-gonic/gin"
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
		dir.GET("/groups", h.ListGroups)
		dir.POST("/sync", h.SyncDirectory)
	}
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

func (h *Handler) ListGroups(c *gin.Context) {
	groups, err := h.directoryService.ListGroups(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to list groups", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Groups listed successfully", groups)
}

func (h *Handler) SyncDirectory(c *gin.Context) {
	// Launch sync in background
	go func() {
		_ = h.directoryService.SyncDirectory(context.Background())
	}()

	response.Success(c, http.StatusAccepted, "Directory sync started", nil)
}
