package commonutil

import "testing"

func TestIsMapValid(t *testing.T) {
	validMap := make(map[string]map[string]string, 1)
	validMap["sample"] = make(map[string]string)

	type args struct {
		target map[string]map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test invalid, empty map",
			args: args{
				target: map[string]map[string]string{},
			},
			wantErr: true,
		},
		{
			name: "Test valid map",
			args: args{
				target: validMap,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := IsMapValid(tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("IsMapValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
