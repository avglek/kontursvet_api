package handler

import (
	"net/http"
	"strconv"

	"kontursvet-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	repo *repository.PortfolioRepository
}

func NewPortfolioHandler(repo *repository.PortfolioRepository) *PortfolioHandler {
	return &PortfolioHandler{repo: repo}
}

func (h *PortfolioHandler) ListCases(c *gin.Context) {
	cases, err := h.repo.ListCases(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cases": cases})
}

func (h *PortfolioHandler) GetCase(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	caseData, err := h.repo.GetCaseWithPhotos(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "case not found"})
		return
	}

	c.JSON(http.StatusOK, caseData)
}
