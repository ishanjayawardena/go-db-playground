# GORM N+1 Query Problem Lab

This module demonstrates the N+1 query problem in GORM and shows different approaches to solve it. It serves as an educational example of database query optimization in Go applications using GORM.

## Overview

The N+1 query problem occurs when an application makes N additional queries to fetch related data after an initial query. This is inefficient and can lead to performance issues, especially with large datasets.

This lab demonstrates:
1. The N+1 problem (what not to do)
2. Using GORM's Preload feature (efficient solution)
3. Using SQL JOINs with aggregation (alternative solution)

## Structure

The code is organized into several sections:

### Types
- `User`: Represents a user with a one-to-many relationship to posts
- `Post`: Represents a post belonging to a user
- `UserWithPostCount`: A composite type used for the JOIN example

### Database Initialization
- `connectToDB()`: Sets up the database connection with logging
- `initDB()`: Handles schema migration and index creation
- `seedData()`: Populates the database with sample data

### Query Demonstrations
- `demonstrateNPlusOneProblem()`: Shows the inefficient N+1 query pattern
- `demonstratePreloadSolution()`: Shows how to use GORM's Preload feature
- `demonstrateJoinsSolution()`: Shows how to use JOINs for aggregation

## Running the Example

```bash
go run main.go
```

The program will:
1. Connect to a SQLite database
2. Initialize the schema and create necessary indexes
3. Seed the database with sample data
4. Demonstrate three different query approaches
5. Show the SQL queries and results for each approach

## Key Learnings

1. **N+1 Problem**
   - Avoid making separate queries for related data
   - Each additional query adds latency and database load

2. **Preload Solution**
   - Uses GORM's Preload feature to fetch related data efficiently
   - Reduces multiple queries to a single JOIN query
   - Maintains clean code structure

3. **JOIN Solution**
   - Useful for aggregation scenarios
   - Performs calculations at the database level
   - Can be more efficient for specific use cases
   - Note: While `GROUP BY users.id, users.name` is valid, in this case `GROUP BY users.id` is sufficient
     since `users.name` is functionally dependent on `users.id` (each user ID has exactly one name).

## Best Practices Demonstrated

- Proper database initialization and setup
- Index creation for performance
- Error handling and logging
- Code organization and documentation
- Query optimization techniques
