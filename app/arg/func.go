// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package arg

import "fmt"

type fset = func(ctx Context, set *FlagSet, val string) error
type Validator = func(arg string) (string, error)

//-----------------------------------------------------------------------------

func (f *FlagSet) Void(name string, help string, call func(ctx Context) error) *Flag {
	checkNameCall(name, call == nil)
	flag := f.Var(new(boolValue), name, help)
	flag.fset = func(ctx Context, set *FlagSet, val string) error {
		if len(val) > 0 {
			return fmt.Errorf("unused argument for flag: -%s", name)
		}
		return call(ctx)
	}
	return flag
}

//-----------------------------------------------------------------------------

func (f *FlagSet) String(name string, help string, call func(ctx Context, arg string) error) *Flag {
	checkNameCall(name, call == nil)
	flag := f.Var(new(stringValue), name, help)
	flag.fset = func(ctx Context, set *FlagSet, val string) error {
		if len(val) > 0 {
			return call(ctx, val)
		}
		arg, ok := set.next()
		if !ok && flag.defv != nil {
			if str, sok := flag.defv.(string); sok {
				arg = str
			}
		}
		if len(arg) == 0 && !flag.hasOpt(argEmpty) {
			return fmt.Errorf("flag needs an argument: -%s", name)
		}
		return call(ctx, arg)
	}
	return flag
}

//-----------------------------------------------------------------------------

func (f *FlagSet) Rest(name string, help string, call func(ctx Context, args []string) error) *Flag {
	checkNameCall(name, call == nil)
	flag := f.Var(new(stringValue), name, help)
	flag.fset = func(ctx Context, set *FlagSet, val string) error {
		if len(val) > 0 {
			return call(ctx, []string{val})
		}
		args := set.Args()
		if !flag.hasOpt(argEmpty) {
			if len(args) == 0 {
				return fmt.Errorf("flag needs an argument: -%s", name)
			}
			for i := range args {
				if len(args[i]) == 0 {
					return fmt.Errorf("empty argument #%d for flag: -%s", i, name)
				}
			}
		}
		return call(ctx, args)
	}
	return flag
}

//-----------------------------------------------------------------------------

func checkNameCall(name string, callIsNil bool) {
	if len(name) == 0 {
		panic("Empty argument name")
	}
	if callIsNil {
		panic("Nil 'call' argument")
	}
}
