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

type IVendorService interface {
	Create(req request.VendorCreateRequest) (entity.Vendor, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.VendorFilters) ([]response.VendorResponse, int, error)
	FindById(id uuid.UUID) (response.VendorResponse, error)
	Update(req request.VendorUpdateRequest) (entity.Vendor, error)
	Delete(id uuid.UUID) (entity.Vendor, error)
	SelectDropdownList(page, pageSize int, search string, options util.SearchOptions) ([]response.SelectDropdownListResponse, int, error)
}

type VendorService struct {
	IVendorRepository repository.IVendorRepository
	validate          *validator.Validate
}

func NewVendorService(repo repository.IVendorRepository, validate *validator.Validate) IVendorService {
	return &VendorService{
		IVendorRepository: repo,
		validate:          validate,
	}
}

func toVendorResponse(v entity.Vendor) response.VendorResponse {
	r := response.VendorResponse{
		Id:        v.Id,
		Code:      v.Code,
		Name:      v.Name,
		Email:     v.Email,
		Phone:     v.Phone,
		Address:   v.Address,
		CoaId:     v.CoaId,
		Status:    v.Status,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
	if v.COA != nil && v.COA.Id != uuid.Nil {
		r.CoaCode = v.COA.Code
		r.CoaName = v.COA.Name
	}
	return r
}

func (s *VendorService) Create(req request.VendorCreateRequest) (entity.Vendor, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Vendor{}, err
	}

	v := entity.Vendor{
		Code:    req.Code,
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
		CoaId:   req.CoaId,
	}

	return s.IVendorRepository.Create(v)
}

func (s *VendorService) FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.VendorFilters) ([]response.VendorResponse, int, error) {
	entities, total, err := s.IVendorRepository.FindAll(page, pageSize, search, options, filters)
	if err != nil {
		return nil, 0, err
	}

	totalPages := (total + pageSize - 1) / pageSize
	if total == 0 || page > totalPages {
		return []response.VendorResponse{}, total, nil
	}

	resps := make([]response.VendorResponse, 0, len(entities))
	for _, v := range entities {
		resps = append(resps, toVendorResponse(v))
	}
	return resps, total, nil
}

func (s *VendorService) FindById(id uuid.UUID) (response.VendorResponse, error) {
	v, err := s.IVendorRepository.FindById(id)
	if err != nil {
		return response.VendorResponse{}, err
	}
	return toVendorResponse(v), nil
}

func (s *VendorService) Update(req request.VendorUpdateRequest) (entity.Vendor, error) {
	v, err := s.IVendorRepository.FindById(req.Id)
	if err != nil {
		return v, err
	}

	if err := s.validate.Struct(req); err != nil {
		return v, err
	}

	v.Code = req.Code
	v.Name = req.Name
	v.Email = req.Email
	v.Phone = req.Phone
	v.Address = req.Address
	v.CoaId = req.CoaId

	if err := s.IVendorRepository.Update(v); err != nil {
		return v, err
	}
	return v, nil
}

func (s *VendorService) Delete(id uuid.UUID) (entity.Vendor, error) {
	v, err := s.IVendorRepository.FindById(id)
	if err != nil {
		return v, err
	}
	return v, s.IVendorRepository.Delete(id)
}

func (s *VendorService) SelectDropdownList(page, pageSize int, search string, options util.SearchOptions) ([]response.SelectDropdownListResponse, int, error) {
	entities, total, err := s.IVendorRepository.FindAll(page, pageSize, search, options, entity.VendorFilters{})
	if err != nil {
		return nil, 0, err
	}

	totalPages := (total + pageSize - 1) / pageSize
	if total == 0 || page > totalPages {
		return []response.SelectDropdownListResponse{}, total, nil
	}

	resps := make([]response.SelectDropdownListResponse, 0, len(entities))
	for _, v := range entities {
		resps = append(resps, response.SelectDropdownListResponse{
			Value: v.Id,
			Label: v.Name,
		})
	}
	return resps, total, nil
}
