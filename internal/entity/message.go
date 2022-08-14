package entity

import "time"

type (
	ChatMessage struct {
		Id                 int
		TextContent        string    `json:"textContent"`
		SenderId           int       `json:"senderId"`
		RecipientId        int       `json:"recipientId"`
		CreatedAt          time.Time `json:"createdAt"`
		UpdatedAt          time.Time `json:"updatedAt"`
		IsRead             bool      `json:"isRead"`
		Edited             bool      `json:"edited"`
		DeletedBySender    bool      `json:"deletedBySender"`
		DeletedByRecipient bool      `json:"deletedByRecipient"`
	}

	Contact struct {
		UserID          int    `json:"userID"`
		Username        string `json:"username"`
		LastMessageText string `json:"lastMessageText"`
		LastMessageRead bool   `json:"lastMessageRead"`
		Unread          int    `json:"unread"`
	}
	
)

func (m ChatMessage) TableName() string {
	return "message"
}
