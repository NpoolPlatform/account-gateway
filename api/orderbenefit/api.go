package orderbenefit

import (
	"context"

	"github.com/NpoolPlatform/message/npool/account/gw/v1/orderbenefit"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	orderbenefit.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	orderbenefit.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return orderbenefit.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
