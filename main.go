package main

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
)


// we need to use json tag in GO, to seralized and convert to json
// So our API can return the struct directly
// And so is convert json into Go struct

// first should be the capital letter because it's is the field which can read outside the module

// json shuold be the lowwer letter because it stand for 'when we serialize json, conver the filed name to lowwer case'
// EX: ID => id

// Or if we get the json object we will look for the lowwer cast and turn into upper case field
// EX: json field of 'title' => TITLE

type book struct {
	ID       string `json:"id"` 
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}


// The 'gin.Context' is all of the information about the request, and it allow to return response

// The 'c' is a variable which can get the specific requst information
// EX: The request header, data payload, query parameter

func getBooks(c *gin.Context) {
	// The 'c.IndentedJSON' which is nicely formmatted the json data
	c.IndentedJSON(http.StatusOK, books) 
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		//gin.H is the shortcut make us to wirte JSON information
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	// The GetQuery will retrun true ^ false in the second variable (which is ok here)
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			//return the pointer of book because we can modify the attribute of the book or the field of the struct from different function
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func createBook(c *gin.Context) {
	var newBook book

	// The 'c.BindJSON' is a part of the requst (because the 'c')
	// With the pointer & to newBook which can that us can directly modify the field values
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
