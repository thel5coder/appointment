package usecase

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/datetime"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/server/requests"
	"time"
)

type DoctorWeekDateScheduleUseCase struct {
	*UcContract
}

//store
func(uc DoctorWeekDateScheduleUseCase) Store(doctorWeekDayScheduleID string,inputs []requests.DoctorScheduleDayRequest) (err error){

	return nil
}

//add
func (uc DoctorWeekDateScheduleUseCase) add(doctorWeekDayScheduleID string, weekDate string, isPresent bool) (err error) {
	repository := actions.NewDoctorWeekDateScheduleRepository(uc.DB)
	now := time.Now().UTC()

	model := models.DoctorWeekDateSchedule{
		WeekDayScheduleID: doctorWeekDayScheduleID,
		WeekDate:          datetime.StrParseToTime(weekDate, "2006-01-02"),
		IsPresent:         isPresent,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	err = repository.Add(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-doctorWeekDateSchedule-add")
		return err
	}

	return nil
}

//delete by
func (uc DoctorWeekDateScheduleUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewDoctorWeekDateScheduleRepository(uc.DB)
	now := time.Now().UTC()

	model := models.DoctorWeekDateSchedule{
		UpdatedAt: now,
		DeletedAt: sql.NullTime{Time: now},
	}
	err = repository.DeleteBy(column, value, operator, model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-doctorWeekDateSchedule-deleteBy")
		return err
	}

	return nil
}
