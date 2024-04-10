package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var config Config

type Config struct {
	Database struct {
		USER                   string `toml:"user"`
		PASSWORD               string `toml:"password"`
		NET                    string `toml:"net"`
		ADDRESS                string `toml:"address"`
		DATABASE_NAME          string `toml:"database_name"`
		ALLOW_NATIVE_PASSWORDS bool   `toml:"allow_native_passwords"`
	}
}

type Employee struct {
	ID         int    `json:"id"`
	NAME       string `json:"name"`
	PRIVILEGES string `json:"privileges"`
}

type Client struct {
	ID   int    `json:"id"`
	NAME string `json:"name"`
}

type Building struct {
	ID        int    `json:"id"`
	CLIENT_ID int    `json:"client_id"`
	ADDRESS   string `json:"address"`
	STATUS    string `json:"status"`
}

type Visit struct {
	ID   int    `json:"id"`
	DATE string `json:"date"`
	CAR  string `json:"car"`
}

type Report struct {
	VISIT_ID  int    `json:"id"`
	DATE      string `json:"date"`
	CAR       string `json:"car"`
	CLIENT_ID int    `json:"client_id"`
	TEAM_IDS  []int  `json:"team_ids"`
	ACTIVITY  []struct {
		ID          int    `json:"activity_id"`
		DESCRIPTION string `json:"activity_description"`
	}
	OBSERVATION []struct {
		ID          int    `json:"observation_id"`
		DESCRIPTION string `json:"observation_description"`
	}
	PENDENCY []struct {
		ID          int    `json:"pendency_id"`
		DESCRIPTION string `json:"pendency_description"`
	}
}

func parseDatbaseConfig(FILE_PATH string) Config {

	var config Config

	data, err := os.ReadFile(FILE_PATH)

	if err != nil {
		log.Fatal(err)
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}

func databaseConnection() *sql.DB {

	cfg := mysql.Config{
		User:                 config.Database.USER,
		Passwd:               config.Database.PASSWORD,
		Net:                  config.Database.NET,
		Addr:                 config.Database.ADDRESS,
		DBName:               config.Database.DATABASE_NAME,
		AllowNativePasswords: config.Database.ALLOW_NATIVE_PASSWORDS,
	}

	var db *sql.DB
	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	if pingErr := db.Ping(); pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("connected!")

	return db
}

func getAllEmployees(context *gin.Context) {

	var employees []Employee

	database := databaseConnection()

	rows, err := database.Query(`SELECT * FROM Employee`)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		log.Fatal(err)
		database.Close()
		return
	}

	for rows.Next() {
		var tempEmployee Employee
		if err := rows.Scan(&tempEmployee.ID, &tempEmployee.NAME, &tempEmployee.PRIVILEGES); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			log.Fatal(err)
			database.Close()
			return
		}

		employees = append(employees, tempEmployee)
	}

	context.JSON(http.StatusOK, employees)

	database.Close()

}

func getEmployeeById(context *gin.Context) {
	var employee Employee
	var employeeId struct {
		ID int `json:"id"`
	}

	if err := context.BindJSON(&employeeId); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		return
	}

	db := databaseConnection()

	preparedStatement, err := db.Prepare(`SELECT * FROM Employee WHERE id = ?`)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		log.Fatal(err)
		db.Close()
		return
	}

	err = preparedStatement.QueryRow(employeeId.ID).Scan(&employee.ID, &employee.NAME, &employee.PRIVILEGES)

	if err != nil {
		return
	}

	context.JSON(http.StatusOK, employee)

	db.Close()
}

func addEmployee(context *gin.Context) {
	var employee Employee

	if err := context.BindJSON(&employee); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		log.Fatal(err)
		return
	}

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`INSERT INTO Employee (name, privileges) VALUES (?, ?)`)

	if err != nil {
		return
	}

	if _, err := preparedStatement.Exec(employee.NAME, employee.PRIVILEGES); err != nil {
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	database.Close()

}

