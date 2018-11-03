package main

import (
	"github.com/mikefaraponov/chatum"
	"github.com/mikefaraponov/clientum/client"
	"github.com/mikefaraponov/clientum/common"
	"github.com/mikefaraponov/clientum/ui"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func main() {
	fx.New(
		common.NewServerAddress(),
		common.NewUser(),
		fx.Provide(chatum.NewChatumClient),
		fx.Provide(grpc.WithInsecure),
		fx.Provide(client.NewGRPCDial),
		fx.Provide(client.NewChatumCommunicateClient),
		fx.Provide(ui.NewUI),
		fx.Invoke(client.Bootstrap),
	).Run()
}
