package service

import (
	"errors"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/request"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ICustomerService interface {
	Create(companyId uuid.UUID, req request.CustomerCreateRequest) (response.CustomerResponse, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.CustomerResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.CustomerResponse, error)
	Update(companyId uuid.UUID, req request.CustomerUpdateRequest) (response.CustomerResponse, error)
	Delete(companyId, id uuid.UUID) error
	SelectDropdownList(companyId uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
}

type CustomerService struct {
	repo     repository.ICustomerRepository
	validate *validator.Validate
}

func NewCustomerService(repo repository.ICustomerRepository, validate *validator.Validate) ICustomerService {
	return &CustomerService{repo: repo, validate: validate}
}

func toCustomerResponse(c entity.Customer) response.CustomerResponse {
	r := response.CustomerResponse{
		Id:        c.Id,
		Code:      c.Code,
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Address:   c.Address,
		CoaId:     c.CoaId,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.COA != nil && c.COA.Id != uuid.Nil {
		r.CoaCode = c.COA.Code
		r.CoaName = c.COA.Name
	}
	return r
}

func (s *CustomerService) Create(companyId uuid.UUID, req request.CustomerCreateRequest) (response.CustomerResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.CustomerResponse{}, err
	}
	c := entity.Customer{
		CompanyId: companyId,
		Code:      req.Code,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CoaId:     req.CoaId,
		Status:    1,
	}
	created, err := s.repo.Create(c)
	if err != nil {
		return response.CustomerResponse{}, err
	}
	return toCustomerResponse(created), nil
}

func (s *CustomerService) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.CustomerResponse, int, error) {
	entities, total, err := s.repo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []response.CustomerResponse{}, 0, nil
	}
	resps := make([]response.CustomerResponse, 0, len(entities))
	for _, c := range entities {
		resps = append(resps, toCustomerResponse(c))
	}
	return resps, total, nil
}

func (s *CustomerService) FindById(companyId, id uuid.UUID) (response.CustomerResponse, error) {
	c, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.CustomerResponse{}, errors.New("customer not found")
	}
	return toCustomerResponse(c), nil
}

func (s *CustomerService) Update(companyId uuid.UUID, req request.CustomerUpdateRequest) (response.CustomerResponse, error) {
	c, err := s.repo.FindById(companyId, req.Id)
	if err != nil {
		return response.CustomerResponse{}, errors.New("customer not found")
	}
	if err := s.validate.Struct(req); err != nil {
		return response.CustomerResponse{}, err
	}
	c.Code = req.Code
	c.Name = req.Name
	c.Email = req.Email
	c.Phone = req.Phone
	c.Address = req.Address
	c.CoaId = req.CoaId
	if err := s.repo.Update(c); err != nil {
		return response.CustomerResponse{}, err
	}
	updated, err := s.repo.FindById(companyId, req.Id)
	if err != nil {
		return response.CustomerResponse{}, err
	}
	return toCustomerResponse(updated), nil
}

func (s *CustomerService) Delete(companyId, id uuid.UUID) error {
	_, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("customer not found")
	}
	return s.repo.Delete(companyId, id)
}

func (s *CustomerService) SelectDropdownList(companyId uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	entities, total, err := s.repo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []response.SelectDropdownListResponse{}, 0, nil
	}
	resps := make([]response.SelectDropdownListResponse, 0, len(entities))
	for _, c := range entities {
		resps = append(resps, response.SelectDropdownListResponse{
			Value: c.Id,
			Label: c.Name,
		})
	}
	return resps, total, nil
}
