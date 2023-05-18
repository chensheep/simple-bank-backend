package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedHostHeader       = "x-forwarded-host"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) ExtractMetadata(ctx context.Context) *Metadata {
	metdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			metdt.UserAgent = userAgents[0]
		}
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			metdt.UserAgent = userAgents[0]
		}
		if clientIps := md.Get(xForwardedHostHeader); len(clientIps) > 0 {
			metdt.ClientIp = clientIps[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		metdt.ClientIp = p.Addr.String()
	}

	return metdt
}
