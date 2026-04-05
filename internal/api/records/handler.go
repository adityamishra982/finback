package records

import (
	"net/http"
	"strconv"

	"github.com/aditya/finback/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// CreateRecord handles record creation (Admin only)
func (h *Handler) CreateRecord(c *gin.Context) {
	var req CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	record := models.Record{
		UserID:   userID,
		Amount:   req.Amount,
		Type:     models.RecordType(req.Type),
		Category: req.Category,
		Date:     req.Date,
		Notes:    req.Notes,
	}

	if err := h.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// ListRecords handles record listing with filtering
func (h *Handler) ListRecords(c *gin.Context) {
	var records []models.Record
	query := h.DB

	// Filtering
	if recordType := c.Query("type"); recordType != "" {
		query = query.Where("type = ?", recordType)
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	// Pagination (optional enhancement)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("date desc").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      records,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateRecord handles record updates (Admin only)
func (h *Handler) UpdateRecord(c *gin.Context) {
	id := c.Param("id")
	var req UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var record models.Record
	if err := h.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	// Manual update mapping to handle nil pointers
	if req.Amount != nil {
		record.Amount = *req.Amount
	}
	if req.Type != nil {
		record.Type = models.RecordType(*req.Type)
	}
	if req.Category != nil {
		record.Category = *req.Category
	}
	if req.Date != nil {
		record.Date = *req.Date
	}
	if req.Notes != nil {
		record.Notes = *req.Notes
	}

	if err := h.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteRecord handles record deletion (Admin only)
func (h *Handler) DeleteRecord(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Record{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}
