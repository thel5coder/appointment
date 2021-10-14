package usecase

import (
	"errors"
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/hashing"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/messages"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"strings"
	"time"
)

type UserUseCase struct {
	*UcContract
}

//browse...
func (uc UserUseCase) Browse(search, order, sort string, page, limit int) (res []viewmodel.UserVm, pagination viewmodel.PaginationVm, err error) {
	repository := actions.NewUserRepository(uc.DB)
	offset, limit, page, order, sort := uc.setPaginationParameter(page, limit, order, sort)

	users, count, err := repository.Browse(search, order, sort, limit, offset)
	if err != nil {
		return res, pagination, err
	}

	for _, user := range users {
		res = append(res, uc.buildBody(&user))
	}
	pagination = uc.setPaginationResponse(page, limit, count)

	return res, pagination, nil
}

//browse all ..
func (uc UserUseCase) BrowseAllBy(column, value string) (res []viewmodel.UserVm, err error) {
	repository := actions.NewUserRepository(uc.DB)
	users, err := repository.BrowseAllBy(column, value)
	if err != nil {
		return res, err
	}

	for _, user := range users {
		res = append(res, uc.buildBody(&user))
	}

	return res, nil
}

//read by ...
func (uc UserUseCase) ReadBy(column, value, operator string) (res viewmodel.UserVm, err error) {
	repository := actions.NewUserRepository(uc.DB)
	user, err := repository.ReadBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-user-readBy")
		return res, err
	}
	res = uc.buildBody(&user)

	return res, nil
}

//profile
func (uc UserUseCase) ReadProfile() (res viewmodel.UserVm, err error) {
	res, err = uc.ReadBy("u.id", uc.UserID, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-user-readByID")
		return res, err
	}

	return res, err
}

//edit ..
func (uc UserUseCase) Edit(ID string, input *requests.UserRequest) (err error) {
	repository := actions.NewUserRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)
	password := ""

	isDuplicate, err := uc.checkDuplication(ID, input.Email, input.MobilePhone)
	if err != nil || isDuplicate {
		return err
	}

	if input.Password != "" {
		password, _ = hashing.HashAndSalt(input.Password)
	}
	formattedMobilePhone := input.MobilePhone
	if string(input.MobilePhone[0]) == "0" {
		formattedMobilePhone = strings.Replace(formattedMobilePhone, "0", "62", 1)
	}
	body := viewmodel.UserVm{
		ID:          ID,
		Name:        input.Name,
		Email:       input.Email,
		MobilePhone: formattedMobilePhone,
		ProfilePicture: viewmodel.FileVm{
			ID: input.ProfilePictureID,
		},
		IsActive:  input.IsActive,
		UpdatedAt: now,
	}
	err = repository.Edit(body, password, uc.TX)
	if err != nil {
		return err
	}

	return nil
}

//edit password ...
func (uc UserUseCase) EditPassword(ID, password string) (err error) {
	repository := actions.NewUserRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)
	hashedPassword, _ := hashing.HashAndSalt(password)

	_, err = repository.EditPassword(ID, hashedPassword, now)
	if err != nil {
		return err
	}

	return nil
}

//edit fcm device token ...
func (uc UserUseCase) EditFcmDeviceToken(ID, fcmDeviceToken string) (err error) {
	repository := actions.NewUserRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)

	_, err = repository.EditFcmDeviceToken(ID, fcmDeviceToken, now)
	if err != nil {
		return err
	}

	return nil
}

//edit activated user...
func (uc UserUseCase) EditActivatedUser(ID string, isActive bool) (err error) {
	repository := actions.NewUserRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)

	_, err = repository.EditActivatedUser(ID, now, now, isActive)
	if err != nil {
		return err
	}

	return nil
}

//add user ...
func (uc UserUseCase) Add(input *requests.UserRequest) (res string, err error) {
	repository := actions.NewUserRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)
	activatedAt := ""
	if input.IsActive {
		activatedAt = now
	}

	isDuplicate, err := uc.checkDuplication("", input.Email, input.MobilePhone)
	if err != nil || isDuplicate {
		return res, err
	}

	password, _ := hashing.HashAndSalt(input.Password)

	formattedMobilePhone := input.MobilePhone
	if string(input.MobilePhone[0]) == "0" {
		formattedMobilePhone = strings.Replace(formattedMobilePhone, "0", "62", 1)
	}
	body := viewmodel.UserVm{
		Name:        input.Name,
		Email:       input.Email,
		MobilePhone: formattedMobilePhone,
		ProfilePicture: viewmodel.FileVm{
			ID: input.ProfilePictureID,
		},
		IsActive:    input.IsActive,
		Role:        viewmodel.UserRoleVm{ID: input.RoleID},
		ActivatedAt: activatedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	res, err = repository.Add(body, password, uc.TX)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc UserUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewUserRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)

	count, err := uc.CountBy("", column, value, operator)
	if err != nil {
		return err
	}

	if count > 0 {
		err = repository.DeleteBy(column, value, operator, now, now, uc.TX)
		if err != nil {
			return err
		}
	}

	return nil
}

//check duplication by email and mobile phone...
func (uc UserUseCase) checkDuplication(ID, email, mobilePhone string) (res bool, err error) {
	countMail, err := uc.CountBy(ID, "email", email, "=")
	if err != nil {
		return false, err
	}
	if countMail > 0 {
		return true, errors.New(messages.EmailAlreadyExist)
	}

	countMobilePhone, err := uc.CountBy(ID, "mobile_phone", mobilePhone, "=")
	if err != nil {
		return false, err
	}
	if countMobilePhone > 0 {
		return true, errors.New(messages.PhoneAlreadyExist)
	}

	return false, nil
}

func (uc UserUseCase) IsPasswordValid(ID, password string) (res bool, err error) {
	repository := actions.NewUserRepository(uc.DB)
	user, err := repository.ReadBy("u.id", ID, "=")
	if err != nil {
		return res, err
	}
	res = hashing.CheckHashString(password, user.Password)

	return res, nil
}

//count by ..
func (uc UserUseCase) CountBy(ID, column, value, operator string) (res int, err error) {
	repository := actions.NewUserRepository(uc.DB)
	res, err = repository.CountBy(ID, column, value, operator)
	if err != nil {
		return res, err
	}

	return res, nil
}

//build body...
func (uc UserUseCase) buildBody(model *models.User) viewmodel.UserVm {
	return viewmodel.UserVm{
		ID:             model.ID,
		Name:           model.Name,
		Email:          model.Email,
		MobilePhone:    model.MobilePhone,
		FcmDeviceToken: model.FcmDeviceToken.String,
		ProfilePicture: viewmodel.FileVm{},
		IsActive:       model.IsActive,
		ActivatedAt:    model.ActivatedAt.String,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		Role: viewmodel.UserRoleVm{
			ID:   model.Role.ID,
			Slug: model.Role.Slug,
			Name: model.Role.Name,
		},
	}
}
