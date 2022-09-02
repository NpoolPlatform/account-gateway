package transfer

import (
	"context"

	"github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	transfer.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	transfer.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return transfer.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
