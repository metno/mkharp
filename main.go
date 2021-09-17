package main

import (
	"flag"
	"log"
	"os"

	"github.com/metno/mkharp/internal/harp/obs"
	"github.com/metno/mkharp/internal/input"
)

func main() {
	out := flag.String("out", "harp.sqlite", "write to the given file")
	obstype := flag.String("obstype", "synop", "write the given observation type")
	sid := flag.Int("sid", 1, "station id to use")
	lon := flag.Float64("lon", 0, "Longitude")
	lat := flag.Float64("lat", 0, "Latitude")
	elev := flag.Int("elevation", 0, "elevation")
	createNew := flag.Bool("create", false, "initialize a new database")
	flag.Parse()

	reader, err := input.Open(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	observations, err := reader.Read()
	if err != nil {
		log.Fatalln(err)
	}

	data := obs.Data{
		SID:          *sid,
		Lon:          float32(*lon),
		Lat:          float32(*lat),
		Elev:         *elev,
		Observations: observations,
	}

	var db *obs.Database
	if *createNew {
		parameters, err := data.Parameters()
		if err != nil {
			log.Fatalln(err)
		}
		db, err = obs.Create(*out, *obstype, parameters)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		db, err = obs.Open(*out, *obstype)
		if err != nil {
			log.Fatalln(err)
		}
	}
	defer db.Close()

	if err := db.Add(data); err != nil {
		log.Fatalln(err)
	}
}
