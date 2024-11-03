package tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/ipv02/chat-server/internal/client/db"
	dbMocks "github.com/ipv02/chat-server/internal/client/db/mocks"
	"github.com/ipv02/chat-server/internal/model"
	"github.com/ipv02/chat-server/internal/repository"
	repoMocks "github.com/ipv02/chat-server/internal/repository/mocks"
	"github.com/ipv02/chat-server/internal/service/chat"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.ChatCreate
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		usersID = []string{
			strconv.FormatInt(gofakeit.Int64(), 10),
		}
		chatName = gofakeit.Name()

		repoErr = fmt.Errorf("repo error")

		req = &model.ChatCreate{
			UsersID:  usersID,
			ChatName: chatName,
		}
	)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		chatRepositoryMock chatRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.ExpectCtxParam1(ctx).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx, req).Return(0, repoErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.ExpectCtxParam1(ctx).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatRepoMock := tt.chatRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)
			service := chat.NewMockService(chatRepoMock, txManagerMock)

			newID, err := service.CreateChat(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
