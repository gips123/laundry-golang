package handlers

import (
	"laundry-go/internal/service"
	"laundry-go/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LaundryHandler struct {
	laundryService service.LaundryService
}

func NewLaundryHandler(laundryService service.LaundryService) *LaundryHandler {
	return &LaundryHandler{laundryService: laundryService}
}

// GetAll handles GET /api/v1/laundries
func (h *LaundryHandler) GetAll(c *gin.Context) {
	search := c.Query("search")
	isOpenStr := c.Query("is_open")
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	sortBy := c.DefaultQuery("sort_by", "")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	var isOpen *bool
	if isOpenStr != "" {
		val := isOpenStr == "true"
		isOpen = &val
	}

	var lat, lng *float64
	if latStr != "" && lngStr != "" {
		if latVal, err := strconv.ParseFloat(latStr, 64); err == nil {
			lat = &latVal
		}
		if lngVal, err := strconv.ParseFloat(lngStr, 64); err == nil {
			lng = &lngVal
		}
	}

	// Get userID if user is logged in (for auto-use user location)
	var userID *string
	if userIDVal, exists := c.Get("user_id"); exists {
		userIDStr := userIDVal.(string)
		userID = &userIDStr
	}

	_ = sortBy // Will be handled by service based on lat/lng availability

	response, err := h.laundryService.GetAll(search, isOpen, lat, lng, userID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", response)
}

// GetByID handles GET /api/v1/laundries/:id
func (h *LaundryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	// Get lat/lng from query if provided
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	var lat, lng *float64
	if latStr != "" && lngStr != "" {
		if latVal, err := strconv.ParseFloat(latStr, 64); err == nil {
			lat = &latVal
		}
		if lngVal, err := strconv.ParseFloat(lngStr, 64); err == nil {
			lng = &lngVal
		}
	}

	response, err := h.laundryService.GetByID(id, lat, lng)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", response)
}

