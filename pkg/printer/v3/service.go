package v3

import (
	"fmt"
	"io"
	"strconv"

	"gotoproto/pkg/models"
)

const (
	newLine = "\n"
	tab     = "	"
)

type ProtoV3Writer struct {
	w io.Writer
}

func NewOutputPrinter(w io.Writer) *ProtoV3Writer {
	return &ProtoV3Writer{w: w}
}

func (p *ProtoV3Writer) WriteResult(result []models.StructInfo) (err error) {
	if err = p.writeString(`syntax = "proto3";`); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, str := range result {
		if err = p.writeString(newLine + newLine); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}

		if err = p.writeString("message "); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}

		if err = p.writeString(str.Name); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}

		if err = p.writeString(" {\n"); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}

		for i, field := range str.Fields {
			if err = p.writeString(tab); err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}

			if field.Type.CustomType != nil {
				if err = p.writeString("reserved " + strconv.Itoa(i+1) + "; //" + *field.Type.CustomType + " \n"); err != nil {
					return fmt.Errorf("failed to write: %w", err)
				}
				continue
			}

			if field.Type.MapType != nil {
				if err = p.writeString(fmt.Sprintf("map<%s, %s> %s = %d;\n", field.Type.MapType.KeyType, field.Type.MapType.ValueType, field.Name, i+1)); err != nil {
					return fmt.Errorf("failed to write: %w", err)
				}
				continue
			}

			if err = p.writeString(field.Type.Name); err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}

			if err = p.writeString(tab); err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}

			if err = p.writeString(field.Name); err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}

			if err = p.writeString(" = "); err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}

			if err = p.writeString(strconv.Itoa(i + 1)); err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}

			if err = p.writeString(";" + newLine); err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}
		}

	}

	if err = p.writeString("}"); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil

}

func (p *ProtoV3Writer) writeString(value string) (err error) {
	_, err = p.w.Write([]byte(value))
	if err != nil {
		return fmt.Errorf("error writing to output: %w", err)
	}

	return nil
}
