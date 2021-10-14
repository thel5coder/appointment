package usecase

import (
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"time"
)

type ClinicUseCase struct {
	*UcContract
}

func (uc ClinicUseCase) Browse(search, order, sort string, page, limit int) (res []viewmodel.ClinicVm, pagination viewmodel.PaginationVm, err error) {
	repository := actions.NewClinicRepository(uc.DB)
	offset, limit, page, order, sort := uc.setPaginationParameter(page, limit, order, sort)

	clinics, count, err := repository.Browse(search, order, sort, limit, offset)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-clinic-browse")
		return res, pagination, err
	}

	for _, clinic := range clinics {
		res = append(res, uc.buildBody(clinic))
	}
	pagination = uc.setPaginationResponse(page, limit, count)

	return res, pagination, err
}

func (uc ClinicUseCase) BrowseAll(search string) (res []viewmodel.ClinicVm, err error) {
	repository := actions.NewClinicRepository(uc.DB)

	clinics, err := repository.BrowseAll(search)
	if err != nil {
		return res, err
	}

	for _, clinic := range clinics {
		res = append(res, uc.buildBody(clinic))
	}

	return res, nil
}

func (uc ClinicUseCase) ReadBy(column, value, operator string) (res viewmodel.ClinicVm, err error) {
	repository := actions.NewClinicRepository(uc.DB)
	clinic, err := repository.ReadBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-clinic-readBy")
		return res, err
	}

	res = uc.buildBody(clinic)

	return res, nil
}

func (uc ClinicUseCase) Edit(input *requests.ClinicRequest, ID string) (err error) {
	repository := actions.NewClinicRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)

	body := viewmodel.ClinicVm{
		ID:          ID,
		Name:        input.Name,
		Address:     input.Address,
		PICName:     input.PICName,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
		UpdatedAt:   now,
	}
	_, err = repository.Edit(body)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-clinic-edit")
		return err
	}

	return nil
}

func (uc ClinicUseCase) Add(input *requests.ClinicRequest) (err error) {
	repository := actions.NewClinicRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)

	body := viewmodel.ClinicVm{
		Name:        input.Name,
		Address:     input.Address,
		PICName:     input.PICName,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err = repository.Add(body)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-clinic-add")
		return err
	}

	return nil
}

func (uc ClinicUseCase) Delete(ID string) (err error) {
	repository := actions.NewClinicRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)

	count, err := uc.CountBy("", "id", ID, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-clinic-countBy")
		return err
	}

	if count > 0 {
		_, err = repository.Delete(ID, now, now)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-clinic-delete")
			return err
		}
	}

	return nil
}

func (uc ClinicUseCase) CountBy(ID, column, value, operator string) (res int, err error) {
	repository := actions.NewClinicRepository(uc.DB)
	res, err = repository.CountBy(ID, column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-clinic-countBy")
		return res, err
	}

	return res, nil
}

func (uc ClinicUseCase) buildBody(model models.Clinic) viewmodel.ClinicVm {
	return viewmodel.ClinicVm{
		ID:          model.ID,
		Name:        model.Name,
		Address:     model.Address,
		PICName:     model.PICName,
		PhoneNumber: model.PhoneNumber,
		Email:       model.Email,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
