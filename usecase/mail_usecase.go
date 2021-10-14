package usecase

import (
	"fmt"
	"io/ioutil"
	"profira-backend/helpers/enums"
	"profira-backend/helpers/str"
	"strings"
)

type MailUseCase struct {
	*UcContract
}

//Send mail...
func (uc MailUseCase) SendMail(payload map[string]interface{}) (res string, err error) {
	if payload["type"] == enums.MailTypeEnums[0] {
		err = uc.SendMailForgotPassword(payload["email"].(string))
	} else if payload["type"] == enums.MailTypeEnums[1] {
		err = uc.SendUserRandomPassword(payload)
	}
	return res, nil
}

//send mail forgot password...
func (uc MailUseCase) SendMailForgotPassword(email string) (err error) {
	userUc := UserUseCase{UcContract: uc.UcContract}
	user, err := userUc.ReadBy("u.email", email, "=")
	if err != nil {
		return err
	}

	//edit password
	password := str.RandomString(6)
	err = userUc.EditPassword(user.ID, password)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//setup template
	htmlFile := "./../statics/forgot-password.html"
	template, err := ioutil.ReadFile(htmlFile)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	templateStr := string(template)
	templateStr = strings.Replace(templateStr, "[password]", password, 1)

	//send new password to email
	err = uc.GoMailConfig.SendGoMail(user.Email, enums.MailSubjectEnums[0], templateStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

//send mail user random password
func (uc MailUseCase) SendUserRandomPassword(payload map[string]interface{}) (err error) {
	//setup template
	htmlFile := "./../statics/user-random-password.html"
	template, err := ioutil.ReadFile(htmlFile)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	templateStr := string(template)
	templateStr = strings.Replace(templateStr, "[password]", payload["password"].(string), 1)
	templateStr = strings.Replace(templateStr, "[email]", payload["email"].(string), 1)

	//send new password to email
	err = uc.GoMailConfig.SendGoMail(payload["email"].(string), enums.MailSubjectEnums[1], templateStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
