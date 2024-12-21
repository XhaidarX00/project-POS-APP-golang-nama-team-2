package orderrepository_test

import (
	"fmt"
	"project_pos_app/helper"
	orderrepository "project_pos_app/repository/order_repository"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetAllPayment(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := zap.NewNop()
	orderRepo := orderrepository.NewOrderRepo(db, log)

	t.Run("Successfully get all payments", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "payments"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).
				AddRow(1, "cash", time.Now()).
				AddRow(2, "bank", time.Now()))

		payments, err := orderRepo.GetAllPayment()

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, payments)
		assert.Len(t, payments, 2)

		assert.Equal(t, "cash", payments[0].Name)
		assert.Equal(t, "bank", payments[1].Name)
	})

	t.Run("No payments found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "payments"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}))

		payments, err := orderRepo.GetAllPayment()

		// Assertions
		assert.NoError(t, err)
		assert.Empty(t, payments)
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "payments"`)).
			WillReturnError(fmt.Errorf("database error"))

		payments, err := orderRepo.GetAllPayment()

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, payments)
		assert.EqualError(t, err, "database error")
	})
}
