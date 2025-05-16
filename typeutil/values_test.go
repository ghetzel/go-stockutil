package typeutil

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ghetzel/testify/require"
)

type Thing struct {
	Name  string
	Value any
}

type subtime struct {
	time.Time
}

func TestIsZero(t *testing.T) {
	var assert = require.New(t)

	var thing Thing
	var things []Thing
	var madeThings = make([]Thing, 0)

	assert.True(IsZero(nil))
	assert.True(IsZero(0))
	assert.True(IsZero(0.0))
	assert.True(IsZero(false))
	assert.True(IsZero(``))
	assert.True(IsZero(thing))
	assert.True(IsZero(Thing{}))
	assert.True(IsZero(things))

	things = append(things, Thing{})

	assert.False(IsZero(1))
	assert.False(IsZero(0.1))
	assert.False(IsZero(true))
	assert.False(IsZero(`value`))
	assert.False(IsZero(Thing{`value`, true}))
	assert.False(IsZero(&Thing{}))
	assert.False(IsZero(things))
	assert.False(IsZero(madeThings))
}

func TestIsEmpty(t *testing.T) {
	var assert = require.New(t)

	var things = make([]Thing, 4)
	var thingmap = make(map[string]Thing)
	var stringmap = map[int]string{
		1: ``,
		2: `    `,
		3: "\t",
	}

	assert.True(IsEmpty(things))
	assert.True(IsEmpty(``))
	assert.True(IsEmpty(` `))
	assert.True(IsEmpty(`     `))
	assert.True(IsEmpty("\t\n  \n\t"))
	assert.True(IsEmpty(thingmap))
	assert.True(IsEmpty(stringmap))
}

func TestIsScalar(t *testing.T) {
	var assert = require.New(t)

	assert.True(IsScalar(1))
	assert.True(IsScalar(true))
	assert.True(IsScalar(3.14))
	assert.True(IsScalar(`four`))
	assert.False(IsScalar([]string{`1`}))
	assert.False(IsScalar(map[string]string{}))
	assert.False(IsScalar(make(chan string)))
	assert.False(IsScalar(time.Time{}))
}

func TestIsStruct(t *testing.T) {
	var assert = require.New(t)

	assert.False(IsStruct(1))
	assert.False(IsStruct(true))
	assert.False(IsStruct(3.14))
	assert.False(IsStruct(`four`))
	assert.False(IsStruct([]string{`1`}))
	assert.False(IsStruct(map[string]string{}))
	assert.False(IsStruct(make(chan string)))
	assert.True(IsStruct(time.Time{}))
}

func TestIsArray(t *testing.T) {
	var assert = require.New(t)

	assert.False(IsArray(nil))

	var a []string
	assert.True(IsArray(a))
	assert.True(IsArray([]string{`1`}))
	assert.True(IsArray(&[]string{`1`}))

	var b any
	b = []string{`1`}

	assert.True(IsArray(b))
	assert.True(IsArray(&b))

	assert.False(IsArray(``))
	assert.False(IsArray(`123`))
	assert.False(IsArray(123))
	assert.False(IsArray(true))
}

func TestIsFunction(t *testing.T) {
	var assert = require.New(t)

	assert.False(IsFunction(nil))
	assert.False(IsFunction(1))
	assert.False(IsFunction(true))
	assert.False(IsFunction(`three`))

	assert.True(IsFunction(func() {}))
	assert.True(IsFunctionArity(func() {}, 0, 0))
	assert.True(IsFunctionArity(func() {}, 0, -1))
	assert.True(IsFunctionArity(func() {}, -1, 0))
	assert.True(IsFunctionArity(func() {}, -1, -1))
	assert.False(IsFunctionArity(func() {}, 99, 0))
	assert.False(IsFunctionArity(func() {}, 0, 99))
	assert.False(IsFunctionArity(func() {}, 99, 99))

	assert.True(IsFunction(func(any) {}))
	assert.True(IsFunctionArity(func(any) {}, 1, 0))
	assert.True(IsFunctionArity(func(any) {}, 1, -1))
	assert.True(IsFunctionArity(func(any) {}, -1, 0))
	assert.True(IsFunctionArity(func(any) {}, -1, -1))
	assert.False(IsFunctionArity(func(any) {}, 99, 0))
	assert.False(IsFunctionArity(func(any) {}, 0, 99))
	assert.False(IsFunctionArity(func(any) {}, 99, 99))

	assert.True(IsFunction(func(any) error { return nil }))
	assert.True(IsFunctionArity(func(any) error { return nil }, 1, 1))
	assert.True(IsFunctionArity(func(any) error { return nil }, 1, -1))
	assert.True(IsFunctionArity(func(any) error { return nil }, -1, 1))
	assert.True(IsFunctionArity(func(any) error { return nil }, -1, -1))
	assert.False(IsFunctionArity(func(any) error { return nil }, 99, 1))
	assert.False(IsFunctionArity(func(any) error { return nil }, 0, 99))
	assert.False(IsFunctionArity(func(any) error { return nil }, 99, 99))

	assert.True(IsFunction(func() error { return nil }))
	assert.True(IsFunctionArity(func() error { return nil }, 0, 1))
	assert.True(IsFunctionArity(func() error { return nil }, 0, -1))
	assert.True(IsFunctionArity(func() error { return nil }, -1, 1))
	assert.True(IsFunctionArity(func() error { return nil }, -1, -1))
	assert.False(IsFunctionArity(func() error { return nil }, 99, 0))
	assert.False(IsFunctionArity(func() error { return nil }, 0, 99))
	assert.False(IsFunctionArity(func() error { return nil }, 99, 99))
}

