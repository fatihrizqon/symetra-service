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
	FindAll(qp *util.QueryParams) ([]response.UserResponse, int, error)
	FindById(reqId uuid.UUID) (response.UserResponse, error)
	Update(req request.UserUpdateRequest) (entity.User, error)
	Delete(reqId uuid.UUID) (entity.User, error)
}

type UserService struct {
	IUserRepository repository.IUserRepository
	validate        *validator.Validate
}

func NewUserService(repo repository.IUserRepository, validate *validator.Validate) IUserService {
	return &UserService{IUserRepository: repo, validate: validate}
}

func (s *UserService) Create(req request.UserCreateRequest) (entity.User, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.User{}, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return entity.User{}, errors.New("failed to hash password")
	}

	u := entity.User{
		Username: strings.ToLower(req.Username),
		Name:     req.Name,
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: string(hashed),
	}

	return s.IUserRepository.Create(u)
}

func (s *UserService) FindAll(qp *util.QueryParams) ([]response.UserResponse, int, error) {
	entities, totalCount, err := s.IUserRepository.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []response.UserResponse{}, 0, nil
	}

	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}

	resps := make([]response.UserResponse, 0, len(entities))
	for _, u := range entities {
		resps = append(resps, response.UserResponse{
			Id:        u.Id,
			Username:  u.Username,
			Name:      u.Name,
			Email:     u.Email,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return resps, totalCount, nil
}

func (s *UserService) FindById(reqId uuid.UUID) (response.UserResponse, error) {
	u, err := s.IUserRepository.FindById(reqId)
	if err != nil {
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		Id:        u.Id,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (s *UserService) Update(req request.UserUpdateRequest) (entity.User, error) {
	u, err := s.IUserRepository.FindById(req.Id)
	if err != nil {
		return u, err
	}

	u.Username = strings.ToLower(req.Username)
	u.Name = req.Name
	u.Email = req.Email

	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
		if err != nil {
			return u, errors.New("failed to generate password")
		}
		u.Password = string(hashed)
	}

	if err := s.IUserRepository.Update(u); err != nil {
		return u, err
	}

	u.Password = ""
	return u, nil
}

func (s *UserService) Delete(reqId uuid.UUID) (entity.User, error) {
	u, err := s.IUserRepository.FindById(reqId)
	if err != nil {
		return u, err
	}
	return u, s.IUserRepository.Delete(reqId)
}
