package platform

import (
	"context"

	"github.com/NpoolPlatform/message/npool/account/gw/v1/platform"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	platform.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	platform.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return platform.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
