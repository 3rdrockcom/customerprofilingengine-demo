package main

import (
	"strconv"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo/models"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/customer/:id", getCustomerAnalysis)
	r.GET("/customer/:id/history", getCustomerAnalysisHistory)

	return r
}

func getCustomerAnalysis(c *gin.Context) {
	var err error

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	analyzer := models.Analyzer{}
	err = db.Where("customers_id = ?", id).Order("date desc").Limit(1).Find(&analyzer).Error
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, analyzer)
}

func getCustomerAnalysisHistory(c *gin.Context) {
	var err error

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	customer := models.Customers{}
	err = db.Where("id = ?", id).Find(&customer).Error
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// loc, _ := time.LoadLocation("America/Los_Angeles")

	query := db.Where("customers_id = ?", id)

	t1, err := time.ParseInLocation(
		"20060102",
		startDate, time.UTC)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	query = query.Where("date >= ?", t1)

	t2, err := time.ParseInLocation(
		"20060102",
		endDate, time.UTC)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	query = query.Where("date < ?", t2)

	analysisHistory := []models.Analyzer{}
	err = query.Order("date asc").Find(&analysisHistory).Error
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return

	}

	c.JSON(200, analysisHistory)
}
