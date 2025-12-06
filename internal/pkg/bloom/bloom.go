/*
Package bloom provides data structures and methods for creating Bloom filters.

A Bloom filter is a representation of a set of _n_ items, where the main
requirement is to make membership queries; _i.e._, whether an item is a
member of a set.

A Bloom filter has two parameters: _m_, a maximum size (typically a reasonably large
multiple of the cardinality of the set to represent) and _k_, the number of hashing
functions on elements of the set. (The actual hashing functions are important, too,
but this is not a parameter for this implementation). A Bloom filter is backed by
a BitSet; a key is represented in the filter by setting the bits at each value of the
hashing functions (modulo _m_). Set membership is done by _testing_ whether the
bits at each value of the hashing functions (again, modulo _m_) are set. If so,
the item is in the set. If the item is actually in the set, a Bloom filter will
never fail (the true positive rate is 1.0); but it is susceptible to false
positives. The art is to choose _k_ and _m_ correctly.

In this implementation, the hashing functions used is murmurhash,
a non-cryptographic hashing function.

This implementation accepts keys for setting as testing as []byte. Thus, to
add a string item, "Love":

	uint n = 1000
	filter := bloom.NewBloomFilter(20*n, 5) // load of 20, 5 keys
	filter.Add([]byte("Love"))

Similarly, to test if "Love" is in bloom:

	if filter.Exists([]byte("Love"))

For numeric data, I recommend that you look into the binary/encoding library. But,
for example, to add a uint32 to the filter:

	i := uint32(100)
	n1 := make([]byte,4)
	binary.BigEndian.PutUint32(n1,i)
	f.Add(n1)

Finally, there is a method to estimate the false positive rate of a
Bloom filter with _m_ bits and _k_ hashing functions for a set of size _n_:

	if bloom.EstimateFalsePositiveRate(20*n, 5, n) > 0.001 ...

You can use it to validate the computed m, k parameters:

	m, k := bloom.EstimateParameters(n, fp)
	ActualfpRate := bloom.EstimateFalsePositiveRate(m, k, n)

or

	f := bloom.NewWithEstimates(n, fp)
	ActualfpRate := bloom.EstimateFalsePositiveRate(f.m, f.k, n)

You would expect ActualfpRate to be close to the desired fp in these cases.

The EstimateFalsePositiveRate function creates a temporary Bloom filter. It is
also relatively expensive and only meant for validation.
*/

// https://pkg.go.dev/github.com/bits-and-blooms/bloom/v3#section-readme
package bloom

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"

	//"github.com/bits-and-blooms/bitset"
	"github.com/redis/go-redis/v9"
)

// A BloomFilter is a representation of a set of _n_ items, where the main
// requirement is to make membership queries; _i.e._, whether an item is a
// member of a set.
type BloomFilter struct {
	m   uint
	k   uint
	b   redis.UniversalClient
	key string
}

const (
	maxRetries      = 5
	minRetryBackoff = 300 * time.Millisecond
	maxRetryBackoff = 500 * time.Millisecond
	dialTimeout     = 5 * time.Second
	readTimeout     = 5 * time.Second
	writeTimeout    = 3 * time.Second
	minIdleConns    = 20
	poolTimeout     = 6 * time.Second
)

// NewBloomFilter creates a new Bloom filter with _m_ bits and _k_ hashing functions
// We force _m_ and _k_ to be at least one to avoid panics.
func NewBloomFilter(m uint, k uint, key string, b redis.UniversalClient) *BloomFilter {
	return &BloomFilter{max(1, m), max(1, k), b, key}
}

// BloomFilterFrom creates a new Bloom filter with len(_data_) * 64 bits and _k_ hashing
// functions. The data slice is not going to be reset.
func BloomFilterFrom(ctx context.Context, data []uint64, k uint, key string, b redis.UniversalClient) *BloomFilter {
	m := uint(len(data) * 64)
	return BloomFilterFromWithM(ctx, data, m, k, key, b)
}

// BloomFilterFromWithM creates a new Bloom filter with _m_ length, _k_ hashing functions.
// The data slice is not going to be reset.
func BloomFilterFromWithM(ctx context.Context, data []uint64, m, k uint, key string, b redis.UniversalClient) *BloomFilter {
	// Create a new BloomFilter with Redis backend
	bf := &BloomFilter{
		m:   m,
		k:   k,
		b:   b,
		key: key, // or you could generate a unique key
	}

	// Initialize the Redis bitmap with the provided data
	// You'll need to convert the uint64 slice to bits in Redis
	for i, word := range data {
		for j := uint(0); j < 64; j++ {
			if word&(1<<j) != 0 {
				pos := uint(i*64) + j
				if pos < m { // Only set bits within the filter size
					bf.b.SetBit(ctx, bf.key, int64(pos), 1)
				}
			}
		}
	}

	return bf
}

