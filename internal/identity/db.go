package identity

type DBKey string

type Record struct {
	Name         string `json:"name"`
	Username     string `json:"username"`
	Password     string `json:"password",store:"false"`
	PasswordHash string `json:"password_hash"`
}

func (r Record) PrimaryKey() DBKey {
	return DBKey(r.Username)
}

func (r Record) ToCSVRecord() []string {
	return []string{
		r.Name,
		r.Username,
		r.PasswordHash,
	}
}

// Interface for the Identity DB.
type IIDB interface {
	Get(key DBKey) (*Record, error)
	Insert(record *Record) error
	Delete(key DBKey)
}
