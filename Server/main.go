package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type URL struct {
	ID          uint   `gorm:"primaryKey"`
	ShortCode   string `gorm:"uniqueIndex"`
	OriginalURL string
	ClickCount  int
	CreatedAt   time.Time
}

func initDB() {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := "localhost"
	dbPort := "5432"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}
	fmt.Println("Database connection successful!")
	db.AutoMigrate(&URL{})
}

// Generate a random short code
func generateShortCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

func createShortURL(c *gin.Context) {
	var request struct {
		OriginalURL string `json:"original_url"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	shortCode := generateShortCode()
	url := URL{ShortCode: shortCode, OriginalURL: request.OriginalURL}
	db.Create(&url)

	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:5000/" + shortCode})
}

func redirectToOriginal(c *gin.Context) {
	shortCode := c.Param("shortCode")
	var url URL
	if err := db.First(&url, "short_code = ?", shortCode).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	url.ClickCount++
	db.Save(&url)
	c.Redirect(http.StatusFound, url.OriginalURL)
}

func getStats(c *gin.Context) {
	shortCode := c.Param("shortCode")
	var url URL
	if err := db.First(&url, "short_code = ?", shortCode).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"original_url": url.OriginalURL,
		"short_code":   url.ShortCode,
		"click_count":  url.ClickCount,
		"created_at":   url.CreatedAt,
	})
}

func main() {
	initDB()
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")                                // Allow your frontend's origin
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allow HTTP methods
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")     // Allow necessary headers
		c.Header("Access-Control-Allow-Credentials", "true")                        // Allow credentials if needed

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.POST("/shorten", createShortURL)
	router.GET("/:shortCode", redirectToOriginal)
	router.GET("/stats/:shortCode", getStats)

	router.Run(":5000")
}
