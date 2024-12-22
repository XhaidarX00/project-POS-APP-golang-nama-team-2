package revenuerepository_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"project_pos_app/helper"
	"project_pos_app/model"
	revenuerepository "project_pos_app/repository/revenue_repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestFindLowStockProducts(t *testing.T) {
	db, mock := helper.SetupTestDB()

	logger := zap.NewNop()
	repo := revenuerepository.NewRevenueRepository(db, logger)

	t.Run("Successfully find low stock products", func(t *testing.T) {
		mockRows := sqlmock.NewRows([]string{"id", "name", "qty"}).
			AddRow(1, "Product A", 3).
			AddRow(2, "Product B", 2)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE qty < $1`)).
			WithArgs(5).
			WillReturnRows(mockRows)

		products, err := repo.FindLowStockProducts(5)

		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.Equal(t, "Product A", products[0].Name)
	})

	t.Run("Fail to find low stock products", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE qty < $1`)).
			WithArgs(5).
			WillReturnError(fmt.Errorf("database error"))

		products, err := repo.FindLowStockProducts(5)

		assert.Error(t, err)
		assert.Nil(t, products)
	})
}

func TestGetProductRevenues(t *testing.T) {
	db, mock := helper.SetupTestDB()

	logger := zap.NewNop()
	repo := revenuerepository.NewRevenueRepository(db, logger)

	t.Run("Successfully get product revenues", func(t *testing.T) {
		mockRows := sqlmock.NewRows([]string{"product_name", "sell_price", "total_revenue", "profit_margin", "revenue_date"}).
			AddRow("Product A", 100.0, 2000.0, 15.0, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "product_revenues" ORDER BY total_revenue DESC`)).
			WillReturnRows(mockRows)

		revenues, err := repo.GetProductRevenues()

		assert.NoError(t, err)
		assert.Len(t, revenues, 1)
		assert.Equal(t, "Product A", revenues[0].ProductName)
	})

	t.Run("Fail to get product revenues", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM product_revenues ORDER BY total_revenue DESC`)).
			WillReturnError(fmt.Errorf("database error"))

		revenues, err := repo.GetProductRevenues()

		assert.Error(t, err)
		assert.Nil(t, revenues)
	})
}

func TestGetTotalRevenueByStatus(t *testing.T) {
	db, mock := helper.SetupTestDB()

	logger := zap.NewNop()
	repo := revenuerepository.NewRevenueRepository(db, logger)

	t.Run("Successfully get total revenue by status", func(t *testing.T) {
		mockRows := sqlmock.NewRows([]string{"status", "revenue"}).
			AddRow("confirmed", 5000.0).
			AddRow("pending", 2000.0)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT status, SUM(revenue) as revenue FROM "order_revenues" GROUP BY "status"`)).
			WillReturnRows(mockRows)

		revenueMap, err := repo.GetTotalRevenueByStatus()

		assert.NoError(t, err)
		assert.Equal(t, 5000.0, revenueMap["confirmed"])
		assert.Equal(t, 2000.0, revenueMap["pending"])
	})

	t.Run("Fail to get total revenue by status", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT status, SUM(revenue) as revenue FROM "order_revenues" GROUP BY status`)).
			WillReturnError(fmt.Errorf("database error"))

		revenueMap, err := repo.GetTotalRevenueByStatus()

		assert.Error(t, err)
		assert.Nil(t, revenueMap)
	})
}

