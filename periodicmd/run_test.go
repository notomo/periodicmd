package periodicmd_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/notomo/periodicmd/periodicmd"
)

func TestRun(t *testing.T) {
	t.Run("execute create and link command", func(t *testing.T) {
		var stdoutWriter bytes.Buffer
		var stderrWriter bytes.Buffer

		if err := periodicmd.Run(
			context.Background(),
			[]periodicmd.Task{
				{
					Frequency: periodicmd.TaskFrequency{
						Weeks: 1,
					},
					StartDate:     "2023-12-04",
					CreateCommand: "echo create-{{.date}}",
					LinkCommand:   "echo link {{.previous.date}}:{{.previous.output}} {{.current.date}}:{{.current.output}}",
				},
			},
			"2023-12-04",
			7,
			false,
			&stdoutWriter,
			&stderrWriter,
		); err != nil {
			t.Fatal(err)
		}

		want := `echo create-2023-12-04
create-2023-12-04
echo create-2023-12-11
create-2023-12-11
echo link : 2023-12-04:create-2023-12-04
link : 2023-12-04:create-2023-12-04
echo link 2023-12-04:create-2023-12-04 2023-12-11:create-2023-12-11
link 2023-12-04:create-2023-12-04 2023-12-11:create-2023-12-11
`
		if got := stdoutWriter.String(); got != want {
			t.Errorf("want %v, but %v:", want, got)
		}
	})

	t.Run("execute with task offset days", func(t *testing.T) {
		var stdoutWriter bytes.Buffer
		var stderrWriter bytes.Buffer

		if err := periodicmd.Run(
			context.Background(),
			[]periodicmd.Task{
				{
					Frequency: periodicmd.TaskFrequency{
						Weeks: 1,
					},
					OffsetDays:    ptr(1),
					StartDate:     "2023-12-04",
					CreateCommand: "echo create-{{.date}}",
				},
			},
			"2023-12-04",
			7,
			false,
			&stdoutWriter,
			&stderrWriter,
		); err != nil {
			t.Fatal(err)
		}

		want := `echo create-2023-12-04
create-2023-12-04
`
		if got := stdoutWriter.String(); got != want {
			t.Errorf("want %v, but %v:", want, got)
		}
	})
}

func ptr[T any](x T) *T {
	return &x
}
