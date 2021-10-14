package mailing

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type Mail struct {

}

func SendMail(to []string,subject,message string) error{
	from := os.Getenv("MAIL_SENDER")
	body := "From : "+from+"\n"+
		"To: "+ strings.Join(to,",")+"\n"+
		"Subject: "+subject+"\n\n"+
		message

	auth := smtp.PlainAuth("",from,os.Getenv("MAIL_SENDER"),os.Getenv("SMTP_HOST"))
	smtpAddress := fmt.Sprintf("%s:%s",os.Getenv("smtp_host"),os.Getenv("smtp_port"))
	err := smtp.SendMail(smtpAddress,auth,from,append(to),[]byte(body))
	if err != nil {
		return err
	}

	return nil
}
