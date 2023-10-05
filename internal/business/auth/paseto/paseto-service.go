package paseto

import (
	"fmt"
	"time"

	"github.com/FadyGamilM/go-websockets/internal/business/auth/token"
	"github.com/FadyGamilM/go-websockets/internal/core"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// Paseto struct is a type which implements the TokenMaker interface
type Paseto struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPaseto(key string) (core.TokenMaker, error) {
	// ensure that the length of the symmetric key is the standard of the paseto library
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symmetric key length")
	}
	return &Paseto{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(key),
	}, nil
}

func (p *Paseto) Create(username string, expiration time.Duration) (string, *token.Payload, error) {
	// create a payload for the token payload
	payload, err := token.NewTokenPayload(username, expiration)
	if err != nil {
		return "", nil, err
	}
	token, err := p.paseto.Encrypt(p.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

func (p *Paseto) Verify(tokenVal string) (*token.Payload, error) {
	payload := new(token.Payload)
	err := p.paseto.Decrypt(tokenVal, p.symmetricKey, &payload, nil)
	if err != nil {
		return nil, fmt.Errorf("error trying to decrypt the token | %v", err.Error())
	}

	// check if the payload is valid or not
	isValid := payload.Valid()

	if !isValid {
		return nil, fmt.Errorf("token is expired")
	}
	return payload, nil
}
