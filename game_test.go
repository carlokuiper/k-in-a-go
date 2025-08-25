package kinago

import (
	"reflect"
	"testing"
)

func TestFinished(t *testing.T) {
	tests := []struct {
		name     string
		board    [][]int
		k        int
		expected []Coordinate
	}{{
		name:     "column unfinished",
		board:    [][]int{{0, 0, 1}, {0, 0, 1}, {0, 0, 0}},
		k:        3,
		expected: nil,
	}, {
		name:     "row finished",
		board:    [][]int{{1, 1, 1}, {0, 0, 0}, {0, 0, 0}},
		k:        3,
		expected: []Coordinate{{0, 0}, {1, 0}, {2, 0}},
	}, {
		name:     "anti diagonal finished",
		board:    [][]int{{0, 1}, {1, 0}},
		k:        2,
		expected: []Coordinate{{0, 1}, {1, 0}},
	}, {
		name:     "off anti diagonal",
		board:    [][]int{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 1}, {0, 0, 1, 0}},
		k:        2,
		expected: []Coordinate{{2, 3}, {3, 2}},
	}, {
		name:     "m large finished",
		board:    [][]int{{0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 1, 0}, {0, 0, 0, 1, 0, 0}},
		k:        3,
		expected: []Coordinate{{3, 2}, {4, 1}, {5, 0}},
	}, {
		name:     "m large unfinished",
		board:    [][]int{{0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 1, 0}, {0, 0, 0, 1, 0, 0}},
		k:        4,
		expected: nil,
	}, {
		name:     "n large finished",
		board:    [][]int{{0, 0}, {0, 0}, {1, 0}, {0, 1}},
		k:        2,
		expected: []Coordinate{{0, 2}, {1, 3}},
	}, {
		name:     "n large unfinished",
		board:    [][]int{{0, 0}, {0, 0}, {1, 0}, {0, 1}},
		k:        3,
		expected: nil,
	}}
	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := finished(test.board, test.k)
			if !reflect.DeepEqual(got, test.expected) {
				t.Log(got)
				t.Log(test.expected)
				t.Fail()
			}
		})
	}
}
