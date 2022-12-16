package api

import (
	"context"

	account "github.com/NpoolPlatform/message/npool/account/gw/v1"

	"github.com/NpoolPlatform/account-gateway/api/goodbenefit"
	"github.com/NpoolPlatform/account-gateway/api/payment"
	"github.com/NpoolPlatform/account-gateway/api/platform"
	"github.com/NpoolPlatform/account-gateway/api/transfer"
	"github.com/NpoolPlatform/account-gateway/api/user"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	account.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	account.RegisterGatewayServer(server, &Server{})
	user.Register(server)
	transfer.Register(server)
	platform.Register(server)
	payment.Register(server)
	goodbenefit.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := account.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	if err := user.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := transfer.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := goodbenefit.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := platform.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := payment.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
