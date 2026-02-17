package domain

type Chat struct {
	ChatID int
}

type ChatField func(chat *Chat) error

func NewChat(fields ...ChatField) (*Chat, error) {
	chat := Chat{}
	for _, field := range fields {
		err := field(&chat)
		if err != nil {
			return nil, err
		}
	}
	return &chat, nil
}

func WithChatChatID(chatID int) ChatField {
	return func(chat *Chat) error {
		chat.ChatID = chatID
		return nil
	}
}
