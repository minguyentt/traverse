package auth

import "golang.org/x/crypto/bcrypt"

type PasswordHash interface {
	generateHash([]byte) error
	compare([]byte) error
}
type password struct {
	text *[]byte
	hash []byte
}

func (p *password) generateHash(text []byte) error {
	hashed, err := bcrypt.GenerateFromPassword(text, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hashed

	return nil
}

func (p *password) compare(input []byte) error {
	return bcrypt.CompareHashAndPassword(p.hash, input)
}
