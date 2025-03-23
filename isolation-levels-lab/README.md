# Database Isolation Levels Lab

This project demonstrates different transaction isolation levels in PostgreSQL and their effects on concurrent transactions. It specifically focuses on the behavior of READ COMMITTED and SERIALIZABLE isolation levels in a money transfer scenario.

## Purpose

The demonstration shows how different isolation levels affect concurrent transactions when multiple operations try to transfer money between two accounts simultaneously. This helps understand:

1. How isolation levels prevent or allow different types of anomalies
2. The trade-off between consistency and concurrency
3. The behavior of concurrent transactions in a real-world scenario
4. Protection against ACID Rain attacks

### ACID Rain Attack

The ACID Rain attack is a type of transaction isolation vulnerability where concurrent transactions can exploit the behavior of certain isolation levels to perform unauthorized operations. In this demonstration:

- Multiple transactions attempt to transfer money simultaneously
- Under READ COMMITTED isolation, transactions can see partial results of other transactions
- This can lead to lost updates and incorrect balance calculations
- SERIALIZABLE isolation prevents these issues by ensuring strict transaction ordering

The attack demonstrates why choosing the right isolation level is crucial for maintaining data integrity in concurrent environments.

### Test Scenario

The demonstration uses a simple money transfer scenario:
- Two accounts: Account A (initial balance: 1000) and Account B (initial balance: 2000)
- Multiple concurrent transactions attempt to transfer 1000 from B to A
- Each transaction only succeeds if B has sufficient balance (>= 1000)
- The total balance remains constant at 3000 throughout

### Expected Outcomes

1. READ COMMITTED:
   - May allow more transfers than expected due to lost updates
   - Each transaction sees the latest committed state
   - Less strict consistency but better concurrency
   - Vulnerable to ACID Rain attacks

2. SERIALIZABLE:
   - Prevents concurrent updates
   - Maintains strict consistency
   - Transactions are serialized to prevent anomalies
   - May result in serialization failures (SQLSTATE 40001)
   - Protects against ACID Rain attacks

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL client (optional, for direct database access)
- Make (optional, for using Makefile commands)

## Setup and Running

### Using Make (Recommended)

1. Start the database:
   ```bash
   make db-up
   ```

2. Run the demonstration:
   ```bash
   make run
   ```

3. View database logs:
   ```bash
   make db-logs
   ```

4. Connect to database:
   ```bash
   make db-shell
   ```

5. Clean up:
   ```bash
   make db-down
   ```

### Manual Setup

1. Start the PostgreSQL database:
   ```bash
   docker-compose up -d
   ```
   This will:
   - Start PostgreSQL with the required configuration
   - Create the database and tables
   - Set up initial account balances

2. Run the demonstration:
   ```bash
   go run main.go
   ```
   The program will:
   - Connect to the database
   - Run tests with both READ COMMITTED and SERIALIZABLE isolation levels
   - Display the results of concurrent transfers
   - Show how different isolation levels affect the outcome
   - Demonstrate protection against ACID Rain attacks

## Database Access

You can connect to the PostgreSQL database using psql in two ways:

1. From inside the container:
   ```bash
   docker exec -it isolation-levels-lab-postgres-1 psql -U isolation_user -d isolation_db
   ```
   Or simply:
   ```bash
   make db-shell
   ```

2. From your host machine:
   ```bash
   psql -h localhost -p 5432 -U isolation_user -d isolation_db
   ```
   When prompted for a password, enter: `isolation_pass`

### Useful psql Commands
- `\dt` - list all tables
- `\d accounts` - describe the accounts table
- `SELECT * FROM accounts;` - view all accounts
- `\q` - quit psql

## Development Commands

The project includes a Makefile with common development commands:

- `make build` - Build the application
- `make run` - Run the application
- `make clean` - Clean build artifacts
- `make test` - Run tests
- `make lint` - Run linter
- `make fmt` - Format code
- `make tidy` - Tidy dependencies
- `make help` - Show all available commands

## Configuration

The demonstration can be configured through environment variables or a config file:
- Database connection details
- Number of concurrent transactions
- Test repetitions
- Timeout duration

## Cleanup

To stop the database:
```bash
make db-down
```
Or manually:
```bash
docker-compose down
```

## Notes

- The demonstration uses connection pooling for better performance
- Each test iteration resets the account balances to their initial state
- The program verifies that the total balance remains constant
- Serialization failures are expected with SERIALIZABLE isolation level
- The ACID Rain attack demonstrates why proper transaction isolation is crucial for financial operations 