package auth

import "golang.org/x/crypto/bcrypt"

type Password interface {
	Hash([]byte) error
	Compare([]byte) error
}

type password struct {
	text *[]byte
	hash []byte
}

func (p *password) Hash(text []byte) error {
	hashed, err := bcrypt.GenerateFromPassword(text, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hashed

	return nil
}

func (p *password) Compare(input []byte) error {
	return bcrypt.CompareHashAndPassword(p.hash, input)
}
