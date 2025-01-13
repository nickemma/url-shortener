package main

import (
	"fmt"
	"math/rand"
	"net/http"
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
	dsn := "host=localhost user=postgres password=tech1234 dbname=url_shortener port=5432 sslmode=disable"
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
	rand.Seed(time.Now().UnixNano())
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

	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8080/" + shortCode})
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
	router.POST("/shorten", createShortURL)
	router.GET("/:shortCode", redirectToOriginal)
	router.GET("/stats/:shortCode", getStats)

	router.Run(":8080")
}
