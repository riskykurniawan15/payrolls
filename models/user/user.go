package user

import "time"

type (
	LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	LoginResponse struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
		User      UserInfo  `json:"user"`
	}

	CreateUserRequest struct {
		Username string `json:"username" validate:"required,min=3,max=50"`
		Password string `json:"password" validate:"required,min=6"`
		Role     string `json:"role" validate:"required,oneof=admin employee"`
	}

	CreateUserResponse struct {
		ID        uint      `json:"id"`
		Username  string    `json:"username"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}

	ProfileResponse struct {
		ID        uint      `json:"id"`
		Username  string    `json:"username"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserInfo struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	User struct {
		ID        uint      `json:"id" gorm:"column:id"`
		Username  string    `json:"username" gorm:"column:username"`
		Password  string    `json:"-" gorm:"column:password"`
		Role      string    `json:"role" gorm:"column:roles"`
		Salary    float64   `json:"salary" gorm:"column:salary;type:decimal(15,2);default:0.00"`
		CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
		UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	}
)
