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

type DoctorTreatmentUseCase struct {
	*UcContract
}

//store
func (uc DoctorTreatmentUseCase) Store(doctorID string, input requests.DoctorTreatmentRequest) (err error) {
	if len(input.Selected) > 0 {
		for _, doctorTreatment := range input.Selected {
			err = uc.add(doctorID, doctorTreatment)
			if err != nil {
				logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-doctorTreatment-add")
				return err
			}
		}
	}

	if len(input.Deleted) > 0 {
		for _, doctorTreatment := range input.Deleted {
			err = uc.DeleteBy("id", doctorTreatment, "=")
			if err != nil {
				logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-doctorTreatment-deleteByID")
				return err
			}
		}
	}

	return nil
}

//add
func (uc DoctorTreatmentUseCase) add(doctorID, treatmentID string) (err error) {
	repository := actions.NewDoctorTreatmentRepository(uc.DB)
	now := time.Now().UTC()
	model := models.DoctorTreatment{
		DoctorID:    doctorID,
		TreatmentID: treatmentID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err = repository.Add(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-doctorTreatment-add")
		return err
	}

	return nil
}

//delete by
func (uc DoctorTreatmentUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewDoctorTreatmentRepository(uc.DB)
	now := time.Now().UTC()

	count, err := uc.CountBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-doctorTreatment-countBy")
		return err
	}

	if count > 0 {
		model := models.DoctorTreatment{
			UpdatedAt: now,
			DeletedAt: sql.NullTime{Time: now},
		}
		err = repository.DeleteBy(column, value, operator, model, uc.TX)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-doctorTreatment-deleteBy")
			return err
		}
	}

	return nil
}

//count by
func (uc DoctorTreatmentUseCase) CountBy(column, value, operator string) (res int, err error) {
	repository := actions.NewDoctorTreatmentRepository(uc.DB)
	res, err = repository.CountBy(column, value, operator)
	if err != nil {
		return res, err
	}

	return res, nil
}
