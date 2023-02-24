package command_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tateexon/cli-exec/command"
)

func TestRunACommandWithOut(t *testing.T) {
	t.Parallel()
	stdout, stderr, err := basicCommandHelper(t, "../scripts/basicCommandWithOut", []string{})
	require.NoError(t, err)
	require.ElementsMatch(t, stdout, []string{"Line 1", "Line 2"})
	require.Empty(t, stderr)
}

func TestRunACommandWithErr(t *testing.T) {
	t.Parallel()
	stdout, stderr, err := basicCommandHelper(t, "../scripts/basicCommandWithErr", []string{})
	require.NoError(t, err)
	require.ElementsMatch(t, stderr, []string{"Error Line 1", "Error Line 2"})
	require.Empty(t, stdout)
}

func TestRunACommandWithMixedOutErr(t *testing.T) {
	t.Parallel()
	stdout, stderr, err := basicCommandHelper(t, "../scripts/basicCommandWithMixedOutErr", []string{})
	require.NoError(t, err)
	require.ElementsMatch(t, stdout, []string{"Line 1", "Line 2"})
	require.ElementsMatch(t, stderr, []string{"Error Line 1", "Error Line 2"})
}

func TestPassInArgs(t *testing.T) {
	t.Parallel()
	stdout, stderr, err := basicCommandHelper(t, "../scripts/basicCommandWithArgs", []string{"1", "2"})
	require.NoError(t, err)
	require.ElementsMatch(t, stdout, []string{"Arg 1: 1", "Arg 2: 2"})
	require.Empty(t, stderr)
}

func TestErrorIsReturned(t *testing.T) {
	t.Parallel()
	stdout, stderr, err := basicCommandHelper(t, "../scripts/basicCommandWithError", []string{})
	require.Error(t, err)
	require.Empty(t, stdout)
	require.Empty(t, stderr)
}

func TestHungAppTimeout(t *testing.T) {
	t.Parallel()
	opts := command.CommandOptions{
		WaitForStdPipe: true,
		StdOutHandler: func(message string) {
			t.Log(message)
		},
		StdErrHandler: func(message string) {
			t.Log(message)
		},
	}
	err := command.ExecuteCommand("../scripts/hang", []string{}, opts)
	require.Error(t, err)
}

func basicCommandHelper(t *testing.T, script string, args []string) ([]string, []string, error) {
	stdout := []string{}
	stderr := []string{}
	opts := command.CommandOptions{
		WaitForStdPipe: true,
		StdOutHandler: func(message string) {
			stdout = append(stdout, message)
			t.Log(message)
		},
		StdErrHandler: func(message string) {
			stderr = append(stderr, message)
			t.Log(message)
		},
	}
	err := command.ExecuteCommand(script, args, opts)
	return stdout, stderr, err
}
