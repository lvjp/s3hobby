package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	chain_of_responsibility "github.com/lvjp/s3hobby/pkg/design-patterns/chain-of-responsibility"
	"github.com/lvjp/s3hobby/pkg/s3/api"
	"github.com/lvjp/s3hobby/pkg/s3/signing"

	"github.com/valyala/fasthttp"
)

type Handler[Input, Output any] = chain_of_responsibility.Handler[*callContext[Input, Output]]

type RequiredBucketInterface interface {
	GetBucket() string
}

type RequiredBucketKeyInterface interface {
	RequiredBucketInterface
	GetKey() string
}

type callContext[Input, Output any] struct {
	context.Context

	Options *Options

	CallInput  Input
	CallOutput Output

	ServerRequest  fasthttp.Request
	ServerResponse *fasthttp.Response
}

type Metadata struct {
	Request  *fasthttp.Request
	Response *fasthttp.Response
}

func PerformCall[
	Input api.HTTPRequestMarshaler,
	OutputPtr interface {
		api.HTTPResponseUnmarshaler
		*OutputBase
	},
	OutputBase any,
](ctx context.Context, options *Options, input Input, optFns ...func(*Options)) (OutputPtr, *Metadata, error) {
	callContext := &callContext[Input, OutputPtr]{
		Context:   ctx,
		Options:   options.With(optFns...),
		CallInput: input,
	}

	chain := chain_of_responsibility.NewChain(
		&httpRequesterHandler[Input, OutputPtr]{},
		&configValidationMiddleware[Input, OutputPtr]{},
		&requiredInputMiddleware[Input, OutputPtr]{},
		&userAgentMiddleware[Input, OutputPtr]{},
		&resolveEndpointMiddleware[Input, OutputPtr]{},
		&transportMiddleware[Input, OutputBase, OutputPtr]{},
		&signerMiddleware[Input, OutputPtr]{},
	)

	err := chain.Handle(callContext)

	metadata := &Metadata{
		Request:  &callContext.ServerRequest,
		Response: callContext.ServerResponse,
	}

	if err != nil {
		return nil, metadata, err
	}

	return callContext.CallOutput, metadata, nil
}

type httpRequesterHandler[Input, Output any] struct{}

func (*httpRequesterHandler[Input, Output]) Handle(ctx *callContext[Input, Output]) error {
	var serverResponse fasthttp.Response

	if err := ctx.Options.HTTPClient.Do(&ctx.ServerRequest, &serverResponse); err != nil {
		return fmt.Errorf("HTTP request error: %v", err)
	}

	ctx.ServerResponse = &serverResponse

	return nil
}

type configValidationMiddleware[Input any, Output any] struct{}

func (*configValidationMiddleware[Input, Output]) Middleware(ctx *callContext[Input, Output], next Handler[Input, Output]) error {
	if err := ctx.Options.Validate(); err != nil {
		return err
	}

	return next.Handle(ctx)
}

type userAgentMiddleware[Input any, Output any] struct{}

func (*userAgentMiddleware[Input, Output]) Middleware(ctx *callContext[Input, Output], next Handler[Input, Output]) error {
	if ctx.Options.UserAgent != nil && *ctx.Options.UserAgent != "" {
		ctx.ServerRequest.Header.SetUserAgent(*ctx.Options.UserAgent)
	}

	return next.Handle(ctx)
}

type resolveEndpointMiddleware[Input any, Output any] struct{}

func (*resolveEndpointMiddleware[Input, Output]) Middleware(ctx *callContext[Input, Output], next Handler[Input, Output]) error {
	params := EndpointParameters{
		Host:         ctx.Options.EndpointHost,
		UseSSL:       ctx.Options.UseSSL,
		UsePathStyle: ctx.Options.UsePathStyle,
	}

	switch v := any(ctx.CallInput).(type) {
	case RequiredBucketKeyInterface:
		params.Bucket = v.GetBucket()
		params.Key = v.GetKey()
	case RequiredBucketInterface:
		params.Bucket = v.GetBucket()
	}

	endpoint, err := ctx.Options.EndpointResolver.ResolveEndpoint(ctx, params)
	if err != nil {
		return fmt.Errorf("cannot resolve endpoint: %v", err)
	}

	ctx.ServerRequest.SetRequestURI(endpoint.URL)

	return next.Handle(ctx)
}

type signerMiddleware[Input any, Output any] struct{}

func (*signerMiddleware[Input, Output]) Middleware(ctx *callContext[Input, Output], next Handler[Input, Output]) error {
	_, _, _, err := ctx.Options.Signer.Sign(signing.SigningArgs{
		Request:     &ctx.ServerRequest,
		Credentials: *ctx.Options.Credentials,
		Region:      ctx.Options.SiginingRegion,
		SigningTime: signing.SigningTimeOf(time.Now()),
	})
	if err != nil {
		return fmt.Errorf("cannot sign the request: %v", err)
	}

	return next.Handle(ctx)
}

type transportMiddleware[
	Input api.HTTPRequestMarshaler,
	OutputBase any,
	OutputPtr interface {
		api.HTTPResponseUnmarshaler
		*OutputBase
	},
] struct{}

func (*transportMiddleware[Input, OutputBase, OutputPtr]) Middleware(ctx *callContext[Input, OutputPtr], next Handler[Input, OutputPtr]) error {
	if err := ctx.CallInput.MarshalHTTP(&ctx.ServerRequest); err != nil {
		return fmt.Errorf("HTTP marshaling error: %v", err)
	}

	if err := next.Handle(ctx); err != nil {
		return err
	}

	var callOutputPtr OutputPtr = new(OutputBase)

	if err := callOutputPtr.UnmarshalHTTP(ctx.ServerResponse); err != nil {
		return err
	}

	ctx.CallOutput = callOutputPtr
	return nil
}

type requiredInputMiddleware[Input any, Output any] struct{}

func (*requiredInputMiddleware[Input, Output]) Middleware(ctx *callContext[Input, Output], next Handler[Input, Output]) error {
	switch v := any(ctx.CallInput).(type) {
	case RequiredBucketKeyInterface:
		if v.GetBucket() == "" {
			return errors.New("bucket is mandatory")
		}

		if v.GetKey() == "" {
			return errors.New("object key is mandatory")
		}
	case RequiredBucketInterface:
		if v.GetBucket() == "" {
			return errors.New("bucket is mandatory")
		}
	}

	return next.Handle(ctx)
}
