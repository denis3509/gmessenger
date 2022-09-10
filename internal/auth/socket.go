package auth

import (
	"encoding/json"
	"messenger/internal/socket"
)

func RegisterHandlers(hub *socket.Hub, service Service) {
	r := resource{service, hub}

	hub.AddHandler("auth:authorize", socket.EventHandler{
		HandleFunc: r.Authorize,
		NeedAuth:   true},
	)
}
type resource struct {
	service Service
	hub     *socket.Hub
}

func (r *resource) Authorize (client *socket.Client, dataRaw []byte) {
	var authData struct{
		Token string `json:"token"`
	}
	json.Unmarshal(dataRaw, &authData)
 
	userId, err := r.service.GetUserId(authData.Token)
	if err != nil {
		errSm := socket.SocketMessage{
			Event: "auth:authorize__error",
			Payload: err.Error(),
		}
		client.Send(errSm)
	} else { 
		r.hub.AddAuthClient(userId, client)
	}

}