func getEmployeesByName(context *gin.Context) {
	var employees []Employee
	var employeeName struct {
		NAME string `json:"name"`
	}

	// fix return
	if err := context.BindJSON(&employeeName); err != nil {
		return
	}

	employeeName.NAME = "%" + employeeName.NAME + "%"

	fmt.Println(employeeName.NAME)

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`SELECT * FROM Employee WHERE name LIKE ?`)

	// fix return
	if err != nil {
		return
	}

	rows, err := preparedStatement.Query(&employeeName.NAME)

	if err != nil {
		return
	}

	for rows.Next() {
		var tempEmployee Employee
		if err := rows.Scan(&tempEmployee); err != nil {
			return
		}
		employees = append(employees, tempEmployee)
	}

	fmt.Println(employees)

	context.JSON(http.StatusOK, employees)

	database.Close()
}

func getAllClients(context *gin.Context) {
	var clients []Client

	db := databaseConnection()
	preparedStatement, err := db.Prepare(`SELECT * FROM Client`)

	if err != nil {
		return
	}

	rows, err := preparedStatement.Query()

	if err != nil {
		return
	}

	for rows.Next() {
		var tempClient Client

		if err := rows.Scan(&tempClient.ID, &tempClient.NAME); err != nil {
			return
		}

		clients = append(clients, tempClient)
	}

	context.JSON(http.StatusOK, clients)

	db.Close()
}

func getClientByName(context *gin.Context) {
	var clients []Client
	var clientName struct {
		NAME string `json:"name"`
	}

	if err := context.BindJSON(&clientName); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		log.Fatal(err)
		return
	}

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`SELECT * FROM Client WHERE name LIKE ?`)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		log.Fatal(err)
		return
	}

	clientName.NAME = "%" + clientName.NAME + "%"

	fmt.Println(clientName.NAME)

	rows, err := preparedStatement.Query(&clientName.NAME)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Errror",
		})
		log.Fatal(err)
		return
	}

	for rows.Next() {
		var tempClient Client
		if err := rows.Scan(&tempClient.ID, &tempClient.NAME); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Errror",
			})
			log.Fatal(err)
			return
		}
		clients = append(clients, tempClient)
	}

	context.JSON(http.StatusOK, clients)

	database.Close()
}

func getClientById(context *gin.Context) {
	var client Client

	if err := context.BindJSON(&client); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		log.Fatal(err)
		return
	}

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`SELECT * FROM Client WHERE id = ?`)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
	}

	err = preparedStatement.QueryRow(&client.ID).Scan(&client.ID, &client.NAME)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		log.Fatal(err)
		return
	}

	context.JSON(http.StatusOK, client)

	database.Close()
}

func addClient(context *gin.Context) {
	var client Client

	// fix return
	if err := context.BindJSON(&client); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		// log.Fatal(err)
		fmt.Println(err)
		return
	}

	fmt.Println(client)

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`INSERT INTO Client (name) VALUES (?)`)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		log.Fatal(err)
		return
	}

	_, err = preparedStatement.Exec(&client.NAME)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		log.Fatal(err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	database.Close()
}

func getAllBuildings(context *gin.Context) {
	var buildings []Building

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`SELECT * FROM Building`)

	if err != nil {
		return
	}

	rows, err := preparedStatement.Query()

	if err != nil {
		return
	}

	for rows.Next() {
		var tempBuilding Building
		if err := rows.Scan(&tempBuilding.ID, &tempBuilding.CLIENT_ID, &tempBuilding.ADDRESS, &tempBuilding.STATUS); err != nil {
			return
		}

		buildings = append(buildings, tempBuilding)
	}

	context.JSON(http.StatusOK, buildings)

	database.Close()
}

func addBuilding(context *gin.Context) {
	var building Building

	if err := context.BindJSON(&building); err != nil {
		return
	}

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`INSERT INTO Building (client_id, address, status) VALUES (?, ?, ?)`)

	if err != nil {
		return
	}

	if _, err := preparedStatement.Exec(&building.CLIENT_ID, &building.ADDRESS, &building.STATUS); err != nil {
		return
	}

	fmt.Println(building)

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	database.Close()
}

