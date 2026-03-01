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
	FindAll(qp *util.QueryParams) ([]response.COAResponse, int, error)
	FindById(reqId uuid.UUID) (response.COAResponse, error)
	Update(req request.COAUpdateRequest) (entity.COA, error)
	Delete(reqId uuid.UUID) (entity.COA, error)
	SelectDropdownList(qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
}

type COAService struct {
	ICOARepository repository.ICOARepository
	validate       *validator.Validate
}

func NewCOAService(repo repository.ICOARepository, validate *validator.Validate) ICOAService {
	return &COAService{ICOARepository: repo, validate: validate}
}

func (e *COAService) Create(req request.COACreateRequest) (entity.COA, error) {
	if err := e.validate.Struct(req); err != nil {
		return entity.COA{}, err
	}
	c := entity.COA{Code: req.Code, Name: req.Name}
	return e.ICOARepository.Create(c)
}

func (e *COAService) FindAll(qp *util.QueryParams) ([]response.COAResponse, int, error) {
	entities, totalCount, err := e.ICOARepository.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []response.COAResponse{}, 0, nil
	}

	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}

	resps := make([]response.COAResponse, 0, len(entities))
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
				Group:     groupResp,
			}
		}

		resps = append(resps, response.COAResponse{
			Id:            value.Id,
			Code:          value.Code,
			Group:         groupResp.Name,
			SubGroup:      subGroupResp.Name,
			Name:          value.Name,
			NormalBalance: groupResp.NormalBalance,
			Status:        value.Status,
			CreatedAt:     value.CreatedAt,
			UpdatedAt:     value.UpdatedAt,
		})
	}

	return resps, totalCount, nil
}

func (e *COAService) FindById(reqId uuid.UUID) (response.COAResponse, error) {
	result, err := e.ICOARepository.FindById(reqId)
	if err != nil {
		return response.COAResponse{}, err
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

func (e *COAService) Update(req request.COAUpdateRequest) (entity.COA, error) {
	c, err := e.ICOARepository.FindById(req.Id)
	if err != nil {
		return c, err
	}
	c.Name = req.Name
	c.Code = req.Code
	return c, e.ICOARepository.Update(c)
}

func (e *COAService) Delete(reqId uuid.UUID) (entity.COA, error) {
	c, err := e.ICOARepository.FindById(reqId)
	if err != nil {
		return c, err
	}
	return c, e.ICOARepository.Delete(reqId)
}

func (e *COAService) SelectDropdownList(qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	entities, totalCount, err := e.ICOARepository.FindAll(qp)
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
