package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates a new in-memory test database for each test
func setupTestDB(t *testing.T) *gorm.DB {
	// Use a unique database name for each test to ensure isolation
	dbName := "file::memory:?cache=shared&_" + t.Name()
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the schema
	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	return db
}

// createTestUser creates a test user for tests
func createTestUser(suffix string) *User {
	return &User{
		Email:     "test-" + suffix + "@example.com",
		Username:  "testuser-" + suffix,
		Name:      "Test User " + suffix,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDBUserRepository(db)

	// Create a test user
	testUser := createTestUser("id-test")
	result := db.Create(testUser)
	assert.NoError(t, result.Error)

	// Test finding by ID
	foundUser, err := repo.GetByID(testUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.ID, foundUser.ID)
	assert.Equal(t, testUser.Email, foundUser.Email)

	// Test not found case
	foundUser, err = repo.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, foundUser)
}

func TestGetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDBUserRepository(db)

	// Create a test user
	testUser := createTestUser("email-test")
	result := db.Create(testUser)
	assert.NoError(t, result.Error)

	// Test finding by email
	foundUser, err := repo.GetByEmail(testUser.Email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.Email, foundUser.Email)
	assert.Equal(t, testUser.Username, foundUser.Username)

	// Test not found case
	foundUser, err = repo.GetByEmail("nonexistent@example.com")
	assert.Error(t, err)
	assert.Nil(t, foundUser)
}

func TestGetByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDBUserRepository(db)

	// Create a test user
	testUser := createTestUser("username-test")
	result := db.Create(testUser)
	assert.NoError(t, result.Error)

	// Test finding by username
	foundUser, err := repo.GetByUsername(testUser.Username)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.Username, foundUser.Username)
	assert.Equal(t, testUser.Email, foundUser.Email)

	// Test not found case
	foundUser, err = repo.GetByUsername("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, foundUser)
}

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDBUserRepository(db)

	// Create a test user
	testUser := createTestUser("create-test")
	err := repo.Create(testUser)
	assert.NoError(t, err)
	assert.NotZero(t, testUser.ID) // Verify ID is generated

	// Verify user exists in DB
	var foundUser User
	result := db.First(&foundUser, "id = ?", testUser.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, testUser.Username, foundUser.Username)
	assert.Equal(t, testUser.Email, foundUser.Email)
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDBUserRepository(db)

	// Create a test user
	testUser := createTestUser("update-test")
	result := db.Create(testUser)
	assert.NoError(t, result.Error)

	// Update user
	testUser.Name = "Updated Name"
	testUser.Email = "updated@example.com"
	err := repo.Update(testUser)
	assert.NoError(t, err)

	// Verify changes persisted
	var foundUser User
	result = db.First(&foundUser, "id = ?", testUser.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Updated Name", foundUser.Name)
	assert.Equal(t, "updated@example.com", foundUser.Email)
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDBUserRepository(db)

	// Create a test user
	testUser := createTestUser("delete-test")
	result := db.Create(testUser)
	assert.NoError(t, result.Error)

	// Delete user
	err := repo.Delete(testUser.ID)
	assert.NoError(t, err)

	// Verify user is deleted
	var foundUser User
	result = db.First(&foundUser, "id = ?", testUser.ID)
	assert.Error(t, result.Error) // Should not find the user

	// Test deleting non-existent user
	err = repo.Delete(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestGetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDBUserRepository(db)

	// Clear any existing data
	db.Exec("DELETE FROM users")

	// Create multiple test users
	expectedCount := 3
	for i := 0; i < expectedCount; i++ {
		user := createTestUser("getall-" + string(rune('0'+i)))
		result := db.Create(user)
		assert.NoError(t, result.Error)
	}

	// Get all users
	users, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, users, expectedCount)
}

func TestUniqueConstraints(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDBUserRepository(db)

	// Create initial test user
	testUser := createTestUser("unique-test")
	err := repo.Create(testUser)
	assert.NoError(t, err)

	// Try to create user with same email
	duplicateEmailUser := createTestUser("unique-email-test")
	duplicateEmailUser.Email = testUser.Email
	err = repo.Create(duplicateEmailUser)
	assert.Error(t, err, "Should fail due to unique email constraint")

	// Try to create user with same username
	duplicateUsernameUser := createTestUser("unique-username-test")
	duplicateUsernameUser.Username = testUser.Username
	err = repo.Create(duplicateUsernameUser)
	assert.Error(t, err, "Should fail due to unique username constraint")
}
