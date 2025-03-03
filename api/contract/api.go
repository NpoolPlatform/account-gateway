package contract

import (
	"context"

	"github.com/NpoolPlatform/message/npool/account/gw/v1/contract"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	contract.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	contract.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return contract.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
