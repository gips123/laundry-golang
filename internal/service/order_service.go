package service

import (
	"errors"
	"laundry-go/internal/models"
	"laundry-go/internal/repository"
	"time"

	"github.com/google/uuid"
)

type OrderService interface {
	Create(userID string, req CreateOrderRequest) (*OrderResponse, error)
	GetByUserID(userID, status string, page, limit int) (*OrderListResponse, error)
	GetByID(userID, orderID string) (*OrderResponse, error)
	CancelOrder(userID, orderID string) (*OrderResponse, error)
	UpdateStatus(laundryOwnerID, orderID string, status string) (*OrderResponse, error)
}

type orderService struct {
	orderRepo        repository.OrderRepository
	orderServiceRepo repository.OrderServiceRepository
	serviceRepo      repository.ServiceRepository
	laundryRepo      repository.LaundryRepository
}

type CreateOrderRequest struct {
	LaundryID        string                `json:"laundry_id"`
	Services         []OrderServiceRequest `json:"services"`
	DeliveryAddress  string                `json:"delivery_address"`
	Notes            string                `json:"notes"`
	EstimatedPickupAt *time.Time          `json:"estimated_pickup_at"`
}

type OrderServiceRequest struct {
	ServiceID string  `json:"service_id"`
	Quantity  float64 `json:"quantity"`
}

type OrderListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	Pagination Pagination      `json:"pagination"`
}

type OrderResponse struct {
	ID                 string                `json:"id"`
	LaundryID          string                `json:"laundry_id"`
	LaundryName        string                `json:"laundry_name"`
	Services           []OrderServiceDetail  `json:"services"`
	TotalPrice         float64               `json:"total_price"`
	Status             string                `json:"status"`
	CreatedAt          time.Time             `json:"created_at"`
	EstimatedPickup    *time.Time            `json:"estimated_pickup"`
	EstimatedDelivery  *time.Time            `json:"estimated_delivery"`
	Address            string                `json:"address"`
	Notes              string                `json:"notes"`
}

