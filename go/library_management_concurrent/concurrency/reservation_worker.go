package concurrency

import (
	"github.com/Ruth652/playground/go/library_management/services"
)

func StartReservationWorker(lib *services.Library) {
	go func() {
		for req := range lib.ReservationChan {
			lib.ProcessReservation(req)
		}
	}()
}
