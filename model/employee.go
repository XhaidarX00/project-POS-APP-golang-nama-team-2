package model

import (
	"fmt"
	"time"
)

type Employee struct {
	ID                uint    `gorm:"primaryKey"`
	UserID            uint    `gorm:"not null"`
	Name              string  `gorm:"type:varchar(50)"`
	PhoneNumber       string  `gorm:"type:varchar(15)"`
	Salary            float64 `gorm:"type:decimal(10,2)"`
	DateOfBirth       time.Time
	ShiftStartTiming  time.Time
	ShiftEndTiming    time.Time
	Address           string `gorm:"type:varchar(255)"`
	AdditionalDetails string `gorm:"type:text"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func alternatingShifts(index int) (time.Time, time.Time) {
	var startTime, endTime time.Time

	if index%2 == 0 { // Shift 1: 08:00 - 17:00
		startTime = time.Date(2024, time.January, 1, 8, 0, 0, 0, time.UTC)
		endTime = startTime.Add(9 * time.Hour) // Tambahkan 9 jam
	} else { // Shift 2: 14:00 - 23:00
		startTime = time.Date(2024, time.January, 1, 14, 0, 0, 0, time.UTC)
		endTime = startTime.Add(9 * time.Hour) // Tambahkan 9 jam
	}

	return startTime, endTime
}

func SeedStaff() []Employee {
	users := []struct {
		UserID uint
		Role   string
	}{
		{2, "Admin"},
		{3, "Admin"},
		{4, "admin"},
		{5, "admin"},
		{6, "Staff"},
		{7, "Staff"},
		{8, "Staff"},
	}

	staff := []Employee{}
	for i, user := range users {
		shiftStart, shiftEnd := alternatingShifts(i)

		staff = append(staff, Employee{
			UserID:            user.UserID,
			Name:              fmt.Sprintf("name %s %d", user.Role, i),
			PhoneNumber:       fmt.Sprintf("123456789%d", i),
			Salary:            4000 + float64(i*500),
			DateOfBirth:       time.Date(1990+i, time.January, 1, 0, 0, 0, 0, time.UTC),
			ShiftStartTiming:  shiftStart,
			ShiftEndTiming:    shiftEnd,
			Address:           fmt.Sprintf("Staff Address %d", i),
			AdditionalDetails: fmt.Sprintf("Details for staff %d", i),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		})
	}

	return staff
}

type ResponseEmployee struct {
	Name  string
	Email string
}
