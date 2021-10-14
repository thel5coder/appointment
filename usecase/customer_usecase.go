package usecase

import (
	"profira-backend/db/models"
	"profira-backend/db/repositories/actions"
	"profira-backend/helpers/amqp"
	"profira-backend/helpers/datetime"
	"profira-backend/helpers/enums"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/hashing"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/str"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"time"
)

type CustomerUseCase struct {
	*UcContract
}

func (uc CustomerUseCase) Browse(search, order, sort string, page, limit int) (res []viewmodel.CustomerVm, pagination viewmodel.PaginationVm, err error) {
	repository := actions.NewCustomerRepository(uc.DB)
	offset, limit, page, order, sort := uc.setPaginationParameter(page, limit, order, sort)

	customers, count, err := repository.Browse(search, order, sort, limit, offset)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-customer-browse")
		return res, pagination, err
	}

	for _, customer := range customers {
		res = append(res, uc.buildBody(customer))
	}
	pagination = uc.setPaginationResponse(page, limit, count)

	return res, pagination, nil
}

func (uc CustomerUseCase) BrowseAllBy(column, value string) (res []viewmodel.CustomerVm, err error) {
	repository := actions.NewCustomerRepository(uc.DB)

	customers, err := repository.BrowseAllBy(column, value)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-customer-browseAllBy")
		return res, err
	}

	for _, customer := range customers {
		res = append(res, uc.buildBody(customer))
	}

	return res, err
}

func (uc CustomerUseCase) ReadBy(column, value, operator string) (res viewmodel.CustomerVm, err error) {
	repository := actions.NewCustomerRepository(uc.DB)
	customer, err := repository.ReadBy(column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-customer-readBy")
		return res, err
	}
	res = uc.buildBody(customer)

	return res, nil
}

func (uc CustomerUseCase) Edit(input *requests.CustomerRequest, ID string) (err error) {
	repository := actions.NewCustomerRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)
	password := ""

	userUc := UserUseCase{UcContract: uc.UcContract}
	if ID != "" {
		user, err := userUc.ReadBy("u.email", input.Email, "=")
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-user-readByEmail")
			return err
		}
		uc.UserID = user.ID
	}

	if input.Password != "" {
		password, _ = hashing.HashAndSalt(input.Password)
	}
	userInput := requests.UserRequest{
		Name:             input.Name,
		Email:            input.Email,
		MobilePhone:      input.MobilePhone,
		Password:         password,
		ProfilePictureID: input.ProfilePictureID,
		IsActive:         input.IsActive,
	}
	err = userUc.Edit(uc.UserID, &userInput)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-user-edit")
		return err
	}

	body := viewmodel.CustomerVm{
		ID:        ID,
		Name:      input.Name,
		Sex:       input.Sex,
		Address:   input.Address,
		BirthDate: input.BirthDate,
		UpdatedAt: now,
		User:      viewmodel.UserCustomerVm{},
	}
	err = repository.Edit(body, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-customer-edit")
		return err
	}

	return nil
}

func (uc CustomerUseCase) EditProfile(input *requests.EditProfileRequest) (err error) {
	userUc := UserUseCase{UcContract: uc.UcContract}
	userEditRequest := requests.UserRequest{
		Name:             input.Name,
		Email:            input.Email,
		MobilePhone:      input.MobilePhone,
		Password:         input.Password,
		ProfilePictureID: input.ProfilePictureID,
		IsActive:         true,
	}
	err = userUc.Edit(uc.UserID, &userEditRequest)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-user-edit")
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	repository := actions.NewCustomerRepository(uc.DB)
	model := models.Customer{
		UserID:    uc.UserID,
		Name:      input.Name,
		BirthDate: datetime.StrParseToTime(input.BirthDate, "2006-01-02"),
		UpdatedAt: now,
	}
	err = repository.EditProfile(model, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-customer-edit-profile")
		return err
	}

	return nil
}