func TestSaveOrderRevenue(t *testing.T) {
	db, mock := helper.SetupTestDB()

	logger := zap.NewNop()
	repo := revenuerepository.NewRevenueRepository(db, logger)

	t.Run("Successfully update existing order revenue", func(t *testing.T) {
		revenueDate := time.Now()
		order := model.OrderRevenue{
			ID:        1,
			Status:    "confirmed",
			Revenue:   1000.0,
			ProductID: 2,
			CreatedAt: revenueDate,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "order_revenues" WHERE id = $1 AND cerate_at = $2 ORDER BY "order_revenues"."id" LIMIT $3`)).
			WithArgs(order.ID, order.CreatedAt, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "status", "revenue", "product_id", "created_at"}).
				AddRow(order.ID, order.Status, order.Revenue, order.ProductID, order.CreatedAt))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "order_revenues" SET`)).
			WithArgs(
				order.ID,
				order.Status,
				order.Revenue,
				order.CreatedAt,
				order.ProductID,
				order.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.SaveOrderRevenue(order)
		assert.NoError(t, err)
	})

	// t.Run("Successfully create new order revenue", func(t *testing.T) {
	// 	revenueDate := time.Now()
	// 	order := model.OrderRevenue{
	// 		ID:        1,
	// 		Status:    "confirmed",
	// 		Revenue:   1000.0,
	// 		ProductID: 2,
	// 		CreatedAt: revenueDate,
	// 	}

	// 	mock.ExpectBegin()
	// 	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "order_revenues" WHERE id = $1 AND cerate_at = $2 ORDER BY "order_revenues"."id" LIMIT $3`)).
	// 		WithArgs(order.ID, order.CreatedAt, 1).
	// 		WillReturnError(gorm.ErrRecordNotFound)

	// 	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "order_revenues"`)).
	// 		WithArgs(
	// 			order.Status,
	// 			order.Revenue,
	// 			order.CreatedAt,
	// 			order.ProductID,
	// 			order.ID,
	// 		).
	// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(order.ID))

	// 	mock.ExpectCommit()

	// 	err := repo.SaveOrderRevenue(order)
	// 	assert.NoError(t, err)
	// })

	t.Run("Fail due to empty status", func(t *testing.T) {
		order := model.OrderRevenue{
			ID:        1,
			Status:    "", // Invalid status
			Revenue:   1000.0,
			CreatedAt: time.Now(),
		}

		err := repo.SaveOrderRevenue(order)
		assert.Error(t, err)
		assert.EqualError(t, err, "order status cannot be empty")
	})

	t.Run("Fail due to negative revenue", func(t *testing.T) {
		order := model.OrderRevenue{
			ID:        1,
			Status:    "confirmed",
			Revenue:   -500.0, // Invalid revenue
			CreatedAt: time.Now(),
		}

		err := repo.SaveOrderRevenue(order)
		assert.Error(t, err)
		assert.EqualError(t, err, "revenue cannot be negative")
	})

	t.Run("Fail due to empty created_at", func(t *testing.T) {
		order := model.OrderRevenue{
			ID:        1,
			Status:    "confirmed",
			Revenue:   1000.0,
			CreatedAt: time.Time{}, // Invalid created_at
		}

		err := repo.SaveOrderRevenue(order)
		assert.Error(t, err)
		assert.EqualError(t, err, "created_at cannot be empty")
	})

	t.Run("Fail due to database error", func(t *testing.T) {
		revenueDate := time.Now()
		order := model.OrderRevenue{
			ID:        1,
			Status:    "confirmed",
			Revenue:   1000.0,
			CreatedAt: revenueDate,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "order_revenues" WHERE id = $1 AND cerate_at = $2 ORDER BY "order_revenues"."id" LIMIT $3`)).
			WithArgs(order.ID, order.CreatedAt, 1).
			WillReturnError(errors.New("database error"))

		mock.ExpectRollback()

		err := repo.SaveOrderRevenue(order)
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}

func TestCalculateOrderRevenue(t *testing.T) {
	db, mock := helper.SetupTestDB()

	logger := zap.NewNop()
	repo := revenuerepository.NewRevenueRepository(db, logger)

	t.Run("Successfully calculate order revenue", func(t *testing.T) {
		mockRows := sqlmock.NewRows([]string{"status", "revenue", "created_at"}).
			AddRow("confirmed", 5000.0, time.Now()).
			AddRow("pending", 2000.0, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT orders.status AS status, SUM(orders.total_amount) AS revenue, CURRENT_DATE AS created_at FROM "orders" LEFT JOIN order_products ON orders.id = order_products.order_id GROUP BY "orders"."status"`)).
			WillReturnRows(mockRows)

		orders, err := repo.CalculateOrderRevenue()

		assert.NoError(t, err)
		assert.Len(t, orders, 2)
		assert.Equal(t, "confirmed", orders[0].Status)
	})

	t.Run("Fail to calculate order revenue due to database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT status, SUM(revenue) as revenue, CURRENT_DATE as created_at FROM "order_revenues" GROUP BY status`)).
			WillReturnError(errors.New("database error"))

		orders, err := repo.CalculateOrderRevenue()

		assert.Error(t, err)
		assert.Nil(t, orders)
	})
}

func TestCalculateProductRevenue(t *testing.T) {
	db, mock := helper.SetupTestDB()

	logger := zap.NewNop()
	repo := revenuerepository.NewRevenueRepository(db, logger)

	t.Run("Successfully calculate product revenue", func(t *testing.T) {
		mockRows := sqlmock.NewRows([]string{"product_name", "sell_price", "total_revenue", "profit_margin", "revenue_date"}).
			AddRow("Product A", 100.0, 2000.0, 15.0, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT products.name AS product_name, products.price AS sell_price, SUM(order_products.qty * products.price) AS total_revenue, CURRENT_DATE AS revenue_date FROM "products" JOIN order_products ON products.id = order_products.product_id JOIN orders ON order_products.order_id = orders.id WHERE orders.status = $1 GROUP BY products.name, products.price`)).
			WillReturnRows(mockRows)

		products, err := repo.CalculateProductRevenue()

		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.Equal(t, "Product A", products[0].ProductName)
	})

	t.Run("Fail to calculate product revenue due to database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT products.name AS product_name, products.price AS sell_price, SUM(order_products.qty * products.price) AS total_revenue, 15.00 AS profit_margin, CURRENT_DATE AS revenue_date FROM "products"`)).
			WillReturnError(errors.New("database error"))

		products, err := repo.CalculateProductRevenue()

		assert.Error(t, err)
		assert.Nil(t, products)
	})
}

