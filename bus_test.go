package command_bus

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testStruct struct {
	bazz string
	foo  int
	List []string
	Map  map[string]int
}

func TestReflect(t *testing.T) {
	a := testStruct{}
	assert.Equal(t, "bus.testStruct", reflect.TypeOf(a).String())
	assert.Equal(t, 4, reflect.Indirect(reflect.ValueOf(a)).NumField())
	for i := 0; i < reflect.Indirect(reflect.ValueOf(a)).NumField(); i++ {
		assert.Equal(t, []string{
			"bazz",
			"foo",
			"List",
			"Map",
		}[i], reflect.Indirect(reflect.ValueOf(a)).Type().Field(i).Name)
	}
}
func BenchmarkCommandBus_Command(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflect.TypeOf(testStruct{})
	}
}

var caller commandName

type s struct {
}

func (a s) foo() {
	baz()
}
func baz() {
	caller = myCaller()
}

func TestCaller(t *testing.T) {
	s{}.foo()
	assert.Equal(t, "chaturbate/src/bus.s.foo", string(caller))
}
