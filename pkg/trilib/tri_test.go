package trilib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMake(t *testing.T) {
	cases := []struct {
		name string
		list []string
		want []Tri[string]
	}{
		{
			name: "empty",
			list: []string{},
			want: []Tri[string]{},
		},
		{
			name: "one",
			list: []string{"a"},
			want: []Tri[string]{
				{
					Current: "a",
				},
			},
		},
		{
			name: "two",
			list: []string{"a", "b"},
			want: []Tri[string]{
				{
					Current: "a",
					Next:    "b",
				},
				{
					Previous: "a",
					Current:  "b",
				},
			},
		},
		{
			name: "three",
			list: []string{"a", "b", "c"},
			want: []Tri[string]{
				{
					Current: "a",
					Next:    "b",
				},
				{
					Previous: "a",
					Current:  "b",
					Next:     "c",
				},
				{
					Previous: "b",
					Current:  "c",
				},
			},
		},
		{
			name: "four",
			list: []string{"a", "b", "c", "d"},
			want: []Tri[string]{
				{
					Current: "a",
					Next:    "b",
				},
				{
					Previous: "a",
					Current:  "b",
					Next:     "c",
				},
				{
					Previous: "b",
					Current:  "c",
					Next:     "d",
				},
				{
					Previous: "c",
					Current:  "d",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Make[string](c.list)
			assert.Equal(t, c.want, got)
		})
	}
}
