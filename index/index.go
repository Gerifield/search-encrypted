package index

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
)

type Index struct {
	hashAlg         func() hash.Hash
	caseInsensitive bool
}

// HashAlg .
type HashAlg func() hash.Hash

// Option .
type Option func(*Index)

// WithHashAlg .
func WithHashAlg(alg HashAlg) Option {
	return func(i *Index) {
		i.hashAlg = alg
	}
}

// WithCaseInsensitive .
func WithCaseInsensitive(enable bool) Option {
	return func(i *Index) {
		i.caseInsensitive = enable
	}
}

// New .
func New(opts ...Option) *Index {
	i := &Index{}
	for _, o := range opts {
		o(i)
	}

	if i.hashAlg == nil {
		i.hashAlg = sha256.New
	}

	return i
}

// HMAC index generation
func (i *Index) HMAC(key string, value string) string {
	sig := hmac.New(i.hashAlg, []byte(key))
	if i.caseInsensitive {
		value = strings.ToUpper(value)
	}
	sig.Write([]byte(value))

	return hex.EncodeToString(sig.Sum(nil))
}

// IntBucket index generation
// This will generate an index using a simple method to put the `num` param in a given bucket
func (i *Index) IntBucket(key string, num int, allBuckets int) string {
	// floor down to the allBuckets-th part using the conversion in the type system
	// of course different values could have different categorization
	bucketIndex := (num / allBuckets) * allBuckets

	return i.HMAC(key, fmt.Sprintf("%d", bucketIndex))
}
