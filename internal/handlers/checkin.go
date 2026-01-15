package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lanhyde/ogenkidesuka-server/internal/database"
	"github.com/lanhyde/ogenkidesuka-server/internal/models"
)

// HealthCheck returns the API health status
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"message":   "Ogenkidesuka API is running",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// CreateCheckIn creates a new wellness check-in
func CreateCheckIn(c *gin.Context) {
	// Get user ID from path parameter
	userID := c.Param("userId")

	var req models.CreateCheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Insert check-in into db
	query := `INSERT INTO check_ins (user_id, check_in_type, step_count, battery_level, checked_at)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, created_at`

	var checkIn models.CheckIn

	parsedUserID, err := parseUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	checkIn.UserID = parsedUserID
	checkIn.CheckInType = req.CheckInType
	checkIn.StepCount = req.StepCount
	checkIn.BatteryLevel = req.BatteryLevel
	checkIn.CheckedAt = time.Now()

	err = database.DB.QueryRow(
		query, checkIn.UserID, checkIn.CheckInType,
		checkIn.StepCount, checkIn.BatteryLevel,
		checkIn.CheckedAt).Scan(&checkIn.ID, &checkIn.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create check-in",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Check-in created successfully",
		"data":    checkIn,
	})
}

// Checks if user has checked in today
func GetTodayCheckIn(c *gin.Context) {
	userID := c.Param("userId")

	query := `
		SELECT id, user_id, check_in_type, step_count, battery_level, checked_at, created_at
		FROM check_ins 
		WHERE user_id = $1 AND DATE(checked_at) = CURRENT_DATE
		ORDER BY checked_at DESC
		LIMIT 1`

	var checkIn models.CheckIn
	err := database.DB.QueryRow(query, userID).Scan(
		&checkIn.ID,
		&checkIn.UserID,
		&checkIn.CheckInType,
		&checkIn.StepCount,
		&checkIn.BatteryLevel,
		&checkIn.CheckedAt,
		&checkIn.CreatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"message": "No check-in today",
			"data":    nil,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch check-in",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "check-in found",
		"data":    checkIn,
	})
}

// Get check-in history for a user
func GetCheckInHistory(c *gin.Context) {
	userID := c.Param("userId")
	limit := c.DefaultQuery("limit", "30") // last 30 days by default

	query := `
		SELECT id, user_id, check_in_type, step_count, battery_level, checked_at, created_at
		FROM check_ins
		WHERE user_id = $1
		ORDER BY checked_at DESC
		LIMIT $2`

	rows, err := database.DB.Query(query, userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch history",
		})
		return
	}
	defer rows.Close()

	var history []models.CheckIn
	for rows.Next() {
		var checkIn models.CheckIn
		err := rows.Scan(
			&checkIn.ID,
			&checkIn.UserID,
			&checkIn.CheckInType,
			&checkIn.StepCount,
			&checkIn.BatteryLevel,
			&checkIn.CheckedAt,
			&checkIn.CreatedAt,
		)
		if err != nil {
			continue
		}
		history = append(history, checkIn)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "History retrieved successfully",
		"data":    history,
		"count":   len(history),
	})
}

func parseUserID(userID string) (int, error) {
	// strconv.Atoi converts string to integer
	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %s", userID)
	}
	return id, nil
}