func (uc CustomerUseCase) Add(input *requests.CustomerRequest) (res string, err error) {
	repository := actions.NewCustomerRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)
	password := input.Password

	userUc := UserUseCase{UcContract: uc.UcContract}
	if input.Password == "" {
		password = str.RandomString(6)
	}

	userInput := requests.UserRequest{
		Name:             input.Name,
		Email:            input.Email,
		MobilePhone:      input.MobilePhone,
		Password:         password,
		RoleID:           defaultCustomerRoleID,
		ProfilePictureID: input.ProfilePictureID,
		IsActive:         input.IsActive,
		ActivatedAt:      now,
	}
	res, err = userUc.Add(&userInput)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-user-add")
		return res, err
	}

	body := viewmodel.CustomerVm{
		Name:      input.Name,
		Sex:       input.Sex,
		Address:   input.Address,
		BirthDate: input.BirthDate,
		CreatedAt: now,
		UpdatedAt: now,
		User:      viewmodel.UserCustomerVm{ID: res},
	}
	err = repository.Add(body, uc.TX)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-customer-add")
		return res, err
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

	return res, nil
}

func (uc CustomerUseCase) DeleteBy(column, value, operator string) (err error) {
	repository := actions.NewCustomerRepository(uc.DB)
	now := time.Now().UTC().Format(time.RFC3339)
	customer, err := uc.ReadBy(`c.`+column, value, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-customer-readBy")
		return err
	}

	if err == nil {
		userUc := UserUseCase{UcContract: uc.UcContract}
		err = userUc.DeleteBy("id", customer.User.ID, operator)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-user-deleteBy")
			return err
		}

		err = repository.DeleteBy(column, customer.ID, operator, now, now, uc.TX)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-customer-deleteBy")
			return err
		}
	}

	return nil
}

//Get Profile
func (uc CustomerUseCase) GetMyProfile() (res viewmodel.CustomerVm, err error) {
	res, err = uc.ReadBy("c.user_id", uc.UserID, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-customer-readByUserID")
		return res, err
	}

	return res, nil
}

//count by
func (uc CustomerUseCase) countBy(ID, column, value, operator string) (res int, err error) {
	repository := actions.NewCustomerRepository(uc.DB)
	res, err = repository.CountBy(ID, column, value, operator)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-customer-countBy")
		return res, err
	}

	return res, nil
}

//build response body
func (uc CustomerUseCase) buildBody(model models.Customer) viewmodel.CustomerVm {
	minioUc := MinioUseCase{UcContract: uc.UcContract}
	var err error
	path := ``
	if model.User.ProfilePictureID.String != "" {
		path, err = minioUc.GetFile(model.User.ProfilePicturePath.String)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "uc-minio-getFile")
		}
	}

	return viewmodel.CustomerVm{
		ID:                 model.ID,
		Name:               model.Name,
		Sex:                model.Sex.String,
		Address:            model.Address.String,
		BirthDate:          model.BirthDate.Format("2006-01-02"),
		MaritalStatus:      model.MaritalStatus.String,
		PhoneNumber:        model.PhoneNumber.String,
		MobilePhoneNumber1: model.MobilePhoneNumber1.String,
		MobilePhoneNumber2: model.MobilePhoneNumber2.String,
		Religion:           model.Religion.String,
		Education:          model.Education.String,
		Hobby:              model.Hobby.String,
		Profession:         model.Profession.String,
		Reference:          model.Reference.String,
		Notes:              model.Notes.String,
		CityID:             model.CityID.String,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
		DeletedAt:          model.DeletedAt.String,
		User: viewmodel.UserCustomerVm{
			ID:          model.UserID,
			Email:       model.User.Email,
			MobilePhone: model.User.MobilePhone,
			RoleID:      model.User.Role.ID,
			ProfilePicture: viewmodel.FileVm{
				ID:   model.User.ProfilePictureID.String,
				Path: path,
			},
			IsActive: model.User.IsActive,
		},
	}
}
