package printer

import "gotoproto/pkg/models"

type ProtoPrinter interface {
	WriteResult(allStructs []models.StructInfo) (err error)
}
