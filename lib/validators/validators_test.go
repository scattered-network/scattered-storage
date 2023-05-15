package validators

import (
	"reflect"
	"regexp"
	"testing"
)

// TestValidateInput tests the ValidateInput function.
func TestValidateInput(t *testing.T) {
	type args struct {
		check *regexp.Regexp
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TestValidateInput",
			args: args{
				check: regexp.MustCompile(`^[a-zA-Z0-9_-]{1,128}$`),
				input: "test",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := ValidateInput(tt.args.check, tt.args.input); got != tt.want {
					t.Errorf("ValidateInput() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

// TestValidateRegex tests the ValidateRegex function.
func TestValidateRegex(t *testing.T) {
	type args struct {
		expression string
	}
	tests := []struct {
		name string
		args args
		want *regexp.Regexp
	}{
		{
			name: "TestRBDNameRegex",
			args: args{
				expression: "^[a-zA-Z0-9-_.]+$",
			},
			want: regexp.MustCompile(`^[a-zA-Z0-9-_.]+$`),
		},
		{
			name: "TestPoolNameRegex",
			args: args{
				expression: "^[a-zA-Z0-9-_]+$",
			},
			want: regexp.MustCompile(`^[a-zA-Z0-9-_]+$`),
		},
		{
			name: "TestRBDSizeRegex",
			args: args{
				expression: "^[0-9]+",
			},
			want: regexp.MustCompile(`^[0-9]+`),
		},
		{
			name: "TestSizeSuffixRegex",
			args: args{
				expression: "^[bBkKmMgGtTpP]$",
			},
			want: regexp.MustCompile(`^[bBkKmMgGtTpP]$`),
		},
		{
			name: "TestValidateDevicePathRegex",
			args: args{
				expression: "^/dev/[[:alnum:]]+$",
			},
			want: regexp.MustCompile(`^/dev/[[:alnum:]]+$`),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := ValidateRegex(tt.args.expression); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ValidateRegex() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
