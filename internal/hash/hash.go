package hash

type Hasher interface {
	Generate(password []byte) ([]byte, error)
	Compare(password, hash []byte) (bool, error)
}
