package cmd

import (
	"reflect"
	"testing"
)

func TestStripDashDash(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"empty", nil, nil},
		{"no dash", []string{"--debug"}, []string{"--debug"}},
		{"leading dash removed", []string{"--", "-v"}, []string{"-v"}},
		{"only dash removed", []string{"--"}, nil},
		{"dash not leading kept", []string{"-v", "--"}, []string{"-v", "--"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripDashDash(tt.args)
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stripDashDash(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
