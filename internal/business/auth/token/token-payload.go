package token

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	// unique id to invalide a token if its leaked
	ID       uuid.UUID `json:"token_id"`
	Username string    `json:"username"`
	ExpireAt time.Time `json:"expire_at"`
	IssuedAt time.Time `json:"issued_at"`
}

// create a new token payload, receiving a username and duration time for expiration
func NewTokenPayload(username string, expiration time.Duration) (*Payload, error) {
	// generate a random uuid for the token id
	tokenID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("error creating a token payload | %v \n", err)
		return nil, err
	}

	// setup the expiration duration
	expireAt := time.Now().Add(expiration)

	// set the token payload params
	payload := &Payload{
		ID:       tokenID,
		Username: username,
		ExpireAt: expireAt,
		IssuedAt: time.Now(),
	}

	// return the token
	return payload, nil
}

func (p *Payload) Valid() bool {
	if time.Now().After(p.ExpireAt) {

		log.Println("the expiration date is not valid")
		return false
	}
	log.Println("the expiration date is  valid")

	return true
}
