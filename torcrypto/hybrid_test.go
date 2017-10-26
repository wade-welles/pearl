package torcrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHybridDecrypt(t *testing.T) {
	k, err := LoadRSAPrivateKeyFromPEMFile("testdata/hybrid_private_key")
	require.NoError(t, err)

	var plain = []byte{
		0x71, 0x45, 0x3e, 0x5d, 0x74, 0x39, 0x5f, 0x21, 0xae, 0x7c, 0x55, 0x29,
		0xbb, 0xa0, 0x22, 0xb6, 0xd2, 0x6c, 0x6c, 0x31, 0x30, 0xd9, 0xad, 0x1b,
		0xaf, 0x25, 0x8b, 0xf1, 0xdd, 0xd4, 0x66, 0x53, 0x9f, 0x3e, 0x35, 0xe9,
		0xdd, 0x12, 0x5a, 0x6a, 0x66, 0xcf, 0x83, 0xbf, 0x52, 0x10, 0x28, 0x19,
		0x24, 0x0d, 0x12, 0x0a, 0x7c, 0x1d, 0xb3, 0x4f, 0x2d, 0xfe, 0x32, 0xc7,
		0xce, 0x99, 0x65, 0x05, 0x6d, 0xe7, 0x18, 0x48, 0xf7, 0x8a, 0x63, 0x24,
		0x9f, 0x8d, 0x04, 0x81, 0x18, 0x5c, 0x09, 0x25, 0xfc, 0xf5, 0x1f, 0xe4,
		0xbc, 0x86, 0x5a, 0x60, 0x1b, 0x91, 0x1d, 0x96, 0xb1, 0x72, 0xef, 0x1f,
		0xf3, 0xdd, 0xff, 0x37, 0x4a, 0xc2, 0xc6, 0x85, 0xbb, 0xec, 0xc4, 0x4a,
		0xc1, 0xaf, 0x56, 0x99, 0xc4, 0xd6, 0x96, 0x3a, 0x18, 0x90, 0xaa, 0x69,
		0xa8, 0xf9, 0x03, 0xd4, 0xa7, 0x11, 0x82, 0xea,
	}

	var cipher = []byte{
		0x95, 0x4a, 0x69, 0xd7, 0x47, 0x4b, 0x9a, 0x4c, 0x57, 0xb3, 0x21, 0xf9,
		0xd3, 0x45, 0xa2, 0xc8, 0x2c, 0x53, 0x5e, 0x42, 0xfe, 0x3b, 0x20, 0x2f,
		0xbe, 0x52, 0xf7, 0x03, 0x34, 0xc6, 0x53, 0xfa, 0xe0, 0xa4, 0xea, 0x3c,
		0xc4, 0x62, 0x66, 0x95, 0xf2, 0x1b, 0x12, 0x44, 0xf5, 0xba, 0x6d, 0x19,
		0xf6, 0xc4, 0x3c, 0xa0, 0x1c, 0xae, 0x77, 0x70, 0x71, 0x86, 0x10, 0x98,
		0x2f, 0x66, 0xf6, 0x9f, 0x56, 0x66, 0x96, 0x94, 0x05, 0xc1, 0x83, 0x5b,
		0x99, 0x33, 0xd2, 0x9b, 0x11, 0x50, 0xea, 0xb5, 0xd1, 0x18, 0xd3, 0x0a,
		0x4c, 0x88, 0xa9, 0x63, 0x53, 0x17, 0x7d, 0x49, 0xfb, 0xb2, 0x30, 0x52,
		0x5d, 0xf2, 0x88, 0x33, 0x7e, 0x4a, 0x91, 0x84, 0x56, 0x09, 0x07, 0x8b,
		0xc4, 0x29, 0x2f, 0xc4, 0xb0, 0xbf, 0x33, 0x6e, 0xd4, 0x31, 0xbb, 0x26,
		0x84, 0x05, 0x11, 0x47, 0x66, 0xda, 0xe6, 0xa9, 0xcf, 0x82, 0x9c, 0xbb,
		0x51, 0x1b, 0x8b, 0xc8, 0x41, 0x69, 0x28, 0x64, 0x95, 0xb3, 0x9c, 0xe6,
		0x14, 0x04, 0xbd, 0x9c, 0x69, 0x59, 0x1b, 0x65, 0x5e, 0x33, 0xd2, 0x3c,
		0x27, 0x52, 0x57, 0xf5, 0x7b, 0x98, 0x41, 0x1a, 0xf9, 0xcd, 0x95, 0x40,
		0xef, 0x60, 0x12, 0x80, 0x39, 0x65, 0xb1, 0xee, 0x64, 0x7f, 0xb8, 0xe8,
		0x90, 0x99, 0xd3, 0x7b, 0x76, 0xd8,
	}

	got := HybridDecrypt(k, cipher)
	assert.Equal(t, plain, got)
}
