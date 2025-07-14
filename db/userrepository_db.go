package db

import (
	"errors"

	"gorm.io/gorm"
)

// UserRepositoryDB implements UserRepository with a GORM database connection
type UserRepositoryDB struct {
	db *gorm.DB
}

// NewDBUserRepository creates a new database-backed user repository
func NewDBUserRepository(db *gorm.DB) *UserRepositoryDB {
	return &UserRepositoryDB{
		db: db,
	}
}

// GetByID finds a user by ID
func (r *UserRepositoryDB) GetByID(id uint) (*User, error) {
	var user User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail finds a user by email
func (r *UserRepositoryDB) GetByEmail(email string) (*User, error) {
	var user User
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByUsername finds a user by username
func (r *UserRepositoryDB) GetByUsername(username string) (*User, error) {
	var user User
	result := r.db.First(&user, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Create adds a new user to the repository
func (r *UserRepositoryDB) Create(user *User) error {
	result := r.db.Create(user)
	return result.Error
}

// Update updates an existing user in the repository
func (r *UserRepositoryDB) Update(user *User) error {
	result := r.db.Save(user)
	return result.Error
}

// Delete removes a user from the repository
func (r *UserRepositoryDB) Delete(id uint) error {
	result := r.db.Delete(&User{}, id)
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return result.Error
}

// GetAll retrieves all users from the database
func (r *UserRepositoryDB) GetAll() ([]*User, error) {
	var users []*User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