func getBuildingsByClientId(context *gin.Context) {
	var buildings []Building
	var clientId Client

	if err := context.BindJSON(&clientId); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		fmt.Println(err)
		return
	}

	database := databaseConnection()

	prepredStatement, err := database.Prepare(`SELECT * FROM Building WHERE client_id=?`)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return
	}

	rows, err := prepredStatement.Query(&clientId.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var tempBuilding Building
		if err := rows.Scan(&tempBuilding.ID, &tempBuilding.CLIENT_ID, &tempBuilding.ADDRESS, &tempBuilding.STATUS); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			fmt.Println(err)
			return
		}

		buildings = append(buildings, tempBuilding)
	}

	context.JSON(http.StatusOK, buildings)

	database.Close()
}

func getBuildingById(context *gin.Context) {
	var building Building

	if err := context.BindJSON(&building); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		fmt.Println(err)
		return
	}

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`SELECT * FROM Building WHERE id=?`)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return
	}

	if err := preparedStatement.QueryRow(&building.ID).Scan(&building.ID, &building.CLIENT_ID, &building.ADDRESS, &building.STATUS); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return

	}

	context.JSON(http.StatusOK, building)

	database.Close()
}

func getAllVisits(context *gin.Context) {
	var visits []Visit

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`SELECT * FROM Visit`)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return
	}

	rows, err := preparedStatement.Query()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var tempVisit Visit

		if err := rows.Scan(&tempVisit.ID, &tempVisit.DATE, &tempVisit.CAR); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			fmt.Println(err)
			return
		}

		visits = append(visits, tempVisit)
	}

	context.JSON(http.StatusOK, visits)

	database.Close()
}

func addVisit(context *gin.Context) {
	var visit Visit

	if err := context.BindJSON(&visit); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		fmt.Println(err)
		return
	}

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`INSERT INTO Visit (date, car) VALUES (?, ?)`)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return
	}

	execResult, err := preparedStatement.Exec(&visit.DATE, &visit.CAR)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return
	}

	fmt.Println(execResult)

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	database.Close()
}

func addReport(context *gin.Context) {
	var report Report
	var visit_id int

	if err := context.BindJSON(&report); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		fmt.Println(err)
		return
	}

	database := databaseConnection()

	if err := database.QueryRow(`SELECT id FROM Visit ORDER BY id DESC`).Scan(&visit_id); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal server error",
		})
		fmt.Println(err)
		return
	}

	report.VISIT_ID = visit_id + 1

	fmt.Println(visit_id)
	fmt.Println(report.VISIT_ID)

	preparedStatement, err := database.Prepare(`INSERT INTO Visit (id, date, car) VALUES (?, ?, ?)`)

	if err != nil {
		return
	}

	if _, err := preparedStatement.Exec(&report.DATE, &report.CAR); err != nil {
		return
	}

	preparedStatement, err = database.Prepare(`INSERT INTO Makes_visit (visit_id, employee_id) VALUES (?, ?)`)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal server error",
		})
		fmt.Println(err)
		return
	}

	for i := 0; i < len(report.TEAM_IDS); i++ {
		_, err = preparedStatement.Exec(&report.TEAM_IDS[i])

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			fmt.Println(err)
			return
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	database.Close()

}

func main() {
	config = parseDatbaseConfig("./configs/database.toml")

	router := gin.Default()

	/* Get endpoints */

	// employees
	router.GET("/employee/getAll", getAllEmployees)
	router.GET("/employee/getById", getEmployeeById)
	router.GET("/employee/getByName", getEmployeesByName)

	// clients
	router.GET("/client/getAll", getAllClients)
	router.GET("/client/getById", getClientById)
	router.GET("/client/getByName", getClientByName)

	// building
	router.GET("/building/getAll", getAllBuildings)
	router.GET("/building/getByClientId", getBuildingsByClientId)
	router.GET("/building/getById", getBuildingById)

	// visit
	router.GET("/visit/getAll", getAllVisits)
	// router.GET("/visit/getById", getVisitById)
	// router.GET("/visit/getByEmployee", getVisitByEmployee)

	/* Post endpoints */

	// employees
	router.POST("/employee/add", addEmployee)

	//clients
	router.POST("/client/add", addClient)

	// building
	router.POST("/building/add", addBuilding)

	// visit
	router.POST("/visit/add", addVisit)
	router.POST("/visit/report/add", addReport)

	// report
	router.POST("/report/test", addReport)

	router.Run("localhost:8080")
}
