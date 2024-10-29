package chat_v1

import (
	"github.com/pkg/errors"
)

// Validate валидация CreateChatRequest
func (req *CreateChatRequest) Validate() error {
	if len(req.UsersId) == 0 {
		return errors.New("validation error: at least one user ID is required")
	}

	for _, userID := range req.UsersId {
		if userID == "" {
			return errors.New("validation error: user ID cannot be empty")
		}
	}

	if req.ChatName == "" {
		return errors.New("validation error: chat name is required")
	}

	return nil
}

// Validate валидация DeleteChatRequest
func (req *DeleteChatRequest) Validate() error {
	if req.Id <= 0 {
		return errors.New("validation error: id must be greater than 0")
	}

	return nil
}

// Validate валидация SendMessageRequest
func (req *SendMessageRequest) Validate() error {
	if req.From == "" {
		return errors.New("validation error: sender ID is required")
	}

	if req.Text == "" {
		return errors.New("validation error: message text is required")
	}

	if req.Timestamp == nil {
		return errors.New("validation error: timestamp is required")
	}

	return nil
}
