package usecase

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"os"
	"profira-backend/helpers/amqp"
	"profira-backend/helpers/enums"
	"profira-backend/helpers/functioncaller"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/messages"
	"profira-backend/server/requests"
	"profira-backend/usecase/viewmodel"
	"profira-backend/usecase/viewmodel/builders"
	"strings"
)

type AuthenticationUseCase struct {
	*UcContract
}

//update session login
func (uc AuthenticationUseCase) UpdateSessionLogin(ID string) (res string, err error) {
	value := uuid.NewV4().String()
	exp := os.Getenv("SESSION_EXP")
	key := "session-" + ID
	resSession := viewmodel.UserSessionVm{}
	resSession.Session = value

	uc.RedisClient.StoreToRedistWithExpired(key, resSession, exp)

	return value, err
}

//generate jwt token
func (uc AuthenticationUseCase) GenerateJwtToken(jwePayload, email, session string) (token, refreshToken, expTokenAt, expRefreshTokenAt string, err error) {
	token, expTokenAt, err = uc.JwtCred.GetToken(session, email, jwePayload)
	if err != nil {
		return token, refreshToken, expTokenAt, expRefreshTokenAt, err
	}

	refreshToken, expRefreshTokenAt, err = uc.JwtCred.GetRefreshToken(session, email, jwePayload)
	if err != nil {
		return token, refreshToken, expTokenAt, expRefreshTokenAt, err
	}

	return token, refreshToken, expTokenAt, expRefreshTokenAt, err
}

func (uc AuthenticationUseCase) Login(input *requests.LoginRequest) (res viewmodel.UserJwtTokenVm, err error) {
	var user viewmodel.UserVm
	var issuer string

	if input.MobilePhone != "" {
		user, err = uc.LoginByByMobilePhone(input.MobilePhone, input.Password, input.FcmDeviceToken)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-loginByMobilePhone")
			return res, err
		}
		issuer = user.MobilePhone
	} else {
		user, err = uc.LoginByEmail(input.Email, input.Password)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-loginByEmail")
			return res, err
		}
		issuer = user.Email
	}

	jwePayload, _ := uc.Jwe.GenerateJwePayload(user.ID)
	session, _ := uc.UpdateSessionLogin(user.ID)
	token, refreshToken, tokenExpiredAt, refreshTokenExpiredAt, err := uc.GenerateJwtToken(jwePayload, issuer, session)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-generateJwtToken")
		return res, err
	}

	res = builders.NewUserTokenJwtVmBuilder().SetToken(token).SetExpiredToken(tokenExpiredAt).SetRefreshToken(refreshToken).
		SetExpiredRefreshToken(refreshTokenExpiredAt).SetIsActive(user.IsActive).GetUserTokenJwtVm()

	return res, nil
}

func (uc AuthenticationUseCase) LoginByEmail(email, password string) (res viewmodel.UserVm, err error) {
	userUc := UserUseCase{UcContract: uc.UcContract}

	count, err := userUc.CountBy("", "email", email, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-user-countByEmail")
		return res, errors.New(messages.CredentialDoNotMatch)
	}
	if count == 0 {
		logruslogger.Log(logruslogger.WarnLevel,messages.CredentialDoNotMatch,functioncaller.PrintFuncName(),"uc-user-countByEmail")
		return res, errors.New(messages.CredentialDoNotMatch)
	}

	res, err = userUc.ReadBy("u.email", email, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-user-readByEmail")
		return res, errors.New(messages.CredentialDoNotMatch)
	}

	isPasswordValid, err := userUc.IsPasswordValid(res.ID, password)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-user-isPasswordValid")
		return res, errors.New(messages.CredentialDoNotMatch)
	}
	if !isPasswordValid{
		logruslogger.Log(logruslogger.WarnLevel,messages.CredentialDoNotMatch,functioncaller.PrintFuncName(),"uc-user-isPasswordValid")
		return res, errors.New(messages.CredentialDoNotMatch)
	}

	return res, err
}

