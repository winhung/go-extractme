package extractor

import (
	logger "go-extractme/cmd/util/logger"
	"testing"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string)  {}
func (m *mockLogger) Debug(msg string) {}
func (m *mockLogger) Warn(msg string)  {}
func (m *mockLogger) Error(msg string) {}
func (m *mockLogger) Fatal(msg string) {}
func createMockLogger() logger.CustomLogger {
	return logger.CustomLogger(&mockLogger{})
}

func TestCreateExtractor(t *testing.T) {
	mockLogger := createMockLogger()

	type args struct {
		conversionType string
		customLogger   logger.CustomLogger
	}
	tests := []struct {
		name         string
		args         args
		expectToFail bool
	}{
		{
			name: "Invalid conversion type",
			args: args{
				conversionType: "invalid type",
				customLogger:   mockLogger,
			},
			expectToFail: true,
		},
		{
			name: "Valid conversion type",
			args: args{
				conversionType: TF2JSON.String(),
				customLogger:   mockLogger,
			},
			expectToFail: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext := CreateExtractor(tt.args.conversionType, tt.args.customLogger)
			if tt.expectToFail && ext != nil {
				t.Errorf("CreateExtractor() ext = %v, expectToFail %v", ext, tt.expectToFail)
				return
			}
		})
	}
}
