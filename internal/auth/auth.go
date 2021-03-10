package auth

import "golang.org/x/crypto/bcrypt"

type Auth struct {
	// pepper is a static salt.
	pepper []byte
	cost   int
}

func New(pepper []byte, cost int) *Auth {
	return &Auth{
		pepper: pepper,
		cost:   cost,
	}
}

func (a *Auth) generateHashFromPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(append(password, a.pepper...), a.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (a *Auth) compareHashAndPassword(hashedPassword, password []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, append(password, a.pepper...)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (a *Auth) Login() error {

	return nil
}

func (a *Auth) Signup() error {
	return nil
}
