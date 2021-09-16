package main

import (
	"log"
	"time"

	"github.com/metno/mkharp/internal/harp/obs"
)

func main() {

	parameters := []obs.Parameter{
		{
			Parameter:  "T2m",
			AccumHours: 0.0,
			Units:      "degC",
		},
		{
			Parameter:  "AccPcp12h",
			AccumHours: 12.0,
			Units:      "kg/m^2",
		},
	}
	data := obs.Data{

		SID:  1492,
		Lon:  10.72,
		Lat:  59.9423,
		Elev: 94,
		Observations: []obs.Observation{
			{
				ValidDate: time.Date(2021, 9, 16, 6, 0, 0, 0, time.UTC),
				Data: map[string]float32{
					"T2m":       12.2,
					"AccPcp12h": 2.1,
				},
			},
		},
	}

	db, err := obs.Create("harp.sqlite", "synop", parameters)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if err := db.Add(data); err != nil {
		log.Fatalln(err)
	}
}
