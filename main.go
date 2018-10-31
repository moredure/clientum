package main

import (
	"github.com/mikefaraponov/chatum"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func main() {
	fx.New(
		NewServerAddress(),
		NewUser(),
		fx.Provide(chatum.NewChatumClient),
		fx.Provide(grpc.WithInsecure),
		fx.Provide(NewGRPCDial),
		fx.Provide(NewChatumCommunicateClient),
		fx.Provide(NewUI),
		fx.Invoke(Bootstrap),
	).Run()
}
