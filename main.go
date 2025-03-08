package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "net/http"
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

// character represents differenct chars from books
type Char struct {
    ID           int     `json:"id"`
    Name         string  `json:"name"`
    Description  string  `json:"description"`
    Book         string  `json:"book"`
    Rating       int     `json:"rating"`
}

var db *sql.DB

func main() {
    // Capture connection properties from environment variables.
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DBHOST"),   // e.g., "127.0.0.1"
        os.Getenv("DBPORT"),   // e.g., "5432"
        os.Getenv("DBUSER"),   // e.g., "postgres"
        os.Getenv("DBPASS"),   // e.g., "yourpassword"
        os.Getenv("DBNAME"),   // e.g., "recordings"
    )

    // Get a database handle.
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Error opening database:", err)
    }
    defer db.Close()

    router := gin.Default()
    router.GET("/chars", getChars)
    router.GET("/char/:name", getCharByName)
    router.POST("/chars", postChars)

    router.Run("0.0.0.0:8080")
}


// getChars responds with the list of all albums as JSON.
func getChars(c *gin.Context) {

    rows, err := db.Query("SELECT * FROM characters")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var characters []Char

    for rows.Next() {
        var cha Char
        if err := rows.Scan(&cha.ID, &cha.Name, &cha.Description, &cha.Book, &cha.Rating); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        characters = append(characters, cha)
    }
    c.JSON(http.StatusOK, characters)
}


// charByID queries for the char with the specified name
func getCharByName(c *gin.Context) {
    // An char to hold data from the returned row.
    name := c.Param("name")
    var cha Char

    row := db.QueryRow("SELECT * FROM characters WHERE name = $1", name)
    if err := row.Scan(&cha.ID, &cha.Name, &cha.Description, &cha.Book, &cha.Rating); err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
    }
    c.JSON(http.StatusOK, cha)
}

// postChars adds an album from JSON received in the request body.
func postChars(c *gin.Context) {
    var newChar Char

    // Call BindJSON to bind the received JSON to newChar.
    if err := c.BindJSON(&newChar); err != nil {
        return
    }

    stmt := `INSERT INTO characters (name, description, book, rating) VALUES ($1, $2, $3, $4) RETURNING id`

    var id int
    // Execute the query and get the ID from the returned result
    err := db.QueryRow(stmt, newChar.Name, newChar.Description, newChar.Book, newChar.Rating).Scan(&id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
    }

    // Update newChar.ID with the generated ID
    newChar.ID = id

    c.IndentedJSON(http.StatusCreated, newChar)
}
