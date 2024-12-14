package seeder

// import (
// 	"log"
// 	"time"

// 	"gorm.io/gorm"
// )

// type Seeder interface {
// 	Seed(db *gorm.DB) error
// }

// func RunSeeders(db *gorm.DB, seeders ...Seeder) {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			log.Fatalf("Transaction rolled back due to panic: %v", r)
// 		}
// 	}()

// 	for _, seeder := range seeders {
// 		if err := seeder.Seed(tx); err != nil {
// 			tx.Rollback()
// 			log.Fatalf("Error running seeder: %v", err)
// 			return
// 		}
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		log.Fatalf("Transaction commit failed: %v", err)
// 	}
// 	log.Println("Seeding completed successfully.")
// }

// // Generalized Seeder Structs
// func dataSeeds() []interface{} {
// 	return []interface{}{
// 		VoucherSeeder{},
// 		UserSeeder{},
// 		// Add other seeder structs here
// 	}
// }

// // VoucherSeeder seeds voucher data
// type VoucherSeeder struct{}

// func (v VoucherSeeder) Seed(db *gorm.DB) error {
// 	vouchers := []models.Voucher{
// 		{
// 			VoucherName:     "10% Discount",
// 			VoucherCode:     "DISCOUNT10",
// 			VoucherType:     "e-commerce",
// 			PointsRequired:  0,
// 			Description:     "10% off for purchases above $100",
// 			VoucherCategory: "Discount",
// 			DiscountValue:   10.0,
// 			MinimumPurchase: 100.0,
// 			PaymentMethods:  []string{"Credit Card", "PayPal"},
// 			StartDate:       time.Now().AddDate(0, 0, -5),
// 			EndDate:         time.Now().AddDate(0, 0, -1),
// 			ApplicableAreas: []string{"US", "Canada"},
// 			Quota:           100,
// 		},
// 		// Add other vouchers here...
// 	}

// 	if err := insertData(db, vouchers); err != nil {
// 		log.Printf("Failed to insert vouchers: %v", err)
// 		return err
// 	}
// 	log.Println("Vouchers seeded successfully.")
// 	return nil
// }

// // UserSeeder seeds user data
// type UserSeeder struct{}

// func (u UserSeeder) Seed(db *gorm.DB) error {
// 	users := []models.User{
// 		{
// 			Name:     "Admin User",
// 			Email:    "admin@example.com",
// 			Password: helper.HashPassword("adminpassword"),
// 			Role:     "admin",
// 		},
// 		{
// 			Name:     "Staff User 1",
// 			Email:    "staff1@example.com",
// 			Password: helper.HashPassword("staffpassword1"),
// 			Role:     "staff",
// 		},
// 		// Add other users here...
// 	}

// 	if err := insertData(db, users); err != nil {
// 		log.Printf("Failed to insert users: %v", err)
// 		return err
// 	}
// 	log.Println("Users seeded successfully.")
// 	return nil
// }

// // insertData inserts any slice of structs into the database
// func insertData(db *gorm.DB, data interface{}) error {
// 	return db.Create(data).Error
// }

// func main() {
// 	db := getDatabaseConnection() // Replace with your database connection logic

// 	seeders := dataSeeds()
// 	for _, seeder := range seeders {
// 		if s, ok := seeder.(Seeder); ok {
// 			RunSeeders(db, s)
// 		}
// 	}
// }

// func getDatabaseConnection() *gorm.DB {
// 	// Replace with actual database connection logic
// 	return &gorm.DB{}
// }
