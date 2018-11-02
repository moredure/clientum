package client

import (
	"context"
	"github.com/mikefaraponov/chatum"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/mikefaraponov/clientum/common"
)

func NewGRPCDial(url common.ServerAddress, d grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(string(url), d)
}

func NewChatumCommunicateClient(client chatum.ChatumClient, user common.User) (chatum.Chatum_CommunicateClient, error) {
	return client.Communicate(metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(common.UsernameField, string(user)),
	))
}
