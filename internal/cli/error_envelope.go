package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type cliErrorEnvelope struct {
	SchemaVersion string                 `json:"schema_version"`
	Error         cliErrorEnvelopeDetail `json:"error"`
}

type cliErrorEnvelopeDetail struct {
	Kind    string         `json:"kind"`
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func commandUsageError(format, output, code, message string, details map[string]any) error {
	err := usageError(code, message, details)
	if format == "json" {
		if werr := writeCLIJSONError(output, err.(*CLIError)); werr != nil {
			return runtimeError(
				"DOCFLOW.RUNTIME.ERROR_ENVELOPE_WRITE_FAILED",
				fmt.Sprintf("nie można zapisać JSON error envelope: %v", werr),
				werr,
				nil,
			)
		}
	}
	return err
}

func commandRuntimeError(format, output, code, message string, cause error, details map[string]any) error {
	err := runtimeError(code, message, cause, details)
	if format == "json" {
		if werr := writeCLIJSONError(output, err.(*CLIError)); werr != nil {
			return runtimeError(
				"DOCFLOW.RUNTIME.ERROR_ENVELOPE_WRITE_FAILED",
				fmt.Sprintf("nie można zapisać JSON error envelope: %v", werr),
				werr,
				nil,
			)
		}
	}
	return err
}

func writeCLIJSONError(output string, cliErr *CLIError) error {
	envelope := cliErrorEnvelope{
		SchemaVersion: "1.0",
		Error: cliErrorEnvelopeDetail{
			Kind:    string(cliErr.Kind),
			Code:    cliErr.Code,
			Message: cliErr.Message,
			Details: cliErr.Details,
		},
	}

	data, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return fmt.Errorf("błąd serializacji error envelope: %w", err)
	}
	if output == "" || output == "-" {
		_, err = os.Stdout.Write(append(data, '\n'))
		return err
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("błąd utworzenia katalogu dla error envelope: %w", err)
	}
	return os.WriteFile(output, data, 0o644)
}
