package user

import (
	"context"

	"github.com/NpoolPlatform/message/npool/account/gw/v1/user"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	user.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	user.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return user.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
