package dashboard

import (
	"net/http"

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

// Summary returns aggregated dashboard data
func (h *Handler) Summary(c *gin.Context) {
	var results []struct {
		Type   string  `json:"type"`
		Amount float64 `json:"amount"`
	}

	// Calculate totals by type
	if err := h.DB.Model(&models.Record{}).Select("type, SUM(amount) as amount").Group("type").Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summary data"})
		return
	}

	var totalIncome, totalExpenses float64
	for _, res := range results {
		if res.Type == string(models.TypeIncome) {
			totalIncome = res.Amount
		} else if res.Type == string(models.TypeExpense) {
			totalExpenses = res.Amount
		}
	}

	// Fetch category-wise totals
	var categoryTotals []struct {
		Category string  `json:"category"`
		Amount   float64 `json:"amount"`
	}
	h.DB.Model(&models.Record{}).Select("category, SUM(amount) as amount").Group("category").Scan(&categoryTotals)

	// Fetch recent activity
	var recentRecords []models.Record
	h.DB.Order("date desc").Limit(5).Find(&recentRecords)

	c.JSON(http.StatusOK, gin.H{
		"total_income":    totalIncome,
		"total_expenses":  totalExpenses,
		"net_balance":     totalIncome - totalExpenses,
		"category_totals": categoryTotals,
		"recent_activity": recentRecords,
	})
}
