package client

import (
	"context"
	"fmt"
)

type EndpointParameters struct {
	Bucket       string
	Key          string
	Host         string
	UseSSL       bool
	UsePathStyle bool
}

type Endpoint struct {
	URL string
}

type EndpointResolver interface {
	ResolveEndpoint(ctx context.Context, params EndpointParameters) (*Endpoint, error)
}

type DefaultEndpointResolver struct{}

func (*DefaultEndpointResolver) ResolveEndpoint(ctx context.Context, params EndpointParameters) (*Endpoint, error) {
	if params.Host == "" {
		return nil, fmt.Errorf("host is required for endpoint resolution")
	}

	url := "http"
	if params.UseSSL {
		url += "s"
	}

	url += "://"

	if !params.UsePathStyle && params.Bucket != "" {
		url += params.Bucket
		url += "."
	}

	url += params.Host

	if params.UsePathStyle && params.Bucket != "" {
		url += "/"
		url += params.Bucket
	}

	if params.Key != "" {
		url += "/"
		url += params.Key
	}

	return &Endpoint{URL: url}, nil
}
