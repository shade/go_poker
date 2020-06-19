package identity

type Record struct {
	Name         string `json:"name"`
	Username     string `json:"username"`
	Password     string `json:"password",store:"false"`
	PasswordHash string `json:"password_hash"`
}

type IdentityGenerator struct {
	db *IdentityDB
	secret string
}

func (i *IdentityGenerator) ValidateRecord(r *Record) {
	// Ensure no user with same username in the db
	if db.Get(r.Username) != nil {
		return nil, errors.New()
	}
	
	// 
}

func (i *IdentityGenerator) 