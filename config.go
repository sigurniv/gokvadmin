package gokvadmin

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Config struct {
	Port string
	Auth *Auth
	TLS  *TLS
}

type TLS struct {
	CertFile string
	KeyFile  string
}

type Auth struct {
	Login        string
	Password     string
	token        string
	tokenExpires time.Time
}

func (a *Auth) GenerateToken() (string, error) {
	var err error
	if a == nil {
		return "token", err
	}

	if len(a.token) > 0 {
		return a.token, err
	}

	b := make([]byte, 8)
	rand.Read(b)

	hash, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	token := string(hash);
	a.token = token;
	a.tokenExpires = time.Now().Local().Add(time.Hour * time.Duration(8))

	return token, err
}

var DefaultConfig = Config{
	Port: "8083",
}