# Go Database Playground

A collection of Go database examples and experiments by [@ishanjayawardena](https://github.com/ishanjayawardena).

## Structure

### Current Labs

#### gorm-nplus1-lab
A lab demonstrating the N+1 query problem in GORM and its solutions:
- N+1 problem demonstration
- Preload solution
- JOIN solution with aggregation

### Planned Labs

More labs will be added over time to demonstrate various database concepts and techniques:

#### Transaction Isolation Levels
- Read Uncommitted
- Read Committed
- Repeatable Read
- Serializable
- Handling race conditions

#### Database Partitioning
- Table partitioning
- Range partitioning
- List partitioning
- Hash partitioning
- Partition pruning

#### Database Sharding
- Horizontal sharding
- Vertical sharding
- Shard key selection
- Cross-shard queries
- Shard rebalancing

#### Connection Pooling
- Pool sizing optimization
- Connection lifecycle
- High concurrency handling
- Connection health checks

#### Query Optimization
- Index strategies
- Query planning
- Execution optimization
- Performance tuning

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/ishanjayawardena/go-db-playground.git
cd go-db-playground
```

2. Run a specific lab:
```bash
# Run the N+1 problem lab
make run-nplus1
```

## Development

### Prerequisites
- Go 1.21 or later
- Make

### Available Make Commands
```bash
make help        # Show available commands
make build       # Build all labs
make test        # Run tests
make clean       # Clean build artifacts
make run-nplus1  # Run the N+1 problem lab
```

### Linting
The project uses golangci-lint for code quality checks. Linting runs automatically on:
- Push to main branch
- Pull requests

To run linting locally:
```bash
golangci-lint run
```

## Contributing

Contributions are welcome! If you'd like to:
1. Add a new lab
2. Improve existing labs
3. Fix bugs
4. Add documentation

Please feel free to submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.