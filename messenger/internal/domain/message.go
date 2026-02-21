package domain

type Message struct {
	UserID int    `json:"user_id"`
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type MessageField func(msg *Message) error

func NewMessage(fields ...MessageField) (*Message, error) {
	msg := Message{}
	for _, field := range fields {
		err := field(&msg)
		if err != nil {
			return nil, err
		}
	}
	return &msg, nil
}

func WithMessageUserID(userID int) MessageField {
	return func(msg *Message) error {
		msg.UserID = userID
		return nil
	}
}

func WithMessageChatID(chatID int) MessageField {
	return func(msg *Message) error {
		msg.ChatID = chatID
		return nil
	}
}

func WithMessageText(text string) MessageField {
	return func(msg *Message) error {
		msg.Text = text
		return nil
	}
}
