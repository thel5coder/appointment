package usecase

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/server/requests"
	"time"
)

type BedTreatmentUseCase struct {
	*UcContract
}

//store
func (uc BedTreatmentUseCase) Store(bedID string, input requests.BedTreatmentRequest) (err error) {
	if len(input.Selected) > 0 {
		for _, treatment := range input.Selected {
			err = uc.add(bedID, treatment)
			if err != nil {
				logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-bedTreatment-add")
				return err
			}
		}
	}

	if len(input.Deleted) > 0 {
		for _, treatment := range input.Deleted {
			err = uc.DeleteBy("id", treatment, "=")
			if err != nil {
				logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-bedTreatment-deleteByID")
				return err
			}
		}
	}

	return nil
}

//add
func (uc BedTreatmentUseCase) add(bedID, treatmentID string) (err error) {
	repository := actions.NewBedTreatmentRepository(uc.DB)
	now := time.Now().UTC()
	model := models.BedTreatment{
		BedID:       bedID,
		TreatmentID: treatmentID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err = repository.Add(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-bedTreatment-add")
		return err
	}

	return nil
}

//delete by
func (uc BedTreatmentUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewBedTreatmentRepository(uc.DB)
	now := time.Now().UTC()
	model := models.BedTreatment{
		UpdatedAt: now,
		DeletedAt: sql.NullTime{Time: now},
	}
	err = repository.DeleteBy(column, value, operator, model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-bedTreatment-deleteBy")
		return err
	}

	return nil
}

//count by
func (uc BedTreatmentUseCase) countBy(column, value, operator string) (res int, err error) {
	repository := actions.NewBedTreatmentRepository(uc.DB)
	res, err = repository.CountBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-bedTreatment-countBy")
		return res, err
	}

	return res, nil
}
