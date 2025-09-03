// defines a token maker used to manage token based auth

package token

import "time"

type Maker interface{
	// creates a token for a user for given duration, returns token as string and error if any
	CreateToken(username string, duration time.Duration) (string, error)

	// checks if the token is valid
	VerifyToken(token string) (*Payload, error)
}