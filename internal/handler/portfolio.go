package handler

import (
	"net/http"
	"strconv"

	_ "kontursvet-api/internal/model"
	"kontursvet-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	repo *repository.PortfolioRepository
}

func NewPortfolioHandler(repo *repository.PortfolioRepository) *PortfolioHandler {
	return &PortfolioHandler{repo: repo}
}

// @Summary List portfolio cases
// @Description Get all portfolio cases
// @Tags portfolio
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/portfolio [get]
func (h *PortfolioHandler) ListCases(c *gin.Context) {
	cases, err := h.repo.ListCases(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cases": cases})
}

// @Summary Get portfolio case by ID
// @Description Get a single portfolio case with photos
// @Tags portfolio
// @Produce json
// @Param id path int true "Case ID"
// @Success 200 {object} model.PortfolioCase
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/portfolio/{id} [get]
func (h *PortfolioHandler) GetCase(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	caseData, err := h.repo.GetCaseWithPhotos(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "case not found"})
		return
	}

	c.JSON(http.StatusOK, caseData)
}
