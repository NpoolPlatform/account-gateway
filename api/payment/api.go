package payment

import (
	"context"

	"github.com/NpoolPlatform/message/npool/account/gw/v1/payment"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	payment.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	payment.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return payment.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
