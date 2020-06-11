package identity

type IdentityGenerator struct {
    secret string
}

func NewIdentityGenerator(secret string) *IdentityGenerator {
    return &IdentityGenerator{
        secret,
    }
}

func (i *IdentityGenerator) NewToken(id string) []byte {
    // TOOD
}

func (i *IdentityGenerator) Validate(t []byte) bool {
    // TODO
    return true
}
