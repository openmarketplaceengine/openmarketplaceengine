package stat

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync/atomic"
	"time"
)

type (
	Context = context.Context
)

type Func = func(ctx Context) (interface{}, error)

type Stat struct {
	name string
	help string
	stat Func
	list *List // sub stats
}

type List struct {
	stat []*Stat
}

type Releaser interface {
	Release()
}

var (
	statNreq uint64
)

//-----------------------------------------------------------------------------
// Stat
//-----------------------------------------------------------------------------

func (s *Stat) Skip() bool {
	return s.stat == nil && (s.list == nil || s.list.Len() == 0)
}

func (s *Stat) Group() bool {
	return s.list != nil
}

func (s *Stat) GroupCount() int {
	if s.list != nil {
		return s.list.Len()
	}
	return 0
}

//-----------------------------------------------------------------------------
// List
//-----------------------------------------------------------------------------

func (ls *List) Add(name string, help string, stat Func) {
	ls.checkName(name, "Stat name")
	if stat == nil {
		panic("Stat func is nil")
	}
	ls.add(&Stat{name: name, help: help, stat: stat})
}

func (ls *List) Group(name string, help string) *List {
	MustJSONField(name, "Group name")
	if stat := ls.Stat(name); stat != nil {
		if stat.list == nil {
			panic(fmt.Sprintf("Group name is in use by a Stat: %q", name))
		}
		return stat.list
	}
	group := &Stat{name: name, list: new(List)}
	if len(help) > 0 {
		group.help = help
	}
	ls.add(group)
	return group.list
}

//-----------------------------------------------------------------------------

func (ls *List) Len() int {
	return len(ls.stat)
}

func (ls *List) Less(i, j int) bool {
	return ls.stat[i].name < ls.stat[j].name
}

func (ls *List) Swap(i, j int) {
	ls.stat[i], ls.stat[j] = ls.stat[j], ls.stat[i]
}

func (ls *List) sort() {
	sort.Sort(ls)
}

//-----------------------------------------------------------------------------

func (ls *List) Find(name string) int {
	a := ls.stat
	n := len(a)
	x := sort.Search(n, func(i int) bool {
		return a[i].name >= name
	})
	if x < n && a[x].name == name {
		return x
	}
	return -1
}

func (ls *List) Stat(name string) *Stat {
	if x := ls.Find(name); x != -1 {
		return ls.stat[x]
	}
	return nil
}

func (ls *List) checkName(name string, errorPrefix string) {
	MustJSONField(name, errorPrefix)
	if x := ls.Find(name); x != -1 {
		panic(fmt.Errorf("%s is already in use: %q", errorPrefix, name))
	}
}

func (ls *List) add(stat *Stat) {
	ls.stat = append(ls.stat, stat)
	ls.sort()
}

//-----------------------------------------------------------------------------

func listJSON(ctx Context, list *List, buf *JSONBuffer) error {
	n := list.Len()
	if n == 0 {
		buf.EmptyObject()
		return nil
	}
	now := time.Now()
	buf.ObjectStart()
	for i := 0; i < n; i++ {
		s := list.stat[i]
		if s.Skip() {
			continue
		}
		buf.Key(s.name, false)
		if s.Group() {
			err := listJSON(ctx, s.list, buf)
			if err != nil {
				return err
			}
			buf.Comma()
			continue
		}
		val, err := s.stat(ctx)
		if err != nil {
			release(val)
			return fmt.Errorf("%q stat() failed: %w", s.name, err)
		}
		err = buf.Value(val)
		if err != nil {
			release(val)
			return fmt.Errorf("%q encoding failed: %w", s.name, err)
		}
		release(val)
		buf.Comma()
	}
	if buf.Lev() == 1 {
		statTime := float64(time.Since(now)) / float64(time.Millisecond)
		buf.Key("stat_time", false)
		_ = buf.Float64(math.Round(statTime*1000) / 1000)
		nreg := atomic.AddUint64(&statNreq, 1)
		buf.Comma()
		buf.Key("stat_nreq", false)
		buf.Uint64(nreg)
	}
	buf.ObjectClose()
	return nil
}

//-----------------------------------------------------------------------------

func release(v interface{}) {
	if r, ok := v.(Releaser); ok {
		r.Release()
	}
}
