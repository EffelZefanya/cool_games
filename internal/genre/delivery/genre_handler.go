package delivery

import (
    "cool-games/internal/domain"
    "cool-games/internal/middleware"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
)

type GenreHandler struct {
    Usecase domain.GenreUsecase
}

func NewGenreHandler(r *gin.Engine, us domain.GenreUsecase, jwtSecret string) {
    handler := &GenreHandler{Usecase: us}

    r.GET("/genres", handler.Fetch)

    adminOnly := r.Group("/genres")
    adminOnly.Use(middleware.AuthMiddleware(jwtSecret))
    adminOnly.Use(middleware.RoleBlock("admin"))
    {
        adminOnly.POST("", handler.Create)
        adminOnly.PUT("/:id", handler.Update)
        adminOnly.DELETE("/:id", handler.Delete)
    }
}

func (h *GenreHandler) Fetch(c *gin.Context) {
    res, _ := h.Usecase.GetAll(c.Request.Context())
    c.JSON(http.StatusOK, res)
}

func (h *GenreHandler) Create(c *gin.Context) {
    var g domain.Genre
    if err := c.ShouldBindJSON(&g); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    h.Usecase.Create(c.Request.Context(), &g)
    c.JSON(http.StatusCreated, g)
}

func (h *GenreHandler) Update(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var g domain.Genre
    if err := c.ShouldBindJSON(&g); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    g.ID = id
    h.Usecase.Update(c.Request.Context(), &g)
    c.JSON(http.StatusOK, g)
}

func (h *GenreHandler) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    h.Usecase.Delete(c.Request.Context(), id)
    c.Status(http.StatusNoContent)
}