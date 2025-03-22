package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Types

// User represents a user in the system.
type User struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `gorm:"size:255;not null"`
	Posts []Post `gorm:"foreignKey:UserID"`
}

// Post represents a blog post in the system.
type Post struct {
	ID     uint   `gorm:"primarykey"`
	Title  string `gorm:"size:255;not null"`
	Body   string `gorm:"type:text"`
	UserID uint   `gorm:"not null;index"`
}

// UserWithPostCount represents a user with their post count.
type UserWithPostCount struct {
	User
	PostCount int
}

func main() {
	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := initDB(db); err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	if err := seedData(db); err != nil {
		log.Fatalf("Failed to seed the database: %v", err)
	}

	// BAD: This will cause N+1 problem
	demonstrateNPlusOneProblem(db)

	// GOOD: Using Preload to avoid N+1 problem
	demonstratePreloadSolution(db)

	// ALTERNATIVE: Using Joins
	demonstrateJoinsSolution(db)
}

// Database initialization functions

// connectToDB establishes a connection to the SQLite database with logging enabled.
func connectToDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	})
}

// initDB initializes the database schema and creates necessary indexes.
func initDB(db *gorm.DB) error {
	log.Println("Auto migrating the schema...")
	if err := db.AutoMigrate(&User{}, &Post{}); err != nil {
		return err
	}

	log.Println("Creating indexes...")
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id)").Error; err != nil {
		log.Printf("Warning: failed to create index: %v", err)
		return err
	}
	return nil
}

// seedData populates the database with sample data.
func seedData(db *gorm.DB) error {
	log.Println("Seeding data...")

	// Clear previous data
	if err := db.Delete(&Post{}, "1=1").Error; err != nil {
		return err
	}
	if err := db.Delete(&User{}, "1=1").Error; err != nil {
		return err
	}

	// Create users
	users := []User{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
	}
	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			return err
		}
	}

	// Create posts
	posts := []Post{
		{Title: "First post", Body: "Content 1", UserID: users[0].ID},
		{Title: "Second post", Body: "Content 2", UserID: users[0].ID},
		{Title: "Hello world", Body: "Content 3", UserID: users[1].ID},
		{Title: "GORM tutorial", Body: "Content 4", UserID: users[2].ID},
		{Title: "Another post", Body: "Content 5", UserID: users[2].ID},
	}
	for _, post := range posts {
		if err := db.Create(&post).Error; err != nil {
			return err
		}
	}
	return nil
}

// Query demonstration functions

// demonstrateNPlusOneProblem shows an example of the N+1 query problem in GORM.
// It performs a query that results in multiple database calls:
// 1. First query fetches all users
// 2. Then for each user, it makes a separate query to fetch their posts
// This is inefficient as it results in N+1 queries where N is the number of users.
// A better approach would be to use Preload or Joins to fetch the data in fewer queries.
func demonstrateNPlusOneProblem(db *gorm.DB) {
	log.Println("===================================================================")
	log.Printf("Database query executed without preloading, potential N+1 problem")
	var users []User
	if err := db.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return
	}
	for _, user := range users {
		log.Printf("User: %s\n", user.Name)
		var posts []Post
		if err := db.Where("user_id = ?", user.ID).Find(&posts).Error; err != nil {
			log.Printf("Error fetching posts for user %s: %v", user.Name, err)
			continue
		}
		for _, post := range posts {
			log.Printf("  - Post: %s\n", post.Title)
		}
	}
}

// demonstratePreloadSolution shows how to efficiently fetch related data using GORM's Preload feature.
// Instead of making N+1 queries, it uses a single query with JOIN to fetch all users and their posts.
// This is much more efficient than the N+1 approach as it reduces the number of database calls to just one.
func demonstratePreloadSolution(db *gorm.DB) {
	log.Println("===================================================================")
	log.Printf("Database query executed with preloading, avoids the N+1 problem")
	var usersWithPosts []User
	if err := db.Preload("Posts").Find(&usersWithPosts).Error; err != nil {
		log.Printf("Error fetching users with posts: %v", err)
		return
	}
	for _, user := range usersWithPosts {
		log.Printf("User: %s\n", user.Name)
		for _, post := range user.Posts {
			log.Printf("  - Post: %s\n", post.Title)
		}
	}
}

// demonstrateJoinsSolution shows an alternative approach to avoid the N+1 problem using SQL JOINs.
// This approach is particularly useful when you need to aggregate data, like counting related records.
// It uses a single query with LEFT JOIN and GROUP BY to fetch users with their post counts.
// This is efficient as it performs the aggregation at the database level rather than in the application.
func demonstrateJoinsSolution(db *gorm.DB) {
	log.Println("===================================================================")
	log.Printf("Database query executed with JOINs, demonstrates aggregation with GROUP BY")
	var usersWithPostCounts []UserWithPostCount

	// Use a more specific SELECT to only fetch needed columns
	err := db.Model(&User{}).
		Select(`users.id, users.name, COUNT(p.id) AS post_count`).
		Joins(`LEFT JOIN posts p ON p.user_id = users.id`).
		Group(`users.id`). // users.name is functionally dependent on users.id
		Limit(100).        // Add pagination for large datasets
		Find(&usersWithPostCounts).
		Error

	if err != nil {
		log.Printf("Error fetching users with post counts: %v", err)
		return
	}

	for _, user := range usersWithPostCounts {
		log.Printf("User: %s has %d posts\n", user.Name, user.PostCount)
	}
}
