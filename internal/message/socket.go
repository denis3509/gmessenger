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
	msg, err := r.service.CreateMessage(client.GetUserId(), data.RecipientId, data.Text)
	if err != nil {
		errMsg := socket.SocketMessage{
			Event:   "message:create__failure",
			Payload: err.Error()}
		client.Send(errMsg)
		return
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		errMsg := socket.SocketMessage{
			Event:   "message:create__failure",
			Payload: err.Error()}
		client.Send(errMsg)
		return
	}
	client.Send(socket.SocketMessage{
		Event:   "message:create__success",
		Payload: string(jsonMsg)})
	newMessage := socket.SocketMessage{
		Event:   "message:new",
		Payload: string(jsonMsg)}
		
	r.hub.SendByUserId(msg.RecipientId, newMessage)

}
