package auth_service

import "simple-rest/data"

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() (int, error) {
	return data.CheckAuth(a.Username, a.Password)
}
