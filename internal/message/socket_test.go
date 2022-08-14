package message

import (
	"encoding/json"
	"fmt"
	// "messenger/internal/entity"
	"messenger/internal/socket"
	"testing"
	// "github.com/brianvoe/gofakeit"
)

func FakeMessage() {

}

func Test_resource_createMessage(t *testing.T) {
	type fields struct {
		service Service
		hub     *socket.Hub
	}
	type args struct {
		client  *socket.Client
		dataRaw []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &resource{
				service: tt.fields.service,
				hub:     tt.fields.hub,
			}
			r.createMessage(tt.args.client, tt.args.dataRaw)
		})
	}
}

func TestMarshal(t *testing.T) {
	faker := Faker{}
	msg := faker.Message()  

	data, _ := json.Marshal(&msg)
	t.Log(string(data))
	fmt.Println()
}
 