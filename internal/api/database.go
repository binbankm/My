package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DatabaseConnection represents a database connection configuration
type DatabaseConnection struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"` // mysql, postgresql
	Host     string    `json:"host"`
	Port     int       `json:"port"`
	Username string    `json:"username"`
	Password string    `json:"password,omitempty"`
	Database string    `json:"database"`
	Status   string    `json:"status"`
	LastTest time.Time `json:"lastTest"`
}

// In-memory storage for database connections
var dbConnections = make(map[string]*DatabaseConnection)

// ListDatabases returns list of database connections
func ListDatabases(c *gin.Context) {
	connections := make([]*DatabaseConnection, 0, len(dbConnections))
	for _, conn := range dbConnections {
		// Don't send passwords in list
		connCopy := *conn
		connCopy.Password = ""
		connections = append(connections, &connCopy)
	}
	c.JSON(http.StatusOK, connections)
}

// CreateDatabase creates a new database connection
func CreateDatabase(c *gin.Context) {
	var conn DatabaseConnection

	if err := c.ShouldBindJSON(&conn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if conn.Name == "" || conn.Type == "" || conn.Host == "" || conn.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Validate type
	if conn.Type != "mysql" && conn.Type != "postgresql" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid database type. Must be 'mysql' or 'postgresql'"})
		return
	}

	// Set default port if not provided
	if conn.Port == 0 {
		if conn.Type == "mysql" {
			conn.Port = 3306
		} else if conn.Type == "postgresql" {
			conn.Port = 5432
		}
	}

	// Generate ID
	conn.ID = fmt.Sprintf("%s-%d", conn.Name, time.Now().Unix())

	// Test connection
	if err := testConnection(&conn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to connect: " + err.Error()})
		return
	}

	conn.Status = "connected"
	conn.LastTest = time.Now()

	// Store connection
	dbConnections[conn.ID] = &conn

	// Don't send password back
	conn.Password = ""
	c.JSON(http.StatusOK, conn)
}

// GetDatabase gets a specific database connection
func GetDatabase(c *gin.Context) {
	id := c.Param("id")

	conn, exists := dbConnections[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database connection not found"})
		return
	}

	// Don't send password
	connCopy := *conn
	connCopy.Password = ""
	c.JSON(http.StatusOK, &connCopy)
}

// DeleteDatabase deletes a database connection
func DeleteDatabase(c *gin.Context) {
	id := c.Param("id")

	if _, exists := dbConnections[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database connection not found"})
		return
	}

	delete(dbConnections, id)
	c.JSON(http.StatusOK, gin.H{"message": "Database connection deleted successfully"})
}

// TestDatabase tests a database connection
func TestDatabase(c *gin.Context) {
	id := c.Param("id")

	conn, exists := dbConnections[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database connection not found"})
		return
	}

	if err := testConnection(conn); err != nil {
		conn.Status = "disconnected"
		c.JSON(http.StatusBadRequest, gin.H{"error": "Connection failed: " + err.Error(), "status": "disconnected"})
		return
	}

	conn.Status = "connected"
	conn.LastTest = time.Now()
	c.JSON(http.StatusOK, gin.H{"message": "Connection successful", "status": "connected", "lastTest": conn.LastTest})
}

// ExecuteQuery executes a SQL query on a database connection
func ExecuteQuery(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Query string `json:"query" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, exists := dbConnections[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database connection not found"})
		return
	}

	db, err := openConnection(conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect: " + err.Error()})
		return
	}
	defer db.Close()

	// Execute query
	rows, err := db.Query(req.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query failed: " + err.Error()})
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get columns: " + err.Error()})
		return
	}

	// Prepare result
	var results []map[string]interface{}
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row: " + err.Error()})
			return
		}

		// Create map for this row
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			// Convert []byte to string
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	c.JSON(http.StatusOK, gin.H{
		"columns": columns,
		"rows":    results,
		"count":   len(results),
	})
}

// ListDatabaseTables lists all tables in a database
func ListDatabaseTables(c *gin.Context) {
	id := c.Param("id")

	conn, exists := dbConnections[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database connection not found"})
		return
	}

	db, err := openConnection(conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect: " + err.Error()})
		return
	}
	defer db.Close()

	var query string
	if conn.Type == "mysql" {
		query = fmt.Sprintf("SHOW TABLES FROM `%s`", conn.Database)
	} else if conn.Type == "postgresql" {
		query = "SELECT tablename FROM pg_tables WHERE schemaname = 'public'"
	}

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tables: " + err.Error()})
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}
		tables = append(tables, tableName)
	}

	c.JSON(http.StatusOK, gin.H{"tables": tables})
}

// Helper function to open a database connection
func openConnection(conn *DatabaseConnection) (*sql.DB, error) {
	var dsn string

	if conn.Type == "mysql" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			conn.Username, conn.Password, conn.Host, conn.Port, conn.Database)
		return sql.Open("mysql", dsn)
	} else if conn.Type == "postgresql" {
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			conn.Host, conn.Port, conn.Username, conn.Password, conn.Database)
		return sql.Open("pgx", dsn)
	}

	return nil, fmt.Errorf("unsupported database type: %s", conn.Type)
}

// Helper function to test a database connection
func testConnection(conn *DatabaseConnection) error {
	db, err := openConnection(conn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Test the connection with a timeout
	db.SetConnMaxLifetime(5 * time.Second)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	return db.Ping()
}
