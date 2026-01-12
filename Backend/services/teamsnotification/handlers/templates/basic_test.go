package templates

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTemplateService is a mock implementation of TemplateService
type MockTemplateService struct {
	mock.Mock
}

func (m *MockTemplateService) Create(ctx context.Context, req *services.CreateTemplateRequest) (*models.Template, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Template), args.Error(1)
}

func (m *MockTemplateService) GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Template), args.Error(1)
}

func (m *MockTemplateService) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Template, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Template), args.Error(1)
}

func (m *MockTemplateService) Update(ctx context.Context, id uuid.UUID, req *services.UpdateTemplateRequest) (*models.Template, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Template), args.Error(1)
}

func (m *MockTemplateService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateTemplate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTemplateService)
	handler := NewHandler(mockService)

	projectID := uuid.New()
	templateID := uuid.New()

	expectedTemplate := &models.Template{
		BaseModel: models.BaseModel{
			ID: templateID,
		},
		ProjectID:            projectID,
		Name:                 "Test Template",
		Description:          "Test Description",
		Variables:            models.JSONBTemplateVariables{},
		DefaultJsonStructure: `{"type":"AdaptiveCard","version":"1.4","body":[{"type":"TextBlock","text":"Hello"}]}`,
	}

	mockService.On("Create", mock.Anything, mock.MatchedBy(func(req *services.CreateTemplateRequest) bool {
		return req.Name == "Test Template" &&
			req.Description == "Test Description" &&
			req.ProjectID == projectID
	})).Return(expectedTemplate, nil)

	router := gin.New()
	router.POST("/projects/:id/templates", handler.CreateTemplate)

	reqBody := map[string]interface{}{
		"name":                 "Test Template",
		"description":          "Test Description",
		"variables":            []interface{}{},
		"defaultJsonStructure": `{"type":"AdaptiveCard","version":"1.4","body":[{"type":"TextBlock","text":"Hello"}]}`,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/projects/"+projectID.String()+"/templates", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestListTemplates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTemplateService)
	handler := NewHandler(mockService)

	projectID := uuid.New()
	templateID := uuid.New()

	expectedTemplates := []models.Template{
		{
			BaseModel: models.BaseModel{
				ID: templateID,
			},
			ProjectID:            projectID,
			Name:                 "Template 1",
			Description:          "Description 1",
			Variables:            models.JSONBTemplateVariables{},
			DefaultJsonStructure: `{"type":"AdaptiveCard"}`,
		},
	}

	mockService.On("GetByProjectID", mock.Anything, projectID).Return(expectedTemplates, nil)

	router := gin.New()
	router.GET("/projects/:id/templates", handler.ListTemplates)

	req, _ := http.NewRequest("GET", "/projects/"+projectID.String()+"/templates", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetTemplate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTemplateService)
	handler := NewHandler(mockService)

	projectID := uuid.New()
	templateID := uuid.New()

	expectedTemplate := &models.Template{
		BaseModel: models.BaseModel{
			ID: templateID,
		},
		ProjectID:            projectID,
		Name:                 "Test Template",
		Description:          "Test Description",
		Variables:            models.JSONBTemplateVariables{},
		DefaultJsonStructure: `{"type":"AdaptiveCard"}`,
	}

	mockService.On("GetByID", mock.Anything, templateID).Return(expectedTemplate, nil)

	router := gin.New()
	router.GET("/projects/:id/templates/:templateId", handler.GetTemplate)

	req, _ := http.NewRequest("GET", "/projects/"+projectID.String()+"/templates/"+templateID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateTemplate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTemplateService)
	handler := NewHandler(mockService)

	projectID := uuid.New()
	templateID := uuid.New()

	expectedTemplate := &models.Template{
		BaseModel: models.BaseModel{
			ID: templateID,
		},
		ProjectID:            projectID,
		Name:                 "Updated Template",
		Description:          "Updated Description",
		Variables:            models.JSONBTemplateVariables{},
		DefaultJsonStructure: `{"type":"AdaptiveCard","version":"1.4"}`,
	}

	mockService.On("Update", mock.Anything, templateID, mock.MatchedBy(func(req *services.UpdateTemplateRequest) bool {
		return req.Name == "Updated Template" &&
			req.Description == "Updated Description"
	})).Return(expectedTemplate, nil)

	router := gin.New()
	router.PUT("/projects/:id/templates/:templateId", handler.UpdateTemplate)

	reqBody := map[string]interface{}{
		"name":                 "Updated Template",
		"description":          "Updated Description",
		"variables":            []interface{}{},
		"defaultJsonStructure": `{"type":"AdaptiveCard","version":"1.4"}`,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/projects/"+projectID.String()+"/templates/"+templateID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteTemplate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTemplateService)
	handler := NewHandler(mockService)

	projectID := uuid.New()
	templateID := uuid.New()

	mockService.On("Delete", mock.Anything, templateID).Return(nil)

	router := gin.New()
	router.DELETE("/projects/:id/templates/:templateId", handler.DeleteTemplate)

	req, _ := http.NewRequest("DELETE", "/projects/"+projectID.String()+"/templates/"+templateID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}
