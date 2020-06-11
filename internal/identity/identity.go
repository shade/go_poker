package identity

type IdentityGenerator struct {
	secret string
}

// NewIdentityGenerator creates an identity generator provided
// a specific token secret for creating HMAC keys.
func NewIdentityGenerator(secret string) *IdentityGenerator {
	return &IdentityGenerator{
		secret,
	}
}

// NewToken Creates a new token provided a unique id,
// note this method is not idempotent
func (i *IdentityGenerator) NewToken(id string) []byte {
	// TODO
	return []byte(i.secret)
}

func (i *IdentityGenerator) Validate(t []byte) bool {
	// TODO
	return true
}
