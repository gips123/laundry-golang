package main

import (
	"fmt"
	"log"
	"laundry-go/internal/config"
	"laundry-go/internal/database"
	"laundry-go/internal/models"
	"laundry-go/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatalf("Transaction rolled back: %v", r)
		}
	}()

	// Seed Users
	users := seedUsers(tx)
	fmt.Println("✓ Seeded users")

	// Seed Laundries
	laundries := seedLaundries(tx, users)
	fmt.Println("✓ Seeded laundries")

	// Seed Services
	services := seedServices(tx, laundries)
	fmt.Println("✓ Seeded services")

	// Seed Orders
	seedOrders(tx, users, laundries, services)
	fmt.Println("✓ Seeded orders")

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("\n✅ All dummy data seeded successfully!")
}

func seedUsers(tx *gorm.DB) map[string]uuid.UUID {
	userMap := make(map[string]uuid.UUID)

	// User 1: John Doe
	passwordHash1, _ := utils.HashPassword("password123")
	user1 := &models.User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "user@laundryhub.com",
		PasswordHash: passwordHash1,
		Phone:        "081234567890",
		Address:      "Jl. Sudirman No. 123, Jakarta Pusat",
		Role:         "customer",
	}
	tx.Create(user1)
	userMap["user-1"] = user1.ID

	// User 2: Admin User
	passwordHash2, _ := utils.HashPassword("admin123")
	user2 := &models.User{
		ID:           uuid.New(),
		Name:         "Admin User",
		Email:        "admin@laundryhub.com",
		PasswordHash: passwordHash2,
		Phone:        "081987654321",
		Address:      "Jl. Thamrin No. 45, Jakarta Pusat",
		Role:         "laundry_owner",
	}
	tx.Create(user2)
	userMap["user-2"] = user2.ID

	// User 3: Test User
	passwordHash3, _ := utils.HashPassword("test123")
	user3 := &models.User{
		ID:           uuid.New(),
		Name:         "Test User",
		Email:        "test@test.com",
		PasswordHash: passwordHash3,
		Phone:        "081111111111",
		Address:      "Jl. Test No. 1, Jakarta",
		Role:         "customer",
	}
	tx.Create(user3)
	userMap["user-3"] = user3.ID

	return userMap
}

func seedLaundries(tx *gorm.DB, users map[string]uuid.UUID) map[string]uuid.UUID {
	laundryMap := make(map[string]uuid.UUID)
	ownerID := users["user-2"] // Admin User owns all laundries

	laundries := []struct {
		id          string
		name        string
		description string
		address     string
		rating      float64
		reviewCount int
		imageURL    string
		isOpen      bool
		openTime    string
		closeTime   string
		lat         float64
		lng         float64
	}{
		{"l1", "Laundry Express Jakarta", "Laundry cepat dan berkualitas dengan layanan pick-up & delivery gratis", "Jl. Sudirman No. 123, Jakarta Pusat", 4.8, 234, "https://images.unsplash.com/photo-1581578731548-c64695cc6952?w=800", true, "08:00", "20:00", -6.2088, 106.8456},
		{"l2", "Clean & Fresh Laundry", "Spesialis dry clean dan cuci premium dengan teknologi terbaru", "Jl. Thamrin No. 45, Jakarta Pusat", 4.9, 189, "https://images.unsplash.com/photo-1628177142898-93e36e4e3a50?w=800", true, "07:00", "21:00", -6.1944, 106.8229},
		{"l3", "Quick Wash Laundry", "Layanan express 6 jam dengan harga terjangkau", "Jl. Gatot Subroto No. 78, Jakarta Selatan", 4.6, 156, "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800", true, "08:00", "22:00", -6.2297, 106.7994},
		{"l4", "Premium Laundry Service", "Layanan premium dengan perawatan khusus untuk pakaian mahal", "Jl. Kemang Raya No. 12, Jakarta Selatan", 4.7, 98, "https://images.unsplash.com/photo-1582735689369-4fe89db7114c?w=800", false, "09:00", "18:00", -6.2603, 106.8106},
		{"l5", "Eco Laundry", "Laundry ramah lingkungan dengan detergen organik", "Jl. Kebayoran Baru No. 56, Jakarta Selatan", 4.5, 201, "https://images.unsplash.com/photo-1600585154340-be6161a56a0c?w=800", true, "07:30", "19:30", -6.2442, 106.7996},
		{"l6", "24/7 Laundry", "Buka 24 jam untuk kenyamanan Anda", "Jl. Senopati No. 34, Jakarta Selatan", 4.4, 312, "https://images.unsplash.com/photo-1556912172-45b7abe8b7e8?w=800", true, "00:00", "23:59", -6.2456, 106.8006},
	}

	for _, l := range laundries {
		laundry := &models.Laundry{
			ID:                 uuid.New(),
			OwnerID:            ownerID,
			Name:               l.name,
			Description:        l.description,
			Address:            l.address,
			Latitude:           &l.lat,
			Longitude:          &l.lng,
			ImageURL:           l.imageURL,
			Rating:             l.rating,
			ReviewCount:        l.reviewCount,
			IsOpen:             l.isOpen,
			OperatingHoursOpen: models.TimeOnly(l.openTime),
			OperatingHoursClose: models.TimeOnly(l.closeTime),
		}
		tx.Create(laundry)
		laundryMap[l.id] = laundry.ID
	}

	return laundryMap
}

