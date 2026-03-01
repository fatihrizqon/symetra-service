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
	FindAll(qp *util.QueryParams) ([]response.COAGroupResponse, int, error)
	FindById(reqId uuid.UUID) (response.COAGroupResponse, error)
	Update(req request.COAGroupUpdateRequest) (entity.COAGroup, error)
	Delete(reqId uuid.UUID) (entity.COAGroup, error)
	SelectDropdownList(qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
}

type COAGroupService struct {
	ICOAGroupRepository repository.ICOAGroupRepository
	validate            *validator.Validate
}

func NewCOAGroupService(repo repository.ICOAGroupRepository, validate *validator.Validate) ICOAGroupService {
	return &COAGroupService{ICOAGroupRepository: repo, validate: validate}
}

func (e *COAGroupService) Create(req request.COAGroupCreateRequest) (entity.COAGroup, error) {
	if err := e.validate.Struct(req); err != nil {
		return entity.COAGroup{}, err
	}
	g := entity.COAGroup{Code: req.Code, Name: req.Name, NormalBalance: req.NormalBalance}
	return e.ICOAGroupRepository.Create(g)
}

func (e *COAGroupService) FindAll(qp *util.QueryParams) ([]response.COAGroupResponse, int, error) {
	entities, totalCount, err := e.ICOAGroupRepository.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []response.COAGroupResponse{}, 0, nil
	}

	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}

	resps := make([]response.COAGroupResponse, 0, len(entities))
	for _, value := range entities {
		resps = append(resps, response.COAGroupResponse{
			Id:            value.Id,
			Code:          value.Code,
			Name:          value.Name,
			NormalBalance: value.NormalBalance,
			Status:        value.Status,
			CreatedAt:     value.CreatedAt,
			UpdatedAt:     value.UpdatedAt,
		})
	}
	return resps, totalCount, nil
}

func (e *COAGroupService) FindById(reqId uuid.UUID) (response.COAGroupResponse, error) {
	result, err := e.ICOAGroupRepository.FindById(reqId)
	if err != nil {
		return response.COAGroupResponse{}, err
	}
	return response.COAGroupResponse{
		Id:            result.Id,
		Code:          result.Code,
		Name:          result.Name,
		NormalBalance: result.NormalBalance,
		Status:        result.Status,
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
	}, nil
}

func (e *COAGroupService) Update(req request.COAGroupUpdateRequest) (entity.COAGroup, error) {
	g, err := e.ICOAGroupRepository.FindById(req.Id)
	if err != nil {
		return g, err
	}
	g.Name = req.Name
	g.Code = req.Code
	g.NormalBalance = req.NormalBalance
	return g, e.ICOAGroupRepository.Update(g)
}

func (e *COAGroupService) Delete(reqId uuid.UUID) (entity.COAGroup, error) {
	g, err := e.ICOAGroupRepository.FindById(reqId)
	if err != nil {
		return g, err
	}
	return g, e.ICOAGroupRepository.Delete(reqId)
}

func (e *COAGroupService) SelectDropdownList(qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	entities, totalCount, err := e.ICOAGroupRepository.FindAll(qp)
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
