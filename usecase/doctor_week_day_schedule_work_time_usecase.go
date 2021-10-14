package usecase

import (
	"database/sql"
	"fmt"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/datetime"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/server/requests"
	"time"
)

type DoctorWeekDayScheduleWorkTimeUseCase struct {
	*UcContract
}

//store
func (uc DoctorWeekDayScheduleWorkTimeUseCase) Store(doctorWeekDayScheduleID string, inputs []requests.DoctorScheduleTimeRequest) (err error) {
	scope := ``

	for _, input := range inputs {
		if input.ID != "" {
			err = uc.edit(input.ID, input.StartAt, input.EndAt)
			scope = `uc-doctorWeekDayScheduleWorkTime-edit`
		} else {
			err = uc.add(doctorWeekDayScheduleID, input.StartAt, input.EndAt)
			scope = `uc-doctorWeekDayScheduleWorkTime-add`
		}
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), scope)
			return err
		}
	}

	return nil
}

//add
func (uc DoctorWeekDayScheduleWorkTimeUseCase) add(doctorScheduleWeekDayID, startAt, endAt string) (err error) {
	repository := actions.NewDoctorWorkDayScheduleWorkTimeRepository(uc.DB)
	now := time.Now().UTC()
	fmt.Println(endAt)

	model := models.DoctorWeekDayScheduleWorkTimes{
		WeekDayScheduleID: doctorScheduleWeekDayID,
		StartAt:           datetime.StrParseToTime(startAt, "15:04"),
		EndAt:             datetime.StrParseToTime(endAt, "15:04"),
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	err = repository.Add(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-doctorWeekDayScheduleWorkTime-add")
		return err
	}

	return nil
}

//edit
func (uc DoctorWeekDayScheduleWorkTimeUseCase) edit(ID, startAt, endAt string) (err error) {
	repository := actions.NewDoctorWorkDayScheduleWorkTimeRepository(uc.DB)
	now := time.Now().UTC()

	model := models.DoctorWeekDayScheduleWorkTimes{
		ID:        ID,
		StartAt:   datetime.StrParseToTime(startAt, "15:04"),
		EndAt:     datetime.StrParseToTime(endAt, "15:04"),
		UpdatedAt: now,
	}
	err = repository.Edit(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-doctorWeekDayScheduleWorkTime-edit")
		return err
	}

	return nil
}

//delete by
func (uc DoctorWeekDayScheduleWorkTimeUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewDoctorWorkDayScheduleWorkTimeRepository(uc.DB)
	now := time.Now().UTC()

	model := models.DoctorWeekDayScheduleWorkTimes{
		UpdatedAt: now,
		DeletedAt: sql.NullTime{Time: now},
	}
	err = repository.DeleteBy(column, value, operator, model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-doctorWeekDayScheduleWorkTime-deleteBy")
		return err
	}

	return nil

}
