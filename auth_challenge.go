package pearl

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
)

// Reference: https://github.com/torproject/torspec/blob/master/tor-spec.txt#L601-L617
//
//	4.3. AUTH_CHALLENGE cells
//
//	   An AUTH_CHALLENGE cell is a variable-length cell with the following
//	   fields:
//	       Challenge [32 octets]
//	       N_Methods [2 octets]
//	       Methods   [2 * N_Methods octets]
//
//	   It is sent from the responder to the initiator. Initiators MUST
//	   ignore unexpected bytes at the end of the cell.  Responders MUST
//	   generate every challenge independently using a strong RNG or PRNG.
//
//	   The Challenge field is a randomly generated string that the
//	   initiator must sign (a hash of) as part of authenticating.  The
//	   methods are the authentication methods that the responder will
//	   accept.  Only one authentication method is defined right now:
//	   see 4.4 below.
//

// AuthMethod represents an authentication method ID.
type AuthMethod uint16

// Defined AuthMethod values.
var (
	AuthMethodRSASHA256TLSSecret   AuthMethod = 1
	AuthMethodEd25519SHA256RFC5705 AuthMethod = 3
)

// String represents the AuthMethod as a string. This is also the TYPE value
// expected in AUTHENTICATE cells.
func (m AuthMethod) String() string {
	return fmt.Sprintf("AUTH%04d", int(m))
}

// AuthChallengeCell represents an AUTH_CHALLENGE cell.
type AuthChallengeCell struct {
	Challenge [32]byte
	Methods   []AuthMethod
}

var _ CellBuilder = new(AuthChallengeCell)

// NewAuthChallengeCell builds an AUTH_CHALLENGE cell with the given method IDs.
// The challenge is generated at random.
func NewAuthChallengeCell(methods []AuthMethod) (*AuthChallengeCell, error) {
	var challenge [32]byte
	_, err := rand.Read(challenge[:])
	if err != nil {
		return nil, errors.Wrap(err, "could not read enough random bytes")
	}
	return &AuthChallengeCell{
		Challenge: challenge,
		Methods:   methods,
	}, nil
}

// NewAuthChallengeCellStandard builds an AUTH_CHALLENGE cell for method 1.
func NewAuthChallengeCellStandard() (*AuthChallengeCell, error) {
	return NewAuthChallengeCell([]AuthMethod{AuthMethodRSASHA256TLSSecret})
}

// Cell constructs the cell bytes.
func (a AuthChallengeCell) Cell(f CellFormat) (Cell, error) {
	// Reference: https://github.com/torproject/torspec/blob/master/tor-spec.txt#L605-L607
	//
	//	       Challenge [32 octets]
	//	       N_Methods [2 octets]
	//	       Methods   [2 * N_Methods octets]
	//
	m := len(a.Methods)
	n := 32 + 2 + 2*m
	c := NewCellEmptyPayload(f, 0, AuthChallenge, uint16(n))
	payload := c.Payload()

	copy(payload, a.Challenge[:])
	binary.BigEndian.PutUint16(payload[32:], uint16(m))
	ptr := 34
	for _, method := range a.Methods {
		binary.BigEndian.PutUint16(payload[ptr:], uint16(method))
		ptr += 2
	}

	return c, nil
}

// AuthenticateCell represents an AUTHENTICATE cell.
type AuthenticateCell struct {
	Method         AuthMethod
	Authentication []byte
}

func ParseAuthenticateCell(c Cell) (*AuthenticateCell, error) {
	if c.Command() != Authenticate {
		return nil, ErrUnexpectedCommand
	}

	payload := c.Payload()
	n := len(payload)

	if n < 4 {
		return nil, errors.New("authenticate cell too short")
	}

	method := binary.BigEndian.Uint16(payload)
	authLen := binary.BigEndian.Uint16(payload[2:])

	// Reference: https://github.com/torproject/torspec/blob/8aaa36d1a062b20ca263b6ac613b77a3ba1eb113/tor-spec.txt#L733-L735
	//
	//	   Responders MUST ignore extra bytes at the end of an AUTHENTICATE
	//	   cell.  Recognized AuthTypes are 1 and 3, described in the next
	//	   two sections.
	//
	if n < int(4+authLen) {
		return nil, errors.New("inconsistent authenticate cell length")
	}

	return &AuthenticateCell{
		Method:         AuthMethod(method),
		Authentication: payload[4 : 4+authLen],
	}, nil
}

func (a AuthenticateCell) Cell(f CellFormat) (Cell, error) {
	// Reference: https://github.com/torproject/torspec/blob/8aaa36d1a062b20ca263b6ac613b77a3ba1eb113/tor-spec.txt#L727-L731
	//
	//	   An AUTHENTICATE cell contains the following:
	//
	//	        AuthType                              [2 octets]
	//	        AuthLen                               [2 octets]
	//	        Authentication                        [AuthLen octets]
	//
	authLen := len(a.Authentication)
	c := NewCellEmptyPayload(f, 0, Authenticate, uint16(4+authLen))
	payload := c.Payload()

	binary.BigEndian.PutUint16(payload, uint16(a.Method))
	binary.BigEndian.PutUint16(payload[2:], uint16(authLen))
	copy(payload[4:], a.Authentication)

	return c, nil
}
