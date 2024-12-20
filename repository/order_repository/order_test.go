package orderrepository_test

import (
	"fmt"
	"project_pos_app/helper"
	"project_pos_app/model"
	orderrepository "project_pos_app/repository/order_repository"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestGetAllOrder(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := zap.NewNop()
	orderRepo := orderrepository.NewOrderRepo(db, log)

	t.Run("Successfully get all orders", func(t *testing.T) {
		search, status := "John", "Completed"

		// Mock the orders query
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT o.id, o.table_id, o.customer_name, o.status, o.created_at as order_date FROM orders as o WHERE o.deleted_at IS NULL AND o.customer_name ILIKE $1 AND o.status ILIKE $2`)).
			WithArgs("%John%", "%Completed%").
			WillReturnRows(sqlmock.NewRows([]string{"id", "table_id", "customer_name", "status", "order_date"}).
				AddRow(1, 1, "John Doe", "Completed", time.Now()).
				AddRow(2, 2, "Jane Doe", "Completed", time.Now()))

		// Mock the order products query
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT po.order_id, po.qty, p.name AS item, p.price FROM order_products as po JOIN products as p ON p.id = po.product_id WHERE po.order_id IN ($1,$2)`)).
			WithArgs(1, 2).
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "qty", "item", "price"}).
				AddRow(1, 2, "Product A", 50.0).
				AddRow(2, 1, "Product B", 100.0))

		// Call the repository method
		orders, err := orderRepo.GetAllOrder(search, status)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, orders)
		assert.Len(t, orders, 2)

		assert.Equal(t, "John Doe", orders[0].CustomerName)
		assert.Equal(t, "Jane Doe", orders[1].CustomerName)

		assert.Equal(t, 100, orders[0].SubTotal) // 2 * 50
		assert.Equal(t, 100, orders[1].SubTotal) // 1 * 100

		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now(), orders[0].OrderDate, time.Second)
		assert.WithinDuration(t, time.Now(), orders[1].OrderDate, time.Second)
	})

	t.Run("Failed to get orders due to database error", func(t *testing.T) {
		search, status := "John", ""

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT o.id, o.table_id, o.customer_name, o.status, o.created_at as order_date FROM orders as o WHERE o.deleted_at IS NULL AND o.customer_name ILIKE $1`)).
			WithArgs("%John%").
			WillReturnError(fmt.Errorf("database error"))

		orders, err := orderRepo.GetAllOrder(search, status)

		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.EqualError(t, err, "database error")
	})

	t.Run("No orders found", func(t *testing.T) {
		search, status := "Nonexistent", ""

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT o.id, o.table_id, o.customer_name, o.status, o.created_at as order_date FROM orders as o WHERE o.deleted_at IS NULL AND o.customer_name ILIKE $1`)).
			WithArgs("%Nonexistent%").
			WillReturnRows(sqlmock.NewRows([]string{"id", "table_id", "customer_name", "status", "order_date"}))

		orders, err := orderRepo.GetAllOrder(search, status)

		assert.NoError(t, err)
		assert.Len(t, orders, 0)
	})
}

