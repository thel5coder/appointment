package usecase

import (
	"database/sql"
	"errors"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/messages"
	"profira-backend/helpers/str"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"strings"
	"time"
)

type TreatmentUseCase struct {
	*UcContract
}

//browse
func (uc TreatmentUseCase) Browse(search, order, sort string, page, limit int) (res []viewmodel.TreatmentVm, pagination viewmodel.PaginationVm, err error) {
	repository := actions.NewTreatmentRepository(uc.DB)
	offset, limit, page, order, sort := uc.setPaginationParameter(page, limit, order, sort)
	treatments, count, err := repository.Browse(search, order, sort, limit, offset)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-treatment-browse")
		return res, pagination, err
	}

	for _, treatment := range treatments {
		res = append(res, uc.buildBody(treatment))
	}
	pagination = uc.setPaginationResponse(page, limit, count)

	return res, pagination, nil
}

//browse all
func (uc TreatmentUseCase) BrowseAll(search string) (res []viewmodel.TreatmentVm, err error) {
	repository := actions.NewTreatmentRepository(uc.DB)
	treatments, err := repository.BrowseAll(search)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-treatment-browseAll")
		return res, err
	}

	for _, treatment := range treatments {
		res = append(res, uc.buildBody(treatment))
	}

	return res, nil
}

//browse by category
func (uc TreatmentUseCase) BrowseByCategory(categoryID string) (res []viewmodel.TreatmentVm, err error) {
	repository := actions.NewTreatmentRepository(uc.DB)
	treatments, err := repository.BrowseByCategory(categoryID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-treatment-browseByCategory")
		return res, err
	}

	for _, treatment := range treatments {
		res = append(res, uc.buildBody(treatment))
	}

	return res, nil
}

//read by
func (uc TreatmentUseCase) ReadBy(column, value, operator string) (res viewmodel.TreatmentVm, err error) {
	repository := actions.NewTreatmentRepository(uc.DB)
	treatment, err := repository.ReadBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-treatment-readBy")
		return res, err
	}
	res = uc.buildBody(treatment)

	return res, nil
}

//edit
func (uc TreatmentUseCase) Edit(input *requests.TreatmentRequest, ID string) (err error) {
	repository := actions.NewTreatmentRepository(uc.DB)
	now := time.Now().UTC()

	count, err := uc.countBy(ID, "t.name", input.Name,"", "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-treatment-countByName")
		return err
	}
	if count > 0 {
		logruslogger.Log(logruslogger.WarnLevel, messages.DataAlreadyExist, functioncaller.PrintFuncName(), "uc-treatment-countByNameExist")
		return errors.New(messages.DataAlreadyExist)
	}

	var duration int32
	for _, procedure := range input.Procedures {
		duration = duration + procedure.Duration
	}
	model := models.Treatment{
		ID:          ID,
		Name:        input.Name,
		Description: input.Description,
		Duration:    duration,
		Price:       sql.NullInt32{Int32: input.Price},
		PhotoID:     sql.NullString{String: input.PhotoID},
		IconID:      sql.NullString{String: input.IconID},
		UpdatedAt:   now,
	}
	err = repository.Edit(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-treatment-edit")
		return err
	}

	treatmentProcedureUc := TreatmentProcedureUseCase{UcContract: uc.UcContract}
	err = treatmentProcedureUc.Store(ID, input.Procedures, input.DeletedProcedures)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-treatmentProcedure-store")
		return err
	}

	return nil
}

//add
func (uc TreatmentUseCase) Add(input *requests.TreatmentRequest) (err error) {
	repository := actions.NewTreatmentRepository(uc.DB)
	now := time.Now().UTC()

	count, err := uc.countBy("", "t.name", input.Name,input.CategoryID, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-treatment-countByName")
		return err
	}
	if count > 0 {
		logruslogger.Log(logruslogger.WarnLevel, messages.DataAlreadyExist, functioncaller.PrintFuncName(), "uc-treatment-countByNameExist")
		return errors.New(messages.DataAlreadyExist)
	}
	var duration int32
	for _, procedure := range input.Procedures {
		duration = duration + procedure.Duration
	}
	model := models.Treatment{
		Name:        input.Name,
		Description: input.Description,
		Duration:    duration,
		Price:       sql.NullInt32{Int32: input.Price},
		PhotoID:     sql.NullString{String: input.PhotoID},
		IconID:      sql.NullString{String: input.IconID},
		CategoryID:  input.CategoryID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	model.ID, err = repository.Add(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-treatment-add")
		return err
	}

	treatmentProcedureUc := TreatmentProcedureUseCase{UcContract: uc.UcContract}
	err = treatmentProcedureUc.Store(model.ID, input.Procedures, input.DeletedProcedures)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-treatmentProcedure-store")
		return err
	}

	return nil
}

//delete
func (uc TreatmentUseCase) Delete(ID string) (err error) {
	repository := actions.NewTreatmentRepository(uc.DB)
	now := time.Now().UTC()

	count, err := uc.countBy("", "t.id", ID,"", "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-treatment-countByID")
		return err
	}
	if count > 0 {
		model := models.Treatment{
			ID:        ID,
			UpdatedAt: now,
			DeletedAt: sql.NullTime{Time: now},
		}
		err = repository.Delete(model, uc.TX)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-treatment-delete")
			return err
		}

		treatmentProcedureUc := TreatmentProcedureUseCase{UcContract: uc.UcContract}
		err = treatmentProcedureUc.DeleteBy("treatment_id", ID, "=")
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-treatmentProcedure-deleteByTreatmentID")
			return err
		}
	}

	return nil
}

//count by
func (uc TreatmentUseCase) countBy(ID, column, value,categoryId, operator string) (res int, err error) {
	repository := actions.NewTreatmentRepository(uc.DB)
	res, err = repository.CountBy(ID, column, value,categoryId, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-treatment-countBy")
		return res, err
	}

	return res, nil
}

//build body
func (uc TreatmentUseCase) buildBody(model models.Treatment) viewmodel.TreatmentVm {
	var err error
	//get icon path & photo path
	fileUc := FileUseCase{UcContract: uc.UcContract}
	var icon viewmodel.FileVm
	var photo viewmodel.FileVm
	if model.IconID.String != "" {
		icon, err = fileUc.Read(model.IconID.String)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-file-readIcon")
		}
	}

	if model.PhotoID.String != "" {
		photo, err = fileUc.Read(model.PhotoID.String)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-file-readPhoto")
		}
	}

	//parse treatment procedure array string
	var treatmentProcedureVm []viewmodel.TreatmentProcedureVm
	if model.TreatmentProcedures.String != "" {
		treatmentProcedures := strings.Split(model.TreatmentProcedures.String, ",")
		for _, treatmentProcedure := range treatmentProcedures {
			procedureArr := strings.Split(treatmentProcedure, ":")
			treatmentProcedureVm = append(treatmentProcedureVm, viewmodel.TreatmentProcedureVm{
				ID:                procedureArr[0],
				ProcedureID:       procedureArr[1],
				ProcedureName:     procedureArr[2],
				ProcedureDuration: str.StringToInt(procedureArr[3]),
			})
		}
	}

	return viewmodel.TreatmentVm{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Duration:    model.Duration,
		Price:       model.Price.Int32,
		PhotoID:     model.PhotoID.String,
		PhotoPath:   photo.Path,
		IconID:      model.IconID.String,
		IconPath:    icon.Path,
		CreatedAt:   model.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   model.UpdatedAt.Format(time.RFC3339),
		Procedures:  treatmentProcedureVm,
	}
}
