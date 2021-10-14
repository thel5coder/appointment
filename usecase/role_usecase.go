package usecase

import (
	"errors"
	"github.com/gosimple/slug"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/messages"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"time"
)

type RoleUseCase struct {
	*UcContract
}

func (uc RoleUseCase) Browse(search, order, sort string, page, limit int) (res []viewmodel.RoleVm, pagination viewmodel.PaginationVm, err error) {
	repository := actions.NewRoleRepository(uc.DB)
	offset, limit, page, order, sort := uc.setPaginationParameter(page, limit, order, sort)

	roles, count, err := repository.Browse(search, order, sort, limit, offset)
	if err != nil {
		return res, pagination, err
	}

	for _, role := range roles {
		res = append(res,uc.buildBody(&role))
	}
	pagination = uc.setPaginationResponse(page, limit, count)

	return res, pagination, err
}

func (uc RoleUseCase) ReadBy(column, value string) (res viewmodel.RoleVm, err error) {
	repository := actions.NewRoleRepository(uc.DB)
	role, err := repository.ReadBy(column, value)
	if err != nil {
		return res, err
	}

	res = uc.buildBody(&role)

	return res, err
}

func (uc RoleUseCase) Edit(ID string, input *requests.RoleRequest) (err error) {
	repository := actions.NewRoleRepository(uc.DB)
	now := time.Now().UTC()
	roleSlug := slug.Make(input.Name)

	count, err := uc.CountBy(ID, "slug", roleSlug)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(messages.DataAlreadyExist)
	}

	body := viewmodel.RoleVm{
		ID:        ID,
		Name:      input.Name,
		Slug:      slug.Make(input.Name),
		UpdatedAt: now.Format(time.RFC3339),
	}
	_, err = repository.Edit(body)

	return err
}

func (uc RoleUseCase) Add(input *requests.RoleRequest) (error error) {
	repository := actions.NewRoleRepository(uc.DB)
	now := time.Now().UTC()
	roleSlug := slug.Make(input.Name)

	count, err := uc.CountBy("", "slug", roleSlug)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(messages.DataAlreadyExist)
	}

	body := viewmodel.RoleVm{
		Name:      input.Name,
		Slug:      roleSlug,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}
	_, err = repository.Add(body)
	if err != nil {
		return err
	}

	return nil
}

func (uc RoleUseCase) Delete(ID string) (err error) {
	repository := actions.NewRoleRepository(uc.DB)
	now := time.Now().UTC()

	count, err := uc.CountBy("", "id", ID)
	if err != nil {
		return errors.New(messages.DataNotFound)
	}

	if count > 0 {
		_, err = repository.Delete(ID, now.Format(time.RFC3339), now.Format(time.RFC3339))
	}

	return nil
}

func (uc RoleUseCase) CountBy(ID, column, value string) (res int, err error) {
	repository := actions.NewRoleRepository(uc.DB)
	res, err = repository.CountBy(ID, column, value)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc RoleUseCase) buildBody(model *models.Role) viewmodel.RoleVm{
	return viewmodel.RoleVm{
		ID:        model.ID,
		Name:      model.Name,
		Slug:      model.Slug,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
