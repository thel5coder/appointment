package fcm

import (
	"github.com/maddevsio/fcm"
	"profira-backend/helpers/interfacepkg"
)

// Connection ...
type Connection struct {
	APIKey string
}

// SendAndroid ...
func (cred *Connection) SendAndroid(to []string, title, body string, data map[string]interface{}) (string, error) {
	c := fcm.NewFCM(cred.APIKey)
	response, err := c.Send(fcm.Message{
		Data:             data,
		RegistrationIDs:  to,
		ContentAvailable: true,
		Priority:         fcm.PriorityHigh,
		Notification: fcm.Notification{
			Title: title,
			Body:  body,
			Sound: "default",
			Badge: "3",
		},
	})
	if err != nil {
		return "", err
	}

	res := interfacepkg.Marshall(response)

	return res, err
}
