// handlers/auth.go
package handlers

import (
	"notes-app-backend/config"
	"notes-app-backend/middleware"
	"notes-app-backend/models"
	"notes-app-backend/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.JSON(c, fiber.StatusBadRequest, "Invalid request", nil)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Where("email = ?", req.Email).First(&models.User{}).Error; err == nil {
		return utils.JSON(c, fiber.StatusBadRequest, "Email already exists", nil)
	}

	config.DB.Create(&user)

	token, _ := middleware.GenerateJWT(user)

	return utils.JSON(c, fiber.StatusCreated, "Register success", fiber.Map{
		"token": token,
		"user":  user,
	})
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.JSON(c, fiber.StatusBadRequest, "Invalid request", nil)
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return utils.JSON(c, fiber.StatusUnauthorized, "Invalid email or password", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return utils.JSON(c, fiber.StatusUnauthorized, "Invalid email or password", nil)
	}

	token, _ := middleware.GenerateJWT(user)

	return utils.JSON(c, fiber.StatusOK, "Login success", fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}