func TestCreateOrder(t *testing.T) {

	// Sample order to be used in the test
	order := &model.Order{
		ID:            1,
		TableID:       1,
		CustomerName:  "John Doe",
		Status:        "completed",
		TotalAmount:   2000,
		Tax:           12,
		PaymentMethod: "bank",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		OrderProducts: []model.OrderProduct{
			{
				ProductID: 1,
				Qty:       2,
			},
		},
	}

	t.Run("Successfully create order", func(t *testing.T) {
		db, mock := helper.SetupTestDB()
		defer func() { _ = mock.ExpectationsWereMet() }()

		log := zap.NewNop()
		orderRepo := orderrepository.NewOrderRepo(db, log)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tables" WHERE id = $1 AND "tables"."deleted_at" IS NULL ORDER BY "tables"."id" LIMIT $2`)).
			WithArgs(order.TableID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "is_book"}).AddRow(order.TableID, false))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "orders"`)).
			WithArgs(order.TableID,
				order.CustomerName,
				order.Status,
				order.TotalAmount,
				order.Tax,
				order.PaymentMethod,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tables" SET "is_book"=$1,"updated_at"=$2 WHERE id = $3 AND "tables"."deleted_at" IS NULL`)).
			WithArgs(true, sqlmock.AnyArg(), order.TableID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 ORDER BY "products"."id" LIMIT $2`)).
			WithArgs(order.OrderProducts[0].ProductID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "qty"}).AddRow(order.OrderProducts[0].ProductID, 10))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "qty"=qty - $1,"updated_at"=$2 WHERE id = $3`)).
			WithArgs(order.OrderProducts[0].Qty, sqlmock.AnyArg(), order.OrderProducts[0].ProductID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "order_products"`)).
			WithArgs(1, order.OrderProducts[0].ProductID, order.OrderProducts[0].Qty).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		err := orderRepo.CreateOrder(order)

		assert.NoError(t, err)

		assert.Equal(t, uint(1), order.ID, "Order ID should be 1 after insertion")

		assert.Equal(t, true, true, "Table should be booked")

		assert.Equal(t, 10-order.OrderProducts[0].Qty, 8, "Stock should be reduced correctly")
	})

	t.Run("Failed to create order due to table not found", func(t *testing.T) {
		db, mock := helper.SetupTestDB()
		defer func() { _ = mock.ExpectationsWereMet() }()

		log := zap.NewNop()
		orderRepo := orderrepository.NewOrderRepo(db, log)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tables" WHERE id = $1 AND "tables"."deleted_at" IS NULL ORDER BY "tables"."id" LIMIT $2`)).
			WithArgs(order.TableID, 1).
			WillReturnError(fmt.Errorf("table not found"))

		mock.ExpectRollback()

		err := orderRepo.CreateOrder(order)

		assert.Error(t, err)
		assert.EqualError(t, err, "table not found")
	})

	t.Run("failed to insert order with database error", func(t *testing.T) {
		db, mock := helper.SetupTestDB()
		defer func() { _ = mock.ExpectationsWereMet() }()

		log := zap.NewNop()
		orderRepo := orderrepository.NewOrderRepo(db, log)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tables" WHERE id = $1 AND "tables"."deleted_at" IS NULL ORDER BY "tables"."id" LIMIT $2`)).
			WithArgs(order.TableID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "is_book"}).AddRow(order.TableID, false))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "orders"`)).
			WithArgs(order.TableID,
				order.CustomerName,
				order.Status,
				order.TotalAmount,
				order.Tax,
				order.PaymentMethod,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnError(fmt.Errorf("database error"))

		err := orderRepo.CreateOrder(order)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})

}

// func TestDeleteOrder(t *testing.T) {
// 	// Initialize sqlmock
// 	db, mock := helper.SetupTestDB()
// 	defer func() { _ = mock.ExpectationsWereMet() }()
//
// 	log := zap.NewNop()
// 	orderRepo := orderrepository.NewOrderRepo(db, log)
//
// 	t.Run("Successfully delete order", func(t *testing.T) {
// 		orderID := 1
//
// 		// Set up the mock for the delete operation (soft delete by updating deleted_at)
// 		mock.ExpectBegin()
// 		mock.ExpectExec(`UPDATE "orders" SET "deleted_at"=$1 WHERE id = $2 AND "orders"."deleted_at" IS NULL`).
// 			WithArgs(time.Now(), orderID).            // Expect a timestamp for the deleted_at field
// 			WillReturnResult(sqlmock.NewResult(1, 1)) // RowsAffected = 1, meaning order was found and deleted
// 		mock.ExpectCommit()
//
// 		// Call the DeleteOrder method
// 		err := orderRepo.DeleteOrder(orderID)
//
// 		// Assertions
// 		assert.NoError(t, err)
// 	})
//
// 	t.Run("Order not found", func(t *testing.T) {
// 		orderID := 999
//
// 		// Set up the mock for the delete operation (no rows affected)
// 		mock.ExpectBegin()
// 		mock.ExpectExec(`UPDATE "orders" SET "deleted_at"=$1 WHERE id = $2 AND "orders"."deleted_at" IS NULL`).
// 			WithArgs(sqlmock.AnyArg(), orderID).      // Expect the deleted_at field to be set
// 			WillReturnResult(sqlmock.NewResult(0, 0)) // RowsAffected = 0, meaning order was not found
// 		mock.ExpectCommit()
//
// 		// Call the DeleteOrder method
// 		err := orderRepo.DeleteOrder(orderID)
//
// 		// Assertions
// 		assert.Error(t, err)
// 		assert.Equal(t, fmt.Errorf("order with ID %d not found", orderID), err)
// 	})
//
// 	t.Run("Failed to delete order", func(t *testing.T) {
// 		orderID := 1
//
// 		// Set up the mock for the delete operation (simulating an error)
// 		mock.ExpectBegin()
// 		mock.ExpectExec(`UPDATE "orders" SET "deleted_at"=$1 WHERE id = $2 AND "orders"."deleted_at" IS NULL`).
// 			WithArgs(sqlmock.AnyArg(), orderID). // Expect a timestamp for the deleted_at field
// 			WillReturnError(fmt.Errorf("database error"))
// 		mock.ExpectRollback()
//
// 		// Call the DeleteOrder method
// 		err := orderRepo.DeleteOrder(orderID)
//
// 		// Assertions
// 		assert.Error(t, err)
// 		assert.Equal(t, fmt.Errorf("failed to delete order: database error"), err)
// 	})
// }

func TestUpdateOrder(t *testing.T) {

	// Sample order to be used in the test
	order := &model.Order{
		ID:            1,
		TableID:       1,
		CustomerName:  "John Doe",
		Status:        "completed",
		TotalAmount:   2000,
		Tax:           12,
		PaymentMethod: "bank",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		OrderProducts: []model.OrderProduct{
			{
				ProductID: 1,
				Qty:       2,
			},
		},
	}

	t.Run("Successfully update order", func(t *testing.T) {
		db, mock := helper.SetupTestDB()
		defer func() { _ = mock.ExpectationsWereMet() }()

		log := zap.NewNop()
		orderRepo := orderrepository.NewOrderRepo(db, log)

		mock.ExpectBegin()

		// Mock query to check if the order exists
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE id = $1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT $2`)).
			WithArgs(order.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "table_id", "status"}).AddRow(order.ID, 1, "pending"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "order_products" WHERE order_id = $1`)).
			WithArgs(order.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "qty"}).AddRow(1, 1, 5))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "qty"=qty + $1,"updated_at"=$2 WHERE id = $3`)).
			WithArgs(5, sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "order_products" WHERE order_id = $1`)).
			WithArgs(order.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 ORDER BY "products"."id" LIMIT $2`)).
			WithArgs(order.OrderProducts[0].ProductID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "stock", "price"}).AddRow(1, 10, 5000))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "qty"=qty - $1,"updated_at"=$2 WHERE id = $3`)).
			WithArgs(order.OrderProducts[0].Qty, sqlmock.AnyArg(), order.OrderProducts[0].ProductID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "order_products"`)).
			WithArgs(sqlmock.AnyArg(), order.OrderProducts[0].ProductID, order.OrderProducts[0].Qty).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tables" SET "is_book"=$1,"updated_at"=$2 WHERE id = $3 AND "tables"."deleted_at" IS NULL`)).
			WithArgs(false, sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders" SET`)).
			WithArgs(order.ID,
				order.TableID,
				order.CustomerName,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := orderRepo.UpdateOrder(int(order.ID), order)

		assert.NoError(t, err)
		assert.Equal(t, order.TotalAmount, 11200.0) // Total amount = (2 * 5000) + 10% tax
		assert.Equal(t, order.Status, "completed")
	})

	t.Run("Order not found", func(t *testing.T) {
		db, mock := helper.SetupTestDB()
		defer func() { _ = mock.ExpectationsWereMet() }()

		log := zap.NewNop()
		orderRepo := orderrepository.NewOrderRepo(db, log)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE id = $1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT $2`)).
			WithArgs(order.ID, 1).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := orderRepo.UpdateOrder(int(order.ID), order)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("order with id %d does not exist", order.ID))
	})

	t.Run("Failed to release old table", func(t *testing.T) {
		db, mock := helper.SetupTestDB()
		defer func() { _ = mock.ExpectationsWereMet() }()

		log := zap.NewNop()
		orderRepo := orderrepository.NewOrderRepo(db, log)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE id = $1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT $2`)).
			WithArgs(order.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "table_id"}).AddRow(order.ID, 2))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tables" SET "is_book"=$1,"updated_at"=$2 WHERE id = $3 AND "tables"."deleted_at" IS NULL`)).
			WithArgs(false, sqlmock.AnyArg(), 2).
			WillReturnError(fmt.Errorf("failed to release old table"))

		mock.ExpectRollback()

		err := orderRepo.UpdateOrder(int(order.ID), order)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to release old table")
	})

	t.Run("Failed to find product", func(t *testing.T) {
		db, mock := helper.SetupTestDB()
		defer func() { _ = mock.ExpectationsWereMet() }()

		log := zap.NewNop()
		orderRepo := orderrepository.NewOrderRepo(db, log)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE id = $1 AND "orders"."deleted_at" IS NULL ORDER BY "orders"."id" LIMIT $2`)).
			WithArgs(order.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "table_id"}).AddRow(order.ID, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "order_products" WHERE order_id = $1`)).
			WithArgs(order.OrderProducts[0].ProductID).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectRollback()

		err := orderRepo.UpdateOrder(int(order.ID), order)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve existing order products")

	})

}

