package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedPort string
		expectError  bool
	}{
		{
			name:         "Default port when no args",
			args:         []string{"./TCPChat"},
			expectedPort: "8989",
			expectError:  false,
		},
		{
			name:         "Valid custom port",
			args:         []string{"./TCPChat", "2525"},
			expectedPort: "2525",
			expectError:  false,
		},
		{
			name:         "Too many arguments",
			args:         []string{"./TCPChat", "8989", "extra"},
			expectedPort: "",
			expectError:  true,
		},
		{
			name:         "Non-numeric port",
			args:         []string{"./TCPChat", "abc"},
			expectedPort: "",
			expectError:  true,
		},
		{
			name:         "Port out of lower bound",
			args:         []string{"./TCPChat", "0"},
			expectedPort: "",
			expectError:  true,
		},
		{
			name:         "Port out of upper bound",
			args:         []string{"./TCPChat", "65536"},
			expectedPort: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = tt.args

			cfg, err := Load()
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}

			if !tt.expectError && cfg.Port != tt.expectedPort {
				t.Errorf("expected port %s, got %s", tt.expectedPort, cfg.Port)
			}
		})
	}
}