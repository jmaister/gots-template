package db

// UserRepository interface for abstracting user database operations
type UserRepository interface {
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
	GetAll() ([]*User, error)
}
