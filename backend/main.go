// main.go
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"notes-app-backend/config"
	"notes-app-backend/handlers"
	"notes-app-backend/middleware"
	"notes-app-backend/models"
)

func main() {
	godotenv.Load()

	config.ConnectDatabase()

	// Auto Migrate (termasuk table logs!)
	config.DB.AutoMigrate(&models.User{}, &models.Note{}, &models.Log{})
	
	if err := os.MkdirAll("/uploads", 0755); err != nil {
        log.Printf("Warning: Failed to create /uploads: %v", err)
    }
	app := fiber.New()
                 // endpoint upload
	app.Static("/uploads", "/uploads")

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Global middleware: Logger dulu, baru JWT di route yang butuh
	app.Use(middleware.Logger())

	// Public routes
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	// Protected routes
	protected := app.Group("/", middleware.JWTProtected())
	{
		protected.Post("/notes", handlers.CreateNote)
		protected.Get("/notes", handlers.GetNotes)
		protected.Get("/notes/:id", handlers.GetNoteByID)
		protected.Delete("/notes/:id", handlers.DeleteNote)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Backend running on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}