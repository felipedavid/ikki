package utils

import "testing"

func TestEnqueue(t *testing.T) {
	q := Queue{}

	for i := range(4096) {
		q.Enqueue(i)
	}

	for i := range(4096) {
		val := q.Dequeu()
		if val != i {
			t.Errorf("unexpected value %v != %v", val, i)
		}
	}
}