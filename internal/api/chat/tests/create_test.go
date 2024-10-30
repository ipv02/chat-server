package tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/ipv02/chat-server/internal/api/chat"
	"github.com/ipv02/chat-server/internal/model"
	"github.com/ipv02/chat-server/internal/service"
	serviceMocks "github.com/ipv02/chat-server/internal/service/mocks"
	"github.com/ipv02/chat-server/pkg/chat_v1"
)

func TestCreate(t *testing.T) {
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *chat_v1.CreateChatRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		usersId = []string{
			strconv.FormatInt(gofakeit.Int64(), 10),
		}
		chatName = gofakeit.Name()

		serviceErr = fmt.Errorf("service error")

		req = &chat_v1.CreateChatRequest{
			UsersId:  usersId,
			ChatName: chatName,
		}

		serviceReq = &model.ChatCreate{
			UsersID:  usersId,
			ChatName: chatName,
		}

		res = &chat_v1.CreateChatResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *chat_v1.CreateChatResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateChatMock.Expect(ctx, serviceReq).Return(id, nil)
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
			err:  serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateChatMock.Expect(ctx, serviceReq).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatServiceMock := tt.chatServiceMock(mc)
			api := chat.NewImplementation(chatServiceMock)

			res, err := api.CreateChat(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