type OrderServiceDetail struct {
	ServiceID   string  `json:"service_id"`
	ServiceName string  `json:"service_name"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	Unit        string  `json:"unit"`
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	orderServiceRepo repository.OrderServiceRepository,
	serviceRepo repository.ServiceRepository,
	laundryRepo repository.LaundryRepository,
) OrderService {
	return &orderService{
		orderRepo:        orderRepo,
		orderServiceRepo: orderServiceRepo,
		serviceRepo:      serviceRepo,
		laundryRepo:      laundryRepo,
	}
}

func (s *orderService) Create(userID string, req CreateOrderRequest) (*OrderResponse, error) {
	// Validation
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	laundryUUID, err := uuid.Parse(req.LaundryID)
	if err != nil {
		return nil, errors.New("invalid laundry ID")
	}

	if len(req.Services) == 0 {
		return nil, errors.New("at least one service is required")
	}

	if req.DeliveryAddress == "" {
		return nil, errors.New("delivery address is required")
	}

	// Verify laundry exists
	laundry, err := s.laundryRepo.FindByID(laundryUUID)
	if err != nil {
		return nil, errors.New("laundry not found")
	}

	// Calculate total price and create order services
	var totalPrice float64
	orderServices := make([]models.OrderService, 0, len(req.Services))
	maxEstimatedHours := 0

	for _, svcReq := range req.Services {
		serviceUUID, err := uuid.Parse(svcReq.ServiceID)
		if err != nil {
			return nil, errors.New("invalid service ID")
		}

		service, err := s.serviceRepo.FindByID(serviceUUID)
		if err != nil {
			return nil, errors.New("service not found")
		}

		if !service.IsActive {
			return nil, errors.New("service is not active")
		}

		if service.LaundryID != laundryUUID {
			return nil, errors.New("service does not belong to this laundry")
		}

		subtotal := service.Price * svcReq.Quantity
		totalPrice += subtotal

		if service.EstimatedTimeHours > maxEstimatedHours {
			maxEstimatedHours = service.EstimatedTimeHours
		}

		orderServices = append(orderServices, models.OrderService{
			ServiceID:   service.ID,
			ServiceName: service.Name,
			Quantity:    svcReq.Quantity,
			UnitPrice:   service.Price,
			Unit:        service.Unit,
			Subtotal:    subtotal,
		})
	}

	// Calculate estimated delivery time
	var estimatedDeliveryAt *time.Time
	if req.EstimatedPickupAt != nil {
		deliveryTime := req.EstimatedPickupAt.Add(time.Duration(maxEstimatedHours) * time.Hour)
		estimatedDeliveryAt = &deliveryTime
	}

	// Create order
	order := &models.Order{
		UserID:            userUUID,
		LaundryID:         laundryUUID,
		Status:            "pending",
		TotalPrice:        totalPrice,
		DeliveryAddress:   req.DeliveryAddress,
		Notes:             req.Notes,
		EstimatedPickupAt: req.EstimatedPickupAt,
		EstimatedDeliveryAt: estimatedDeliveryAt,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, errors.New("failed to create order")
	}

	// Create order services
	for i := range orderServices {
		orderServices[i].OrderID = order.ID
	}

	if err := s.orderServiceRepo.CreateBatch(orderServices); err != nil {
		return nil, errors.New("failed to create order services")
	}

	// Convert to response
	serviceDetails := make([]OrderServiceDetail, 0, len(orderServices))
	for _, os := range orderServices {
		serviceDetails = append(serviceDetails, OrderServiceDetail{
			ServiceID:   os.ServiceID.String(),
			ServiceName: os.ServiceName,
			Quantity:    os.Quantity,
			Price:       os.UnitPrice,
			Unit:        os.Unit,
		})
	}

	return &OrderResponse{
		ID:                order.ID.String(),
		LaundryID:         order.LaundryID.String(),
		LaundryName:       laundry.Name,
		Services:          serviceDetails,
		TotalPrice:        order.TotalPrice,
		Status:            order.Status,
		CreatedAt:         order.CreatedAt,
		EstimatedPickup:   order.EstimatedPickupAt,
		EstimatedDelivery: order.EstimatedDeliveryAt,
		Address:           order.DeliveryAddress,
		Notes:             order.Notes,
	}, nil
}

func (s *orderService) GetByUserID(userID, status string, page, limit int) (*OrderListResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	orders, total, err := s.orderRepo.FindByUserID(userUUID, status, page, limit)
	if err != nil {
		return nil, errors.New("failed to fetch orders")
	}

	orderResponses := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, *s.toOrderResponse(&order))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &OrderListResponse{
		Orders: orderResponses,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *orderService) GetByID(userID, orderID string) (*OrderResponse, error) {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	order, err := s.orderRepo.FindByID(orderUUID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Verify ownership
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	if order.UserID != userUUID {
		return nil, errors.New("unauthorized")
	}

	return s.toOrderResponse(order), nil
}

func (s *orderService) CancelOrder(userID, orderID string) (*OrderResponse, error) {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	order, err := s.orderRepo.FindByID(orderUUID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Verify ownership
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	if order.UserID != userUUID {
		return nil, errors.New("unauthorized")
	}

	// Check if order can be cancelled
	if order.Status != "pending" && order.Status != "confirmed" {
		return nil, errors.New("order cannot be cancelled at this stage")
	}

	order.Status = "cancelled"
	if err := s.orderRepo.Update(order); err != nil {
		return nil, errors.New("failed to cancel order")
	}

	return s.toOrderResponse(order), nil
}

func (s *orderService) UpdateStatus(laundryOwnerID, orderID string, status string) (*OrderResponse, error) {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	order, err := s.orderRepo.FindByID(orderUUID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Verify laundry ownership
	laundryOwnerUUID, err := uuid.Parse(laundryOwnerID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	laundry, err := s.laundryRepo.FindByID(order.LaundryID)
	if err != nil {
		return nil, errors.New("laundry not found")
	}

	if laundry.OwnerID != laundryOwnerUUID {
		return nil, errors.New("unauthorized")
	}

	// Validate status
	validStatuses := []string{"pending", "confirmed", "picked-up", "washing", "drying", "ironing", "ready", "delivered", "completed", "cancelled"}
	valid := false
	for _, vs := range validStatuses {
		if status == vs {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errors.New("invalid status")
	}

	order.Status = status
	if err := s.orderRepo.Update(order); err != nil {
		return nil, errors.New("failed to update order status")
	}

	return s.toOrderResponse(order), nil
}

func (s *orderService) toOrderResponse(order *models.Order) *OrderResponse {
	serviceDetails := make([]OrderServiceDetail, 0, len(order.OrderServices))
	for _, os := range order.OrderServices {
		serviceDetails = append(serviceDetails, OrderServiceDetail{
			ServiceID:   os.ServiceID.String(),
			ServiceName: os.ServiceName,
			Quantity:    os.Quantity,
			Price:       os.UnitPrice,
			Unit:        os.Unit,
		})
	}

	laundryName := ""
	if order.Laundry.Name != "" {
		laundryName = order.Laundry.Name
	}

	return &OrderResponse{
		ID:                order.ID.String(),
		LaundryID:         order.LaundryID.String(),
		LaundryName:       laundryName,
		Services:          serviceDetails,
		TotalPrice:        order.TotalPrice,
		Status:            order.Status,
		CreatedAt:         order.CreatedAt,
		EstimatedPickup:   order.EstimatedPickupAt,
		EstimatedDelivery: order.EstimatedDeliveryAt,
		Address:           order.DeliveryAddress,
		Notes:             order.Notes,
	}
}

