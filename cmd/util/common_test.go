package commonutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsFileExist(t *testing.T) {
	currDir, err := os.Getwd()
	if err != nil {
		t.Fatal("Unable to get current working directory")
	}
	nonExistentFileFullpath := filepath.Join(currDir, "test", "nothing.txt")
	folderFullpath := filepath.Join(currDir, "test", "sample")
	validFile := filepath.Join(currDir, "test", "exist.txt")

	type args struct {
		fullPath string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test with non-existent file",
			args: args{
				fullPath: nonExistentFileFullpath,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Test with directory",
			args: args{
				fullPath: folderFullpath,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Test with existent file",
			args: args{
				fullPath: validFile,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsFileExist(tt.args.fullPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsFileExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsFileExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
