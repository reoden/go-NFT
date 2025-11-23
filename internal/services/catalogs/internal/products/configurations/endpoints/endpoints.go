package endpoints

import (
	"github.com/reoden/go-echo-template/pkg/core/web/route"
)

func RegisterEndpoints(endpoints []route.Endpoint) error {
	for _, endpoint := range endpoints {
		endpoint.MapEndpoint()
	}

	return nil
}
