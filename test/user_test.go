package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindById(id uuid.UUID) (entity.User, error) {
	args := m.Called(id)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) Create(entity entity.User) (entity.User, error) {
	args := m.Called(entity)
	return entity, args.Error(1)
}

func (m *MockUserRepository) Update(entity entity.User) error {
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll() ([]entity.User, error) {
	args := m.Called()
	return args.Get(0).([]entity.User), args.Error(1)
}

func TestUserCreate(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)

	// Define route
	app.Post("/users", func(c *fiber.Ctx) error {
		var user entity.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).SendString("invalid request")
		}
		createdUser, err := mockRepo.Create(user)
		if err != nil {
			return c.Status(500).SendString("failed to create user")
		}
		return c.Status(201).JSON(createdUser)
	})

	// Mock data
	user1 := entity.User{Id: uuid.New(), Username: "admin"}

	// Test Create
	mockRepo.On("Create", user1).Return(user1, nil)

	reqBody, _ := json.Marshal(user1)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestUserFindById(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		parsedId, err := uuid.Parse(id)
		if err != nil {
			util.HandleError(c, 500, err)
			return nil
		}

		user, err := mockRepo.FindById(parsedId)
		if err != nil {
			return c.Status(404).SendString("record not found")
		}
		return c.Status(200).JSON(user)
	})

	// Mock data
	user1 := entity.User{Id: uuid.New(), Username: "admin"}

	// Test FindById
	mockRepo.On("FindById", user1.Id).Return(user1, nil)
	req := httptest.NewRequest("GET", "/users/"+user1.Id.String(), nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
