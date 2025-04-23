package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"swift-codes-api/internal/config"
	"swift-codes-api/models"
	"swift-codes-api/repositories/interfaces"
)

type SwiftCodesHandler struct {
	cfg  config.Config
	repo interfaces.SwiftRepository
}

func NewSwiftHandler(cfg config.Config, repo interfaces.SwiftRepository) *SwiftCodesHandler {
	return &SwiftCodesHandler{
		cfg:  cfg,
		repo: repo,
	}
}

func (h *SwiftCodesHandler) GetSwiftCode(c *gin.Context) {
	code := c.Param("swift-code")

	result, err := h.repo.FindByCode(context.TODO(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "SWIFT code not found"})
		return
	}

	if result.IsHeadquarter {
		branches, err := h.repo.FindBranchesByPrefix(context.TODO(), result.SwiftPrefix)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"address":       result.Address,
				"bankName":      result.BankName,
				"countryISO2":   result.CountryISO2,
				"countryName":   result.CountryName,
				"isHeadquarter": result.IsHeadquarter,
				"swiftCode":     result.SwiftCode,
				"branches":      branches,
			})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *SwiftCodesHandler) GetSwiftCodesByCountry(c *gin.Context) {
	countryISO2 := c.Param("countryISO2code")

	if len(countryISO2) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid country code format. Must be a 2-letter ISO country code",
		})
		return
	}

	countryISO2 = strings.ToUpper(countryISO2)

	swiftCodes, countryName, err := h.repo.FindByCountryISO2(context.TODO(), countryISO2)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No SWIFT codes found for this country"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"countryISO2": countryISO2,
		"countryName": countryName,
		"swiftCodes":  swiftCodes,
	})
}

func (h *SwiftCodesHandler) AddSwiftCode(c *gin.Context) {
	var swiftCode models.SwiftCode

	if err := c.ShouldBindJSON(&swiftCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request format"})
		return
	}

	if swiftCode.SwiftCode == "" || swiftCode.BankName == "" ||
		swiftCode.CountryISO2 == "" || swiftCode.CountryName == "" ||
		swiftCode.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields"})
		return
	}

	if len(swiftCode.CountryISO2) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid country code format. Must be a 2-letter ISO country code",
		})
		return
	}

	swiftCode.CountryISO2 = strings.ToUpper(swiftCode.CountryISO2)

	swiftCodeLength := len(swiftCode.SwiftCode)
	if swiftCodeLength != 8 && swiftCodeLength != 11 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid SWIFT code format. Must be 8 or 11 characters",
		})
		return
	}

	err := h.repo.AddSwiftCode(context.TODO(), swiftCode)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add SWIFT code"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "SWIFT code added successfully"})
}
