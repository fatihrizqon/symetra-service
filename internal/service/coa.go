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

type ICOAService interface {
	Create(req request.COACreateRequest) (entity.COA, error)
	FindAll(page, pageSize int, search string, options util.SearchOptions) ([]response.COAResponse, int, error)
	FindById(reqId uuid.UUID) (response.COAResponse, error)
	Update(req request.COAUpdateRequest) (entity.COA, error)
	Delete(reqId uuid.UUID) (entity.COA, error)

	SelectDropdownList(page, pageSize int, search string, options util.SearchOptions) ([]response.SelectDropdownListResponse, int, error)
}

type COAService struct {
	ICOARepository repository.ICOARepository
	validate       *validator.Validate
}

func NewCOAService(repo repository.ICOARepository, validate *validator.Validate) ICOAService {
	return &COAService{
		ICOARepository: repo,
		validate:       validate,
	}
}

// Create implements ICOAService.
func (e *COAService) Create(req request.COACreateRequest) (entity.COA, error) {
	entity := entity.COA{
		Code: req.Code,
		Name: req.Name,
	}

	if err := e.validate.Struct(req); err != nil {
		return entity, err
	}

	entity, err := e.ICOARepository.Create(entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}

// FindAll implements ICOAService with pagination.
func (e *COAService) FindAll(page, pageSize int, search string, options util.SearchOptions) ([]response.COAResponse, int, error) {
	var resps []response.COAResponse
	entities, totalCount, err := e.ICOARepository.FindAll(page, pageSize, search, options)

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

		var groupResp *response.COAGroupResponse
		if value.SubGroup.Group.Id != uuid.Nil {
			groupResp = &response.COAGroupResponse{
				Id:            value.SubGroup.Group.Id,
				Code:          value.SubGroup.Group.Code,
				Name:          value.SubGroup.Group.Name,
				NormalBalance: value.SubGroup.Group.NormalBalance,
				Status:        value.SubGroup.Group.Status,
				CreatedAt:     value.SubGroup.Group.CreatedAt,
				UpdatedAt:     value.SubGroup.Group.UpdatedAt,
			}
		}

		var subGroupResp *response.COASubGroupResponse
		if value.SubGroup.Id != uuid.Nil {
			subGroupResp = &response.COASubGroupResponse{
				Id:        value.SubGroup.Id,
				Code:      value.SubGroup.Code,
				Name:      value.SubGroup.Name,
				Status:    value.SubGroup.Status,
				CreatedAt: value.SubGroup.CreatedAt,
				UpdatedAt: value.SubGroup.UpdatedAt,
				GroupId:   value.SubGroup.GroupId,
				Group:     groupResp, // <- masukkan group ke sini
			}
		}

		resp := response.COAResponse{
			Id: value.Id,
			// SubgroupId: value.SubgroupId,
			// SubGroup:   subGroupResp,
			Code:          value.Code,
			Group:         groupResp.Name,
			SubGroup:      subGroupResp.Name,
			Name:          value.Name,
			NormalBalance: groupResp.NormalBalance,
			Status:        value.Status,
			CreatedAt:     value.CreatedAt,
			UpdatedAt:     value.UpdatedAt,
		}

		resps = append(resps, resp)
	}

	return resps, totalCount, nil
}

// FindById implements ICOAService.
func (e *COAService) FindById(reqId uuid.UUID) (response.COAResponse, error) {
	var res response.COAResponse
	result, err := e.ICOARepository.FindById(reqId)

	if err != nil {
		return res, err
	}

	return response.COAResponse{
		Id:        result.Id,
		Name:      result.Name,
		Code:      result.Code,
		Status:    result.Status,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

// Update implements ICOAService.
func (e *COAService) Update(req request.COAUpdateRequest) (entity.COA, error) {
	entity, err := e.ICOARepository.FindById(req.Id)
	if err != nil {
		return entity, err
	}

	entity.Name = req.Name
	entity.Code = req.Code

	err = e.ICOARepository.Update(entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}

// Delete implements ICOAService.
func (e *COAService) Delete(reqId uuid.UUID) (entity.COA, error) {
	entity, err := e.ICOARepository.FindById(reqId)
	if err != nil {
		return entity, err
	}

	err = e.ICOARepository.Delete(reqId)
	if err != nil {
		return entity, err
	}

	return entity, nil
}

// SelectDropdownList implements ICOAService.
func (e *COAService) SelectDropdownList(page int, pageSize int, search string, options util.SearchOptions) ([]response.SelectDropdownListResponse, int, error) {
	var resps []response.SelectDropdownListResponse
	entities, totalCount, err := e.ICOARepository.FindAll(page, pageSize, search, options)

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
		resp := response.SelectDropdownListResponse{
			Value: value.Id,
			Label: value.Name,
		}
		resps = append(resps, resp)
	}

	return resps, totalCount, nil
}
