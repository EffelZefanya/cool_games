package delivery

import (
    "cool-games/internal/domain"
    "cool-games/internal/middleware"
    "net/http"
    "github.com/gin-gonic/gin"
)

type OrderHandler struct {
    Usecase domain.OrderUsecase 
}

func NewOrderHandler(r *gin.Engine, us domain.OrderUsecase, jwtSecret string) {
    handler := &OrderHandler{Usecase: us} 

    protected := r.Group("/orders")
    protected.Use(middleware.AuthMiddleware(jwtSecret))
    {
        protected.POST("/buy", middleware.RoleBlock("customer"), handler.Purchase)
		protected.POST("/topup", middleware.RoleBlock("customer"), handler.TopUp)
        
        protected.GET("/sales-report", middleware.RoleBlock("publisher"), handler.GetSalesReport)
		protected.GET("/library", middleware.RoleBlock("customer"), handler.GetLibrary)
    }
}

func (h *OrderHandler) Purchase(c *gin.Context) {
    var req domain.PurchaseRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.MustGet("user_id").(int)
    err := h.Usecase.BuyGame(c.Request.Context(), userID, req.GameID) 
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Purchase successful!"})
}

func (h *OrderHandler) GetSalesReport(c *gin.Context) {
    publisherID := c.MustGet("user_id").(int)

    report, err := h.Usecase.GetPublisherSalesReport(c.Request.Context(), publisherID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report: " + err.Error()})
        return
    }

    if report == nil {
        report = []domain.SalesReportEntry{}
    }

    c.JSON(http.StatusOK, report)
}

func (h *OrderHandler) TopUp(c *gin.Context) {
    var req struct {
        Amount float64 `json:"amount" binding:"required,gt=0"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
        return
    }

    userID := c.MustGet("user_id").(int)
    
    err := h.Usecase.AddBalance(c.Request.Context(), userID, req.Amount)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Balance updated successfully!"})
}

func (h *OrderHandler) GetLibrary(c *gin.Context) {
    userID := c.MustGet("user_id").(int)

    games, err := h.Usecase.GetCustomerLibrary(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch library"})
        return
    }

    if games == nil {
        games = []domain.Game{}
    }

    c.JSON(http.StatusOK, games)
}