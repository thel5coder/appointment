package usecase

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"profira-backend/helpers/amqp"
	"profira-backend/helpers/messages"
	"profira-backend/helpers/str"
	"profira-backend/server/requests"
	"strings"
	"time"
)

type OtpUseCase struct {
	*UcContract
}

func (uc OtpUseCase) RequestOTP(mobilePhone string) (err error) {
	formattedMobilePhone := mobilePhone
	if string(mobilePhone[0]) == "0" {
		formattedMobilePhone = strings.Replace(formattedMobilePhone, "0", "62", 1)
	}
	err = uc.LimitRetryByKey(`limitDailyOTP-`+formattedMobilePhone, 3)
	if err != nil {
		if err.Error() == messages.MaxRetryKey {
			return errors.New(messages.MaxRequestOTPDaily)
		}

		return err
	}

	queueBody := map[string]interface{}{
		"qid":         uc.ReqID,
		"mobilePhone": formattedMobilePhone,
	}
	err = uc.PushToQueue(queueBody, amqp.SMSIncoming, amqp.SMSDeadLetter)
	if err != nil {
		return err
	}

	return nil
}

func (uc OtpUseCase) SendOTP(mobilePhone string) (res map[string]interface{}, err error) {
	otp := str.RandomNumberString(4)
	res = map[string]interface{}{
		"otp": otp,
	}

	formattedMobilePhone := mobilePhone
	if string(mobilePhone[0]) == "0" {
		formattedMobilePhone = strings.Replace(formattedMobilePhone, "0", "62", 1)
	}

	err = uc.RedisClient.StoreToRedistWithExpired(`+`+formattedMobilePhone, res, "1m")
	if err != nil {
		return res, err
	}

	err = uc.TwilioClient.SendSMS(uc.TwilioClient.DefaultSender, `+`+formattedMobilePhone, otp)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc OtpUseCase) SubmitOTP(input *requests.SubmitOTPRequest) (err error) {
	otpRes := map[string]interface{}{}
	err = uc.RedisClient.GetFromRedis(input.MobilePhone, &otpRes)
	if err != nil {
		return errors.New(messages.InvalidKey)
	}

	if otpRes["otp"] != input.OTP {
		err = uc.LimitRetryByKey(input.MobilePhone+`retrySubmitOTP`, 3)
		if err != nil {
			return err
		}

		return errors.New(messages.InvalidOTP)
	}

	return nil
}

func (uc OtpUseCase) RequestXAPIKey() string {
	var sha1 = sha1.New()
	now := time.Now().UTC().Format("2006-01-02T15:04Z07:00")

	fmt.Println(now)

	sha1.Write([]byte(now))
	sha1EncryptedStr := sha1.Sum(nil)
	encrypted := fmt.Sprintf("%x", sha1EncryptedStr)

	return encrypted
}
