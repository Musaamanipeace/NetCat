# Net-Cat: Modular TCP Chat Server


**A blazing-fast, concurrent TCP chat server built in Go—refactored using clean architectural guidelines.**


## 📂 Project Structure

This project follows a clean, modular Go layout pattern:

net/
├── cmd/
│   └── server/
│       └── main.go       # Application entry point
├── internal/
│   └── chat/
│       ├── ascii.go      # Graphic assets & assets management
│       ├── client.go     # Client logic and session context
│       ├── message.go    # Message structures & text layout formatter
│       └── server.go     # TCP infrastructure & broadcasting multiplexer
├── pkg/
│   └── config/
│       └── config.go     # Arguments parser & parameters configuration
├── Makefile              # Project workflow automation
├── go.mod                # Module definitions
└── README.md             # This comprehensive architecture breakdown
🛠 Automation Workflow (Makefile)
A Makefile is included to streamline operations.

Clean & Build Binary
Bash
make build
This generates a thread-safe static binary file named net-cat.

Execute Unit Test Files
Bash
make test
Runs the test suites for configuration validation parser logic, message layout formatters, and network connectivity state maps with active race condition tracking.

Launch Server Instance
Bash
# Launch on default port (8989)
make run

# Launch on explicit custom port
make run port=2525
Clean Work Directory
Bash
make clean
✨ System Features
Concurrent Engine: Utilizes goroutines alongside isolated channels to coordinate messaging structures without tracking blocks.

Thread-Safety Controls: Mutex flags protect modifications across standard arrays and internal active client memory pools.

Historic Catch-Up: Active connections receive raw buffered elements to capture all historical communications inside the execution pool.

Strict Parameter Constraints: Rejects empty usernames, drops active duplicate credentials automatically, and strictly limits structural concurrency to 10 endpoints.