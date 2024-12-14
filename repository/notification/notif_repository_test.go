package notification_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"project_pos_app/helper"
	"project_pos_app/model"
	"project_pos_app/repository/notification"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// SetupTestDB membuat mock database untuk pengujian
func SetupTestDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("failed to create sqlmock: %v", err))
	}
	return db, mock
}

func TestCreateNotification(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	repo := notification.NewNotifRepo(db, &log)

	t.Run("Successfully create a notification", func(t *testing.T) {
		now := time.Now()
		notification := &model.Notification{
			Title:     "Testing Notification",
			Message:   "This is a test notification",
			Status:    "new",
			CreatedAt: now,
			UpdatedAt: now,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "notifications" ("title","message","status","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
			WithArgs(
				notification.Title,
				notification.Message,
				notification.Status,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.Create(*notification)

		assert.NoError(t, err)
	})

	t.Run("Failed to create notification - Insertion failure", func(t *testing.T) {
		now := time.Now()
		notification := &model.Notification{
			Title:     "Testing Notification",
			Message:   "This is a test notification",
			Status:    "new",
			CreatedAt: now,
			UpdatedAt: now,
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "notifications" ("title","message","status","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
			WithArgs(
				notification.Title,
				notification.Message,
				notification.Status,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(fmt.Errorf("failed to insert notification"))
		mock.ExpectRollback()

		err := repo.Create(*notification)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create notification")
	})
}

func TestGetAllNotifications(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	repo := notification.NewNotifRepo(db, &log)

	t.Run("Successfully get all notifications", func(t *testing.T) {
		mockRows := sqlmock.NewRows([]string{"id", "title", "message", "status", "created_at", "updated_at"}).
			AddRow(1, "Notification 1", "Content 1", "new", time.Now(), time.Now()).
			AddRow(2, "Notification 2", "Content 2", "read", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "notifications" WHERE status = $1 ORDER BY created_at desc`)).
			WithArgs("new").
			WillReturnRows(mockRows)

		var notifications []model.Notification
		err := repo.GetAll(&notifications, "new")

		assert.NoError(t, err)
		assert.Len(t, notifications, 2)
		assert.Equal(t, "Notification 1", notifications[0].Title)
		assert.Equal(t, "Notification 2", notifications[1].Title)
	})

	t.Run("Failed to get all notifications - Database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "notifications" WHERE status = $1 ORDER BY created_at desc`)).
			WithArgs("new").
			WillReturnError(fmt.Errorf("database error"))

		var notifications []model.Notification
		err := repo.GetAll(&notifications, "new")

		assert.Error(t, err)
		assert.Empty(t, notifications)
	})
}

func TestFindNotificationByID(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	repo := notification.NewNotifRepo(db, &log)

	t.Run("Successfully find notification by ID", func(t *testing.T) {
		mockRow := sqlmock.NewRows([]string{"id", "title", "message", "status", "created_at", "updated_at"}).
			AddRow(1, "Notification 1", "Content 1", "new", time.Now(), time.Now())

		// Update query expectation to match the actual SQL query
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "notifications" WHERE "notifications"."id" = $1 ORDER BY "notifications"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(mockRow)

		notification, err := repo.FindByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, notification)
		assert.Equal(t, 1, notification.ID)
	})

	t.Run("Failed to find notification by ID - Not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "notifications" WHERE "notifications"."id" = $1`)).
			WithArgs(2).
			WillReturnError(fmt.Errorf("record not found"))

		notification, err := repo.FindByID(2)

		assert.Error(t, err)
		assert.Equal(t, 0, notification.ID)
		assert.EqualError(t, err, "notification not found")
	})
}

func TestUpdateNotification(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	repo := notification.NewNotifRepo(db, &log)

	t.Run("Successfully update notification", func(t *testing.T) {
		mockRow := sqlmock.NewRows([]string{"id", "title", "message", "status", "created_at", "updated_at"}).
			AddRow(1, "Notification 1", "Content 1", "new", time.Now(), time.Now())

		// Update query expectation to match the actual SQL query
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "notifications" WHERE "notifications"."id" = $1 ORDER BY "notifications"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(mockRow)

		// Mock untuk UPDATE
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "notifications" SET "title"=$1,"message"=$2,"status"=$3,"created_at"=$4,"updated_at"=$5 WHERE "id" = $6`)).
			WithArgs("Notification 1", "Content 1", "readed", sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		notification := model.Notification{
			ID:      1,
			Title:   "testing",
			Message: "testing ke sekian kali",
			Status:  "readed",
		}

		err := repo.Update(&notification, 1)
		assert.NoError(t, err)
	})

	t.Run("Failed to update notification - Not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "notifications" SET "status" = \$1 WHERE "id" = \$2`)).
			WithArgs("read", 2).
			WillReturnError(fmt.Errorf("record not found"))
		mock.ExpectRollback()

		notification := model.Notification{
			ID:     2,
			Status: "read",
		}

		err := repo.Update(&notification, 2)
		assert.Error(t, err)
		assert.EqualError(t, err, "error update notif : notification not found")
	})
}

func TestDeleteNotification(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	repo := notification.NewNotifRepo(db, &log)

	t.Run("Successfully delete notification", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "notifications" WHERE "notifications"."id" = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("Failed to delete notification - Not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "notifications" WHERE "notifications"."id" = $1`)).
			WithArgs(2).
			WillReturnError(fmt.Errorf("record not found"))
		mock.ExpectRollback()

		err := repo.Delete(2)
		assert.Error(t, err)
		assert.EqualError(t, err, "error delete notif : record not found")
	})
}

func TestMarkAllAsRead(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	repo := notification.NewNotifRepo(db, &log)

	t.Run("Successfully mark all notifications as read", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "notifications" SET "status"=$1,"updated_at"=$2 WHERE status = $3`)).
			WithArgs("readed", sqlmock.AnyArg(), "new").
			WillReturnResult(sqlmock.NewResult(1, 2)) // Mengembalikan 2 baris yang terpengaruh
		mock.ExpectCommit()

		err := repo.MarkAllAsRead()
		assert.NoError(t, err)
	})

	t.Run("Failed to mark all notifications as read - Database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "notifications" SET "status"=$1,"updated_at"=$2 WHERE status = $3`)).
			WithArgs("readed", time.Now(), "new").
			WillReturnError(fmt.Errorf("database error"))
		mock.ExpectRollback()

		err := repo.MarkAllAsRead()
		assert.Error(t, err)
		assert.EqualError(t, err, "error update status all notif")
	})
}