func TestSaveProductRevenue(t *testing.T) {
	db, mock := helper.SetupTestDB()

	logger := zap.NewNop()
	repo := revenuerepository.NewRevenueRepository(db, logger)

	t.Run("Successfully save product revenue", func(t *testing.T) {
		// Gunakan waktu yang tetap untuk konsistensi dalam pengujian
		revenueDate := time.Date(2024, 12, 18, 0, 0, 0, 0, time.UTC)

		productRevenue := model.ProductRevenue{
			ProductName:  "Product A",
			SellPrice:    100,
			Profit:       1900,
			ProfitMargin: 15,
			TotalRevenue: 2000,
			RevenueDate:  revenueDate,
		}

		// Mock untuk skenario record tidak ditemukan
		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT \* FROM "product_revenues" WHERE product_name = \$1 AND revenue_date = \$2 ORDER BY "product_revenues"."id" LIMIT \$3`).
			WithArgs(productRevenue.ProductName, productRevenue.RevenueDate, 1).
			WillReturnError(gorm.ErrRecordNotFound) // Simulasi record tidak ditemukan

		// Ekspektasi insert untuk record baru
		// Dalam unit test
		mock.ExpectQuery(`INSERT INTO "product_revenues"`).
			WithArgs(
				productRevenue.ProductName,
				productRevenue.SellPrice,
				productRevenue.Profit,
				productRevenue.ProfitMargin,
				productRevenue.TotalRevenue,
				productRevenue.RevenueDate,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		// Jalankan fungsi SaveProductRevenue
		err := repo.SaveProductRevenue(productRevenue)
		assert.NoError(t, err)

		// Verifikasi bahwa semua ekspektasi dipenuhi
		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Update existing product revenue", func(t *testing.T) {
		revenueDate := time.Date(2024, 12, 18, 0, 0, 0, 0, time.UTC)

		existingProductRevenue := model.ProductRevenue{
			ID:           1,
			ProductName:  "Product A",
			SellPrice:    80,
			Profit:       1500,
			ProfitMargin: 12,
			TotalRevenue: 1800,
			RevenueDate:  revenueDate,
		}

		updatedProductRevenue := model.ProductRevenue{
			ProductName:  "Product A",
			SellPrice:    100,
			Profit:       1900,
			ProfitMargin: 15,
			TotalRevenue: 2000,
			RevenueDate:  revenueDate,
		}

		// Mock untuk skenario record sudah ada
		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT \* FROM "product_revenues" WHERE product_name = \$1 AND revenue_date = \$2 ORDER BY "product_revenues"."id" LIMIT \$3`).
			WithArgs(updatedProductRevenue.ProductName, updatedProductRevenue.RevenueDate, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product_name", "sell_price", "profit", "profit_margin", "total_revenue", "revenue_date"}).
				AddRow(existingProductRevenue.ID, existingProductRevenue.ProductName, existingProductRevenue.SellPrice,
					existingProductRevenue.Profit, existingProductRevenue.ProfitMargin, existingProductRevenue.TotalRevenue, existingProductRevenue.RevenueDate))

		// Ekspektasi update
		mock.ExpectExec(`UPDATE "product_revenues" SET`).
			WithArgs(
				updatedProductRevenue.ProductName,
				updatedProductRevenue.SellPrice,
				updatedProductRevenue.Profit,
				updatedProductRevenue.ProfitMargin,
				updatedProductRevenue.TotalRevenue,
				updatedProductRevenue.RevenueDate,
				existingProductRevenue.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Jalankan fungsi SaveProductRevenue
		err := repo.SaveProductRevenue(updatedProductRevenue)
		assert.NoError(t, err)

		// Verifikasi bahwa semua ekspektasi dipenuhi
		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	// Test case untuk validasi input
	t.Run("Validation - Empty Product Name", func(t *testing.T) {
		productRevenue := model.ProductRevenue{
			ProductName: "",
			SellPrice:   100,
			RevenueDate: time.Now(),
		}

		err := repo.SaveProductRevenue(productRevenue)
		assert.Error(t, err)
		assert.EqualError(t, err, "product name cannot be empty")
	})

	t.Run("Validation - Invalid Sell Price", func(t *testing.T) {
		productRevenue := model.ProductRevenue{
			ProductName: "Product A",
			SellPrice:   0,
			RevenueDate: time.Now(),
		}

		err := repo.SaveProductRevenue(productRevenue)
		assert.Error(t, err)
		assert.EqualError(t, err, "sell price must be positive")
	})

	t.Run("Validation - Empty Revenue Date", func(t *testing.T) {
		productRevenue := model.ProductRevenue{
			ProductName: "Product A",
			SellPrice:   100,
			RevenueDate: time.Time{},
		}

		err := repo.SaveProductRevenue(productRevenue)
		assert.Error(t, err)
		assert.EqualError(t, err, "revenue date cannot be empty")
	})

}
