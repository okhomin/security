package bcrypt

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct {
	// pepper is a static salt, min 32 chars, always in prefix.
	pepper []byte

	// cost is a work factor. 2^cost iterations for hashing.
	cost int
}

func New(pepper []byte, cost int) *Bcrypt {
	return &Bcrypt{
		pepper: pepper,
		cost:   cost,
	}
}

func (b *Bcrypt) Generate(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(append(b.pepper, password...), b.cost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (b *Bcrypt) Compare(hashedPassword, password []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, append(b.pepper, password...)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
