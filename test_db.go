package main

import (
"fmt"
"log"

"github.com/binbankm/My/internal/model"
"github.com/binbankm/My/internal/util"
)

func main() {
// Initialize database
if err := model.InitDB(); err != nil {
log.Fatalf("Failed to initialize database: %v", err)
}

// Fetch the admin user
var user model.User
if err := model.DB.Where("username = ?", "admin").First(&user).Error; err != nil {
fmt.Printf("Failed to find admin user: %v\n", err)
return
}

fmt.Printf("User found: %+v\n", user)
fmt.Printf("Password hash: %s\n", user.Password)

// Test password check
password := "admin123"
if util.CheckPassword(password, user.Password) {
fmt.Println("Password check: SUCCESS")
} else {
fmt.Println("Password check: FAILED")
}
}
