package crawler_test

import (
	"gawr/crawler"
	"testing"
)

func TestQueue_String(t *testing.T) {
	input := "foo"
	expected := "foo"

	q := crawler.Queue[string]{}
	q.Push(input)
	result, err := q.Pop()
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestQueue_IsEmpty(t *testing.T) {
	tbl := []struct {
		name     string
		input    []int
		expected bool
	}{
		{
			"should be true",
			nil,
			true,
		},
		{
			"should not be true",
			[]int{1, 2, 3, 4, 5},
			false,
		},
	}

	for _, tc := range tbl {
		t.Run(tc.name, func(t *testing.T) {
			q := crawler.Queue[int]{}
			for _, e := range tc.input {
				q.Push(e)
			}

			res := q.IsEmpty()
			if res != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, res)
			}
		})
	}
}
