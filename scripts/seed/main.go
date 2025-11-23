package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/abneribeiro/goapi/internal/config"
	"github.com/abneribeiro/goapi/internal/database"
	"github.com/abneribeiro/goapi/internal/model"
	"github.com/abneribeiro/goapi/internal/repository"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	ctx := context.Background()

	userRepo := repository.NewUserRepository(db)
	equipmentRepo := repository.NewEquipmentRepository(db)

	password, _ := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)

	usersData := []struct {
		Email    string
		Name     string
		Phone    string
		Role     model.UserRole
	}{
		{"owner@example.com", "John Owner", "+1234567890", model.RoleOwner},
		{"renter@example.com", "Jane Renter", "+0987654321", model.RoleRenter},
		{"owner2@example.com", "Bob Owner", "+1122334455", model.RoleOwner},
	}

	users := make([]*model.User, len(usersData))

	fmt.Println("Creating/fetching users...")
	for i, data := range usersData {
		existingUser, err := userRepo.GetByEmail(ctx, data.Email)
		if err == nil {
			users[i] = existingUser
			fmt.Printf("Found existing user: %s (ID: %s)\n", existingUser.Email, existingUser.ID)
			continue
		}

		user := &model.User{
			Email:        data.Email,
			PasswordHash: string(password),
			Name:         data.Name,
			Phone:        data.Phone,
			Role:         data.Role,
			Verified:     true,
		}

		if err := userRepo.Create(ctx, user); err != nil {
			fmt.Printf("Warning: Could not create user %s: %v\n", data.Email, err)
			continue
		}
		users[i] = user
		fmt.Printf("Created user: %s (ID: %s)\n", user.Email, user.ID)
	}

	priceHour := 15.0
	priceDay := 50.0
	priceWeek := 200.0
	priceDay2 := 75.0
	priceWeek2 := 350.0
	priceDay3 := 100.0
	priceDay4 := 30.0

	equipmentList := []*model.Equipment{
		{
			OwnerID:      users[0].ID,
			Name:         "Professional Camera Canon EOS R5",
			Description:  "High-end mirrorless camera perfect for professional photography and video production. Includes 24-70mm lens.",
			Category:     "Photography",
			PricePerHour: &priceHour,
			PricePerDay:  &priceDay,
			PricePerWeek: &priceWeek,
			Location:     "New York, NY",
			AutoApprove:  false,
		},
		{
			OwnerID:      users[0].ID,
			Name:         "DJI Mavic 3 Pro Drone",
			Description:  "Professional drone with 4/3 CMOS sensor. Perfect for aerial photography and videography.",
			Category:     "Drones",
			PricePerDay:  &priceDay2,
			PricePerWeek: &priceWeek2,
			Location:     "New York, NY",
			AutoApprove:  true,
		},
		{
			OwnerID:      users[2].ID,
			Name:         "Sony A7 IV Full Frame Camera",
			Description:  "Versatile full-frame camera for both photos and video. Great for content creators.",
			Category:     "Photography",
			PricePerDay:  &priceDay,
			PricePerWeek: &priceWeek,
			Location:     "Los Angeles, CA",
			AutoApprove:  false,
		},
		{
			OwnerID:      users[2].ID,
			Name:         "Professional Lighting Kit",
			Description:  "Complete lighting setup with 3 LED panels, softboxes, and stands. Ideal for studio work.",
			Category:     "Lighting",
			PricePerDay:  &priceDay4,
			Location:     "Los Angeles, CA",
			AutoApprove:  true,
		},
		{
			OwnerID:      users[0].ID,
			Name:         "MacBook Pro 16\" M3 Max",
			Description:  "Powerful laptop for video editing and 3D rendering. 64GB RAM, 2TB SSD.",
			Category:     "Computers",
			PricePerDay:  &priceDay3,
			Location:     "New York, NY",
			AutoApprove:  false,
		},
	}

	fmt.Println("\nCreating equipment...")
	for _, equipment := range equipmentList {
		if err := equipmentRepo.Create(ctx, equipment); err != nil {
			fmt.Printf("Warning: Could not create equipment %s: %v\n", equipment.Name, err)
		} else {
			fmt.Printf("Created equipment: %s (ID: %s)\n", equipment.Name, equipment.ID)
		}
	}

	reservationRepo := repository.NewReservationRepository(db)

	reservations := []*model.Reservation{
		{
			EquipmentID: equipmentList[0].ID,
			RenterID:    users[1].ID,
			StartDate:   time.Now().AddDate(0, 0, 7),
			EndDate:     time.Now().AddDate(0, 0, 10),
			Status:      model.StatusApproved,
			TotalPrice:  150.0,
		},
		{
			EquipmentID: equipmentList[1].ID,
			RenterID:    users[1].ID,
			StartDate:   time.Now().AddDate(0, 0, 14),
			EndDate:     time.Now().AddDate(0, 0, 16),
			Status:      model.StatusPending,
			TotalPrice:  150.0,
		},
	}

	fmt.Println("\nCreating reservations...")
	for _, reservation := range reservations {
		if err := reservationRepo.Create(ctx, reservation); err != nil {
			fmt.Printf("Warning: Could not create reservation: %v\n", err)
		} else {
			fmt.Printf("Created reservation: %s (Status: %s)\n", reservation.ID, reservation.Status)
		}
	}

	fmt.Println("\n=== Seed completed successfully! ===")
	fmt.Println("\nTest credentials:")
	fmt.Println("  Owner: owner@example.com / Password123")
	fmt.Println("  Renter: renter@example.com / Password123")
	fmt.Println("  Owner 2: owner2@example.com / Password123")
}