// LoginByByMobilePhone login by mobile phone
func (uc AuthenticationUseCase) LoginByByMobilePhone(mobilePhone, password, fcmDeviceToken string) (res viewmodel.UserVm, err error) {
	userUc := UserUseCase{UcContract: uc.UcContract}

	formattedMobilePhone := mobilePhone
	if string(mobilePhone[0]) == "0" {
		formattedMobilePhone = strings.Replace(formattedMobilePhone, "0", "62", 1)
	}

	count, err := userUc.CountBy("", "mobile_phone", formattedMobilePhone, "=")
	if err != nil  {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-user-countByMobilePhone")
		return res, errors.New(messages.CredentialDoNotMatch)
	}
	if count == 0 {
		logruslogger.Log(logruslogger.WarnLevel,messages.CredentialDoNotMatch,functioncaller.PrintFuncName(),"uc-user-countByMobilePhone")
		return res,errors.New(messages.CredentialDoNotMatch)
	}

	res, err = userUc.ReadBy("u.mobile_phone", formattedMobilePhone, "=")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,messages.CredentialDoNotMatch,functioncaller.PrintFuncName(),"uc-user-readByMobilePhone")
		return res, errors.New(messages.CredentialDoNotMatch)
	}

	err = userUc.EditFcmDeviceToken(res.ID, fcmDeviceToken)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,messages.CredentialDoNotMatch,functioncaller.PrintFuncName(),"uc-user-editFcmToken")
		return res, err
	}

	isPasswordValid, err := userUc.IsPasswordValid(res.ID, password)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel,err.Error(),functioncaller.PrintFuncName(),"uc-user-isPasswordValid")
		return res, errors.New(messages.CredentialDoNotMatch)
	}
	if !isPasswordValid {
		logruslogger.Log(logruslogger.WarnLevel,messages.CredentialDoNotMatch,functioncaller.PrintFuncName(),"uc-user-isPasswordValid")
		return res, errors.New(messages.CredentialDoNotMatch)
	}

	return res, err
}

// Registration registration
func (uc AuthenticationUseCase) Registration(input *requests.RegisterRequest) (res viewmodel.UserJwtTokenVm, err error) {
	customerUc := CustomerUseCase{UcContract: uc.UcContract}
	var customerInput requests.CustomerRequest
	customerInput = requests.CustomerRequest{
		Name:        input.Name,
		BirthDate:   input.BirthDate,
		Email:       input.Email,
		MobilePhone: input.MobilePhone,
		Password:    input.Password,
		Sex:         defaultSex,
		IsActive:    false,
	}

	userID, err := customerUc.Add(&customerInput)
	if err != nil {
		return res, err
	}

	jwePayload, _ := uc.Jwe.GenerateJwePayload(userID)
	session, _ := uc.UpdateSessionLogin(userID)
	token, refreshToken, tokenExpiredAt, refreshTokenExpiredAt, err := uc.GenerateJwtToken(jwePayload, input.Email, session)
	if err != nil {
		return res, err
	}

	res = viewmodel.UserJwtTokenVm{
		Token:           token,
		ExpTime:         tokenExpiredAt,
		RefreshToken:    refreshToken,
		ExpRefreshToken: refreshTokenExpiredAt,
		IsActive:        false,
	}

	return res, nil
}

// ActivationCustomer activation user customer
func (uc AuthenticationUseCase) ActivationCustomer(input *requests.ActivationCustomerRequest) (err error) {
	userUc := UserUseCase{UcContract: uc.UcContract}
	user, err := userUc.ReadBy("u.id", uc.UserID, "=")
	if err != nil {
		return err
	}

	otpUc := OtpUseCase{UcContract: uc.UcContract}
	otpInput := requests.SubmitOTPRequest{
		Type:        enums.OTPSubmitTypeEnums[0],
		MobilePhone: `+`+user.MobilePhone,
		OTP:         input.OTP,
	}
	err = otpUc.SubmitOTP(&otpInput)
	if err != nil {
		return err
	}

	err = userUc.EditActivatedUser(uc.UserID, true)
	if err != nil {
		return err
	}

	return nil
}

// ForgotPassword forgot password
func (uc AuthenticationUseCase) ForgotPassword(email string) (err error) {
	userUc := UserUseCase{UcContract: uc.UcContract}

	count, err := userUc.CountBy("", "email", email, "=")
	if count > 0 && err == nil {
		queueBody := map[string]interface{}{
			"qid": uc.ReqID,
			"payload": map[string]interface{}{
				"email": email,
				"type":  enums.MailTypeEnums[0],
			},
		}
		err = uc.PushToQueue(queueBody, amqp.MailIncoming, amqp.MailDeadLetter)
		if err != nil {
			return err
		}
	} else {
		return errors.New(messages.EmailNotFound)
	}

	return nil
}
