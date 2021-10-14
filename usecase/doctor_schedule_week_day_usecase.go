package usecase

import (
	"database/sql"
	"fmt"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/datetime"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/str"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"strings"
	"time"
)

type DoctorScheduleWeekDayUseCase struct {
	*UcContract
}

// BrowseBy browse by
func (uc DoctorScheduleWeekDayUseCase) BrowseBy(filters map[string]interface{}) (res []viewmodel.DoctorScheduleVm, err error) {
	repository := actions.NewDoctorWeekDayScheduleRepository(uc.DB)
	doctorWeekDaySchedules, err := repository.BrowseBy(filters)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-doctorScheduleWeekDay-browseBy")
		return res, err
	}

	for _, doctorWeekDaySchedule := range doctorWeekDaySchedules {
		res = append(res, uc.buildBody(doctorWeekDaySchedule))
	}

	return res, nil
}

// Store store
func (uc DoctorScheduleWeekDayUseCase) Store(input *requests.DoctorScheduleRequest) (err error) {
	time.Now().Weekday()
	//delete week day schedules
	if len(input.DeletedScheduleDay) > 0 {
		for _, ID := range input.DeletedScheduleDay {
			err = uc.DeleteBy("id", ID, "=")
			if err != nil {
				logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-doctorScheduleWeekDay-deleteByID")
				return err
			}
		}
	}

	//delete work time schedule
	doctorWeekDayScheduleWorkTimeUc := DoctorWeekDayScheduleWorkTimeUseCase{UcContract: uc.UcContract}
	if len(input.DeletedScheduleTimes) > 0 {
		for _, ID := range input.DeletedScheduleTimes {
			err = doctorWeekDayScheduleWorkTimeUc.DeleteBy("id", ID, "=")
			if err != nil {
				logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-doctorWeekDayScheduleWorkTime-deleteByID")
				return err
			}
		}
	}

	//add schedule day
	err = uc.add(input.Schedules, input.DoctorID, input.ClinicID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-doctorScheduleWeekDay-add")
		return err
	}

	return nil
}

//add
func (uc DoctorScheduleWeekDayUseCase) add(inputs []requests.DoctorScheduleDayRequest, doctorID, clinicID string) (err error) {
	repository := actions.NewDoctorWeekDayScheduleRepository(uc.DB)
	now := time.Now().UTC()
	doctorWeekDayScheduleWorkTimeUc := DoctorWeekDayScheduleWorkTimeUseCase{UcContract: uc.UcContract}

	for _, input := range inputs {
		ID := input.ID
		if input.ID == "" {
			model := models.DoctorWeekDaySchedule{
				DoctorID:  doctorID,
				ClinicID:  clinicID,
				Day:       input.Day,
				CreatedAt: now,
				UpdatedAt: now,
			}
			ID, err = repository.Add(model, uc.TX)
			if err != nil {
				logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-doctorScheduleWeekDay-add")
				return err
			}
		}

		//add work times
		err = doctorWeekDayScheduleWorkTimeUc.Store(ID, input.ScheduleTimes)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-doctorWeekDayScheduleWorkTime-store")
			return err
		}
	}

	return nil
}

// DeleteBy delete by
func (uc DoctorScheduleWeekDayUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewDoctorWeekDayScheduleRepository(uc.DB)
	now := time.Now().UTC()

	model := models.DoctorWeekDaySchedule{
		UpdatedAt: now,
		DeletedAt: sql.NullTime{Time: now},
	}
	err = repository.DeleteBy(column, value, operator, model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-doctorScheduleWeekDay-deleteBy")
		return err
	}

	//delete work time
	doctorWeekDayScheduleWorkTimeUc := DoctorWeekDayScheduleWorkTimeUseCase{UcContract: uc.UcContract}
	err = doctorWeekDayScheduleWorkTimeUc.DeleteBy("week_day_schedule_id", value, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-doctorWeekDayScheduleWorkTime-deleteByWeekDayScheduleID")
		return err
	}

	return nil
}

//build body doctor week day work times
func (uc DoctorScheduleWeekDayUseCase) buildBodyDoctorWeekDayScheduleWorkTime(workTime string) (res []models.DoctorWeekDayScheduleWorkTimes) {
	workTimes := str.Unique(strings.Split(workTime, ","))
	for _, workTime := range workTimes {
		workTimeArr := strings.Split(workTime, "|")
		fmt.Println(workTimeArr[3])
		fmt.Println(datetime.StrParseToTime(workTimeArr[3], "15:04:05"))
		res = append(res, models.DoctorWeekDayScheduleWorkTimes{
			ID:                workTimeArr[0],
			WeekDayScheduleID: workTimeArr[1],
			StartAt:           datetime.StrParseToTime(workTimeArr[2], "15:04:05"),
			EndAt:             datetime.StrParseToTime(workTimeArr[3], "15:04:05"),
		})
	}

	return res
}

//build body doctor week day
func (uc DoctorScheduleWeekDayUseCase) buildBodyDoctorWeekDay(weekDaySchedule string, weekDayWorkTimes []models.DoctorWeekDayScheduleWorkTimes) (res []viewmodel.DoctorScheduleDaysVm) {
	weekDaySchedules := str.Unique(strings.Split(weekDaySchedule, ","))
	for _, weekDaySchedule := range weekDaySchedules {
		weekDayScheduleArr := strings.Split(weekDaySchedule, "|")
		res = append(res, viewmodel.DoctorScheduleDaysVm{
			ID:        weekDayScheduleArr[0],
			Day:       weekDayScheduleArr[1],
			CreatedAt: weekDayScheduleArr[2],
			UpdatedAt: weekDayScheduleArr[3],
			WorkTimes: nil,
		})
	}

	//add work time to work day
	for i := 0; i < len(res); i++ {
		for _, workTime := range weekDayWorkTimes {
			if res[i].ID == workTime.WeekDayScheduleID {
				res[i].WorkTimes = append(res[i].WorkTimes, viewmodel.DoctorScheduleWorkTimesVm{
					ID:      workTime.ID,
					StartAt: workTime.StartAt.Format("15:04"),
					EndAt:   workTime.EndAt.Format("15:04"),
				})
			}
		}
	}

	return res
}

//build body
func (uc DoctorScheduleWeekDayUseCase) buildBody(model models.DoctorWeekDaySchedule) viewmodel.DoctorScheduleVm {
	var doctorScheduleWorkTimeModel []models.DoctorWeekDayScheduleWorkTimes
	if model.ScheduleTimes.String != "" {
		doctorScheduleWorkTimeModel = uc.buildBodyDoctorWeekDayScheduleWorkTime(model.ScheduleTimes.String)
	}

	var doctorScheduleWeekDay []viewmodel.DoctorScheduleDaysVm
	if model.ScheduleDays.String != "" {
		doctorScheduleWeekDay = uc.buildBodyDoctorWeekDay(model.ScheduleDays.String, doctorScheduleWorkTimeModel)
	}

	return viewmodel.DoctorScheduleVm{
		DoctorID:   model.DoctorID,
		DoctorName: model.Staff.Name,
		ClinicID:   model.ClinicID,
		ClinicName: model.Clinic.Name,
		Schedules:  doctorScheduleWeekDay,
	}
}
