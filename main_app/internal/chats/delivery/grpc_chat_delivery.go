package delivery

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	grpcChat "github.com/go-park-mail-ru/2024_2_EaglesDesigner/protos/gen/go/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	grpcChat.UnimplementedChatServiceServer // Хитрая штука, о ней ниже
	chatRepository                          ChatRepository
}

type ChatRepository interface {
	GetUserChats(ctx context.Context, userId string) (chatIds []string, err error)

	GetUsersFromChat(ctx context.Context, chatId string) (userIds []string, err error)
}

func Register(gRPCServer *grpc.Server, chatRepository ChatRepository) {
	grpcChat.RegisterChatServiceServer(gRPCServer, &serverAPI{chatRepository: chatRepository})
}

func (s *serverAPI) GetUserChats(ctx context.Context, in *grpcChat.UserChatsRequest) (*grpcChat.UserChatsResponse, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Infof("Пришел запрос на получение чатов пользователя: %v", in.UserId)

	chatIds, err := s.chatRepository.GetUserChats(ctx, in.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get chats")
	}

	return &grpcChat.UserChatsResponse{ChatIds: chatIds}, nil
}

func (s *serverAPI) GetUsersFromChat(ctx context.Context, in *grpcChat.UsersFromChatRequest) (*grpcChat.UsersFromChatResponse, error) {
	// TODO
	userIds, err := s.chatRepository.GetUsersFromChat(ctx, in.ChatId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get users")
	}

	return &grpcChat.UsersFromChatResponse{UserIds: userIds}, nil
}
