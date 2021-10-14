package usecase

import (
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/str"
	"profira-backend/usecase/viewmodel"
	"strings"
)

type DoctorUseCase struct {
	*UcContract
}

//browse
func (uc DoctorUseCase) Browse(roleID, search, sort, order string, page, limit int) (res []viewmodel.DoctorVm, pagination viewmodel.PaginationVm, err error) {
	repository := actions.NewStaffRepository(uc.DB)
	offset, limit, page, order, sort := uc.setPaginationParameter(page, limit, order, sort)

	staffs, count, err := repository.BrowseByRole(roleID, search, order, sort, limit, offset)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-staff-browseByRole")
		return res, pagination, err
	}

	for _, staff := range staffs {
		res = append(res, uc.buildBody(staff))
	}
	pagination = uc.setPaginationResponse(page, limit, count)

	return res, pagination, nil
}

//browse all doctor by clinic
func (uc DoctorUseCase) BrowseAllByClinic(clinicId, roleId string) (res []viewmodel.DoctorAvailabilityVm, err error) {
	repository := actions.NewStaffRepository(uc.DB)

	staffs, err := repository.BrowseStaffDoctorByClinic(clinicId, roleId)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-staff-browseStaffDoctorByClinic")
		return res, err
	}

	for _, staff := range staffs {
		doctor := viewmodel.NewDoctorAvailabilityBuilderVm().SetID(staff.ID).SetName(staff.Name).Build()
		res = append(res, doctor)
	}

	return res, nil
}

//browse all
func (uc DoctorUseCase) BrowseAll(roleID, search string) (res []viewmodel.DoctorVm, err error) {
	repository := actions.NewStaffRepository(uc.DB)

	staffs, err := repository.BrowseAllByRole(roleID, search)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-staff-browseAllByRole")
		return res, err
	}

	for _, staff := range staffs {
		res = append(res, uc.buildBody(staff))
	}

	return res, nil
}

//read by
func (uc DoctorUseCase) ReadBy(column, value, operator string) (res viewmodel.DoctorVm, err error) {
	repository := actions.NewStaffRepository(uc.DB)

	staff, err := repository.ReadBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-staff-readBy")
		return res, err
	}
	res = uc.buildBody(staff)

	return res, nil
}

//build doctor clinics
func (uc DoctorUseCase) buildClinic(clinics string) (res []viewmodel.StaffClinicVm) {
	if clinics != "" {
		clinicsArr := str.Unique(strings.Split(clinics, ","))
		for _, clinic := range clinicsArr {
			clinicStringArr := strings.Split(clinic, ":")
			res = append(res, viewmodel.StaffClinicVm{
				ID:          clinicStringArr[0],
				Name:        clinicStringArr[1],
				Address:     clinicStringArr[2],
				PICName:     clinicStringArr[3],
				PhoneNumber: clinicStringArr[4],
				Email:       clinicStringArr[5],
			})
		}
	}

	return res
}

//build doctor treatments
func (uc DoctorUseCase) buildTreatments(treatment string) (res []viewmodel.DoctorTreatmentVm) {
	if treatment != "" {
		treatments := str.Unique(strings.Split(treatment, ","))
		for _, treatment := range treatments {
			treatmentArr := strings.Split(treatment, ":")
			res = append(res, viewmodel.DoctorTreatmentVm{
				ID:          treatmentArr[0],
				TreatmentID: treatmentArr[1],
				Name:        treatmentArr[2],
			})
		}
	}

	return res
}

//build body
func (uc DoctorUseCase) buildBody(model models.Staff) viewmodel.DoctorVm {
	//doctor clinics
	clinics := uc.buildClinic(model.Clinics)
	//doctor treatment
	treatments := uc.buildTreatments(model.Treatments.String)
	//get file
	fileUc := FileUseCase{UcContract: uc.UcContract}
	file, err := fileUc.Read(model.User.ProfilePictureID.String)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-file-read")
	}

	return viewmodel.DoctorVm{
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
		Clinics:    clinics,
		Treatments: treatments,
	}
}
