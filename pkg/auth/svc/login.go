package svc

import (
	"fmt"
	"net/http/cookiejar"
)

func Login(username string, password string) (*cookiejar.Jar, error) {

	fmt.Println("yay")
	return nil, nil

}
