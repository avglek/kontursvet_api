package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"kontursvet-api/internal/model"
	"kontursvet-api/internal/service"

	"github.com/gin-gonic/gin"
)

type LeadHandler struct {
	leadService   *service.LeadService
	maxBotService *service.MaxBotService
}

func NewLeadHandler(leadService *service.LeadService, maxBotService *service.MaxBotService) *LeadHandler {
	return &LeadHandler{
		leadService:   leadService,
		maxBotService: maxBotService,
	}
}

// @Summary Create lead
// @Description Create a new lead and send notification to MAX
// @Tags leads
// @Accept json
// @Produce json
// @Param lead body model.Lead true "Lead data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/leads [post]
func (h *LeadHandler) Create(c *gin.Context) {
	lead := &model.Lead{}
	if err := c.ShouldBindJSON(lead); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.leadService.Create(ctx, lead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.maxBotService.SendLeadNotification(lead); err != nil {
		// Log error but don't fail the request
		_ = err
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "id": lead.ID})
}

// @Summary List leads
// @Description Get all leads with pagination
// @Tags leads
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/leads [get]
func (h *LeadHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	leads, total, err := h.leadService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leads":  leads,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// @Summary Upload file
// @Description Upload a file for a lead
// @Tags uploads
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/upload [post]
func (h *LeadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	data, _ := io.ReadAll(src)

	filename, err := h.leadService.SaveUploadFile(file.Filename, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"filename": filename})
}

// @Summary Upload JSON with lead data
// @Description Upload JSON with lead data and attachments
// @Tags uploads
// @Accept json
// @Produce json
// @Param data body model.LeadMessage true "Lead message with attachments"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/upload-json [post]
func (h *LeadHandler) UploadJSON(c *gin.Context) {
	jsonData, _ := io.ReadAll(c.Request.Body)
	var msg model.LeadMessage
	if err := json.Unmarshal(jsonData, &msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lead := &model.Lead{
		Name:         msg.Text.Name,
		PhoneDigital: msg.Text.PhoneDigital,
		PhoneFormat:  msg.Text.PhoneFormat,
		HomeType:     msg.Text.HomeType,
		Location:     msg.Text.Location,
		Message:      msg.Text.Message,
	}

	ctx := c.Request.Context()
	if err := h.leadService.Create(ctx, lead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.maxBotService.SendLeadNotification(lead); err != nil {
		_ = err
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "id": lead.ID})
}
