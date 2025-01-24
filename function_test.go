package main

import (
	"io"
	"os"
	"testing"

	fuzz "github.com/google/gofuzz"
)

func TestRunSimple(t *testing.T) {
	cases := []struct {
		command string
		args    []string
		payload []byte
		stdout  string
		stderr  string
	}{
		{"echo", []string{"test"}, []byte("{}"), "test\n", ""},
		{"sh", []string{"-c", "echo test >&2"}, []byte("{}"), "", "test\n"},
	}

	for i, c := range cases {
		fn := NewFunction(c.command, c.args, c.payload)

		rstdout, wstdout, _ := os.Pipe()
		rstderr, wstderr, _ := os.Pipe()

		fn.SetStdout(wstdout)
		fn.SetStderr(wstderr)

		stdout, stderr, err := fn.Run()

		wstdout.Close()
		wstderr.Close()

		if err != nil {
			t.Errorf("simple case(%d) failed): %s", i, err)
		}

		if stdout != c.stdout {
			t.Errorf("simple case(%d) stdout want: %s, got: %s", i, c.stdout, stdout)
		}

		rso, _ := io.ReadAll(rstdout)
		if string(rso) != c.stdout {
			t.Errorf("simple case(%d) rstdout want: %s, got: %s", i, c.stdout, rso)
		}

		if stderr != c.stderr {
			t.Errorf("simple case(%d) stderr want: %s, got: %s", i, c.stderr, stderr)
		}

		rse, _ := io.ReadAll(rstderr)
		if string(rse) != c.stderr {
			t.Errorf("simple case(%d) rstderr want: %s, got: %s", i, c.stderr, rse)
		}
	}
}

func TestRunFuzzing(t *testing.T) {
	for i := 1; i <= 10; i++ {
		f := fuzz.New()
		var str string
		f.Fuzz(&str)

		fn := NewFunction("echo", []string{"-n", str}, []byte("{}"))

		fn.SetStdout(io.Discard)
		fn.SetStderr(io.Discard)

		stdout, _, _ := fn.Run()

		if stdout != str {
			t.Errorf("fuzzing stdout want: %s, got: %s", str, stdout)
		}
	}
}
