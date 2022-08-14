package message

import (
	"errors"
	"messenger/internal/entity"
	"time"
)

type Emitter interface {
	Emit(event string, data []byte)
}
type Service interface {
	CreateMessage(senderId int, recipientId int, text string) (*entity.ChatMessage, error)
	GetContactList(userId int) (*ContactList, error)
	Dialog(userId int, contactId int, limit int, page int) []entity.ChatMessage
}
type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return service{repository}
}

func (s service) CreateMessage(senderId int, recipientId int, text string) (*entity.ChatMessage, error) {
	now := time.Now()
	msg := entity.ChatMessage{
		TextContent:        text,
		SenderId:           senderId,
		RecipientId:        recipientId,
		CreatedAt:          now,
		UpdatedAt:          now,
		IsRead:             false,
		Edited:             false,
		DeletedBySender:    false,
		DeletedByRecipient: false,
	}

	if senderId == recipientId {
		return nil, errors.New("sending message to the same user is not allowed")
	}

	id, err := s.repository.Create(msg)
	if err != nil {
		return nil, err
	}
	msg, err = s.repository.Get(id)
	return &msg, err

}

type ContactList struct {
	TotalUnread int
	contacts    []entity.Contact
}

func (s service) GetContactList(userId int) (*ContactList, error) {
	contacts, err := s.repository.Contacts(userId)
	if err != nil {
		return nil, err
	}
	totalUnread := 0
	for _, v := range contacts {
		totalUnread += v.Unread
	}
	return &ContactList{totalUnread, contacts}, nil

}

func (s service) Dialog(userId int, contactId int, limit int, page int) []entity.ChatMessage {
	return []entity.ChatMessage{}
	
}
