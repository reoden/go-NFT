package contracts

import (
	"context"
	"testing"

	postgres "github.com/reoden/go-NFT/pkg/postgrespgx"
)

type PostgresPgxContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*PostgresContainerOptions,
	) (*postgres.PostgresPgxOptions, error)
	Cleanup(ctx context.Context) error
}
