package bloom

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/reoden/go-NFT/pkg/logger"
	"go.uber.org/fx"
)

var (
	// Module provided to fxlog
	Module = fx.Module(
		"bloomfx",
		bloomProviders,
		bloomInvokes,
	) //nolint:gochecknoglobals

	bloomProviders = fx.Options(fx.Provide(
		// Provide a function that creates a BloomFilter factory
		NewBloomFilterFactory,
	)) //nolint:gochecknoglobals

	bloomInvokes = fx.Options(
		fx.Invoke(registerBloomHooks),
	) //nolint:gochecknoglobals
)

// BloomFilterFactory holds the redis client and provides methods to create bloom filters
type BloomFilterFactory struct {
	client redis.UniversalClient
}

// NewBloomFilterFactory creates a new BloomFilterFactory with the provided redis client
func NewBloomFilterFactory(client redis.UniversalClient) *BloomFilterFactory {
	return &BloomFilterFactory{
		client: client,
	}
}

// NewBloomFilter creates a new BloomFilter with the specified parameters
func (f *BloomFilterFactory) New(m uint, k uint, key string) *BloomFilter {
	// Create BloomFilter with pre-configured redis client
	return NewBloomFilter(m, k, key, f.client)
}

// NewWithEstimates creates a new BloomFilter for about n items with fp false positive rate
func (f *BloomFilterFactory) NewWithEstimates(n uint, fp float64, key string) *BloomFilter {
	m, k := EstimateParameters(n, fp)
	return f.New(m, k, key)
}

// BloomFilterFrom creates a new Bloom filter with len(_data_) * 64 bits and _k_ hashing functions
func (f *BloomFilterFactory) From(ctx context.Context, data []uint64, k uint, key string) *BloomFilter {
	m := uint(len(data) * 64)
	return f.FromWithM(ctx, data, m, k, key)
}

// BloomFilterFromWithM creates a new Bloom filter with _m_ length, _k_ hashing functions
func (f *BloomFilterFactory) FromWithM(ctx context.Context, data []uint64, m, k uint, key string) *BloomFilter {
	return BloomFilterFromWithM(ctx, data, m, k, key, f.client)
}

func registerBloomHooks(
	lc fx.Lifecycle,
	logger logger.Logger,
) {
	// No special hooks needed for bloom filter since it depends on redis
	// which already has its own lifecycle management
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("bloom filter module initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("bloom filter module stopped")
			return nil
		},
	})
}
