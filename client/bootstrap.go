package client

import (
	"context"
	"github.com/marcusolsson/tui-go"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"syscall"
	"github.com/mikefaraponov/clientum/common"
)

func Bootstrap(lc fx.Lifecycle, conn *grpc.ClientConn, ui tui.UI) {
	ui.SetKeybinding(common.Esc, func() {
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