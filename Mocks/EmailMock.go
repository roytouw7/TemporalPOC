package Mocks

import "log"

func SendEmail(msg string) (bool, error) {
	log.Println(msg)
	return true, nil
}
