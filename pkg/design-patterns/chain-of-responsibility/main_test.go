package chain_of_responsibility

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockRequest struct {
	context.Context
	MyData string
}

func newMockMiddleware(v string) Middleware[*mockRequest] {
	return MiddlewareFunc[*mockRequest](func(request *mockRequest, next Handler[*mockRequest]) error {
		fmt.Println(">>>", v)
		fmt.Println("  request.MyData:", request.MyData)
		request.MyData = v
		err := next.Handle(request)
		fmt.Println("  err:", err)
		fmt.Println("<<<", v)

		return errors.New("error-" + v)
	})
}

func TestNewChain(t *testing.T) {
	handlerInput := "handler-input"
	handlerError := errors.New("handler-error")

	handler := HandlerFunc[*mockRequest](func(request *mockRequest) error {
		require.Equal(t, handlerInput, request.MyData)
		return handlerError
	})

	chain := NewChain(handler)

	actualError := chain.Handle(&mockRequest{
		Context: t.Context(),
		MyData:  handlerInput,
	})
	require.Same(t, handlerError, actualError)
}

func ExampleNewChain() {
	m1 := newMockMiddleware("m1")
	m2 := newMockMiddleware("m2")
	m3 := newMockMiddleware("m3")

	h := HandlerFunc[*mockRequest](func(request *mockRequest) error {
		fmt.Println(">>> handler")
		fmt.Println("  request.MyData:", request.MyData)
		fmt.Println("<<< handler")
		return errors.New("error-handler")
	})

	chain := NewChain(h, m1, m2, m3)
	err := chain.Handle(&mockRequest{
		Context: context.Background(),
		MyData:  "m0",
	})
	fmt.Println("err:", err)

	// Output:
	// >>> m1
	//   request.MyData: m0
	// >>> m2
	//   request.MyData: m1
	// >>> m3
	//   request.MyData: m2
	// >>> handler
	//   request.MyData: m3
	// <<< handler
	//   err: error-handler
	// <<< m3
	//   err: error-m3
	// <<< m2
	//   err: error-m2
	// <<< m1
	// err: error-m1
}