func seedServices(tx *gorm.DB, laundries map[string]uuid.UUID) map[string]uuid.UUID {
	serviceMap := make(map[string]uuid.UUID)

	// Service definitions
	services := []struct {
		id            string
		name          string
		description   string
		price         float64
		unit          string
		estimatedTime int
		category      string
		laundryIDs    []string // Which laundries have this service
	}{
		{"s1", "Cuci Reguler", "Cuci dan setrika pakaian biasa", 8000, "kg", 24, "regular", []string{"l1", "l2", "l3", "l5", "l6"}},
		{"s2", "Cuci Express", "Cuci cepat 6 jam", 12000, "kg", 6, "express", []string{"l1", "l3", "l5", "l6"}},
		{"s3", "Dry Clean", "Dry clean untuk pakaian khusus", 25000, "pcs", 48, "dry-clean", []string{"l1", "l2", "l4"}},
		{"s4", "Setrika Saja", "Hanya setrika tanpa cuci", 5000, "kg", 12, "ironing", []string{"l1", "l2", "l3", "l5", "l6"}},
		{"s5", "Cuci Karpet", "Cuci karpet ukuran kecil", 50000, "pcs", 72, "regular", []string{"l2", "l4"}},
	}

	for _, s := range services {
		for _, laundryID := range s.laundryIDs {
			service := &models.Service{
				ID:               uuid.New(),
				LaundryID:        laundries[laundryID],
				Name:             s.name,
				Description:      s.description,
				Price:            s.price,
				Unit:             s.unit,
				EstimatedTimeHours: s.estimatedTime,
				Category:         s.category,
				IsActive:         true,
			}
			tx.Create(service)
			// Store first occurrence of each service type for order mapping
			if _, exists := serviceMap[s.id]; !exists {
				serviceMap[s.id] = service.ID
			}
		}
	}

	return serviceMap
}