// baseHashes returns the four hash values of data that are used to create k
// hashes
func baseHashes(data []byte) [4]uint64 {
	var d digest128 // murmur hashing
	hash1, hash2, hash3, hash4 := d.sum256(data)
	return [4]uint64{
		hash1, hash2, hash3, hash4,
	}
}

// location returns the ith hashed location using the four base hash values
func location(h [4]uint64, i uint) uint64 {
	ii := uint64(i)
	return h[ii%2] + ii*h[2+(((ii+(ii%2))%4)/2)]
}

// location returns the ith hashed location using the four base hash values
func (f *BloomFilter) location(h [4]uint64, i uint) uint {
	return uint(location(h, i) % uint64(f.m))
}

// EstimateParameters estimates requirements for m and k.
// Based on https://bitbucket.org/ww/bloom/src/829aa19d01d9/bloom.go
// used with permission.
func EstimateParameters(n uint, p float64) (m uint, k uint) {
	m = uint(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	k = uint(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return
}

// NewWithEstimates creates a new Bloom filter for about n items with fp
// false positive rate
func NewWithEstimates(n uint, fp float64, key string, b redis.UniversalClient) *BloomFilter {
	m, k := EstimateParameters(n, fp)
	return NewBloomFilter(m, k, key, b)
}

// Cap returns the capacity, _m_, of a Bloom filter
func (f *BloomFilter) Cap() uint {
	return f.m
}

// K returns the number of hash functions used in the BloomFilter
func (f *BloomFilter) K() uint {
	return f.k
}

// BitSet returns the underlying bitset for this filter.
func (f *BloomFilter) BitSet() redis.UniversalClient {
	return f.b
}

// Add data to the Bloom Filter. Returns the filter (allows chaining)
func (f *BloomFilter) Add(ctx context.Context, data []byte) *BloomFilter {
	h := baseHashes(data)
	for i := uint(0); i < f.k; i++ {
		f.b.SetBit(ctx, f.key, int64(f.location(h, i)), 1).Result()
	}
	return f
}

// Merge the data from two Bloom Filters.
func (f *BloomFilter) Merge(g *BloomFilter) error {
	// Make sure the m's and k's are the same, otherwise merging has no real use.
	if f.m != g.m {
		return fmt.Errorf("m's don't match: %d != %d", f.m, g.m)
	}

	if f.k != g.k {
		return fmt.Errorf("k's don't match: %d != %d", f.k, g.k)
	}

	err := f.b.BitOpOr(context.Background(), f.key, f.key, g.key).Err()
	if err != nil {
		return fmt.Errorf("failed to merge Bloom filters: %w", err)
	}
	return nil
}

// Copy creates a copy of a Bloom filter.
func (f *BloomFilter) Copy(ctx context.Context) *BloomFilter {
	fc := NewBloomFilter(f.m, f.k, fmt.Sprintf("%s_%d", f.key, time.Now().UnixNano()), f.b)
	err := fc.b.BitOpOr(ctx, fc.key, f.key, f.key).Err()
	if err != nil {
		// Handle error appropriately
		// In a real implementation, you might want to return an error
		return nil
	}

	return fc
}

// AddString to the Bloom Filter. Returns the filter (allows chaining)
func (f *BloomFilter) AddString(ctx context.Context, data string) *BloomFilter {
	return f.Add(ctx, []byte(data))
}

// Exists returns true if the data is in the BloomFilter, false otherwise.
// If true, the result might be a false positive. If false, the data
// is definitely not in the set.
func (f *BloomFilter) Exists(ctx context.Context, data []byte) bool {
	h := baseHashes(data)
	for i := uint(0); i < f.k; i++ {
		pos := f.location(h, i)
		// Use the key field here
		bit, _ := f.b.GetBit(ctx, f.key, int64(pos)).Result()
		if bit == 0 {
			return false
		}
	}
	return true
}

// ExistsString returns true if the string is in the BloomFilter, false otherwise.
// If true, the result might be a false positive. If false, the data
// is definitely not in the set.
func (f *BloomFilter) ExistsString(ctx context.Context, data string) bool {
	return f.Exists(ctx, []byte(data))
}

// TestLocations returns true if all locations are set in the BloomFilter, false
// otherwise.
func (f *BloomFilter) TestLocations(locs []uint64) bool {
	for i := 0; i < len(locs); i++ {
		if bit, _ := f.b.GetBit(context.Background(), f.key, int64(uint(locs[i]%uint64(f.m)))).Result(); bit == 0 {
			return false
		}
	}
	return true
}

// TestAndAdd is equivalent to calling Exists(data) then Add(data).
// The filter is written to unconditionnally: even if the element is present,
// the corresponding bits are still set. See also TestOrAdd.
// Returns the result of Exists.
func (f *BloomFilter) TestAndAdd(data []byte) bool {
	present := true
	h := baseHashes(data)
	for i := uint(0); i < f.k; i++ {
		l := f.location(h, i)
		if bit, _ := f.b.GetBit(context.Background(), f.key, int64(l)).Result(); bit == 0 {
			present = false
		}
		f.b.SetBit(context.Background(), f.key, int64(l), 1)
	}
	return present
}

// TestAndAddString is the equivalent to calling Exists(string) then Add(string).
// The filter is written to unconditionnally: even if the string is present,
// the corresponding bits are still set. See also TestOrAdd.
// Returns the result of Exists.
func (f *BloomFilter) TestAndAddString(data string) bool {
	return f.TestAndAdd([]byte(data))
}

// TestOrAdd is equivalent to calling Exists(data) then if not present Add(data).
// If the element is already in the filter, then the filter is unchanged.
// Returns the result of Exists.
func (f *BloomFilter) TestOrAdd(data []byte) bool {
	present := true
	h := baseHashes(data)
	for i := uint(0); i < f.k; i++ {
		l := f.location(h, i)
		if bit, _ := f.b.GetBit(context.Background(), f.key, int64(l)).Result(); bit == 0 {
			present = false
			f.b.SetBit(context.Background(), f.key, int64(l), 1)
		}
	}
	return present
}

// TestOrAddString is the equivalent to calling Exists(string) then if not present Add(string).
// If the string is already in the filter, then the filter is unchanged.
// Returns the result of Exists.
func (f *BloomFilter) TestOrAddString(data string) bool {
	return f.TestOrAdd([]byte(data))
}

// ClearAll clears all the data in a Bloom filter, removing all keys
func (f *BloomFilter) ClearAll(ctx context.Context) *BloomFilter {
	f.b.Del(ctx, f.key)
	return f
}

// EstimateFalsePositiveRate returns, for a BloomFilter of m bits
// and k hash functions, an estimation of the false positive rate when
//
//	storing n entries. This is an empirical, relatively slow
//
// test using integers as keys.
// This function is useful to validate the implementation.
func EstimateFalsePositiveRate(ctx context.Context, m, k, n uint, key string, b redis.UniversalClient) (fpRate float64) {
	rounds := uint32(100000)
	// We construct a new filter.
	f := NewBloomFilter(m, k, key, b)
	n1 := make([]byte, 4)
	// We populate the filter with n values.
	for i := uint32(0); i < uint32(n); i++ {
		binary.BigEndian.PutUint32(n1, i)
		f.Add(ctx, n1)
	}
	fp := 0
	// test for number of rounds
	for i := uint32(0); i < rounds; i++ {
		binary.BigEndian.PutUint32(n1, i+uint32(n)+1)
		if f.Exists(ctx, n1) {
			fp++
		}
	}
	fpRate = float64(fp) / (float64(rounds))
	return
}

// Approximating the number of items
// https://en.wikipedia.org/wiki/Bloom_filter#Approximating_the_number_of_items_in_a_Bloom_filter
func (f *BloomFilter) ApproximatedSize() uint32 {
	ctx := context.Background()
	count, err := f.b.BitCount(ctx, f.key, nil).Result()
	if err != nil {
		// Handle error appropriately
		count = 0
	}
	x := float64(count)
	m := float64(f.Cap())
	k := float64(f.K())
	size := -1 * m / k * math.Log(1-x/m) / math.Log(math.E)
	return uint32(math.Floor(size + 0.5)) // round
}

// bloomFilterJSON is an unexported type for marshaling/unmarshaling BloomFilter struct.
//type bloomFilterJSON struct {
//	M uint           `json:"m"`
//	K uint           `json:"k"`
//	B *bitset.BitSet `json:"b"`
//}

// MarshalJSON implements json.Marshaler interface.
//func (f BloomFilter) MarshalJSON() ([]byte, error) {
//	return json.Marshal(bloomFilterJSON{f.m, f.k, f.b})
//}

// UnmarshalJSON implements json.Unmarshaler interface.
//func (f *BloomFilter) UnmarshalJSON(data []byte) error {
//	var j bloomFilterJSON
//	err := json.Unmarshal(data, &j)
//	if err != nil {
//		return err
//	}
//	f.m = j.M
//	f.k = j.K
//	f.b = j.B
//	return nil
//}

// WriteTo writes a binary representation of the BloomFilter to an i/o stream.
// It returns the number of bytes written.
//
// Performance: if this function is used to write to a disk or network
// connection, it might be beneficial to wrap the stream in a bufio.Writer.
// E.g.,
//
//	      f, err := os.Create("myfile")
//		       w := bufio.NewWriter(f)
//func (f *BloomFilter) WriteTo(stream io.Writer) (int64, error) {
//	err := binary.Write(stream, binary.BigEndian, uint64(f.m))
//	if err != nil {
//		return 0, err
//	}
//	err = binary.Write(stream, binary.BigEndian, uint64(f.k))
//	if err != nil {
//		return 0, err
//	}
//	numBytes, err := f.b.WriteTo(stream)
//	return numBytes + int64(2*binary.Size(uint64(0))), err
//}

// ReadFrom reads a binary representation of the BloomFilter (such as might
// have been written by WriteTo()) from an i/o stream. It returns the number
// of bytes read.
//
// Performance: if this function is used to read from a disk or network
// connection, it might be beneficial to wrap the stream in a bufio.Reader.
// E.g.,
//
//	f, err := os.Open("myfile")
//	r := bufio.NewReader(f)
//func (f *BloomFilter) ReadFrom(stream io.Reader) (int64, error) {
//	var m, k uint64
//	err := binary.Read(stream, binary.BigEndian, &m)
//	if err != nil {
//		return 0, err
//	}
//	err = binary.Read(stream, binary.BigEndian, &k)
//	if err != nil {
//		return 0, err
//	}
//	b := &bitset.BitSet{}
//	numBytes, err := b.ReadFrom(stream)
//	if err != nil {
//		return 0, err
//	}
//	f.m = uint(m)
//	f.k = uint(k)
//	f.b = b
//	return numBytes + int64(2*binary.Size(uint64(0))), nil
//}

// GobEncode implements gob.GobEncoder interface.
//func (f *BloomFilter) GobEncode() ([]byte, error) {
//	var buf bytes.Buffer
//	_, err := f.WriteTo(&buf)
//	if err != nil {
//		return nil, err
//	}
//
//	return buf.Bytes(), nil
//}

// GobDecode implements gob.GobDecoder interface.
//func (f *BloomFilter) GobDecode(data []byte) error {
//	buf := bytes.NewBuffer(data)
//	_, err := f.ReadFrom(buf)
//
//	return err
//}

// MarshalBinary implements binary.BinaryMarshaler interface.
//func (f *BloomFilter) MarshalBinary() ([]byte, error) {
//	var buf bytes.Buffer
//	_, err := f.WriteTo(&buf)
//	if err != nil {
//		return nil, err
//	}
//
//	return buf.Bytes(), nil
//}

// UnmarshalBinary implements binary.BinaryUnmarshaler interface.
//func (f *BloomFilter) UnmarshalBinary(data []byte) error {
//	buf := bytes.NewBuffer(data)
//	_, err := f.ReadFrom(buf)
//
//	return err
//}

// Equal tests for the equality of two Bloom filters
func (f *BloomFilter) Equal(g *BloomFilter) bool {
	if g == nil || f == nil {
		return g == f
	}
	if f.m != g.m {
		return false
	}
	if f.k != g.k {
		return false
	}

	// For Redis-based implementation, we need to compare the bitmaps in Redis
	ctx := context.Background()

	// Get the bit counts for both filters
	countF, errF := f.b.BitCount(ctx, f.key, nil).Result()
	countG, errG := g.b.BitCount(ctx, g.key, nil).Result()

	// If bit counts differ, the filters are not equal
	if errF != nil || errG != nil || countF != countG {
		log.Println("bit count differ")
		return false
	}

	// Perform a bitwise comparison using BITOP XOR
	// If all bits are the same, XOR result should be all zeros
	tempKey := fmt.Sprintf("temp_xor_%d", time.Now().UnixNano())
	defer g.b.Del(ctx, tempKey) // Clean up temporary key

	err := f.b.BitOpXor(ctx, tempKey, f.key, g.key).Err()
	if err != nil {
		log.Println("error performing BITOP XOR")
		return false
	}

	// Check if the XOR result has any set bits
	xorCount, err := f.b.BitCount(ctx, tempKey, nil).Result()
	if err != nil {
		log.Println("error getting XOR bit count")
		return false
	}

	// If XOR result has zero bits set, the filters are equal
	return xorCount == 0
}

// Locations returns a list of hash locations representing a data item.
func Locations(data []byte, k uint) []uint64 {
	locs := make([]uint64, k)

	// calculate locations
	h := baseHashes(data)
	for i := uint(0); i < k; i++ {
		locs[i] = location(h, i)
	}

	return locs
}
