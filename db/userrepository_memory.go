package db

import (
	"errors"
	"sync"
)

// UserRepositoryMemory implements UserRepository using in-memory storage
type UserRepositoryMemory struct {
	users map[uint]*User
	mu    sync.RWMutex
	id    uint
}

// NewMemoryUserRepository creates a new memory-backed user repository
func NewMemoryUserRepository() *UserRepositoryMemory {
	return &UserRepositoryMemory{
		users: make(map[uint]*User),
	}
}

// GetByID finds a user by ID
func (r *UserRepositoryMemory) GetByID(id uint) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetByEmail finds a user by email
func (r *UserRepositoryMemory) GetByEmail(email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// GetByUsername finds a user by username
func (r *UserRepositoryMemory) GetByUsername(username string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// Create adds a new user to the repository
func (r *UserRepositoryMemory) Create(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.id++
	user.ID = r.id
	r.users[user.ID] = user
	return nil
}

// Update updates an existing user in the repository
func (r *UserRepositoryMemory) Update(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	r.users[user.ID] = user
	return nil
}

// Delete removes a user from the repository
func (r *UserRepositoryMemory) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}

// GetAll retrieves all users
func (r *UserRepositoryMemory) GetAll() ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}
