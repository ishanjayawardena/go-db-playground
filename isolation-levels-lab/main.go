package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"isolation-levels-lab/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

// runConfig represents a run configuration
type RunConfig struct {
	Concurrency     int
	MaxConnections  int
	TestRepetitions int
	Timeout         time.Duration
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()
	log.Printf("Connection pool created %v %v", pool.Config().MaxConns, pool.Config().MinConns)

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	log.Printf("Connected to PostgreSQL database %q", cfg.Database.DBName)
	nProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(nProcs)
	log.Printf("Number of CPUs: %d", nProcs)

	runCfg := &RunConfig{
		Concurrency:     nProcs * 2,
		MaxConnections:  nProcs * 2,
		TestRepetitions: 100,
		Timeout:         30 * time.Second,
	}
	demonstrateIsolationLevels(pool, runCfg)
}

// readBalance reads the balance of an account within a transaction
func readBalance(ctx context.Context, tx pgx.Tx, accountID int) (int64, error) {
	var balance int64
	err := tx.QueryRow(ctx, "SELECT balance FROM accounts WHERE id = $1", accountID).Scan(&balance) // <= FOR UPDATE is removed
	if err != nil {
		return 0, fmt.Errorf("error reading balance: %w", err)
	}
	return balance, nil
}

// updateBalance updates the balance of an account within a transaction
func updateBalance(ctx context.Context, tx pgx.Tx, accountID int, amount int64) error {
	_, err := tx.Exec(ctx, "UPDATE accounts SET balance = $1 WHERE id = $2", amount, accountID)
	if err != nil {
		return fmt.Errorf("error updating balance: %w", err)
	}
	return nil
}

// transfer attempts to transfer 1000 from account 2 to account 1 if sufficient balance exists
func transfer(ctx context.Context, pool *pgxpool.Pool, isolationLevel pgx.TxIsoLevel) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: isolationLevel})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Ensure rollback if commit doesn't happen

	var balance2, balance1 int64
	if balance2, err = readBalance(ctx, tx, 2); err != nil {
		return err
	}
	if balance1, err = readBalance(ctx, tx, 1); err != nil {
		return err
	}
	if balance2 >= 1000 {
		if err = updateBalance(ctx, tx, 1, balance1+1000); err != nil {
			return err
		}
		if err = updateBalance(ctx, tx, 2, balance2-1000); err != nil {
			return err
		}
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return tx.Commit(ctx)
	}
}

// resetBalances resets the balances of both accounts to their initial values
func resetBalances(ctx context.Context, pool *pgxpool.Pool, isolationLevel pgx.TxIsoLevel) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: isolationLevel})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := updateBalance(ctx, tx, 1, 1000); err != nil {
		return err
	}
	if err := updateBalance(ctx, tx, 2, 2000); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return tx.Commit(ctx)
	}
}

// verifyBalances checks that the total balance remains constant at 3000
func verifyBalances(ctx context.Context, pool *pgxpool.Pool, isolationLevel pgx.TxIsoLevel) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: isolationLevel})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	var balance1, balance2 int64
	if balance1, err = readBalance(ctx, tx, 1); err != nil {
		return fmt.Errorf("error reading balance 1: %w", err)
	}
	if balance2, err = readBalance(ctx, tx, 2); err != nil {
		return fmt.Errorf("error reading balance 2: %w", err)
	}
	if !(balance1 >= 0 && balance2 >= 0 && balance1+balance2 == 3000) {
		return fmt.Errorf("balances are not correct: %d %d", balance1, balance2)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return tx.Commit(ctx)
	}
}

// bulkTransfer executes multiple concurrent transfers with the specified isolation level
func bulkTransfer(ctx context.Context, pool *pgxpool.Pool, runCfg *RunConfig, isolationLevel pgx.TxIsoLevel) error {
	if err := resetBalances(ctx, pool, isolationLevel); err != nil {
		return fmt.Errorf("error resetting balances: %w", err)
	}
	defer func() {
		if err := verifyBalances(ctx, pool, isolationLevel); err != nil {
			log.Printf("error verifying balances: %v", err)
		} else {
			//log.Printf("Balances verified. Test completed successfully.")
		}
	}()

	pool.Config().MaxConns = int32(runCfg.MaxConnections)
	eg, ctx0 := errgroup.WithContext(ctx)
	eg.SetLimit(runCfg.Concurrency)
	for range make([]int, runCfg.Concurrency) {
		eg.Go(func() error {
			return transfer(ctx0, pool, isolationLevel)
		})
	}
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error in transfer: %w", err)
	}
	return nil
}

// demonstrateIsolationLevels runs tests for different isolation levels
func demonstrateIsolationLevels(pool *pgxpool.Pool, runCfg *RunConfig) {
	// Test Scenario:
	// 1. Initial state:
	//    - Account A: 1000
	//    - Account B: 2000
	//    - Total: 3000
	// 2. Concurrent transfers:
	//    - Each transaction attempts to transfer 1000 from B to A
	//    - Only succeeds if B has sufficient balance (>= 1000)
	// 3. Expected outcomes:
	//    - Only two transactions should succeed (B has 2000 initially) at most
	//    - Remaining transactions fail due to insufficient balance
	//    - Total balance remains 3000 throughout
	// 4. Isolation level differences:
	//    - READ COMMITTED: May allow more transfers due to lost updates
	//    - SERIALIZABLE: Prevents concurrent updates, maintains strict consistency

	ctx, cancel := context.WithTimeout(context.Background(), runCfg.Timeout)
	defer cancel()

	log.Println("============= Demonstrating READ COMMITTED =============")
	demonstrateReadCommitted(ctx, pool, runCfg)

	log.Println("============= Demonstrating SERIALIZABLE =============")
	demonstrateSerializable(ctx, pool, runCfg)
}

// demonstrateReadCommitted runs tests with READ COMMITTED isolation level
func demonstrateReadCommitted(ctx context.Context, pool *pgxpool.Pool, runCfg *RunConfig) {
	runTest(ctx, pool, runCfg, pgx.ReadCommitted)
}

// demonstrateSerializable runs tests with SERIALIZABLE isolation level
func demonstrateSerializable(ctx context.Context, pool *pgxpool.Pool, runCfg *RunConfig) {
	runTest(ctx, pool, runCfg, pgx.Serializable)
}

// runTest executes multiple test iterations with the specified isolation level
func runTest(ctx context.Context, pool *pgxpool.Pool, runCfg *RunConfig, isolationLevel pgx.TxIsoLevel) {
	log.Printf("Running %d transactions repeatedly...", runCfg.Concurrency)
	for range make([]int, runCfg.TestRepetitions) {
		err := bulkTransfer(ctx, pool, runCfg, isolationLevel)
		if err != nil {
			// For SERIALIZABLE: returns "ERROR: could not serialize access due to concurrent update (SQLSTATE 40001)"
			// but account balances remain consistent (total always 3000)
			//log.Printf("error in bulk transfer: %v", err)
		}
	}
	log.Printf("Test completed...")
}
