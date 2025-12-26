package service

import (
	"errors"
	"fmt"
	"laundry-go/internal/repository"
	"laundry-go/internal/utils"

	"github.com/google/uuid"
)

type LaundryService interface {
	GetAll(search string, isOpen *bool, lat, lng *float64, userID *string, page, limit int) (*LaundryListResponse, error)
	GetByID(id string, lat, lng *float64) (*LaundryDetailResponse, error)
}

type laundryService struct {
	laundryRepo repository.LaundryRepository
	serviceRepo repository.ServiceRepository
	userRepo    repository.UserRepository
}

type LaundryListResponse struct {
	Laundries   []LaundryListItem `json:"laundries"`
	Pagination  Pagination        `json:"pagination"`
	UserLocation *UserLocation    `json:"user_location,omitempty"`
}

type UserLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LaundryListItem struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Address         string              `json:"address"`
	Rating          float64             `json:"rating"`
	ReviewCount     int                 `json:"review_count"`
	Image           string              `json:"image"`
	PriceRange      string              `json:"price_range"`
	Distance        *float64            `json:"distance,omitempty"`
	IsOpen          bool                `json:"is_open"`
	OperatingHours OperatingHours      `json:"operating_hours"`
}

type LaundryDetailResponse struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Address         string              `json:"address"`
	Rating          float64             `json:"rating"`
	ReviewCount     int                 `json:"review_count"`
	Image           string              `json:"image"`
	PriceRange      string              `json:"price_range"`
	Distance        *float64            `json:"distance,omitempty"`
	IsOpen          bool                `json:"is_open"`
	OperatingHours OperatingHours      `json:"operating_hours"`
	Services        []ServiceResponse   `json:"services"`
}

type OperatingHours struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

type ServiceResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Unit            string  `json:"unit"`
	EstimatedTime   int     `json:"estimated_time"`
	Category        string  `json:"category"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func NewLaundryService(laundryRepo repository.LaundryRepository, serviceRepo repository.ServiceRepository, userRepo repository.UserRepository) LaundryService {
	return &laundryService{
		laundryRepo: laundryRepo,
		serviceRepo: serviceRepo,
		userRepo:    userRepo,
	}
}

func (s *laundryService) GetAll(search string, isOpen *bool, lat, lng *float64, userID *string, page, limit int) (*LaundryListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Prioritas penggunaan lat/lng:
	// 1. Query params (lat, lng) - PRIORITAS TERTINGGI
	// 2. User profile lat/lng (jika user sudah login dan punya lokasi)
	// 3. Null/0 (jika tidak ada sama sekali)
	var finalLat, finalLng *float64
	var userLocation *UserLocation

	if lat != nil && lng != nil {
		// Gunakan query params jika ada
		finalLat = lat
		finalLng = lng
		userLocation = &UserLocation{
			Latitude:  *lat,
			Longitude: *lng,
		}
	} else if userID != nil {
		// Coba ambil dari user profile
		userUUID, err := uuid.Parse(*userID)
		if err == nil {
			user, err := s.userRepo.FindByID(userUUID)
			if err == nil && user.Latitude != nil && user.Longitude != nil {
				finalLat = user.Latitude
				finalLng = user.Longitude
				userLocation = &UserLocation{
					Latitude:  *user.Latitude,
					Longitude: *user.Longitude,
				}
			}
		}
	}

	laundries, total, err := s.laundryRepo.FindAll(search, isOpen, page, limit)
	if err != nil {
		return nil, errors.New("failed to fetch laundries")
	}

	items := make([]LaundryListItem, 0, len(laundries))
	for _, laundry := range laundries {
		// Get price range
		minPrice, maxPrice, _ := s.serviceRepo.GetPriceRange(laundry.ID)
		priceRange := formatPriceRange(minPrice, maxPrice)

		// Calculate distance if lat/lng provided
		var distance *float64
		if finalLat != nil && finalLng != nil && laundry.Latitude != nil && laundry.Longitude != nil {
			dist := utils.CalculateDistance(*finalLat, *finalLng, *laundry.Latitude, *laundry.Longitude)
			distance = &dist
		}

		items = append(items, LaundryListItem{
			ID:              laundry.ID.String(),
			Name:            laundry.Name,
			Description:     laundry.Description,
			Address:         laundry.Address,
			Rating:          laundry.Rating,
			ReviewCount:     laundry.ReviewCount,
			Image:           laundry.ImageURL,
			PriceRange:      priceRange,
			Distance:        distance,
			IsOpen:          laundry.IsOpen,
			OperatingHours: OperatingHours{
				Open:  string(laundry.OperatingHoursOpen),
				Close: string(laundry.OperatingHoursClose),
			},
		})
	}

	// Sort by distance if lat/lng tersedia, else sort by rating
	if finalLat != nil && finalLng != nil {
		// Sort by distance (nearest first)
		for i := 0; i < len(items)-1; i++ {
			for j := i + 1; j < len(items); j++ {
				if items[i].Distance != nil && items[j].Distance != nil {
					if *items[i].Distance > *items[j].Distance {
						items[i], items[j] = items[j], items[i]
					}
				} else if items[i].Distance == nil && items[j].Distance != nil {
					// Items without distance go to the end
					items[i], items[j] = items[j], items[i]
				}
			}
		}
	} else {
		// Sort by rating (highest first)
		for i := 0; i < len(items)-1; i++ {
			for j := i + 1; j < len(items); j++ {
				if items[i].Rating < items[j].Rating {
					items[i], items[j] = items[j], items[i]
				}
			}
		}
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &LaundryListResponse{
		Laundries:    items,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
		UserLocation: userLocation,
	}, nil
}

func (s *laundryService) GetByID(id string, lat, lng *float64) (*LaundryDetailResponse, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid laundry ID")
	}

	laundry, err := s.laundryRepo.FindByID(uuid)
	if err != nil {
		return nil, errors.New("laundry not found")
	}

	// Get price range
	minPrice, maxPrice, _ := s.serviceRepo.GetPriceRange(laundry.ID)
	priceRange := formatPriceRange(minPrice, maxPrice)

	// Calculate distance if lat/lng provided
	var distance *float64
	if lat != nil && lng != nil && laundry.Latitude != nil && laundry.Longitude != nil {
		dist := utils.CalculateDistance(*lat, *lng, *laundry.Latitude, *laundry.Longitude)
		distance = &dist
	}

	// Convert services
	services := make([]ServiceResponse, 0, len(laundry.Services))
	for _, service := range laundry.Services {
		services = append(services, ServiceResponse{
			ID:            service.ID.String(),
			Name:          service.Name,
			Description:   service.Description,
			Price:         service.Price,
			Unit:          service.Unit,
			EstimatedTime: service.EstimatedTimeHours,
			Category:      service.Category,
		})
	}

	return &LaundryDetailResponse{
		ID:              laundry.ID.String(),
		Name:            laundry.Name,
		Description:     laundry.Description,
		Address:         laundry.Address,
		Rating:          laundry.Rating,
		ReviewCount:     laundry.ReviewCount,
		Image:           laundry.ImageURL,
		PriceRange:      priceRange,
		Distance:        distance,
		IsOpen:          laundry.IsOpen,
		OperatingHours: OperatingHours{
			Open:  string(laundry.OperatingHoursOpen),
			Close: string(laundry.OperatingHoursClose),
		},
		Services: services,
	}, nil
}

func formatPriceRange(minPrice, maxPrice float64) string {
	if minPrice == 0 && maxPrice == 0 {
		return "Rp 0"
	}
	if minPrice == maxPrice {
		return fmt.Sprintf("Rp %.0f", minPrice)
	}
	return fmt.Sprintf("Rp %.0f - Rp %.0f", minPrice, maxPrice)
}

