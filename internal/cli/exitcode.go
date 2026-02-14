package cli

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// exitCodeDomain marks domain failures (e.g. strict validation violations).
	exitCodeDomain = 1
	// exitCodeUsage marks CLI usage/flag validation failures.
	exitCodeUsage = 2
	// exitCodeRuntime marks runtime/IO/configuration failures.
	exitCodeRuntime = 3
)

type CLIErrorKind string

const (
	CLIErrorKindUsage   CLIErrorKind = "usage"
	CLIErrorKindRuntime CLIErrorKind = "runtime"
)

type CLIError struct {
	Kind    CLIErrorKind
	Code    string
	Message string
	Details map[string]any

	exitCode int
	cause    error
}

func (e *CLIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return ""
}

func (e *CLIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *CLIError) ExitCode() int {
	if e == nil {
		return exitCodeDomain
	}
	return e.exitCode
}

func usageError(code, message string, details map[string]any) error {
	return &CLIError{
		Kind:     CLIErrorKindUsage,
		Code:     code,
		Message:  message,
		Details:  details,
		exitCode: exitCodeUsage,
	}
}

func usageErrorf(format string, args ...any) error {
	return usageError("DOCFLOW.CLI.USAGE", fmt.Sprintf(format, args...), nil)
}

func runtimeError(code, message string, cause error, details map[string]any) error {
	return &CLIError{
		Kind:     CLIErrorKindRuntime,
		Code:     code,
		Message:  message,
		Details:  details,
		exitCode: exitCodeRuntime,
		cause:    cause,
	}
}

// ExitCode maps command execution errors to process exit codes.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	var coder interface {
		ExitCode() int
	}
	if errors.As(err, &coder) {
		code := coder.ExitCode()
		if code > 0 {
			return code
		}
	}

	// Cobra/pflag argument parsing errors are usage failures.
	msg := err.Error()
	if strings.HasPrefix(msg, "required flag(s)") ||
		strings.HasPrefix(msg, "unknown flag:") ||
		strings.HasPrefix(msg, "unknown shorthand flag:") ||
		strings.HasPrefix(msg, "flag needs an argument:") ||
		strings.HasPrefix(msg, "invalid argument ") {
		return exitCodeUsage
	}

	return exitCodeDomain
}
