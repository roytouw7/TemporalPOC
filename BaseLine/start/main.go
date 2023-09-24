package main

import "log"

func main() {
	service, cleanup, err := initialiseApplication()
	defer cleanup()

	if err != nil {
		panic(err)
	}

	ok, err := service.Upgrade()
	if err != nil {
		log.Fatal("failed upgrading", err)
	}
	if ok {
		log.Println("Upgraded reservation!")
	}
}
