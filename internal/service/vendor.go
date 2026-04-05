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
	Create(companyID uuid.UUID, req request.VendorCreateRequest) (entity.Vendor, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.VendorResponse, int, error)
	FindById(companyID, id uuid.UUID) (response.VendorResponse, error)
	Update(companyID uuid.UUID, req request.VendorUpdateRequest) (entity.Vendor, error)
	Delete(companyID, id uuid.UUID) (entity.Vendor, error)
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
}

type VendorService struct {
	IVendorRepository repository.IVendorRepository
	validate          *validator.Validate
}

func NewVendorService(repo repository.IVendorRepository, validate *validator.Validate) IVendorService {
	return &VendorService{IVendorRepository: repo, validate: validate}
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

func (s *VendorService) Create(companyID uuid.UUID, req request.VendorCreateRequest) (entity.Vendor, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Vendor{}, err
	}
	v := entity.Vendor{
		CompanyId: companyID,
		Code:      req.Code,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CoaId:     req.CoaId,
	}
	return s.IVendorRepository.Create(v)
}

func (s *VendorService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.VendorResponse, int, error) {
	entities, total, err := s.IVendorRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}

	totalPages := (total + qp.PageSize - 1) / qp.PageSize
	if total == 0 || qp.Page > totalPages {
		return []response.VendorResponse{}, total, nil
	}

	resps := make([]response.VendorResponse, 0, len(entities))
	for _, v := range entities {
		resps = append(resps, toVendorResponse(v))
	}
	return resps, total, nil
}

func (s *VendorService) FindById(companyID, id uuid.UUID) (response.VendorResponse, error) {
	v, err := s.IVendorRepository.FindById(companyID, id)
	if err != nil {
		return response.VendorResponse{}, err
	}
	return toVendorResponse(v), nil
}

func (s *VendorService) Update(companyID uuid.UUID, req request.VendorUpdateRequest) (entity.Vendor, error) {
	v, err := s.IVendorRepository.FindById(companyID, req.Id)
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

func (s *VendorService) Delete(companyID, id uuid.UUID) (entity.Vendor, error) {
	v, err := s.IVendorRepository.FindById(companyID, id)
	if err != nil {
		return v, err
	}
	return v, s.IVendorRepository.Delete(companyID, id)
}

func (s *VendorService) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	entities, total, err := s.IVendorRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}

	totalPages := (total + qp.PageSize - 1) / qp.PageSize
	if total == 0 || qp.Page > totalPages {
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
