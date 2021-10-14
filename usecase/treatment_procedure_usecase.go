package usecase

import (
	"database/sql"
	"fmt"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/server/requests"
	"time"
)

type TreatmentProcedureUseCase struct {
	*UcContract
}

func (uc TreatmentProcedureUseCase) edit(ID, procedureID string) (err error) {
	repository := actions.NewTreatmentProcedureRepository(uc.DB)
	now := time.Now().UTC()

	model := models.TreatmentProcedure{
		ID:          ID,
		ProcedureID: procedureID,
		UpdatedAt:   now,
	}
	err = repository.Edit(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-treatmentProcedure-edit")
		return err
	}

	return nil
}

func (uc TreatmentProcedureUseCase) add(treatmentID, procedureID string) (err error) {
	repository := actions.NewTreatmentProcedureRepository(uc.DB)
	now := time.Now().UTC()

	model := models.TreatmentProcedure{
		TreatmentID: treatmentID,
		ProcedureID: procedureID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err = repository.Add(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-treatmentProcedure-add")
		return err
	}

	return nil
}

func (uc TreatmentProcedureUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewTreatmentProcedureRepository(uc.DB)
	now := time.Now().UTC()

	model := models.TreatmentProcedure{
		UpdatedAt: now,
		DeletedAt: sql.NullTime{Time: now},
	}
	err = repository.DeleteBy(column, value, operator, model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-treatmentProcedure-deleteBy")
		return err
	}

	return nil
}

func (uc TreatmentProcedureUseCase) Store(treatmentID string, inputs []requests.TreatmentProcedureRequest, deletedTreatmentProcedures []string) (err error) {
	var scopes string
	fmt.Println(inputs)
	if len(inputs) > 0 {
		for _, input := range inputs {
			if input.ID == "" {
				err = uc.add(treatmentID, input.ProcedureID)
				scopes = "uc-treatmentProcedure-add"
			} else {
				err = uc.edit(input.ID, input.ProcedureID)
				scopes = "uc-treatmentProcedure-edit"
			}
			if err != nil {
				logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), scopes)
				return err
			}
		}
	}

	if len(deletedTreatmentProcedures) > 0 {
		for _, ID := range deletedTreatmentProcedures {
			err = uc.DeleteBy("id", ID, "=")
			if err != nil {
				logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-treatmentProcedure-deleteBy")
				return err
			}
		}
	}

	return nil
}
