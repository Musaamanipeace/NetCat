// Package config_test tests our configuration loading logic.
//
// These tests make sure our config.Load() function behaves correctly
// for all the different ways a user might try to start the server.
package config

// "os" - Changes the command-line arguments for testing
// "testing" - Provides testing framework (t.Run, t.Errorf, t.Fatalf)
import (
	"os"      // For temporarily changing command-line arguments during tests
	"testing" // For running tests and checking if things work correctly
)

// TestLoad runs through various scenarios of how users might invoke our program.
// Each scenario (called a "test case") checks if Load() returns the right configuration
// or returns an error when appropriate.
//
// Table-driven tests are a Go idiom where we list all test cases in a table
// and loop through them. It's cleaner than writing separate test functions.
func TestLoad(t *testing.T) {
	// Define our test cases as a list of scenarios.
	// Each struct represents a different way someone might run the program.
	// The underscore _ in tests := []struct{...} creates an anonymous struct type.
	tests := []struct {
		name        string // Describes what this test case is testing
		args        []string // The fake command-line arguments we'll use
		expectedPort string // What we expect the port to be after Load()
		expectError  bool   // Should Load() return an error?
	}{
		// Test case 1: No arguments at all
		// When user types "./net-cat" with NO extra arguments, we use port 8989.
		{
			name:         "Default port when no args", // Describes this scenario
			args:         []string{"./TCPChat"},      // Fake command line with just program name
			expectedPort: "8989",                    // We expect the default port
			expectError:  false,                     // No error should happen
		},
		// Test case 2: Valid custom port
		// When user types "./net-cat 2525", we use port 2525.
		{
			name:         "Valid custom port",
			args:         []string{"./TCPChat", "2525"}, // Program + custom port
			expectedPort: "2525",
			expectError:  false,
		},
		// Test case 3: Too many arguments
		// When user types "./net-cat 8989 extra", that's wrong!
		{
			name:         "Too many arguments",
			args:         []string{"./TCPChat", "8989", "extra"}, // 3 args = error
			expectedPort: "",  // Port doesn't matter when there's an error
			expectError:  true, // We expect an error
		},
		// Test case 4: Non-numeric port
		// When user types "./net-cat abc", "abc" isn't a number!
		{
			name:         "Non-numeric port",
			args:         []string{"./TCPChat", "abc"}, // "abc" can't be a port
			expectedPort: "",
			expectError:  true,
		},
		// Test case 5: Port below valid range
		// Port 0 is invalid - port numbers start from 1.
		{
			name:         "Port out of lower bound",
			args:         []string{"./TCPChat", "0"}, // Invalid port
			expectedPort: "",
			expectError:  true,
		},
		// Test case 6: Port above valid range
		// Port 65536 is too big - maximum is 65535.
		{
			name:         "Port out of upper bound",
			args:         []string{"./TCPChat", "65536"}, // Too large
			expectedPort: "",
			expectError:  true,
		},
	}

	// Loop through each test case and run it.
	// The "for range" construct iterates through the slice.
	for _, tt := range tests {
		// t.Run creates a subtest with a name, so when tests fail we know which one.
		t.Run(tt.name, func(t *testing.T) {
			// Save the original os.Args so we can restore it later.
			// This is important because changing os.Args affects ALL tests.
			oldArgs := os.Args

			// defer schedules this function call to run AFTER the test finishes.
			// We restore os.Args to avoid breaking other tests.
			defer func() { os.Args = oldArgs }()

			// Set up our fake command-line arguments for this test.
			os.Args = tt.args

			// Load the configuration with our fake arguments.
			cfg, err := Load()

			// Check if we got an error when we expected one, or vice versa.
			// The "!=" operator means "not equal".
			// So "(err != nil) != tt.expectError" means:
			// "we got an error" is NOT equal to "we expected an error"
			if (err != nil) != tt.expectError {
				// t.Fatalf stops the test immediately with an error message.
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}

			// If there was no error expected, check if the port is correct.
			if !tt.expectError && cfg.Port != tt.expectedPort {
				// t.Errorf records an error but CONTINUES the test.
				// Useful when you want to check multiple things and see all failures.
				t.Errorf("expected port %s, got %s", tt.expectedPort, cfg.Port)
			}
		})
	}
}