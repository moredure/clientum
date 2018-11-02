package common

import (
	"go.uber.org/fx"
	"os"
)

type (
	ServerAddress string
	User string
)

func NewServerAddress() fx.Option {
	url, ok := os.LookupEnv("SERVER_URL")
	if !ok {
		return fx.Error(ServerAddressErr)
	}
	return fx.Provide(func() ServerAddress {
		return ServerAddress(url)
	})
}

func NewUser() fx.Option {
	user, ok := os.LookupEnv("USER")
	if !ok {
		return fx.Error(UserMissingErr)
	}
	return fx.Provide(func() User {
		return User(user)
	})
}
