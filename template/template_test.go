package template

import (
	"errors"
	"sort"
	"testing"
)

// MockSSMClient is a mock implementation of SSMParameterInterface for testing
type MockSSMClient struct {
	parameters map[string]string
	shouldFail bool
}

func NewMockSSMClient() *MockSSMClient {
	return &MockSSMClient{
		parameters: map[string]string{
			"/myapp/db_host":     "db.example.com",
			"/myapp/db_name":     "myapp_db",
			"/myapp/db_user":     "dbuser",
			"/myapp/db_password": "secret-password",
			"/myapp/api_key":     "api-key-12345",
		},
	}
}

func (m *MockSSMClient) GetParameter(name string) (string, error) {
	if m.shouldFail {
		return "", errors.New("mock SSM error")
	}

	value, ok := m.parameters[name]
	if !ok {
		return "", errors.New("parameter not found")
	}
	return value, nil
}

func TestProcess(t *testing.T) {
	mockClient := NewMockSSMClient()
	processor := NewProcessor(mockClient)

	tests := []struct {
		name     string
		template string
		expected string
		wantErr  bool
	}{
		{
			name:     "No placeholders",
			template: "KEY=value\nANOTHER=value2",
			expected: "KEY=value\nANOTHER=value2",
			wantErr:  false,
		},
		{
			name:     "With placeholders",
			template: "DB_HOST={{SSM:/myapp/db_host}}\nDB_PASSWORD={{SSM:/myapp/db_password}}",
			expected: "DB_HOST=db.example.com\nDB_PASSWORD=secret-password",
			wantErr:  false,
		},
		{
			name:     "Mixed content",
			template: "APP_ENV=dev\nDB_HOST={{SSM:/myapp/db_host}}\nPORT=3000",
			expected: "APP_ENV=dev\nDB_HOST=db.example.com\nPORT=3000",
			wantErr:  false,
		},
		{
			name:     "With comments",
			template: "# This is a comment\nAPP_ENV=dev\n# Another comment\nDB_HOST={{SSM:/myapp/db_host}}",
			expected: "# This is a comment\nAPP_ENV=dev\n# Another comment\nDB_HOST=db.example.com",
			wantErr:  false,
		},
		{
			name:     "With empty lines",
			template: "APP_ENV=dev\n\nDB_HOST={{SSM:/myapp/db_host}}\n\nPORT=3000",
			expected: "APP_ENV=dev\n\nDB_HOST=db.example.com\n\nPORT=3000",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.Process(tt.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("Process() = %v, want %v", result, tt.expected)
			}
		})
	}

	// Test error case
	t.Run("SSM error", func(t *testing.T) {
		mockClient.shouldFail = true
		_, err := processor.Process("DB_HOST={{SSM:/myapp/db_host}}")
		if err == nil {
			t.Errorf("Process() expected error, got nil")
		}
		mockClient.shouldFail = false
	})

	// Test parameter not found
	t.Run("Parameter not found", func(t *testing.T) {
		_, err := processor.Process("KEY={{SSM:/nonexistent/key}}")
		if err == nil {
			t.Errorf("Process() expected error, got nil")
		}
	})
}

func TestExtractParameters(t *testing.T) {
	tests := []struct {
		name     string
		template string
		expected []string
	}{
		{
			name:     "No placeholders",
			template: "KEY=value\nANOTHER=value2",
			expected: []string{},
		},
		{
			name:     "With placeholders",
			template: "DB_HOST={{SSM:/myapp/db_host}}\nDB_PASSWORD={{SSM:/myapp/db_password}}",
			expected: []string{"/myapp/db_host", "/myapp/db_password"},
		},
		{
			name:     "With duplicate placeholders",
			template: "DB_HOST={{SSM:/myapp/db_host}}\nDB_HOST_REPLICA={{SSM:/myapp/db_host}}",
			expected: []string{"/myapp/db_host"},
		},
		{
			name:     "Mixed content",
			template: "APP_ENV=dev\nDB_HOST={{SSM:/myapp/db_host}}\nAPI_KEY={{SSM:/myapp/api_key}}",
			expected: []string{"/myapp/db_host", "/myapp/api_key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractParameters(tt.template)
			if err != nil {
				t.Errorf("ExtractParameters() error = %v", err)
				return
			}
			
			// Sort both slices for comparison
			sort.Strings(result)
			sort.Strings(tt.expected)
			
			if len(result) != len(tt.expected) {
				t.Errorf("ExtractParameters() returned %d parameters, want %d", len(result), len(tt.expected))
				return
			}
			
			for i, param := range result {
				if param != tt.expected[i] {
					t.Errorf("ExtractParameters()[%d] = %v, want %v", i, param, tt.expected[i])
				}
			}
		})
	}
}
