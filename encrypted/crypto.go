package encrypted

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
)

const (
	PrivKeyLen   = 64
	signatureLen = 64
	PubKeyLen    = 32
	seedLen      = 32
	addressLen   = 30
)

type PublicKey struct {
	key ed25519.PublicKey
}

func (p *PublicKey) Address() Address {
	return Address{
		value: p.key[len(p.key)-addressLen:],
	}
}

func PublicKeyFromBytes(b []byte) *PublicKey {
	if len(b) != PubKeyLen {
		panic("invalid public key length")
	}
	return &PublicKey{
		key: ed25519.PublicKey(b),
	}
}

type PrivateKey struct {
	key ed25519.PrivateKey
}

// Genarate New Private Key from string
func newPrivateKeyFromString(s string) *PrivateKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return newPrivateKeyFromSeed(b)
}

// Genarate New Private Key from seed
func newPrivateKeyFromSeed(seed []byte) *PrivateKey {
	if len(seed) != seedLen {
		panic("Invalid length")
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func GeneratePrivateKey() *PrivateKey {
	seed := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		panic(err)
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func (p *PrivateKey) Bytes() []byte {
	return p.key
}

func (p *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		ed25519.Sign(p.key, msg),
	}
}

func (p *PrivateKey) Public() *PublicKey {
	b := make([]byte, PubKeyLen)
	copy(b, p.key[32:])

	return &PublicKey{
		key: b,
	}
}

func (p *PublicKey) Bytes() []byte {
	return p.key
}

type Signature struct {
	value []byte
}

func SignatureFromBytes(b []byte) *Signature {
	if len(b) != signatureLen {
		panic("length of the bytes nor equal to 64")
	}
	return &Signature{
		value: b,
	}
}

func (s *Signature) Bytes() []byte {
	return s.value
}

func (s *Signature) Verify(pubKey *PublicKey, msg []byte) bool {
	return ed25519.Verify(pubKey.key, msg, s.value)
}

type Address struct {
	value []byte
}

func (a Address) Bytes() []byte {
	return a.value
}

func (a Address) String() string {
	return hex.EncodeToString(a.value)
}