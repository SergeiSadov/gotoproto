package parser

import (
	"reflect"
	"testing"

	"gotoproto/pkg/models"
)

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []models.StructInfo
		wantErr    bool
	}{
		{
			name: "basic",
			args: args{
				s: `
type User struct{
	ID int
	PhoneNumbers []string
	Contacts map[string]string
}`,
			},
			wantResult: []models.StructInfo{
				{
					Name: "User",
					Fields: []models.Field{
						{
							Name: "ID",
							Type: models.Type{
								Name: "int64",
							},
						},
						{
							Name: "PhoneNumbers",
							Type: models.Type{
								Name: "string[]",
							},
						},
						{
							Name: "Contacts",
							Type: models.Type{
								Name: "",
								MapType: &models.MapType{
									KeyType:   "string",
									ValueType: "string",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should return error",
			args: args{
				s: `input`,
			},
			wantResult: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := Parse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Parse() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
