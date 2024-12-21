package model

type AccessPermission struct {
	ID           uint `gorm:"primaryKey"`
	UserID       uint `json:"user_id"`
	PermissionID uint `json:"permission_id"`
	Status       bool
}

type Permission struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type ResponseAccess struct {
	Permission string `json:"permission"`
	Status     bool
	Role       string
	UserID     int
}

func SeedPermissions() []Permission {
	return []Permission{
		{Name: "Dashboard"},
		{Name: "Revenue"},
		{Name: "Inventory"},
		{Name: "Order"},
		{Name: "Superadmin"},
		{Name: "Product"},
		{Name: "Notification"},
	}
}

func SeedAccessPermissions() []AccessPermission {

	accessPermissions := []AccessPermission{
		{UserID: 1, PermissionID: 1, Status: true},
		{UserID: 1, PermissionID: 2, Status: false},
		{UserID: 1, PermissionID: 3, Status: true},
		{UserID: 1, PermissionID: 4, Status: false},
		{UserID: 1, PermissionID: 5, Status: true},
		{UserID: 1, PermissionID: 6, Status: true},

		{UserID: 2, PermissionID: 1, Status: false},
		{UserID: 2, PermissionID: 2, Status: true},
		{UserID: 2, PermissionID: 3, Status: true},
		{UserID: 2, PermissionID: 4, Status: true},
		{UserID: 2, PermissionID: 5, Status: false},
		{UserID: 2, PermissionID: 6, Status: true},

		{UserID: 3, PermissionID: 1, Status: true},
		{UserID: 3, PermissionID: 2, Status: false},
		{UserID: 3, PermissionID: 3, Status: true},
		{UserID: 3, PermissionID: 4, Status: true},
		{UserID: 3, PermissionID: 5, Status: true},
		{UserID: 3, PermissionID: 6, Status: false},

		{UserID: 4, PermissionID: 1, Status: true},
		{UserID: 4, PermissionID: 2, Status: false},
		{UserID: 4, PermissionID: 3, Status: true},
		{UserID: 4, PermissionID: 4, Status: false},
		{UserID: 4, PermissionID: 5, Status: true},
		{UserID: 4, PermissionID: 6, Status: true},

		{UserID: 5, PermissionID: 1, Status: true},
		{UserID: 5, PermissionID: 2, Status: false},
		{UserID: 5, PermissionID: 3, Status: true},
		{UserID: 5, PermissionID: 4, Status: true},
		{UserID: 5, PermissionID: 5, Status: false},
		{UserID: 5, PermissionID: 6, Status: true},

		{UserID: 6, PermissionID: 1, Status: true},
		{UserID: 6, PermissionID: 2, Status: false},
		{UserID: 6, PermissionID: 3, Status: true},
		{UserID: 6, PermissionID: 4, Status: false},
		{UserID: 6, PermissionID: 5, Status: false},
		{UserID: 6, PermissionID: 6, Status: true},

		{UserID: 7, PermissionID: 1, Status: true},
		{UserID: 7, PermissionID: 2, Status: true},
		{UserID: 7, PermissionID: 3, Status: true},
		{UserID: 7, PermissionID: 4, Status: true},
		{UserID: 7, PermissionID: 5, Status: true},
		{UserID: 7, PermissionID: 6, Status: true},

		{UserID: 8, PermissionID: 1, Status: true},
		{UserID: 8, PermissionID: 2, Status: true},
		{UserID: 8, PermissionID: 3, Status: true},
		{UserID: 8, PermissionID: 4, Status: false},
		{UserID: 8, PermissionID: 5, Status: true},
		{UserID: 8, PermissionID: 6, Status: true},

		{UserID: 9, PermissionID: 1, Status: true},
		{UserID: 9, PermissionID: 2, Status: true},
		{UserID: 9, PermissionID: 3, Status: true},
		{UserID: 9, PermissionID: 4, Status: true},
		{UserID: 9, PermissionID: 5, Status: true},
		{UserID: 9, PermissionID: 6, Status: false},

		{UserID: 10, PermissionID: 1, Status: true},
		{UserID: 10, PermissionID: 2, Status: true},
		{UserID: 10, PermissionID: 3, Status: true},
		{UserID: 10, PermissionID: 4, Status: false},
		{UserID: 10, PermissionID: 5, Status: true},
		{UserID: 10, PermissionID: 6, Status: true},
	}

	return accessPermissions
}
