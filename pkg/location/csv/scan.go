package csv

import (
	"bufio"
	"io"
	"strings"
)

type Scanner struct {
	s *bufio.Scanner
}

func NewScan(r io.Reader) *Scanner {
	s := bufio.NewScanner(r)
	return &Scanner{
		s: s,
	}
}

func (s *Scanner) NextLocation() (*Location, error) {

	for {
		scan := s.s.Scan()
		if !scan {
			return nil, s.s.Err()
		}

		line := s.s.Text()
		isHeader := strings.Contains(line, "driver_id")
		if isHeader || strings.TrimSpace(line) == "" {
			continue
		}
		return Parse(line)
	}
}