func TestSetValueInt(t *testing.T) {
	var assert = require.New(t)

	// INT
	// ---------------------------------------------------------------------------------------------
	var i int
	// int* -> int
	assert.NoError(SetValue(&i, int(42)))
	assert.Equal(int(42), i)

	assert.NoError(SetValue(&i, int8(42)))
	assert.Equal(int(42), i)

	assert.NoError(SetValue(&i, int16(42)))
	assert.Equal(int(42), i)

	assert.NoError(SetValue(&i, int32(42)))
	assert.Equal(int(42), i)

	assert.NoError(SetValue(&i, int64(42)))
	assert.Equal(int(42), i)

	// uint* -> int
	assert.NoError(SetValue(&i, uint(42)))
	assert.Equal(int(42), i)

	assert.NoError(SetValue(&i, uint8(42)))
	assert.Equal(int(42), i)

	assert.NoError(SetValue(&i, uint16(42)))
	assert.Equal(int(42), i)

	assert.NoError(SetValue(&i, uint32(42)))
	assert.Equal(int(42), i)

	assert.NoError(SetValue(&i, uint64(42)))
	assert.Equal(int(42), i)

	// ---------------------------------------------------------------------------------------------
	var i8 int8
	// int* -> int
	assert.NoError(SetValue(&i8, int(42)))
	assert.Equal(int8(42), i8)

	assert.NoError(SetValue(&i8, int8(42)))
	assert.Equal(int8(42), i8)

	assert.NoError(SetValue(&i8, int16(42)))
	assert.Equal(int8(42), i8)

	assert.NoError(SetValue(&i8, int32(42)))
	assert.Equal(int8(42), i8)

	assert.NoError(SetValue(&i8, int64(42)))
	assert.Equal(int8(42), i8)

	// uint* -> int
	assert.NoError(SetValue(&i8, uint(42)))
	assert.Equal(int8(42), i8)

	assert.NoError(SetValue(&i8, uint8(42)))
	assert.Equal(int8(42), i8)

	assert.NoError(SetValue(&i8, uint16(42)))
	assert.Equal(int8(42), i8)

	assert.NoError(SetValue(&i8, uint32(42)))
	assert.Equal(int8(42), i8)

	assert.NoError(SetValue(&i8, uint64(42)))
	assert.Equal(int8(42), i8)

	assert.NotNil(SetValue(time.Time{}, time.Now()))

	// var i16 int16
	// var i32 int32
	// var u uint
	// var u8 uint8
	// var u16 uint16
	// var u32 uint32
	// var u64 uint64
}

type testEnum string

const (
	Value1 testEnum = `value-1`
	Value2 testEnum = `value-2`
	Value3 testEnum = `value-3`
)

type testSettable struct {
	Name          string
	Type          testEnum
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	Other         *string
	CreatedAtNano int64
}

func TestSetValueStruct(t *testing.T) {
	var assert = require.New(t)

	var t1 = &testSettable{
		Name: `t1`,
		Type: Value2,
	}

	// -------------------------------------------------------------------------
	assert.Equal(Value2, t1.Type)

	assert.NoError(SetValue(
		reflect.ValueOf(t1).Elem().Field(0),
		42,
	))

	assert.NoError(SetValue(
		reflect.ValueOf(t1).Elem().Field(1),
		Value3,
	))

	assert.Equal(`42`, t1.Name)
	assert.Equal(Value3, t1.Type)

	// -------------------------------------------------------------------------
	assert.NoError(SetValue(
		reflect.ValueOf(t1).Elem().Field(1),
		`value-4`,
	))

	assert.Equal(testEnum(`value-4`), t1.Type)

	// -------------------------------------------------------------------------
	var tm = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	assert.NoError(SetValue(
		reflect.ValueOf(t1).Elem().Field(2),
		tm,
	))

	assert.True(t1.CreatedAt.Equal(tm))

	// -------------------------------------------------------------------------
	var tmI any
	tmI = tm

	assert.NoError(SetValue(
		reflect.ValueOf(t1).Elem().Field(2),
		tmI,
	))

	assert.True(t1.CreatedAt.Equal(tm))

	// -------------------------------------------------------------------------
	assert.NoError(SetValue(
		reflect.ValueOf(t1).Elem().Field(3),
		&tm,
	))

	assert.True(t1.UpdatedAt.Equal(tm), fmt.Sprintf("%v", t1.UpdatedAt))

	// -------------------------------------------------------------------------
	assert.Error(SetValue(
		reflect.ValueOf(t1).Elem().Field(3),
		tm,
	))

	var st = subtime{}
	assert.NoError(SetValue(&st, tm))
	assert.True(st.Time.Equal(tm))

	// -------------------------------------------------------------------------
	tm = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	assert.NoError(SetValue(
		reflect.ValueOf(t1).Elem().Field(5),
		tm,
	))

	assert.Equal(tm.UnixNano(), t1.CreatedAtNano)
}

func TestSetValueSlice(t *testing.T) {
	var assert = require.New(t)

	var s []string

	assert.NoError(SetValue(&s, []string{`hey`}))
	assert.Equal([]string{`hey`}, s)

	assert.NoError(SetValue(&s, []string{`there`}))
	assert.Equal([]string{`there`}, s)

}
