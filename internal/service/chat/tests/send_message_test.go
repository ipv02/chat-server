package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ipv02/chat-server/internal/model"
	"github.com/ipv02/chat-server/internal/repository"
	repoMocks "github.com/ipv02/chat-server/internal/repository/mocks"
	"github.com/ipv02/chat-server/internal/service/chat"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx context.Context
		req *model.ChatSendMessage
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		from      = gofakeit.Name()
		text      = gofakeit.City()
		timestamp = gofakeit.Date()

		repoErr = fmt.Errorf("repo error")

		req = &model.ChatSendMessage{
			From:      from,
			Text:      text,
			Timestamp: timestamppb.New(timestamp),
		}
	)

	tests := []struct {
		name               string
		args               args
		want               error
		err                error
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, req).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repoErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, req).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatRepoMock := tt.chatRepositoryMock(mc)
			service := chat.NewMockService(chatRepoMock)

			err := service.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, nil)
		})
	}
}
