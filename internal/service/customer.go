package service

import (
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/request"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ICustomerService interface {
	Create(req request.CustomerCreateRequest) (entity.Customer, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.CustomerFilters) ([]response.CustomerResponse, int, error)
	FindById(id uuid.UUID) (response.CustomerResponse, error)
	Update(req request.CustomerUpdateRequest) (entity.Customer, error)
	Delete(id uuid.UUID) (entity.Customer, error)
	SelectDropdownList(page, pageSize int, search string, options util.SearchOptions) ([]response.SelectDropdownListResponse, int, error)
}

type CustomerService struct {
	ICustomerRepository repository.ICustomerRepository
	validate            *validator.Validate
}

func NewCustomerService(repo repository.ICustomerRepository, validate *validator.Validate) ICustomerService {
	return &CustomerService{
		ICustomerRepository: repo,
		validate:            validate,
	}
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

func (s *CustomerService) Create(req request.CustomerCreateRequest) (entity.Customer, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Customer{}, err
	}

	c := entity.Customer{
		Code:    req.Code,
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
		CoaId:   req.CoaId,
	}

	return s.ICustomerRepository.Create(c)
}

func (s *CustomerService) FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.CustomerFilters) ([]response.CustomerResponse, int, error) {
	entities, total, err := s.ICustomerRepository.FindAll(page, pageSize, search, options, filters)
	if err != nil {
		return nil, 0, err
	}

	totalPages := (total + pageSize - 1) / pageSize
	if total == 0 || page > totalPages {
		return []response.CustomerResponse{}, total, nil
	}

	resps := make([]response.CustomerResponse, 0, len(entities))
	for _, c := range entities {
		resps = append(resps, toCustomerResponse(c))
	}
	return resps, total, nil
}

func (s *CustomerService) FindById(id uuid.UUID) (response.CustomerResponse, error) {
	c, err := s.ICustomerRepository.FindById(id)
	if err != nil {
		return response.CustomerResponse{}, err
	}
	return toCustomerResponse(c), nil
}

func (s *CustomerService) Update(req request.CustomerUpdateRequest) (entity.Customer, error) {
	c, err := s.ICustomerRepository.FindById(req.Id)
	if err != nil {
		return c, err
	}

	if err := s.validate.Struct(req); err != nil {
		return c, err
	}

	c.Code = req.Code
	c.Name = req.Name
	c.Email = req.Email
	c.Phone = req.Phone
	c.Address = req.Address
	c.CoaId = req.CoaId

	if err := s.ICustomerRepository.Update(c); err != nil {
		return c, err
	}
	return c, nil
}

func (s *CustomerService) Delete(id uuid.UUID) (entity.Customer, error) {
	c, err := s.ICustomerRepository.FindById(id)
	if err != nil {
		return c, err
	}
	return c, s.ICustomerRepository.Delete(id)
}

func (s *CustomerService) SelectDropdownList(page, pageSize int, search string, options util.SearchOptions) ([]response.SelectDropdownListResponse, int, error) {
	entities, total, err := s.ICustomerRepository.FindAll(page, pageSize, search, options, entity.CustomerFilters{})
	if err != nil {
		return nil, 0, err
	}

	totalPages := (total + pageSize - 1) / pageSize
	if total == 0 || page > totalPages {
		return []response.SelectDropdownListResponse{}, total, nil
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