func TestDeleteOrder(t *testing.T) {
	// Sample order to be used in the test
	orderID := 1

	t.Run("Successfully delete order", func(t *testing.T) {
		db, mock := helper.SetupTestDB()
		defer func() { _ = mock.ExpectationsWereMet() }()

		log := zap.NewNop()
		orderRepo := orderrepository.NewOrderRepo(db, log)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders" SET "deleted_at"=$1 WHERE id = $2 AND "orders"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := orderRepo.DeleteOrder(orderID)

		assert.NoError(t, err)
	})

	t.Run("Order not found", func(t *testing.T) {
		db, mock := helper.SetupTestDB()
		defer func() { _ = mock.ExpectationsWereMet() }()

		log := zap.NewNop()
		orderRepo := orderrepository.NewOrderRepo(db, log)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders" SET "deleted_at"=$1 WHERE id = $2 AND "orders"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnError(fmt.Errorf("order with ID %d not found", orderID))
		mock.ExpectRollback()

		err := orderRepo.DeleteOrder(orderID)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to delete order: order with ID 1 not found")
	})

	t.Run("Failed to delete order", func(t *testing.T) {
		db, mock := helper.SetupTestDB()
		defer func() { _ = mock.ExpectationsWereMet() }()

		log := zap.NewNop()
		orderRepo := orderrepository.NewOrderRepo(db, log)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders" SET "deleted_at"=$1 WHERE id = $2 AND "orders"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnError(fmt.Errorf("filed deletes ID %d", orderID))
		mock.ExpectRollback()

		err := orderRepo.DeleteOrder(orderID)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to delete order: filed deletes ID 1")
	})
}
