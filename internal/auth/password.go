package auth

import "golang.org/x/crypto/bcrypt"

type Password struct {
	Pass *[]byte
	Hash []byte
}

func (p *Password) Set(pass []byte) error{
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

    p.Pass = &pass
    p.Hash = hashed

    return nil
}

func (p *Password) Compare(input []byte) error {
	return bcrypt.CompareHashAndPassword(p.Hash, input)
}
