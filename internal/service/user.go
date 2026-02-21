package service

import (
	"errors"
	"strings"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/request"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Create(req request.UserCreateRequest) (entity.User, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.UserFilters) ([]response.UserResponse, int, error)
	FindById(reqId uuid.UUID) (response.UserResponse, error)
	Update(req request.UserUpdateRequest) (entity.User, error)
	Delete(reqId uuid.UUID) (entity.User, error)
}
type UserService struct {
	IUserRepository repository.IUserRepository
	validate        *validator.Validate
}

func NewUserService(repo repository.IUserRepository, validate *validator.Validate) IUserService {
	return &UserService{
		IUserRepository: repo,
		validate:        validate,
	}
}

// Create implements IUserService.
func (e *UserService) Create(req request.UserCreateRequest) (entity.User, error) {
	var user entity.User

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return user, errors.New("failed to hash password")
	}

	entity := entity.User{
		Username: strings.ToLower(req.Username),
		Name:     req.Name,
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: string(hashed),
	}

	if err := e.validate.Struct(req); err != nil {
		return entity, err
	}

	entity, err = e.IUserRepository.Create(entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}

// FindAll implements IUserService with pagination.
func (e *UserService) FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.UserFilters) ([]response.UserResponse, int, error) {
	var resps []response.UserResponse
	entities, totalCount, err := e.IUserRepository.FindAll(page, pageSize, search, options, filters)

	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return resps, totalCount, nil
	}

	totalPages := (totalCount + pageSize - 1) / pageSize
	if page > totalPages {
		return nil, totalCount, nil
	}

	for _, value := range entities {
		resp := response.UserResponse{
			Id:        value.Id,
			Username:  value.Username,
			Name:      value.Name,
			Email:     value.Email,
			Status:    value.Status,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
		resps = append(resps, resp)
	}

	return resps, totalCount, nil
}

// FindById implements IUserService.
func (e *UserService) FindById(reqId uuid.UUID) (response.UserResponse, error) {
	var res response.UserResponse
	result, err := e.IUserRepository.FindById(reqId)

	if err != nil {
		return res, err
	}

	return response.UserResponse{
		Id:        result.Id,
		Username:  result.Username,
		Name:      result.Name,
		Email:     result.Email,
		Status:    result.Status,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

// Update implements IUserService.
func (e *UserService) Update(req request.UserUpdateRequest) (entity.User, error) {
	entity, err := e.IUserRepository.FindById(req.Id)
	if err != nil {
		return entity, err
	}

	entity.Username = strings.ToLower(req.Username)
	entity.Name = req.Name
	entity.Email = req.Email

	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
		if err != nil {
			return entity, errors.New("failed to generate password")
		}
		entity.Password = string(hashed)
	}

	err = e.IUserRepository.Update(entity)
	if err != nil {
		return entity, err
	}

	entity.Password = ""
	return entity, nil
}

// Delete implements IUserService.
func (e *UserService) Delete(reqId uuid.UUID) (entity.User, error) {
	entity, err := e.IUserRepository.FindById(reqId)
	if err != nil {
		return entity, err
	}

	err = e.IUserRepository.Delete(reqId)
	if err != nil {
		return entity, err
	}

	return entity, nil
}
