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

type ICOAGroupService interface {
	Create(req request.COAGroupCreateRequest) (entity.COAGroup, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.COAGroupFilters) ([]response.COAGroupResponse, int, error)
	FindById(reqId uuid.UUID) (response.COAGroupResponse, error)
	Update(req request.COAGroupUpdateRequest) (entity.COAGroup, error)
	Delete(reqId uuid.UUID) (entity.COAGroup, error)
}

type COAGroupService struct {
	ICOAGroupRepository repository.ICOAGroupRepository
	validate            *validator.Validate
}

func NewCOAGroupService(repo repository.ICOAGroupRepository, validate *validator.Validate) ICOAGroupService {
	return &COAGroupService{
		ICOAGroupRepository: repo,
		validate:            validate,
	}
}

// Create implements ICOAGroupService.
func (e *COAGroupService) Create(req request.COAGroupCreateRequest) (entity.COAGroup, error) {
	entity := entity.COAGroup{
		Code: req.Code,
		Name: req.Name,
	}

	if err := e.validate.Struct(req); err != nil {
		return entity, err
	}

	entity, err := e.ICOAGroupRepository.Create(entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}

// FindAll implements ICOAGroupService with pagination.
func (e *COAGroupService) FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.COAGroupFilters) ([]response.COAGroupResponse, int, error) {
	var resps []response.COAGroupResponse
	entities, totalCount, err := e.ICOAGroupRepository.FindAll(page, pageSize, search, options, filters)

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
		resp := response.COAGroupResponse{
			Id:        value.Id,
			Code:      value.Code,
			Name:      value.Name,
			Status:    value.Status,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
		resps = append(resps, resp)
	}

	return resps, totalCount, nil
}

// FindById implements ICOAGroupService.
func (e *COAGroupService) FindById(reqId uuid.UUID) (response.COAGroupResponse, error) {
	var res response.COAGroupResponse
	result, err := e.ICOAGroupRepository.FindById(reqId)

	if err != nil {
		return res, err
	}

	return response.COAGroupResponse{
		Id:        result.Id,
		Name:      result.Name,
		Code:      result.Code,
		Status:    result.Status,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

// Update implements ICOAGroupService.
func (e *COAGroupService) Update(req request.COAGroupUpdateRequest) (entity.COAGroup, error) {
	entity, err := e.ICOAGroupRepository.FindById(req.Id)
	if err != nil {
		return entity, err
	}

	entity.Name = req.Name
	entity.Code = req.Code

	err = e.ICOAGroupRepository.Update(entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}

// Delete implements ICOAGroupService.
func (e *COAGroupService) Delete(reqId uuid.UUID) (entity.COAGroup, error) {
	entity, err := e.ICOAGroupRepository.FindById(reqId)
	if err != nil {
		return entity, err
	}

	err = e.ICOAGroupRepository.Delete(reqId)
	if err != nil {
		return entity, err
	}

	return entity, nil
}
