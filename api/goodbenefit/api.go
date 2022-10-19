package goodbenefit

import (
	"context"

	"github.com/NpoolPlatform/message/npool/account/gw/v1/goodbenefit"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	goodbenefit.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	goodbenefit.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return goodbenefit.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
