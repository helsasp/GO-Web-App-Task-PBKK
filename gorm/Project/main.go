package main

import (
    "database/sql"
    "log"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// User represents the schema for the 'users' table
type User struct {
    ID           uint           `gorm:"primaryKey"`  // Primary key
    Name         string         // Regular string field
    Email        *string        // Nullable string
    Age          uint8          // Unsigned 8-bit integer
    Birthday     *time.Time     // Nullable date field
    MemberNumber sql.NullString // Nullable string
    ActivatedAt  sql.NullTime   // Nullable time
    CreatedAt    time.Time      // Auto-managed creation time
    UpdatedAt    time.Time      // Auto-managed update time
}

// Blog represents the schema for the 'blogs' table
type Blog struct {
    ID      int64  `gorm:"primaryKey"` // Primary key
    Name    string                      // Blog name
    Email   string                      // Author email
    Upvotes int32                       // Upvote count
}

func main() {
    // Data Source Name (DSN) for MySQL connection
    dsn := "root:helsasp@tcp(127.0.0.1:3306)/gormcoba1?charset=utf8mb4&parseTime=True&loc=Local"

    // Initialize GORM connection to the MySQL database
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect database: %v", err)
    }

    // Automatically migrate the schema (create tables if they don't exist)
    err = db.AutoMigrate(&User{}, &Blog{})
    if err != nil {
        log.Fatalf("Failed to migrate schema: %v", err)
    }

    // Creating multiple users with proper initialization
    now := time.Now()
    users := []*User{
        {Name: "Jinzhu", Age: 18, Birthday: &now},
        {Name: "Jackson", Age: 19, Birthday: &now},
    }

    // Insert users into the database
    result := db.Create(users)
    if result.Error != nil {
        log.Fatalf("Error while inserting users: %v", result.Error)
    }

    log.Printf("Inserted %d users", result.RowsAffected)

    // Example of inserting a single user
    user := User{Name: "Alice", Age: 20, Birthday: &now}
    db.Create(&user)

    // Nullable fields example
    memberNumber := sql.NullString{String: "12345", Valid: true}
    activatedAt := sql.NullTime{Time: now, Valid: true}

    userWithNullFields := User{
        Name:         "Bob",
        Age:          25,
        Birthday:     &now,
        MemberNumber: memberNumber,
        ActivatedAt:  activatedAt,
    }
    db.Create(&userWithNullFields)

    // Select specific fields and create a new user
    db.Select("Name", "Age", "CreatedAt").Create(&User{Name: "John", Age: 22, CreatedAt: time.Now()})

    // Check the result
    log.Printf("Inserted user: %+v", user)
}
