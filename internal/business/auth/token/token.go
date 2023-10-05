package token

import (
	"time"
)

// token-maker has two methods, the create token and verify token method
// this is a specs for any type of token auth impl
type TokenMaker interface {
	/*
		@Params:
			username => to be included in the payload
			expiration => to set a short-life time for the token
	*/
	Create(username string, expiration time.Duration) (string, *Payload, error)
	Verify(token string) (*Payload, error)
}
