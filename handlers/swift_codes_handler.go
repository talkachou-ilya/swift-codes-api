package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"swift-codes-api/internal/config"
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
