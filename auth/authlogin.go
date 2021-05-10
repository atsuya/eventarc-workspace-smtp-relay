package auth

import (
	"errors"
	"net/smtp"
)

// GoogleAccountAuth represents data required for AUTH LOGIN
type GoogleAccountAuth struct {
	Username string
	Password string
}

// AuthLogin returns Auth implementation that performs AUTH LOGIN
func AuthLogin(username string, password string) smtp.Auth {
	return &GoogleAccountAuth{Username: username, Password: password}
}

// Start implements smtp.Auth.Start
func (g *GoogleAccountAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

// Next implements smtp.Auth.Next
func (g *GoogleAccountAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(g.Username), nil
		case "Password:":
			return []byte(g.Password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}

	return nil, nil
}
