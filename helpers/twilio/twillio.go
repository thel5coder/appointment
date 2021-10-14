package twilio

import (
	"fmt"
	"github.com/sfreiberg/gotwilio"
)

type Client struct {
	sid           string
	token         string
	DefaultSender string
	twilio        *gotwilio.Twilio
}

func NewTwilioClient(sid, token, defaultSender string) *Client {
	return &Client{
		sid:           sid,
		token:         token,
		DefaultSender: defaultSender,
		twilio:        gotwilio.NewTwilioClient(sid, token),
	}
}

func (client Client) SendSMS(from, to, message string) (err error) {
	preMessage := `My Profira : Kode Otorisasi Anda Adalah `
	smsRes, exception, err := client.twilio.SendSMS(client.DefaultSender, to, preMessage+message, "", "MG1702d8ee40e1bb607922846635abfc02")
	if smsRes != nil {
		fmt.Println("sms response")
		fmt.Println(smsRes)
	}
	if exception != nil {
		fmt.Println("exception")
		fmt.Println(exception)
	}
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
