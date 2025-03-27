package diff

import (
	"reflect"
	t "testing"
)

type TestCase struct {
	l1  []string
	l2  []string
	exp []Op[string]
}

func TestDiff(t *t.T) {
	A := "A"
	B := "B"
	C := "C"
	testCases := []TestCase{
		{[]string{}, []string{}, []Op[string]{}},
		{[]string{}, []string{"foo"}, []Op[string]{{OpInsert, 0, 0, "foo"}}},
		{[]string{"foo"}, []string{}, []Op[string]{{OpDelete, 0, -1, "foo"}}},
		{[]string{"foo", "bar", "baz"}, []string{"foo", "bar", "baz"}, []Op[string]{}},
		{[]string{"foo", "bar", "baz"}, []string{"foo", "baz"}, []Op[string]{{OpDelete, 1, -1, "bar"}}},
		{[]string{"baz"}, []string{"foo", "baz"}, []Op[string]{{OpInsert, 0, 0, "foo"}}},
		{[]string{"bar", "baz"}, []string{"foo", "baz"}, []Op[string]{{OpDelete, 0, -1, "bar"}, {OpInsert, 1, 0, "foo"}}},
		{[]string{"foo", "bar", "baz"}, []string{"foo", "bar"}, []Op[string]{{OpDelete, 2, -1, "baz"}}},
		{
			[]string{A, B, C, A, B, B, A},
			[]string{C, B, A, B, A, C},
			[]Op[string]{{OpDelete, 0, -1, A}, {OpInsert, 1, 0, C}, {OpDelete, 2, -1, C}, {OpDelete, 5, -1, B}, {OpInsert, 7, 5, C}},
		},
		{
			[]string{C, A, B, A, B, A, B, A, B, A, B, A, B, C},
			[]string{B, A, B, A, B, A, B, A, B, A, B, A, B, A},
			[]Op[string]{{OpDelete, 0, -1, C}, {OpInsert, 1, 0, B}, {OpDelete, 13, -1, C}, {OpInsert, 14, 13, A}},
		},
		{
			[]string{B},
			[]string{A, B, C, B, A},
			[]Op[string]{{OpInsert, 0, 0, A}, {OpInsert, 0, 1, B}, {OpInsert, 0, 2, C}, {OpInsert, 1, 4, A}},
		},
	}
	for _, c := range testCases {
		act := myersDiffStr(c.l1, c.l2, true)
		if !reflect.DeepEqual(c.exp, act) {
			t.Errorf("Failed diff, expected %v actual %v\n", c.exp, act)
		}
	}
}
