package clients

import "context"

type Provider interface {
	Name() string
	Capabilities() Capability

	Generate(ctx context.Context, req *Request) (*Response, error)
	Stream(ctx context.Context, req *Request) (Stream, error)
}
