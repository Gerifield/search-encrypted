package index

import (
	"crypto/sha512"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlgOption(t *testing.T) {
	i := New(WithHashAlg(sha512.New))

	assert.IsType(t, sha512.New(), i.hashAlg())
}

func TestCasingOption(t *testing.T) {
	i := New(WithCaseInsensitive(true))

	assert.True(t, i.caseInsensitive)
}

func TestHMACVersions(t *testing.T) {
	i := New()

	assert.Equal(t, "2ae2e0eed2d0aa7dc7c18af891e5557a890b28dc23d790b25efedb69b091f520", i.HMAC("testKey", "testVal"))
	assert.Equal(t, "c22a9edc46210388a797ba52cb33b12106415039b2613bd8942264ea369974ee", i.HMAC("testKey", "testVal1"))
	assert.Equal(t, "a1a910d022aa0a65f64fca7212150c2ce860091a4479b4f420115fb81a3eb246", i.HMAC("testKey1", "testVal"))

	i = New(WithCaseInsensitive(true))
	assert.Equal(t, "e8f0317eb8458168afef734be5edcf87e8ef71f5f47be6361abcbdd6f465c9a5", i.HMAC("testKey", "testVal"))
	assert.Equal(t, "2ea9cf8551abc3b5040a3ae62f6e4f5219a159d2f997f1f78c7d140ffa02b6f4", i.HMAC("testKey", "testVal1"))
	assert.Equal(t, "61b820c580cea34c34059050a58900713459474dde7e0f8896da8ff98f4c329d", i.HMAC("testKey1", "testVal"))
}

func TestBuckets(t *testing.T) {
	i := New()

	assert.Equal(t, "659186e84c8d6cc98dbcf9f64a02b9177124e9d1a36b4fa8eb85d462f6ac757b", i.IntBucket("testKey", 0, 10))
	assert.Equal(t, "659186e84c8d6cc98dbcf9f64a02b9177124e9d1a36b4fa8eb85d462f6ac757b", i.IntBucket("testKey", 9, 10))
	assert.Equal(t, "d1011078104565a31f7ebdbb83ef3de42c68d86fa38ec13531b657ebd83f7952", i.IntBucket("testKey", 10, 10))
	assert.Equal(t, "d1011078104565a31f7ebdbb83ef3de42c68d86fa38ec13531b657ebd83f7952", i.IntBucket("testKey", 19, 10))

	assert.Equal(t, "659186e84c8d6cc98dbcf9f64a02b9177124e9d1a36b4fa8eb85d462f6ac757b", i.IntBucket("testKey", 0, 5))
	assert.Equal(t, "659186e84c8d6cc98dbcf9f64a02b9177124e9d1a36b4fa8eb85d462f6ac757b", i.IntBucket("testKey", 4, 5))
	assert.Equal(t, "5504da0b7dcf7a712902bcdce7dcbfe87f1f7395dab848ddf51f0bd23aba6183", i.IntBucket("testKey", 5, 5))
}
