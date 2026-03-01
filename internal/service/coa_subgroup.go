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
	FindAll(qp *util.QueryParams) ([]response.COASubGroupResponse, int, error)
	FindById(reqId uuid.UUID) (response.COASubGroupResponse, error)
	Update(req request.COASubGroupUpdateRequest) (entity.COASubGroup, error)
	Delete(reqId uuid.UUID) (entity.COASubGroup, error)
	SelectDropdownList(qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
}

type COASubGroupService struct {
	ICOASubGroupRepository repository.ICOASubGroupRepository
	validate               *validator.Validate
}

func NewCOASubGroupService(repo repository.ICOASubGroupRepository, validate *validator.Validate) ICOASubGroupService {
	return &COASubGroupService{ICOASubGroupRepository: repo, validate: validate}
}

func (e *COASubGroupService) Create(req request.COASubGroupCreateRequest) (entity.COASubGroup, error) {
	if err := e.validate.Struct(req); err != nil {
		return entity.COASubGroup{}, err
	}
	sg := entity.COASubGroup{GroupId: req.GroupId, Code: req.Code, Name: req.Name}
	return e.ICOASubGroupRepository.Create(sg)
}

func (e *COASubGroupService) FindAll(qp *util.QueryParams) ([]response.COASubGroupResponse, int, error) {
	entities, totalCount, err := e.ICOASubGroupRepository.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []response.COASubGroupResponse{}, 0, nil
	}

	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}

	resps := make([]response.COASubGroupResponse, 0, len(entities))
	for _, value := range entities {
		var groupResp *response.COAGroupResponse
		if value.Group.Id != uuid.Nil {
			groupResp = &response.COAGroupResponse{
				Id:            value.Group.Id,
				Code:          value.Group.Code,
				Name:          value.Group.Name,
				NormalBalance: value.Group.NormalBalance,
				Status:        value.Group.Status,
				CreatedAt:     value.Group.CreatedAt,
				UpdatedAt:     value.Group.UpdatedAt,
			}
		}
		resps = append(resps, response.COASubGroupResponse{
			Id:        value.Id,
			GroupId:   value.GroupId,
			Group:     groupResp,
			Code:      value.Code,
			Name:      value.Name,
			Status:    value.Status,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		})
	}
	return resps, totalCount, nil
}

func (e *COASubGroupService) FindById(reqId uuid.UUID) (response.COASubGroupResponse, error) {
	result, err := e.ICOASubGroupRepository.FindById(reqId)
	if err != nil {
		return response.COASubGroupResponse{}, err
	}
	return response.COASubGroupResponse{
		Id:        result.Id,
		GroupId:   result.GroupId,
		Code:      result.Code,
		Name:      result.Name,
		Status:    result.Status,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (e *COASubGroupService) Update(req request.COASubGroupUpdateRequest) (entity.COASubGroup, error) {
	sg, err := e.ICOASubGroupRepository.FindById(req.Id)
	if err != nil {
		return sg, err
	}
	sg.GroupId = req.GroupId
	sg.Name = req.Name
	sg.Code = req.Code
	return sg, e.ICOASubGroupRepository.Update(sg)
}

func (e *COASubGroupService) Delete(reqId uuid.UUID) (entity.COASubGroup, error) {
	sg, err := e.ICOASubGroupRepository.FindById(reqId)
	if err != nil {
		return sg, err
	}
	return sg, e.ICOASubGroupRepository.Delete(reqId)
}

func (e *COASubGroupService) SelectDropdownList(qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	entities, totalCount, err := e.ICOASubGroupRepository.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}

	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if totalCount == 0 || qp.Page > totalPages {
		return []response.SelectDropdownListResponse{}, totalCount, nil
	}

	resps := make([]response.SelectDropdownListResponse, 0, len(entities))
	for _, value := range entities {
		resps = append(resps, response.SelectDropdownListResponse{
			Value: value.Id,
			Label: value.Name,
		})
	}
	return resps, totalCount, nil
}
