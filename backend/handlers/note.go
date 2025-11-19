// handlers/note.go
package handlers

import (
	"os"
	"fmt"
	"notes-app-backend/config"
	"notes-app-backend/models"
	"notes-app-backend/utils"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CreateNoteRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url,omitempty"`
}

func CreateNote(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// Handle upload gambar kalau ada
	var imageURL string
	file, err := c.FormFile("image")
	if err == nil {
		// Buat nama unik
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		savePath := filepath.Join("/uploads", filename)

		// Simpan file
		if err := c.SaveFile(file, savePath); err != nil {
			return utils.JSON(c, fiber.StatusInternalServerError, "Failed to save image", nil)
		}

		// Generate URL fleksibel
        baseURL := os.Getenv("APP_URL")
        if baseURL == "" {
            baseURL = "http://localhost:8000"
        }
        imageURL = fmt.Sprintf("%s/uploads/%s", baseURL, filename)
    }

	// Parse title & content
	title := c.FormValue("title")
	content := c.FormValue("content")

	note := models.Note{
		Title:    title,
		Content:  content,
		ImageURL: imageURL,
		UserID:   userID,
	}

	if err := config.DB.Create(&note).Error; err != nil {  // TAMBAH ERROR CHECK INI
    // Cleanup file kalau DB gagal (opsional, biar rapi)
    if imageURL != "" {
        os.Remove(filepath.Join("/uploads", filepath.Base(imageURL)))  // Hapus file sementara
    }
    return utils.JSON(c, fiber.StatusInternalServerError, "Failed to save note to database: "+err.Error(), nil)
	}

	return utils.JSON(c, fiber.StatusCreated, "Note created", note)
}

func GetNotes(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var notes []models.Note
	config.DB.Where("user_id = ?", userID).Find(&notes)

	return utils.JSON(c, fiber.StatusOK, "Notes fetched", notes)
}

func GetNoteByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var note models.Note
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		return utils.JSON(c, fiber.StatusNotFound, "Note not found", nil)
	}

	return utils.JSON(c, fiber.StatusOK, "Note found", note)
}

func DeleteNote(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var note models.Note
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		return utils.JSON(c, fiber.StatusNotFound, "Note not found or not owner", nil)
	}

	config.DB.Delete(&note)
	return utils.JSON(c, fiber.StatusOK, "Note deleted", nil)
}