// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package arg

import (
	"context"
	"fmt"
	"io"
	"os"
)

type (
	Context = context.Context
)

type ArgSet struct {
	FlagSet
}

//-----------------------------------------------------------------------------

func NewArgSet() *ArgSet {
	return new(ArgSet).Init()
}

//-----------------------------------------------------------------------------

func (s *ArgSet) Init() *ArgSet {
	s.SetOutput(io.Discard)
	s.Usage = func() {}
	return s
}

//-----------------------------------------------------------------------------

func (s *ArgSet) Parse(ctx Context, args []string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	err := s.FlagSet.parse(ctx, args)
	if err == nil {
		return nil
	}
	switch s.errorHandling {
	case ContinueOnError:
		if s.checkHelp(err) {
			return nil
		}
		return err
	case ExitOnError:
		if s.checkHelp(err) {
			os.Exit(0)
		}
		_, _ = fmt.Fprintln(s.Output(), err)
		os.Exit(2)
	case PanicOnError:
		if s.checkHelp(err) {
			os.Exit(0)
		}
		panic(err)
	}
	return err
}

func (s *ArgSet) checkHelp(err error) bool {
	if err == ErrHelp {
		s.PrintHelp()
		return true
	}
	return false
}

//-----------------------------------------------------------------------------

func (s *ArgSet) PrintHelp() {
	s.SetOutput(os.Stdout)
	fmt.Println("Usage:")
	s.PrintDefaults()
}
