package identity

type Identity struct {
	unique []byte
	data   Payload
}

type IdentityDB []Identity

 
func (*IdentityDB) Get(unique string) (Identity, bool) {

}

func (*IdentityDB) Insert(payload Identity, unique string) bool {
	id := Identity {
		unique, 
		payload
	}
}

func (*IdentityDB) Sync(path string) bool {
}