func seedOrders(tx *gorm.DB, users map[string]uuid.UUID, laundries map[string]uuid.UUID, services map[string]uuid.UUID) {
	user1ID := users["user-1"]

	// Order 1
	order1ID := uuid.New()
	order1 := &models.Order{
		ID:                order1ID,
		UserID:            user1ID,
		LaundryID:         laundries["l1"],
		Status:            "washing",
		TotalPrice:        34000,
		DeliveryAddress:   "Jl. Contoh No. 123, Jakarta",
		Notes:             "Mohon hati-hati dengan pakaian putih",
		EstimatedPickupAt: timePtr(time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)),
		EstimatedDeliveryAt: timePtr(time.Date(2024, 1, 16, 18, 0, 0, 0, time.UTC)),
		CreatedAt:         time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}
	tx.Create(order1)

	// Order 1 Services
	order1Service1 := &models.OrderService{
		ID:          uuid.New(),
		OrderID:     order1ID,
		ServiceID:   services["s1"],
		ServiceName: "Cuci Reguler",
		Quantity:    3,
		UnitPrice:   8000,
		Unit:        "kg",
		Subtotal:    24000,
	}
	tx.Create(order1Service1)

	order1Service2 := &models.OrderService{
		ID:          uuid.New(),
		OrderID:     order1ID,
		ServiceID:   services["s4"],
		ServiceName: "Setrika Saja",
		Quantity:    2,
		UnitPrice:   5000,
		Unit:        "kg",
		Subtotal:    10000,
	}
	tx.Create(order1Service2)

	// Order 2
	order2ID := uuid.New()
	order2 := &models.Order{
		ID:                order2ID,
		UserID:            user1ID,
		LaundryID:         laundries["l2"],
		Status:            "ready",
		TotalPrice:        24000,
		DeliveryAddress:   "Jl. Contoh No. 123, Jakarta",
		EstimatedPickupAt: timePtr(time.Date(2024, 1, 14, 9, 0, 0, 0, time.UTC)),
		EstimatedDeliveryAt: timePtr(time.Date(2024, 1, 14, 15, 0, 0, 0, time.UTC)),
		CreatedAt:         time.Date(2024, 1, 14, 8, 0, 0, 0, time.UTC),
	}
	tx.Create(order2)

	// Order 2 Services - l2 needs Cuci Express, but l2 doesn't have it in services list
	// According to dummy data, l2 should have express service, so we'll find/create it
	// For now, we'll use a service from l2 that exists (Cuci Reguler) but with express pricing
	var l2ExpressServiceID uuid.UUID
	tx.Model(&models.Service{}).
		Where("laundry_id = ? AND name = ?", laundries["l2"], "Cuci Express").
		Select("id").
		First(&l2ExpressServiceID)
	
	// If express doesn't exist for l2, create it
	if l2ExpressServiceID == uuid.Nil {
		expressService := &models.Service{
			ID:               uuid.New(),
			LaundryID:        laundries["l2"],
			Name:             "Cuci Express",
			Description:      "Cuci cepat 6 jam",
			Price:            12000,
			Unit:             "kg",
			EstimatedTimeHours: 6,
			Category:         "express",
			IsActive:         true,
		}
		tx.Create(expressService)
		l2ExpressServiceID = expressService.ID
	}

	order2Service1 := &models.OrderService{
		ID:          uuid.New(),
		OrderID:     order2ID,
		ServiceID:   l2ExpressServiceID,
		ServiceName: "Cuci Express",
		Quantity:    2,
		UnitPrice:   12000,
		Unit:        "kg",
		Subtotal:    24000,
	}
	tx.Create(order2Service1)

	// Order 3
	order3ID := uuid.New()
	order3 := &models.Order{
		ID:                order3ID,
		UserID:            user1ID,
		LaundryID:         laundries["l3"],
		Status:            "delivered",
		TotalPrice:        25000,
		DeliveryAddress:   "Jl. Contoh No. 123, Jakarta",
		EstimatedPickupAt: timePtr(time.Date(2024, 1, 13, 13, 0, 0, 0, time.UTC)),
		EstimatedDeliveryAt: timePtr(time.Date(2024, 1, 15, 16, 0, 0, 0, time.UTC)),
		CreatedAt:         time.Date(2024, 1, 13, 12, 0, 0, 0, time.UTC),
	}
	tx.Create(order3)

	// Order 3 Services - l3 needs Dry Clean, but l3 doesn't have it in services list
	// According to dummy data, we'll create it for l3
	var l3DryCleanServiceID uuid.UUID
	tx.Model(&models.Service{}).
		Where("laundry_id = ? AND name = ?", laundries["l3"], "Dry Clean").
		Select("id").
		First(&l3DryCleanServiceID)
	
	// If dry clean doesn't exist for l3, create it
	if l3DryCleanServiceID == uuid.Nil {
		dryCleanService := &models.Service{
			ID:               uuid.New(),
			LaundryID:        laundries["l3"],
			Name:             "Dry Clean",
			Description:      "Dry clean untuk pakaian khusus",
			Price:            25000,
			Unit:             "pcs",
			EstimatedTimeHours: 48,
			Category:         "dry-clean",
			IsActive:         true,
		}
		tx.Create(dryCleanService)
		l3DryCleanServiceID = dryCleanService.ID
	}

	order3Service1 := &models.OrderService{
		ID:          uuid.New(),
		OrderID:     order3ID,
		ServiceID:   l3DryCleanServiceID,
		ServiceName: "Dry Clean",
		Quantity:    1,
		UnitPrice:   25000,
		Unit:        "pcs",
		Subtotal:    25000,
	}
	tx.Create(order3Service1)
}

func timePtr(t time.Time) *time.Time {
	return &t
}

