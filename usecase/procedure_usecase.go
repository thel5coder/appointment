package usecase

import (
	"database/sql"
	"errors"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/messages"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"time"
)

type ProcedureUseCase struct {
	*UcContract
}

func (uc ProcedureUseCase) Browse(search, order, sort string, page, limit int) (res []viewmodel.ProcedureVm, pagination viewmodel.PaginationVm, err error) {
	repository := actions.NewProcedureRepository(uc.DB)

	offset, limit, page, order, sort := uc.setPaginationParameter(page, limit, order, sort)
	procedures, count, err := repository.Browse(search, order, sort, limit, offset)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-browse-procedure")
		return res, pagination, err
	}

	for _, procedure := range procedures {
		res = append(res, uc.buildBody(procedure))
	}
	pagination = uc.setPaginationResponse(page, limit, count)

	return res, pagination, nil
}

func (uc ProcedureUseCase) BrowseAll(search string) (res []viewmodel.ProcedureVm, err error) {
	repository := actions.NewProcedureRepository(uc.DB)

	procedures, err := repository.BrowseAll(search)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-browseAll-procedure")
		return res, err
	}

	for _, procedure := range procedures {
		res = append(res, uc.buildBody(procedure))
	}

	return res, nil
}

func (uc ProcedureUseCase) ReadBy(column, value, operator string) (res viewmodel.ProcedureVm, err error) {
	repository := actions.NewProcedureRepository(uc.DB)

	procedure, err := repository.ReadBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-readBy-procedure")
		return res, err
	}
	res = uc.buildBody(procedure)

	return res, nil
}

func (uc ProcedureUseCase) Edit(input *requests.ProcedureRequest, ID string) (err error) {
	repository := actions.NewProcedureRepository(uc.DB)
	now := time.Now().UTC()

	count, err := uc.countBy(ID, "name", input.Name, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-countBy-procedure")
		return err
	}
	if count > 0 {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-countBy-procedure")
		return errors.New(messages.DataAlreadyExist)
	}

	model := models.Procedure{
		ID:        ID,
		Name:      input.Name,
		Duration:  input.Duration,
		UpdatedAt: now,
	}
	_, err = repository.Edit(model)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-edit-procedure")
		return err
	}

	return nil
}

func (uc ProcedureUseCase) Add(input *requests.ProcedureRequest) (err error) {
	repository := actions.NewProcedureRepository(uc.DB)
	now := time.Now().UTC()

	count, err := uc.countBy("", "name", input.Name, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-countBy-procedure")
		return err
	}
	if count > 0 {
		logruslogger.Log(logruslogger.WarnLevel, messages.DataAlreadyExist, functioncaller.PrintFuncName(), "uc-countBy-procedure")
		return errors.New(messages.DataAlreadyExist)
	}

	model := models.Procedure{
		Name:      input.Name,
		Duration:  input.Duration,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err = repository.Add(model)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-add-procedure")
		return err
	}

	return nil
}

func (uc ProcedureUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewProcedureRepository(uc.DB)
	now := time.Now().UTC()

	count, err := uc.countBy("", column, value, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-countBy-procedure")
		return err
	}
	if count > 0 {
		model := models.Procedure{
			UpdatedAt: now,
			DeletedAt: sql.NullTime{
				Time: now,
			},
		}
		_, err = repository.DeleteBy(column, value, operator, model)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-deleteBy-procedure")
			return err
		}

		return nil
	}

	return nil
}

func (uc ProcedureUseCase) countBy(ID, column, value, operator string) (res int, err error) {
	repository := actions.NewProcedureRepository(uc.DB)

	res, err = repository.CountBy(ID, column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-procedure-countBy")
		return res, err
	}

	return res, nil
}

func (uc ProcedureUseCase) buildBody(model models.Procedure) viewmodel.ProcedureVm {
	return viewmodel.ProcedureVm{
		ID:        model.ID,
		Name:      model.Name,
		Duration:  model.Duration,
		CreatedAt: model.CreatedAt.Format(time.RFC3339),
		UpdatedAt: model.UpdatedAt.Format(time.RFC3339),
	}
}
