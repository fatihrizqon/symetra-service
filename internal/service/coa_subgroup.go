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

type ICOASubGroupService interface {
	Create(req request.COASubGroupCreateRequest) (entity.COASubGroup, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.COASubGroupFilters) ([]response.COASubGroupResponse, int, error)
	FindById(reqId uuid.UUID) (response.COASubGroupResponse, error)
	Update(req request.COASubGroupUpdateRequest) (entity.COASubGroup, error)
	Delete(reqId uuid.UUID) (entity.COASubGroup, error)
}

type COASubGroupService struct {
	ICOASubGroupRepository repository.ICOASubGroupRepository
	validate               *validator.Validate
}

func NewCOASubGroupService(repo repository.ICOASubGroupRepository, validate *validator.Validate) ICOASubGroupService {
	return &COASubGroupService{
		ICOASubGroupRepository: repo,
		validate:               validate,
	}
}

// Create implements ICOASubGroupService.
func (e *COASubGroupService) Create(req request.COASubGroupCreateRequest) (entity.COASubGroup, error) {
	entity := entity.COASubGroup{
		GroupId: req.GroupId,
		Code:    req.Code,
		Name:    req.Name,
	}

	if err := e.validate.Struct(req); err != nil {
		return entity, err
	}

	entity, err := e.ICOASubGroupRepository.Create(entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}

// FindAll implements ICOASubGroupService with pagination.
func (e *COASubGroupService) FindAll(page, pageSize int, search string, options util.SearchOptions, filters entity.COASubGroupFilters) ([]response.COASubGroupResponse, int, error) {
	var resps []response.COASubGroupResponse
	entities, totalCount, err := e.ICOASubGroupRepository.FindAll(page, pageSize, search, options, filters)

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
		resp := response.COASubGroupResponse{
			Id:        value.Id,
			GroupId:   value.GroupId,
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

// FindById implements ICOASubGroupService.
func (e *COASubGroupService) FindById(reqId uuid.UUID) (response.COASubGroupResponse, error) {
	var res response.COASubGroupResponse
	result, err := e.ICOASubGroupRepository.FindById(reqId)

	if err != nil {
		return res, err
	}

	return response.COASubGroupResponse{
		Id:        result.Id,
		GroupId:   result.GroupId,
		Name:      result.Name,
		Code:      result.Code,
		Status:    result.Status,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

// Update implements ICOASubGroupService.
func (e *COASubGroupService) Update(req request.COASubGroupUpdateRequest) (entity.COASubGroup, error) {
	entity, err := e.ICOASubGroupRepository.FindById(req.Id)
	if err != nil {
		return entity, err
	}

	entity.GroupId = req.GroupId
	entity.Name = req.Name
	entity.Code = req.Code

	err = e.ICOASubGroupRepository.Update(entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}

// Delete implements ICOASubGroupService.
func (e *COASubGroupService) Delete(reqId uuid.UUID) (entity.COASubGroup, error) {
	entity, err := e.ICOASubGroupRepository.FindById(reqId)
	if err != nil {
		return entity, err
	}

	err = e.ICOASubGroupRepository.Delete(reqId)
	if err != nil {
		return entity, err
	}

	return entity, nil
}
