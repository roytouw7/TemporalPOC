package Mocks

import "log"

// GetRoomToUpgrade trivial mock which should represent getting the room to upgrade to for given reservation id
func GetRoomToUpgrade(reservationId string) (room string, err error) {
	log.Println("The hotel is notified")
	return "Suite", nil
}
