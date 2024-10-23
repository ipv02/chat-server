package chat

import "context"

func (s *serv) DeleteChat(ctx context.Context, id int64) error {
	err := s.chatRepository.DeleteChat(ctx, id)

	return err
}
