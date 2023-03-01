package server

import (
	"context"

	"github.com/XHXHXHX/medical_marketing_proto/gen/go/proto/v1api"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type ReportServer struct {
	c *Common
	v1api.ApiReportServiceServer
}

func NewReportServer(c *Common) *ReportServer {
	return &ReportServer{
		c: c,
	}
}

func (s *ReportServer) RegisterGRPC(svr grpc.ServiceRegistrar) {
	v1api.RegisterApiReportServiceServer(svr, s)
}

func (s *ReportServer) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return v1api.RegisterApiReportServiceHandler(ctx, mux, conn)
}

func (s *ReportServer) ReportCreate(ctx context.Context, req *v1api.ReportCreateRequest) (*v1api.ReportCreateResponse, error) {

	return &v1api.ReportCreateResponse{}, nil
}