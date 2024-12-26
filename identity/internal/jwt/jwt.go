package jwt

type JWT string

type Pair struct {
	Access  JWT
	Refresh JWT
}
