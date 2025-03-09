package converter

import (
	"bytes"
	"fmt"

	"gotoproto/pkg/parser"
	v3 "gotoproto/pkg/printer/v3"
)

func Convert(input string) string {
	res, err := parser.Parse(input)
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err).Error()
	}

	if len(res) == 0 {
		return "invalid input"
	}

	var out bytes.Buffer

	err = v3.NewOutputPrinter(&out).WriteResult(res)
	if err != nil {
		return fmt.Errorf("failed to write: %w", err).Error()
	}

	return out.String()
}
