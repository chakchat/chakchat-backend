package jwt

type InternalJWT string

type JWT string

type Pair struct {
	Access  JWT
	Refresh JWT
}
