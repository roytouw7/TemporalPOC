package email

import (
	"github.com/google/uuid"
)

type ReservationService interface {
	Upgrade() (success bool, err error)
}

type reservationService struct {
	id           string
	emailService EmailService
}

func NewReservationService(emailService EmailService) ReservationService {
	return &reservationService{
		id:           uuid.NewString(),
		emailService: emailService,
	}
}

func (r *reservationService) Upgrade() (success bool, err error) {
	return r.emailService.SendUpgradeEmail(r.id)
}
