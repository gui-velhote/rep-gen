package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// json structs definitions

type User struct {
	ID         int    `json:"id"`
	NAME       string `json:"name"`
	PRIVILEGES string `json:"privileges"`
}

type Client struct {
	ID      int    `json:"id"`
	NAME    string `json:"name"`
	ADDRESS string `json:"address"`
}

type Report struct {
	ID          int      `json:"id"`
	CLIENT_ID   int      `json:"client_id"`
	DATE        string   `json:"date"`
	CAR_LICENSE string   `json:"car_license"`
	TEAM        []int    `json:"team"`
	ACTIVITIES  []string `json:"activities"`
	PENDENCIES  []string `json:"pendencies"`
	OBERVATIONS []string `json:"observations"`
}

// Change database connection properties access
func dbConnection() *sql.DB {
	var db *sql.DB

	// capture connection properties
	cfg := mysql.Config{
		User:                 "golang",
		Passwd:               "2002",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "repgen",
		AllowNativePasswords: true,
	}
	// Get a database handle
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("connected!")

	return db
}

/*
 * Function to get all users data from database
 */
func getAllUsers(c *gin.Context) {

	var users []User

	data_base := dbConnection()
	rows, err := data_base.Query(`SELECT * FROM User`)

	if err != nil {
		data_base.Close()
		c.JSON(http.StatusNoContent, nil)
	}
	defer rows.Close()

	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.NAME, &user.PRIVILEGES)
		if err != nil {
			c.JSON(http.StatusNoContent, nil)
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)

	data_base.Close()
}

/*
 * Function that adds a user on the database
 */
func addUser(c *gin.Context) {

	var newUser User

	data_base := dbConnection()

	if err := c.BindJSON(&newUser); err != nil {
		data_base.Close()
		return
	}

	query := `INSERT INTO User (name, privileges) VALUES (?, ?)`

	_, err := data_base.ExecContext(context.Background(), query, newUser.NAME, newUser.PRIVILEGES)

	if err != nil {
		data_base.Close()
		return
	}

	// find a way to not need the insertResult variable
	// fmt.Println(insertResult)

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	data_base.Close()
}

func getUserById(c *gin.Context) {

	var user User

	// Connect to database
	db := dbConnection()

	// Collect GET method data
	if err := c.BindJSON(&user); err != nil {
		return
	}

	// Query database for user

	err := db.QueryRow(`SELECT * FROM User WHERE id=?`, &user.ID).Scan(&user.ID, &user.NAME, &user.PRIVILEGES)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, user)

	// Close database connection
	db.Close()

}

func getUserByName(c *gin.Context) {

	var queryUser User
	var users []User

	if err := c.BindJSON(&queryUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad request",
		})
		fmt.Println(err)
		return
	}

	db := dbConnection()

	rows, err := db.Query(`SELECT * FROM User WHERE name=?`, queryUser.NAME)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad request",
		})
		return
	}

	for rows.Next() {

		var tempUser User

		if err := rows.Scan(&tempUser.ID, &tempUser.NAME, &tempUser.PRIVILEGES); err != nil {
			fmt.Println(err)
			return
		}

		users = append(users, tempUser)

	}

	c.JSON(http.StatusOK, users)

	db.Close()
}

func getUserByNameQuery(c *gin.Context) {

	var queryUser User
	var users []User

	if err := c.BindJSON(&queryUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad request",
		})
		fmt.Println(err)
		return
	}

	db := dbConnection()

	userName := "%" + queryUser.NAME + "%"

	rows, err := db.Query(`SELECT * FROM User WHERE name LIKE ?`, userName)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "No data match",
		})
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var tempUser User

		if err := rows.Scan(&tempUser.ID, &tempUser.NAME, &tempUser.PRIVILEGES); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			fmt.Println(err)
			return
		}

		users = append(users, tempUser)
	}

	c.JSON(http.StatusOK, users)

	db.Close()

}

/*
 * function to add new clients to database
 */
func addClient(c *gin.Context) {

	var newClient Client

	// Get data from post
	if err := c.BindJSON(&newClient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad data post",
		})
		return
	}

	// Connect to database
	data_base := dbConnection()

	// prepare statement
	prepared_statement, err := data_base.PrepareContext(context.Background(), `INSERT INTO Client (name, address) VALUES(?, ?)`)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad data post",
		})
		data_base.Close()
		return
	}

	// Execute query
	insertResult, err := prepared_statement.Exec(newClient.NAME, newClient.ADDRESS)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad data post",
		})
		data_base.Close()
		return
	}

	// figure out a way to not need this
	fmt.Println(insertResult)

	// Return success status
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	data_base.Close()
}

