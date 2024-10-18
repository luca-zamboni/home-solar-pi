package worker

import (
	"home-solar-pi/pkg/db"
	"home-solar-pi/pkg/device"
	"log"
	"math/rand"
	"time"
)

const debug = true

type WorkerSevice struct {
	inverterService *device.InterverService
	heaterService   *device.HeaterService
	logger          *log.Logger
	dbService       *db.DbService
	threshold       int
}

func NewWorkerSevice(inverterService *device.InterverService, heaterService *device.HeaterService,
	logger *log.Logger, dbService *db.DbService, threshold int) WorkerSevice {
	return WorkerSevice{
		inverterService: inverterService,
		heaterService:   heaterService,
		logger:          logger,
		dbService:       dbService,
		threshold:       threshold,
	}
}

func (w *WorkerSevice) UpdateInverter(interval time.Duration) {

	for {

		reading := w.getInverterReading()

		w.logger.Printf("Power %+v\n", reading)

		w.dbService.InsertReading(reading)

		if reading > w.threshold {
			_, err := w.heaterService.PowerOn()

			if err != nil {
				w.logger.Printf("Error heater := %s\n", err.Error())
			}

			w.logger.Println("Heater activated")
			w.dbService.InsertHeaterlog(reading, true)
		}

		time.Sleep(interval)
	}

}

func (w *WorkerSevice) getInverterReading() int {
	var power *device.InverterResponse

	if debug {
		power = &device.InverterResponse{}

		power.Body.Data.PAC.Values = map[string]int{}
		power.Body.Data.PAC.Values["1"] = rand.Int()%300 + 400
		power.Body.Data.PAC.Unit = "W"

	} else {
		var err error
		power, err = w.inverterService.GetCurrentPower()

		if err != nil {
			w.logger.Printf("Error retrieving current power\n")
		}
	}

	var reading int = power.Body.Data.PAC.Values["1"]

	return reading
}
