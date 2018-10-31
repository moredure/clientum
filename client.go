package main

import (
	"context"
	"github.com/mikefaraponov/chatum"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewGRPCDial(url ServerAddress, d grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(string(url), d)
}

func NewChatumCommunicateClient(client chatum.ChatumClient, user User) (chatum.Chatum_CommunicateClient, error) {
	return client.Communicate(metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(UsernameField, string(user)),
	))
}
