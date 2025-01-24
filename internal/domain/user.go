package domain

const UserCollection = "users"

type User struct {
	ID               string `bson:"_id,omitempty" json:"id,omitempty"`
	Email            string `json:"email" validate:"email,required"`
	Password         string `json:"password" validate:"required"`
	TelegramUsername string `bson:"telegram_username" json:"telegram_username,omitempty"`
	FirstName        string `bson:"first_name" json:"first_name,omitempty"`
	LastName         string `bson:"last_name" json:"last_name,omitempty"`
	PhoneNumber      string `bson:"phone_number" json:"phone_number" validate:"required"`
	IsActive         bool   `bson:"is_active" json:"is_active,omitempty"`
	IsAdmin          bool   `bson:"is_admin" json:"is_admin,omitempty"`
}

type UserFilter struct {
	Email            string `bson:"email" json:"email,omitempty"`
	TelegramUsername string `bson:"telegram_username" json:"telegram_username,omitempty"`
	PhoneNumber      string `bson:"phone_number" json:"phone_number,omitempty"`
}

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetUsers(filter UserFilter) ([]User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id string) error
	GetUserByID(id string) (User, error)
}

type UserUsecase interface {
	CreateUser(user User) (User, error)
	GetUsers() ([]User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id string) error
	GetUserByID(id string) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserByTgUserName(username string) (User, error)
	ActivateUser(id string) error
	DeactivateUser(id string) error
	PromoteUser(id string) error
	DemoteUser(id string) error
	Login(email, password string) (User, error)
}
