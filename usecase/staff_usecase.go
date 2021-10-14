package usecase

import (
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/amqp"
	"profira-backend/helpers/enums"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/str"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"strings"
	"time"
)

type StaffUseCase struct {
	*UcContract
}

//browse pagination by role
func (uc StaffUseCase) BrowseByRole(roleID, search, sort, order string, page, limit int) (res []viewmodel.StaffVm, pagination viewmodel.PaginationVm, err error) {
	repository := actions.NewStaffRepository(uc.DB)
	offset, limit, page, order, sort := uc.setPaginationParameter(page, limit, order, sort)

	staffs, count, err := repository.BrowseByRole(roleID, search, order, sort, limit, offset)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-staff-browseByRole")
		return res,pagination,err
	}

	for _, staff := range staffs {
		res = append(res, uc.buildBody(staff))
	}
	pagination = uc.setPaginationResponse(page, limit, count)

	return res, pagination, nil
}

//browse all by role
func (uc StaffUseCase) BrowseAllByRole(roleID, search string) (res []viewmodel.StaffVm, err error) {
	repository := actions.NewStaffRepository(uc.DB)

	staffs, err := repository.BrowseAllByRole(roleID, search)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-staff-browseAllByRole")
		return res, err
	}

	for _, staff := range staffs {
		res = append(res, uc.buildBody(staff))
	}

	return res, nil
}

//read by
func (uc StaffUseCase) ReadBy(column, value, operator string) (res viewmodel.StaffVm, err error) {
	repository := actions.NewStaffRepository(uc.DB)

	staff, err := repository.ReadBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-staff-readBy")
		return res, err
	}
	res = uc.buildBody(staff)

	return res, nil
}

//edit
func (uc StaffUseCase) Edit(input *requests.StaffRequest, ID string) (err error) {
	repository := actions.NewStaffRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)

	staff, err := uc.ReadBy("s.id", ID, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-staff-readByID")
		return err
	}

	//edit user
	userUc := UserUseCase{UcContract: uc.UcContract}
	userInput := requests.UserRequest{
		Name:             input.Name,
		Email:            input.Email,
		MobilePhone:      input.MobilePhone,
		Password:         input.Password,
		ProfilePictureID: input.ProfilePictureID,
		IsActive:         input.IsActive,
	}
	err = userUc.Edit(staff.User.ID, &userInput)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-user-edit")
		return err
	}

	//edit staff
	body := viewmodel.StaffVm{
		ID:          ID,
		Name:        input.Name,
		Description: input.Description,
		UpdatedAt:   now,
	}
	err = repository.Edit(body, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-staff-edit")
		return err
	}

	//add staff clinic
	if len(input.ClinicIDs) > 0 {
		staffClinicUc := StaffClinicUseCase{UcContract: uc.UcContract}
		err = staffClinicUc.Store(ID, input.ClinicIDs)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-staffClinic-store")
			return err
		}
	}

	//add treatment
	doctorTreatmentUc := DoctorTreatmentUseCase{UcContract:uc.UcContract}
	err = doctorTreatmentUc.Store(ID,input.Treatment)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-doctorTreatment-store")
		return err
	}

	return nil
}

//add
func (uc StaffUseCase) Add(input *requests.StaffRequest, roleID string) (err error) {
	repository := actions.NewStaffRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)

	//user add
	userUc := UserUseCase{UcContract: uc.UcContract}
	password := input.Password
	if input.Password == "" {
		password = str.RandomString(6)
	}

	userInput := requests.UserRequest{
		Name:             input.Name,
		Email:            input.Email,
		MobilePhone:      input.MobilePhone,
		Password:         password,
		RoleID:           roleID,
		ProfilePictureID: input.ProfilePictureID,
		IsActive:         true,
	}
	userID, err := userUc.Add(&userInput)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-user-add")
		return err
	}

	//add staff
	body := viewmodel.StaffVm{
		Name:        input.Name,
		Description: input.Description,
		User:        viewmodel.UserStaffVm{ID: userID},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	staffID, err := repository.Add(body, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-staff-add")
		return err
	}

	//add staff clinic
	staffClinicUc := StaffClinicUseCase{UcContract: uc.UcContract}
	err = staffClinicUc.Store(staffID, input.ClinicIDs)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-staffClinic-store")
		return err
	}

	//add treatment
	doctorTreatmentUc := DoctorTreatmentUseCase{UcContract:uc.UcContract}
	err = doctorTreatmentUc.Store(staffID,input.Treatment)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-doctorTreatment-store")
		return err
	}

	queueBody := map[string]interface{}{
		"qid": uc.ReqID,
		"payload": map[string]interface{}{
			"email":    input.Email,
			"password": password,
			"type":     enums.MailTypeEnums[1],
		},
	}
	err = uc.PushToQueue(queueBody, amqp.MailIncoming, amqp.MailDeadLetter)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-queue-pushToQueue")
		return err
	}

	return nil
}

func (uc StaffUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewStaffRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)

	count, err := uc.countBy("", column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-staff-countBy")
		return err
	}

	if count > 0 {
		staff, err := uc.ReadBy(column, value, operator)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-staff-readBy")
			return err
		}

		//delete user
		userUc := UserUseCase{UcContract: uc.UcContract}
		err = userUc.DeleteBy("id", staff.User.ID, operator)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-user-DeleteByID")
			return err
		}

		//delete staff
		err = repository.DeleteBy("id", value, operator, now, now, uc.TX)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"query-staff-delete")
			return err
		}

		//delete staff clinic
		staffClinicUc := StaffClinicUseCase{UcContract: uc.UcContract}
		err = staffClinicUc.Delete("staff_id", value, operator)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-staffClinic-delete")
			return err
		}

		//delete treatment if doctor
		doctorTreatmentUc := DoctorTreatmentUseCase{UcContract:uc.UcContract}
		err = doctorTreatmentUc.DeleteBy("doctor_id",value,"=")
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-doctorTreatment-deleteByDoctorID")
			return err
		}

	}

	return nil
}

func (uc StaffUseCase) countBy(ID, column, value, operator string) (res int, err error) {
	repository := actions.NewStaffRepository(uc.DB)
	res, err = repository.CountBy(ID, column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-staff-countBy")
		return res, err
	}

	return res, nil
}

func (uc StaffUseCase) buildBody(model models.Staff) viewmodel.StaffVm {
	//parse staff clinic
	var staffClinics []viewmodel.StaffClinicVm
	clinicsArr := strings.Split(model.Clinics, ",")
	for _, clinic := range clinicsArr {
		clinicStringArr := strings.Split(clinic, ":")
		staffClinics = append(staffClinics, viewmodel.StaffClinicVm{
			ID:          clinicStringArr[0],
			Name:        clinicStringArr[1],
			Address:     clinicStringArr[2],
			PICName:     clinicStringArr[3],
			PhoneNumber: clinicStringArr[4],
			Email:       clinicStringArr[5],
		})
	}

	//get file
	fileUc := FileUseCase{UcContract: uc.UcContract}
	file, err := fileUc.Read(model.User.ProfilePictureID.String)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-file-read")
	}

	return viewmodel.StaffVm{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		User: viewmodel.UserStaffVm{
			ID:             model.UserID,
			Email:          model.User.Email,
			MobilePhone:    model.User.MobilePhone,
			RoleID:         model.User.Role.ID,
			ProfilePicture: file,
			IsActive:       model.User.IsActive,
		},
		Clinics: staffClinics,
	}
}
