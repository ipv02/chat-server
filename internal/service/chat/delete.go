package chat

import "context"

func (s *serv) DeleteChat(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.chatRepository.DeleteChat(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	return err
}
