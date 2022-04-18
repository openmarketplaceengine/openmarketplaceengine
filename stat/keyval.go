package stat

import (
	"database/sql"
	"sync"
)

type KVOption int

const (
	KVEscapeKeys KVOption = 1 << iota
	KVReleaseVal          = 1 << iota
)

//-----------------------------------------------------------------------------
// IntKeyVal
//-----------------------------------------------------------------------------

type IntKeyVal struct {
	Key []string
	Val []int64
	Opt KVOption
}

//-----------------------------------------------------------------------------

var _intKeyValPool = sync.Pool{New: func() interface{} {
	return new(IntKeyVal).Alloc(16)
}}

func GetIntKeyVal(opts ...KVOption) *IntKeyVal {
	kv, _ := _intKeyValPool.Get().(*IntKeyVal)
	kv.Opt = maskOpts(opts)
	return kv
}

//-----------------------------------------------------------------------------

func (kv *IntKeyVal) Get(index int) (string, int64) {
	return kv.Key[index], kv.Val[index]
}

//-----------------------------------------------------------------------------

func (kv *IntKeyVal) Add(key string, val int64) {
	kv.Key = append(kv.Key, key)
	kv.Val = append(kv.Val, val)
}

//-----------------------------------------------------------------------------

func (kv *IntKeyVal) Alloc(size int) *IntKeyVal {
	kv.Key = make([]string, 0, size)
	kv.Val = make([]int64, 0, size)
	return kv
}

//-----------------------------------------------------------------------------

func (kv *IntKeyVal) Len() int {
	return len(kv.Key)
}

func (kv *IntKeyVal) Empty() bool {
	return len(kv.Key) == 0
}

//-----------------------------------------------------------------------------

func (kv *IntKeyVal) Scan(rows *sql.Rows) error {
	var key string
	var val int64
	for rows.Next() {
		err := rows.Scan(&key, &val)
		if err != nil {
			return err
		}
		kv.Add(key, val)
	}
	return nil
}

//-----------------------------------------------------------------------------

func (kv *IntKeyVal) Total() int64 {
	var tot int64
	for i := 0; i < len(kv.Val); i++ {
		tot += kv.Val[i]
	}
	return tot
}

//-----------------------------------------------------------------------------

func (kv *IntKeyVal) Release() {
	kv.Reset()
	_intKeyValPool.Put(kv)
}

//-----------------------------------------------------------------------------

func (kv *IntKeyVal) Reset() {
	if n := len(kv.Key); n > 0 {
		for i := 0; i < n; i++ {
			kv.Key[i] = ""
		}
		kv.Key = kv.Key[:0]
		kv.Val = kv.Val[:0]
	}
	kv.Opt = 0
}

func (kv *IntKeyVal) WriteJSON(b *JSONBuffer) error {
	if kv.Empty() {
		b.EmptyObject()
		return nil
	}
	escKey := (kv.Opt&KVEscapeKeys != 0)
	b.ObjectStart()
	for i := 0; i < kv.Len(); i++ {
		key, val := kv.Get(i)
		b.Key(key, escKey)
		b.Int64(val)
		b.Comma()
	}
	b.ObjectClose()
	return nil
}

//-----------------------------------------------------------------------------
// Options
//-----------------------------------------------------------------------------

func maskOpts(opts []KVOption) (mask KVOption) {
	for i := 0; i < len(opts); i++ {
		mask |= opts[i]
	}
	return
}
