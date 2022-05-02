package main

import (
	"flag"
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/detector"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/tollgate/yaml"
	"log"
	"os"
	"path/filepath"
)

func main() {

	var t string
	var l string
	flag.StringVar(&t, "t", "", "path to tollgates yaml file")
	flag.StringVar(&l, "l", "", "path to locations csv file")
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

	//lFile, err := os.Open(lPath)
	//if err != nil {
	//	log.Fatalf("open locations file error: %s", err)
	//}
	//defer tFile.Close()

	tolls, err := yaml.ReadYaml(tFile)
	if err != nil {
		log.Fatalf("readin tollgates error: %s", err)
	}

	fmt.Printf("yaml %v\n", tolls)

	tollgates := transform(tolls)
	fmt.Printf("transformed %v\n", tollgates)

	//d := detector.NewDetector(tollgates, detector.NewMapStorage())

	//d.DetectCrossing(context.Background(), m)
}

func transform(tollgates []yaml.Tollgate) []*detector.Tollgate {
	var r []*detector.Tollgate
	for _, t := range tollgates {
		var boxes []*detector.BBox

		for _, b := range t.BBoxes.Boxes {
			boxes = append(boxes, &detector.BBox{
				LonMin: b.LonMin,
				LatMin: b.LatMin,
				LonMax: b.LonMax,
				LatMax: b.LatMax,
			})
		}

		r = append(r, &detector.Tollgate{
			ID: t.ID,
			Line: &detector.Line{
				Lon1: t.Line.Lon1,
				Lat1: t.Line.Lat1,
				Lon2: t.Line.Lon2,
				Lat2: t.Line.Lat2,
			},
			BBoxes:         boxes,
			BBoxesRequired: t.BBoxes.BoxesRequired,
		})
	}
	return r
}
