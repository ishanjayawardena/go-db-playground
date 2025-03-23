# Go Database Playground

A collection of Go projects demonstrating various database concepts and patterns.

## Projects

### 1. Isolation Levels Lab
A demonstration of PostgreSQL transaction isolation levels and their effects on concurrent transactions. This lab specifically focuses on:
- READ COMMITTED vs SERIALIZABLE isolation levels
- ACID Rain attack vulnerability
- Money transfer scenario with concurrent transactions
- Connection pooling and transaction management

See [isolation-levels-lab/README.md](isolation-levels-lab/README.md) for detailed information.

### 2. GORM N+1 Query Lab
A demonstration of the N+1 query problem using GORM and different solutions to prevent it.

See [gorm-nplus1-lab/README.md](gorm-nplus1-lab/README.md) for detailed information.

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL client (optional)

## Getting Started

Each lab is a self-contained module with its own setup instructions. Please refer to the individual lab READMEs for specific setup and running instructions.

## Project Structure

```
.
├── isolation-levels-lab/    # Transaction isolation levels demonstration
├── gorm-nplus1-lab/        # GORM N+1 query problem demonstration
└── README.md              # This file
```

## Contributing

Feel free to contribute by:
1. Creating new labs
2. Improving existing labs
3. Fixing bugs
4. Adding documentation

## License

This project is licensed under the MIT License - see the LICENSE file for details.