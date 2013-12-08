package interactive

import (
	"testing"
	"reflect"
)

func TestCompositeArrayEmpty(t *testing.T) {
	type Alice [0]int

	env := makeEnv()
	env.Types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { }
	expr := "Alice {}"

	expectResult(t, expr, env, expected)
}

func TestCompositeArrayValues(t *testing.T) {
	type Alice [3]int

	env := makeEnv()
	env.Types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { 1, 2, 3 }
	expr := "Alice { 1, 2, 3 }"

	expectResult(t, expr, env, expected)
}

func TestCompositeArrayKeyValues(t *testing.T) {
	type Alice [3]int

	env := makeEnv()
	env.Types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { 1: 1, 2 }
	expr := "Alice { 1: 1, 2 }"

	expectResult(t, expr, env, expected)
}

func TestCompositeArrayIncompleteValues(t *testing.T) {
	type Alice [3]int

	env := makeEnv()
	env.Types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { 1, 2 }
	expr := "Alice { 1, 2 }"

	expectResult(t, expr, env, expected)
}

func TestCompositeSliceEmpty(t *testing.T) {
	type Alice []int

	env := makeEnv()
	env.Types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { }
	expr := "Alice { }"

	expectResult(t, expr, env, expected)
}

func TestCompositeSliceValues(t *testing.T) {
	type Alice []int

	env := makeEnv()
	env.Types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { 1, 2, 3 }
	expr := "Alice { 1, 2, 3 }"

	expectResult(t, expr, env, expected)
}

func TestCompositeSliceKeyValues(t *testing.T) {
	type Alice []int

	env := makeEnv()
	env.Types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { 1, 10: 1 }
	expr := "Alice { 1, 10: 1 }"

	expectResult(t, expr, env, expected)
}

func TestCompositeAnonArray(t *testing.T) {
	env := makeEnv()

	expected := [3]int { 1, 2 }
	expr := "[3]int { 1, 2 }"

	expectResult(t, expr, env, expected)
}

func TestCompositeAnonAnonArray(t *testing.T) {
	env := makeEnv()

	print("\n\n")
	expected := [3][3]int { [3]int { 1, 2 }, [3]int { 3, 4 } }
	expr := "[3][3]int { [3]int { 1, 2 }, [3]int { 3, 4 } }"

	expectResult(t, expr, env, expected)
}

func TestCompositeStructValues(t *testing.T) {
	type Alice struct {
		Bob int
	}

	env := makeEnv()
	env.Types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { 10 }
	expr := "Alice{ 10 }"

	expectResult(t, expr, env, expected)
}

func TestCompositeStructKeyValues(t *testing.T) {
	type Alice struct {
		Bob int
	}

	env := makeEnv()
	env.Types["Alice"] = reflect.TypeOf(Alice{})

	expected := Alice { Bob: 10 }
	expr := "Alice{ Bob: 10 }"

	expectResult(t, expr, env, expected)
}

