package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/adapters/dtos"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/core/domain"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/core/ports"
)

type UserHandler struct {
	userService ports.UserService // เรียกผ่าน Interface
}

func NewUserHandler(r *gin.Engine, svc ports.UserService) *UserHandler {
	handler := &UserHandler{userService: svc}
	return handler
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.Register(domain.User{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "สมัครสมาชิกไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ลงทะเบียนสำเร็จ!"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	user, err := h.userService.GetByID(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ดึงข้อมูลไม่สำเร็จ"})
		return
	}
	userResponse := dtos.ToUserResponse(user)
	c.JSON(http.StatusOK, userResponse)
}
