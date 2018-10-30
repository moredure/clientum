package main

import (
	"context"
	"github.com/marcusolsson/tui-go"
	"github.com/mikefaraponov/chatum"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"syscall"
)

func NewGRPCDial(url ServerAddress, d grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(string(url), d)
}

func NewChatumCommunicateClient(client chatum.ChatumClient, user User) (chatum.Chatum_CommunicateClient, error) {
	return client.Communicate(metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("username", string(user)),
	))
}

func Register(lc fx.Lifecycle, conn *grpc.ClientConn, ui tui.UI) {
	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	})

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go ui.Run()
			return nil
		},
		OnStop: func(context.Context) error {
			ui.Quit()
			return conn.Close()
		},
	})
}