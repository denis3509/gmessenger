package message

import (
	"encoding/json"
	// "fmt"
	"messenger/internal/socket"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func RegisterHandlers(hub *socket.Hub, service Service) {
	r := resource{service, hub}

	hub.AddHandler("message:create", socket.EventHandler{
		HandleFunc: r.createMessage,
		NeedAuth:   true},
	)
}

type resource struct {
	service Service
	hub     *socket.Hub
}

type CreateMessageData struct {
	RecipientId int    `json:"recipientId"`
	Text        string `json:"text"`
}

func (m CreateMessageData) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.RecipientId, validation.Required, validation.Min(1)),
		validation.Field(&m.Text, validation.Required, validation.Length(0, 1024)),
	)
}
func (r *resource) createMessage(client *socket.Client, dataRaw []byte) {

	var data CreateMessageData
	json.Unmarshal(dataRaw, &data)
	msg, err := r.service.CreateMessage(client.UserId, data.RecipientId, data.Text)
	if err != nil {
		client.Send("message:create__failure", []byte(err.Error()))
		return
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		client.Send("message:create__failure", []byte(err.Error()))
		return
	}
	client.Send("message:create__success", jsonMsg)
	recipient := r.hub.ClientById(msg.RecipientId)
	if recipient != nil {
		recipient.Send("message:new", jsonMsg)
	}

}
