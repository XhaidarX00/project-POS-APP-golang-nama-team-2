package model

import "time"

type Reservation struct {
	ID              uint      `gorm:"primaryKey" json:"id,omitempty"`
	TableNumber     uint      `json:"tableNumber,omitempty"`
	Pax             int       `json:"pax,omitempty"`
	Date            string    `gorm:"-" json:"date,omitempty"`
	Time            string    `gorm:"-" json:"time,omitempty"`
	ReservationDate time.Time `gorm:"type:timestamp" json:"-"`
	DepositFee      int       `json:"depositFee,omitempty"`
	Status          string    `json:"status,omitempty" gorm:"default:'Confirmed'"`
	Title           string    `json:"title,omitempty"`
	Fullname        string    `json:"fullName,omitempty"`
	PhoneNumber     string    `json:"phoneNumber,omitempty"`
	Email           string    `json:"email,omitempty"`
}
type FormUpdate struct {
	TableNumber uint
	Status      string
}

//	func SeedReservations() []Reservation {
//		return []Reservation{
//			{TableID: 1, Pax: 5, Date: "2024-12-21", Time: "11:11:00", DepositFee: 60, Title: "Mr", Fullname: "Watson Joyce", PhoneNumber: "086566156", Email: "watsonjoyce@mail.com", PaymentMethodID: 2},
//			{TableID: 1, Pax: 7, Date: "2024-12-21", Time: "12:11:00", DepositFee: 60, Title: "Mr", Fullname: "John Doe", PhoneNumber: "086566156", Email: "johndoe@mail.com", PaymentMethodID: 2},
//			{TableID: 1, Pax: 5, Date: "2024-12-21", Time: "13:11:00", DepositFee: 60, Title: "Mr", Fullname: "Will Smith", PhoneNumber: "086566156", Email: "willsmith@mail.com", PaymentMethodID: 2},
//		}
//	}
func SeedReservations() []Reservation {
	return []Reservation{
		{TableNumber: 1, Pax: 5, ReservationDate: time.Date(2024, time.December, 21, 11, 12, 0, 0, time.UTC), DepositFee: 60, Title: "Mr", Fullname: "Watson Joyce", PhoneNumber: "086566156", Email: "watsonjoyce@mail.com"},
		{TableNumber: 1, Pax: 7, ReservationDate: time.Date(2024, time.December, 21, 12, 12, 0, 0, time.UTC), DepositFee: 60, Title: "Mr", Fullname: "John Doe", PhoneNumber: "086566156", Email: "johndoe@mail.com"},
		{TableNumber: 1, Pax: 5, ReservationDate: time.Date(2024, time.December, 22, 13, 12, 0, 0, time.UTC), DepositFee: 60, Title: "Mr", Fullname: "Will Smith", PhoneNumber: "086566156", Email: "willsmith@mail.com"},
	}
}
