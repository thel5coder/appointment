package pusher

import (
	"github.com/pusher/pusher-http-go"
)

var (
	// channelDefault ...
	channelDefault = "qibla"

	// Event New & Updated Stage ...
	EventStory = "story"

	//Event New Board
	EventBoard = "board"
)

// Credential ...
type Credential struct {
	AppID   string
	Key     string
	Secret  string
	Cluster string
}

// Send ...
func (cred *Credential) Send(eventName string, data interface{}) (err error) {
	client := pusher.Client{
		AppID:   cred.AppID,
		Key:     cred.Key,
		Secret:  cred.Secret,
		Cluster: cred.Cluster,
		Secure:  true,
	}

	err = client.Trigger(channelDefault, eventName, data)

	return err
}