func getAllClients(c *gin.Context) {

	var clients []Client

	db := dbConnection()

	rows, err := db.Query(`SELECT * FROM Client`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		db.Close()
		return
	}

	for rows.Next() {
		var tempClient Client

		if err := rows.Scan(&tempClient.ID, &tempClient.NAME, &tempClient.ADDRESS); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			fmt.Println(err)
			return
		}

		clients = append(clients, tempClient)
	}

	c.JSON(http.StatusOK, clients)

	db.Close()
}

func getClientById(c *gin.Context) {

	var queryClient Client

	// get data from request
	if err := c.BindJSON(&queryClient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad request",
		})
		return
	}

	// connect to database
	db := dbConnection()

	// Query for client from id
	err := db.QueryRow(`SELECT * FROM Client WHERE id=?`, queryClient.ID).Scan(&queryClient.ID, &queryClient.NAME, &queryClient.ADDRESS)

	fmt.Println(queryClient)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad request",
		})
		fmt.Println(err)
		db.Close()
		return
	}

	// http the request
	c.JSON(http.StatusOK, queryClient)

	// close database
	db.Close()
}

func getClientByName(c *gin.Context) {

	var queryClient Client

	// if elements, err := c.GetRawData(); err != nil {
	//
	// } else {
	//     for _, element := range elements {
	//         fmt.Print(string(element))
	//     }
	// }

	// get json from request
	if err := c.BindJSON(&queryClient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad request",
		})
		fmt.Println(err)
		return
	}

	// connect to database
	db := dbConnection()

	// query for client info
	err := db.QueryRow(`SELECT * FROM Client WHERE name LIKE ?`, queryClient.NAME).Scan(&queryClient.ID, &queryClient.NAME, &queryClient.ADDRESS)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad request",
		})
		fmt.Println(err)
		db.Close()
		return
	}

	// return http with info
	c.JSON(http.StatusOK, queryClient)

	// close database
	db.Close()

}

func getClientByNameQuery(c *gin.Context) {

	var clients []Client
	var queryClient Client

	db := dbConnection()

	if err := c.BindJSON(&queryClient); err != nil {
		fmt.Println(err)
		return
	}

	queryClient.NAME = "%" + queryClient.NAME + "%"

	rows, err := db.Query(`SELECT * FROM Client WHERE name LIKE ?`, queryClient.NAME)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad request",
		})
		fmt.Println(err)
		db.Close()
		return
	}

	// cicle through rows to get data
	for rows.Next() {
		var tempClient Client

		if err := rows.Scan(&tempClient.ID, &tempClient.NAME, &tempClient.ADDRESS); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"status": "No data match",
			})
			fmt.Println(err)
			db.Close()
			return
		}

		clients = append(clients, tempClient)

	}

	c.JSON(http.StatusOK, clients)

	db.Close()
}

func addReport(c *gin.Context) {

	var report Report

	// Add the request json to the variable
	if err := c.BindJSON(&report); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(report)

	db := dbConnection()

	query := `INSERT INTO Report (team, client_id, date) VALUES (?, ?, ?)`

	_, err := db.ExecContext(context.Background(), query, report.TEAM_ID, report.CLIENT_ID, report.DATE)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad request",
		})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stauts": "success",
	})

	db.Close()
}

// func getAllReports(c *gin.Context) {
//
//     var reports []Report
//
//     db := dbConnection()
//
//     rows, err := db.Query(`SELECT `, args ...any)
//
//     db.Close()
// }

func addReport(c *gin.Context) {
	var report Report

}

func main() {
	// Creates the default gin router
	router := gin.Default()

	/* Create endpoints here */
	/* Example: router.GET("/endpoint", function) */

	/* Get endpoints */

	// Users endpoints
	router.GET("/user/getAll", getAllUsers)
	router.GET("/user/getById", getUserById)
	router.GET("/user/getByName", getUserByName)
	router.GET("/user/getByNameQuery", getUserByNameQuery)

	// Clients endpoints
	router.GET("/client/getAll", getAllClients)
	router.GET("/client/getByID", getClientById)
	router.GET("/client/getByName", getClientByName)
	router.GET("/client/getByNameQuery", getClientByNameQuery)

	// Post endpoints
	router.POST("/user/new", addUser)
	router.POST("/client/new", addClient)
	router.POST("/user/addReport", addReport)

	// Run the server
	router.Run("localhost:8080")
}
