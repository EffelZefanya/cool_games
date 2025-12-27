package delivery

import (
	"cool-games/internal/domain"
	"cool-games/internal/middleware"
	"net/http"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	CustomerUsecase domain.CustomerUsecase
}

func NewCustomerHandler(r *gin.Engine, cu domain.CustomerUsecase, jwtSecret string) {
	handler := &CustomerHandler{CustomerUsecase: cu}

	customerGroup := r.Group("/me")
	customerGroup.Use(middleware.AuthMiddleware(jwtSecret))
	customerGroup.Use(middleware.RoleBlock("customer"))
	{
		customerGroup.GET("/profile", handler.GetProfile)
	}
}

func (h *CustomerHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	profile, err := h.CustomerUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}