package delivery

import (
	"cool-games/internal/domain"
	"cool-games/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	GameUsecase domain.GameUsecase
}

func NewGameHandler(r *gin.Engine, us domain.GameUsecase, jwtSecret string) {
	handler := &GameHandler{GameUsecase: us}

	r.GET("/games", handler.Fetch)
	r.GET("/games/:id", handler.GetByID)

	protected := r.Group("/games")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protected.GET("/my-games", middleware.RoleBlock("publisher"), handler.GetMyGames)
		protected.POST("", middleware.RoleBlock("publisher"), handler.Create)
		protected.PUT("/:id", middleware.RoleBlock("publisher"), handler.Update)
		protected.DELETE("/:id", middleware.RoleBlock("publisher", "admin"), handler.Delete)
		protected.PATCH("/:id/restock", middleware.RoleBlock("publisher"), handler.Restock)
	}
}

func (h *GameHandler) Create(c *gin.Context) {
	var g domain.Game
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(int)
	if err := h.GameUsecase.Create(c.Request.Context(), &g, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, g)
}

func (h *GameHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.MustGet("user_id").(int)
	role := c.MustGet("role").(string)

	var g domain.Game
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.GameUsecase.Update(c.Request.Context(), id, &g, userID, role); err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrUnauthorizedAction { status = http.StatusForbidden }
		if err == domain.ErrGameNotFound { status = http.StatusNotFound }
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *GameHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := h.GameUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GameHandler) Fetch(c *gin.Context) {
	search := c.Query("search")
	minPrice, _ := strconv.ParseFloat(c.DefaultQuery("min_price", "0"), 64)
	maxPrice, _ := strconv.ParseFloat(c.DefaultQuery("max_price", "0"), 64)

	res, err := h.GameUsecase.GetAll(c.Request.Context(), search, minPrice, maxPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GameHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	role := c.MustGet("role").(string)
	userID := c.MustGet("user_id").(int)

	if err := h.GameUsecase.Delete(c.Request.Context(), id, userID, role); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *GameHandler) GetMyGames(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	res, err := h.GameUsecase.GetByPublisher(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GameHandler) Restock(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.MustGet("user_id").(int)
	var req domain.RestockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.GameUsecase.Restock(c.Request.Context(), id, userID, req.Amount); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}