package pearl

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLinkSpecTCPIPv4(t *testing.T) {
	s := NewLinkSpecTCP(net.IPv4(127, 0, 0, 1), 5002)
	assert.Equal(t, s.Type, LinkSpecTLSTCPIPv4)
	assert.Equal(t, []byte{127, 0, 0, 1, 0x13, 0x8a}, s.Spec)
}

func TestNewLinkSpecTCPIPv6(t *testing.T) {
	ip := net.IP([]byte{
		0, 1, 2, 3, 0, 1, 2, 3,
		0, 1, 2, 3, 0, 1, 2, 3,
	})
	s := NewLinkSpecTCP(ip, 0x1337)
	assert.Equal(t, s.Type, LinkSpecTLSTCPIPv6)
	expect := []byte{
		0, 1, 2, 3, 0, 1, 2, 3,
		0, 1, 2, 3, 0, 1, 2, 3,
		0x13, 0x37,
	}
	assert.Equal(t, expect, s.Spec)
}

func TestNewLinkSpecTCPUnexpected(t *testing.T) {
	assert.Panics(t, func() {
		NewLinkSpecTCP(net.IP(make([]byte, 7)), 0x1337)
	})
}

func TestLinkSpecAddressIPv4(t *testing.T) {
	s := LinkSpec{
		Type: LinkSpecTLSTCPIPv4,
		Spec: []byte{23, 123, 43, 1, 1, 1},
	}
	addr, err := s.Address()
	require.NoError(t, err)
	require.NotNil(t, addr)
	assert.Equal(t, "23.123.43.1:257", addr.String())
}

func TestLinkSpecAddressIPv6(t *testing.T) {
	s := LinkSpec{
		Type: LinkSpecTLSTCPIPv6,
		Spec: []byte{
			0, 1, 2, 3, 4, 5, 6, 7,
			0, 1, 2, 3, 4, 5, 6, 7,
			1, 1,
		},
	}
	addr, err := s.Address()
	require.NoError(t, err)
	require.NotNil(t, addr)
	assert.Equal(t, "[1:203:405:607:1:203:405:607]:257", addr.String())
}

func TestLinkSpecAddressNil(t *testing.T) {
	for _, lt := range []LinkSpecType{
		LinkSpecLegacyIdentity,
		LinkSpecEd25519Identity,
	} {
		s := LinkSpec{Type: lt}
		addr, err := s.Address()
		assert.NoError(t, err)
		assert.Nil(t, addr)
	}
}

func TestLinkSpecAddressBadLength(t *testing.T) {
	s := LinkSpec{
		Type: LinkSpecTLSTCPIPv4,
		Spec: []byte{0, 1, 2, 3, 4, 5, 6},
	}
	_, err := s.Address()
	assert.Error(t, err)
}

func AssertLinkSpecEqual(t *testing.T, expect, got LinkSpec) {
	assert.Equal(t, expect.Type, got.Type)
	assert.Equal(t, expect.Spec, got.Spec)
}

func TestExtend2UnmarshalBinary(t *testing.T) {
	payload := []byte{
		0x0e, 0x00, 0x00, 0x00, 0x00, 0x57, 0xd4, 0x8d, 0x22, 0x00, 0x77, 0x02,
		0x00, 0x06, 0x7f, 0x00, 0x00, 0x01, 0x13, 0x8a, 0x02, 0x14, 0x8f, 0xd0,
		0xc0, 0xef, 0x1c, 0x8a, 0xdc, 0x3c, 0x52, 0x9b, 0xf5, 0xe1, 0x9f, 0xc7,
		0x86, 0xca, 0x91, 0xde, 0x80, 0xd5, 0x00, 0x02, 0x00, 0x54, 0x8f, 0xd0,
		0xc0, 0xef, 0x1c, 0x8a, 0xdc, 0x3c, 0x52, 0x9b, 0xf5, 0xe1, 0x9f, 0xc7,
		0x86, 0xca, 0x91, 0xde, 0x80, 0xd5, 0x9d, 0x9c, 0x6c, 0x68, 0xa3, 0xe7,
		0x6f, 0x1e, 0xdf, 0xd3, 0x21, 0xa8, 0x53, 0x68, 0x8d, 0xf1, 0x30, 0xe2,
		0xf2, 0x49, 0x5c, 0x42, 0xd4, 0x2b, 0x00, 0xc6, 0xdb, 0x11, 0x0f, 0xbf,
		0x95, 0x11, 0x45, 0xe3, 0xad, 0xf4, 0x91, 0x4d, 0x6c, 0x74, 0xff, 0xb3,
		0x22, 0x83, 0xb7, 0x2a, 0xa3, 0xa1, 0x7f, 0x3c, 0x26, 0x31, 0x6a, 0x17,
		0xda, 0x63, 0x2b, 0x8b, 0x84, 0x00, 0xdc, 0x8e, 0xd8, 0x41, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
	}
	rc := NewRelayCellFromBytes(payload)
	data, err := rc.RelayData()
	require.NoError(t, err)
	e := Extend2Payload{}
	err = e.UnmarshalBinary(data)
	require.NoError(t, err)

	assert.Len(t, e.LinkSpecs, 2)
	AssertLinkSpecEqual(t, NewLinkSpecTCP(net.IPv4(127, 0, 0, 1), 5002), e.LinkSpecs[0])
	AssertLinkSpecEqual(t, NewLinkSpecLegacyID([]byte{
		0x8f, 0xd0, 0xc0, 0xef, 0x1c, 0x8a, 0xdc, 0x3c, 0x52, 0x9b,
		0xf5, 0xe1, 0x9f, 0xc7, 0x86, 0xca, 0x91, 0xde, 0x80, 0xd5,
	}), e.LinkSpecs[1])
	assert.Equal(t, data[31:], e.HandshakeData)
}
