package main 


import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)


//declaring book struct with json format
type book struct {
	ID string  `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Quantity int `json:"quantity"`
}


//book variable of slice type Of Books struct
var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}



//like controller in Express Js
func getBooks(c *gin.Context){ //gin context is essentially all information about the request
	c.IndentedJSON(http.StatusOK , books)
}


//creating a new book
func createBook(c *gin.Context){ 
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		return 
	}
    books = append(books, newBook)
    c.IndentedJSON(http.StatusCreated, newBook)
}




//====== fetch book buy add : 

//get the book and either pointer to the book or error
func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}


//the route controller 
func getBook(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}




//updating a book by grabbing the book id and updating it using body data sent 
func updateBook(c *gin.Context) {
    id := c.Param("id")
    var updatedBook book

	//bing the received data with the new variable updatedBook of type book
    if err := c.BindJSON(&updatedBook); err!= nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body."})
        return
    }


	//getting the book
    book, err := getBookById(id)
    if err!= nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
        return
    }


	//updating 
    book.Title = updatedBook.Title
    book.Author = updatedBook.Author
    book.Quantity = updatedBook.Quantity

	//return the updated book 
    c.IndentedJSON(http.StatusOK, book)
}



//deleting a book 
func deleteBook(c *gin.Context) {
	id := c.Param("id")
	book , err := getBookById(id)

	fmt.Println(book)

	if err!= nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}


	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book deleted successfully."})
}


//main function 
func main (){
	router := gin.Default()
    router.GET("/books", getBooks)


	//to test POST /books
	//use curl command : 
	//curl localhost:6699/books --include --header "Content-Type : application/json" -d @body.json --request "POST"
	router.POST("/books", createBook)
    router.GET("/books/:id", getBook)
    router.PUT("/books/:id", updateBook)
    router.DELETE("/books/:id", deleteBook)
    router.Run(":6699") 
}