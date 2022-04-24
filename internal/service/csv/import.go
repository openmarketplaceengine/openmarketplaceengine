package csv

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ImportLocations(rd io.Reader) error {
	s := bufio.NewScanner(rd)

	var line int
	for {
		scan := s.Scan()
		line++
		if !scan {
			err := s.Err()
			if err != nil {
				return fmt.Errorf("scan error:%w", err)
			}
			return nil
		}
		text := s.Text()
		isHeader := strings.Contains(text, "driver_id")
		if isHeader {
			continue
		}
		location, err := Parse(text)
		if err != nil {
			return fmt.Errorf("line %v parse location error:%w", line, err)
		}
		fmt.Printf("%v\n", location)
	}
}
