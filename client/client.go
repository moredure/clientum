package client

import (
	"context"
	"github.com/mikefaraponov/chatum"
	"github.com/mikefaraponov/clientum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewGRPCDial(env *common.Environment, d grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(env.ServerAddress, d)
}

func NewContext(env *common.Environment) context.Context {
	return metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(common.UsernameField, env.User),
	)
}

func NewChatumCommunicateClient(client chatum.ChatumClient, ctx context.Context) (chatum.Chatum_CommunicateClient, error) {
	return client.Communicate(ctx)
}
