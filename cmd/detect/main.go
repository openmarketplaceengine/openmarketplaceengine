package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/location/csv"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/tollgate/yaml"
	"golang.org/x/net/context"
)

func main() {
	var t string
	var l string
	flag.StringVar(&t, "t", "pkg/tollgate/yaml/tollgates.yaml", "path to tollgates yaml file")
	flag.StringVar(&l, "l", "pkg/location/csv/testdata/coopdrive-gps-pings-2022.04.06.csv", "path to locations csv file")
	flag.Parse()

	if t == "" || l == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	tPath, err := filepath.Abs(t)
	if err != nil {
		log.Fatalf("%s to absolute path error: %s", t, err)
	}

	lPath, err := filepath.Abs(l)
	if err != nil {
		log.Fatalf("%s to absolute path error: %s", l, err)
	}
	fmt.Printf("tollgates file: %s\n", tPath)
	fmt.Printf("locations file: %s\n", lPath)

	tFile, err := os.Open(tPath)
	if err != nil {
		log.Fatalf("open tollgates file error: %s", err)
	}
	defer tFile.Close()

	lFile, err := os.Open(lPath)
	if err != nil {
		log.Fatalf("open locations file error: %s", err)
	}
	defer tFile.Close()

	tolls, err := yaml.ReadYaml(tFile)
	if err != nil {
		log.Fatalf("readin tollgates error: %s", err)
	}

	tollgates := transform(tolls)

	d := detector.NewDetector(tollgates, detector.NewMapStorage())

	scan := csv.NewScan(lFile)

	var crossings []*detector.Crossing
	var from *detector.Location
	for {
		location, err := scan.NextLocation()

		if err != nil {
			log.Fatalf("scan next location error: %s", err)
		}

		if location == nil {
			break
		}

		to := &detector.Location{
			Lon: location.Lon,
			Lat: location.Lat,
		}

		if from == nil {
			from = to
			continue
		}

		crossing, err := d.DetectCrossing(context.Background(), location.DriverID, &detector.Movement{
			From: from,
			To:   to,
		})

		if err != nil {
			log.Fatalf("detect  crossing error: %s", err)
		}
		from = to
		if crossing != nil {
			crossings = append(crossings, crossing)
		}
	}

	for i, c := range crossings {
		fmt.Printf("detected crossing [%v] %v\n", i, c)
	}
}

func transform(tollgates []yaml.Tollgate) []*detector.Tollgate {
	r := make([]*detector.Tollgate, len(tollgates))
	for i, t := range tollgates {
		var boxes []*detector.BBox

		for _, b := range t.BBoxes.Boxes {
			boxes = append(boxes, &detector.BBox{
				LonMin: b.LonMin,
				LatMin: b.LatMin,
				LonMax: b.LonMax,
				LatMax: b.LatMax,
			})
		}

		r[i] = &detector.Tollgate{
			ID: t.ID,
			Line: &detector.Line{
				Lon1: t.Line.Lon1,
				Lat1: t.Line.Lat1,
				Lon2: t.Line.Lon2,
				Lat2: t.Line.Lat2,
			},
			BBoxes:         boxes,
			BBoxesRequired: t.BBoxes.BoxesRequired,
		}
	}
	return r
}
