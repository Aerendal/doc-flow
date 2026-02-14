package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestExitCodeMapping(t *testing.T) {
	plainErr := errors.New("boom")
	usageErr := usageErrorf("bad args")
	runtimeErr := runtimeError("DOCFLOW.RUNTIME.TEST", "runtime failed", plainErr, nil)

	cases := []struct {
		name string
		err  error
		want int
	}{
		{name: "nil", err: nil, want: 0},
		{name: "domain", err: plainErr, want: 1},
		{name: "usage", err: usageErr, want: 2},
		{name: "runtime", err: runtimeErr, want: 3},
		{name: "wrapped usage", err: fmt.Errorf("wrapped: %w", usageErr), want: 2},
		{name: "required flag parse", err: errors.New(`required flag(s) "rules" not set`), want: 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ExitCode(tc.err); got != tc.want {
				t.Fatalf("ExitCode() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestComplianceUsageErrorEmitsJSONEnvelope(t *testing.T) {
	cmd := complianceCmd()
	cmd.SetArgs([]string{"--format", "json", "--output", "-", "--strict"})
	stdout, err := runCommandCaptureStdout(cmd)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := ExitCode(err); got != 2 {
		t.Fatalf("ExitCode() = %d, want 2; err=%v", got, err)
	}

	var envelope struct {
		SchemaVersion string `json:"schema_version"`
		Error         struct {
			Kind string `json:"kind"`
			Code string `json:"code"`
		} `json:"error"`
	}
	if jerr := json.Unmarshal(stdout, &envelope); jerr != nil {
		t.Fatalf("invalid json envelope: %v\nstdout=%s", jerr, string(stdout))
	}
	if envelope.SchemaVersion != "1.0" {
		t.Fatalf("unexpected schema_version: %s", envelope.SchemaVersion)
	}
	if envelope.Error.Kind != "usage" {
		t.Fatalf("unexpected error kind: %s", envelope.Error.Kind)
	}
	if envelope.Error.Code != "DOCFLOW.CLI.MISSING_RULES" {
		t.Fatalf("unexpected error code: %s", envelope.Error.Code)
	}
}

func TestComplianceRuntimeErrorEmitsJSONEnvelope(t *testing.T) {
	cmd := complianceCmd()
	cmd.SetArgs([]string{"--format", "json", "--output", "-", "--rules", "/tmp/definitely/missing-rules.yaml"})
	stdout, err := runCommandCaptureStdout(cmd)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := ExitCode(err); got != 3 {
		t.Fatalf("ExitCode() = %d, want 3; err=%v", got, err)
	}

	var envelope struct {
		Error struct {
			Kind string `json:"kind"`
			Code string `json:"code"`
		} `json:"error"`
	}
	if jerr := json.Unmarshal(stdout, &envelope); jerr != nil {
		t.Fatalf("invalid json envelope: %v\nstdout=%s", jerr, string(stdout))
	}
	if envelope.Error.Kind != "runtime" {
		t.Fatalf("unexpected error kind: %s", envelope.Error.Kind)
	}
	if envelope.Error.Code != "DOCFLOW.RUNTIME.RULES_CHECKSUM_FAILED" {
		t.Fatalf("unexpected error code: %s", envelope.Error.Code)
	}
}

func TestValidateUsageErrorEmitsJSONEnvelope(t *testing.T) {
	cmd := validateCmd()
	cmd.SetArgs([]string{"--format", "json", "--output", "-", "--strict", "--warn"})
	stdout, err := runCommandCaptureStdout(cmd)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := ExitCode(err); got != 2 {
		t.Fatalf("ExitCode() = %d, want 2; err=%v", got, err)
	}

	var envelope struct {
		Error struct {
			Kind string `json:"kind"`
			Code string `json:"code"`
		} `json:"error"`
	}
	if jerr := json.Unmarshal(stdout, &envelope); jerr != nil {
		t.Fatalf("invalid json envelope: %v\nstdout=%s", jerr, string(stdout))
	}
	if envelope.Error.Kind != "usage" {
		t.Fatalf("unexpected error kind: %s", envelope.Error.Kind)
	}
	if envelope.Error.Code != "DOCFLOW.CLI.CONFLICTING_FLAGS" {
		t.Fatalf("unexpected error code: %s", envelope.Error.Code)
	}
}

func TestComplianceUsageErrorsAreExit2(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{
			name: "missing rules",
			args: []string{"--strict"},
		},
		{
			name: "strict and warn conflict",
			args: []string{"--strict", "--warn", "--rules", "docs/_meta/GOVERNANCE_RULES.yaml"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := complianceCmd()
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if got := ExitCode(err); got != 2 {
				t.Fatalf("ExitCode() = %d, want 2; err=%v", got, err)
			}
		})
	}
}

func runCommandCaptureStdout(cmd interface{ Execute() error }) ([]byte, error) {
	stdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	err := cmd.Execute()

	_ = wOut.Close()
	os.Stdout = stdout

	out, _ := io.ReadAll(rOut)
	_ = rOut.Close()
	return out, err
}

func TestValidateUsageErrorsAreExit2(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{
			name: "strict and warn conflict",
			args: []string{"--strict", "--warn"},
		},
		{
			name: "unsupported format",
			args: []string{"--format", "yaml"},
		},
		{
			name: "fail-on new requires against",
			args: []string{"--format", "json", "--output", "-", "--fail-on", "new"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := validateCmd()
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if got := ExitCode(err); got != 2 {
				t.Fatalf("ExitCode() = %d, want 2; err=%v", got, err)
			}
		})
	}
}

func TestComplianceFailOnNewRequiresAgainst(t *testing.T) {
	cmd := complianceCmd()
	cmd.SetArgs([]string{"--format", "json", "--output", "-", "--rules", "docs/_meta/GOVERNANCE_RULES.yaml", "--fail-on", "new"})

	stdout, err := runCommandCaptureStdout(cmd)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := ExitCode(err); got != 2 {
		t.Fatalf("ExitCode() = %d, want 2; err=%v", got, err)
	}

	var envelope struct {
		Error struct {
			Kind string `json:"kind"`
			Code string `json:"code"`
		} `json:"error"`
	}
	if jerr := json.Unmarshal(stdout, &envelope); jerr != nil {
		t.Fatalf("invalid json envelope: %v\nstdout=%s", jerr, string(stdout))
	}
	if envelope.Error.Kind != "usage" || envelope.Error.Code != "DOCFLOW.CLI.MISSING_AGAINST" {
		t.Fatalf("unexpected error envelope: %+v", envelope.Error)
	}
}
