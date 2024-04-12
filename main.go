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
	"golang.org/x/net/context"
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
	Network struct {
		ADDRESS string `toml:"network_address"`
		PORT    string `toml:"network_port"`
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
	ID          int    `json:"id"`
	DATE        string `json:"date"`
	CAR         string `json:"car"`
	CLIENT_ID   int    `json:"client_id"`
	CLIENT_NAME string `json:"client_name"`
}

type Report struct {
	VISIT_ID    int      `json:"id"`
	DATE        string   `json:"date"`
	CAR         string   `json:"car"`
	CLIENT_ID   int      `json:"client_id"`
	CLIENT_NAME string   `json:"client_name"`
	BUILDING_ID int      `json:"building_id"`
	TEAM_IDS    []int    `json:"team_ids"`
	TEAM_NAMES  []string `json:"team_names"`
	ACTIVITY    []struct {
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

func parseError(err error, httpCode int, c *gin.Context) {
	fmt.Errorf("CreateOrder: %v", err)
	c.JSON(httpCode, gin.H{
		"status": "Error",
	})
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

func getAllEmployees(c *gin.Context) {

	var employees []Employee

	database := databaseConnection()

	rows, err := database.Query(`SELECT * FROM Employee`)

	if err != nil {
		parseError(err, http.StatusBadRequest, c)
	}

	for rows.Next() {
		var tempEmployee Employee
		if err := rows.Scan(&tempEmployee.ID, &tempEmployee.NAME, &tempEmployee.PRIVILEGES); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			log.Fatal(err)
			database.Close()
			return
		}

		employees = append(employees, tempEmployee)
	}

	c.JSON(http.StatusOK, employees)

	database.Close()

}

func getEmployeeById(c *gin.Context) {
	var employee Employee
	var employeeId struct {
		ID int `json:"id"`
	}

	if err := c.BindJSON(&employeeId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		return
	}

	db := databaseConnection()

	preparedStatement, err := db.Prepare(`SELECT * FROM Employee WHERE id = ?`)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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

	c.JSON(http.StatusOK, employee)

	db.Close()
}

func addEmployee(c *gin.Context) {
	var employee Employee

	if err := c.BindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	database.Close()

}

func getEmployeesByName(c *gin.Context) {
	var employees []Employee
	var employeeName struct {
		NAME string `json:"name"`
	}

	// fix return
	if err := c.BindJSON(&employeeName); err != nil {
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

	c.JSON(http.StatusOK, employees)

	database.Close()
}

func getAllClients(c *gin.Context) {
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

	c.JSON(http.StatusOK, clients)

	db.Close()
}

func getClientByName(c *gin.Context) {
	var clients []Client
	var clientName struct {
		NAME string `json:"name"`
	}

	if err := c.BindJSON(&clientName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		log.Fatal(err)
		return
	}

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`SELECT * FROM Client WHERE name LIKE ?`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		log.Fatal(err)
		return
	}

	clientName.NAME = "%" + clientName.NAME + "%"

	fmt.Println(clientName.NAME)

	rows, err := preparedStatement.Query(&clientName.NAME)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Errror",
		})
		log.Fatal(err)
		return
	}

	for rows.Next() {
		var tempClient Client
		if err := rows.Scan(&tempClient.ID, &tempClient.NAME); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Errror",
			})
			log.Fatal(err)
			return
		}
		clients = append(clients, tempClient)
	}

	c.JSON(http.StatusOK, clients)

	database.Close()
}

func getClientById(c *gin.Context) {
	var client Client

	if err := c.BindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		log.Fatal(err)
		return
	}

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`SELECT * FROM Client WHERE id = ?`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
	}

	err = preparedStatement.QueryRow(&client.ID).Scan(&client.ID, &client.NAME)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		log.Fatal(err)
		return
	}

	c.JSON(http.StatusOK, client)

	database.Close()
}

func addClient(c *gin.Context) {
	var client Client

	// fix return
	if err := c.BindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		log.Fatal(err)
		return
	}

	_, err = preparedStatement.Exec(&client.NAME)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		log.Fatal(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	database.Close()
}

func getAllBuildings(c *gin.Context) {
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

	c.JSON(http.StatusOK, buildings)

	database.Close()
}

func addBuilding(c *gin.Context) {
	var building Building

	if err := c.BindJSON(&building); err != nil {
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

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	database.Close()
}

func getBuildingsByClientId(c *gin.Context) {
	var buildings []Building
	var clientId Client

	if err := c.BindJSON(&clientId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		fmt.Println(err)
		return
	}

	database := databaseConnection()

	prepredStatement, err := database.Prepare(`SELECT * FROM Building WHERE client_id=?`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return
	}

	rows, err := prepredStatement.Query(&clientId.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var tempBuilding Building
		if err := rows.Scan(&tempBuilding.ID, &tempBuilding.CLIENT_ID, &tempBuilding.ADDRESS, &tempBuilding.STATUS); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			fmt.Println(err)
			return
		}

		buildings = append(buildings, tempBuilding)
	}

	c.JSON(http.StatusOK, buildings)

	database.Close()
}

func getBuildingById(c *gin.Context) {
	var building Building

	if err := c.BindJSON(&building); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		fmt.Println(err)
		return
	}

	database := databaseConnection()

	preparedStatement, err := database.Prepare(`SELECT * FROM Building WHERE id=?`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return
	}

	if err := preparedStatement.QueryRow(&building.ID).Scan(&building.ID, &building.CLIENT_ID, &building.ADDRESS, &building.STATUS); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return

	}

	c.JSON(http.StatusOK, building)

	database.Close()
}

func getAllVisits(c *gin.Context) {
	var visits []Visit

	database := databaseConnection()

	rows, err := database.Query(`SELECT v.id, v.date, v.car, c.id, c.name FROM Visit as v INNER JOIN Recieves_visit AS rv ON rv.visit_id=v.id INNER JOIN Building AS b ON b.id=rv.building_id INNER JOIN Client AS c ON c.id=b.client_id`)

	if err != nil {
		return
	}

	for rows.Next() {
		var tempVisit Visit

		err = rows.Scan(&tempVisit.ID, &tempVisit.DATE, &tempVisit.CAR, &tempVisit.CLIENT_ID, &tempVisit.CLIENT_NAME)

		if err != nil {
			return
		}

		visits = append(visits, tempVisit)
	}

	c.JSON(http.StatusOK, visits)

	database.Close()
}

func getReport(c *gin.Context) {
	var report Report

	fmt.Println(c.Request.Body)

	if err := c.BindJSON(&report); err != nil {
		fmt.Println(report)
		fmt.Println(err)
		return
	}

	database, err := databaseConnection().BeginTx(context.Background(), nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer database.Rollback()

	if err := database.QueryRow(
		`SELECT c.id, c.name, v.id, v.date, v.car, b.id FROM Client AS c
		INNER JOIN Building AS b ON b.client_id=c.id
		INNER JOIN Recieves_visit AS rv ON rv.building_id=b.id
		INNER JOIN Visit AS v ON v.id=rv.visit_id WHERE v.id=?`, &report.VISIT_ID).Scan(&report.CLIENT_ID, &report.CLIENT_NAME, &report.VISIT_ID, &report.DATE, &report.CAR, &report.BUILDING_ID); err != nil {
		fmt.Println(err)
		return
	}

	rows, err := database.Query(`SELECT e.id, e.name FROM Employee as e
	INNER JOIN Makes_visit AS mv ON e.id=mv.employee_id
	INNER JOIN Visit AS v ON mv.visit_id=v.id WHERE v.id=?`, &report.VISIT_ID)

	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.ID, &employee.NAME)

		if err != nil {
			fmt.Println(err)
			return
		}

		report.TEAM_NAMES = append(report.TEAM_NAMES, employee.NAME)
		report.TEAM_IDS = append(report.TEAM_IDS, employee.ID)
	}

	rows, err = database.Query(`SELECT a.id, a.description, o.id, o.description, p.id, p.description FROM Activity AS a INNER JOIN Observation AS o ON a.visit_id=o.visit_id INNER JOIN Pendency AS p ON a.visit_id=p.visit_id WHERE a.visit_id=?`, &report.VISIT_ID)

	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var descriptions struct {
			ACTIVITY struct {
				ID          int
				DESCRIPTION string
			}
			OBSERVATION struct {
				ID          int
				DESCRIPTION string
			}
			PENDENCY struct {
				ID          int
				DESCRIPTION string
			}
		}

		err := rows.Scan(&descriptions.ACTIVITY.ID, &descriptions.ACTIVITY.DESCRIPTION, &descriptions.OBSERVATION.ID, &descriptions.OBSERVATION.DESCRIPTION, &descriptions.PENDENCY.ID, &descriptions.PENDENCY.DESCRIPTION)

		if err != nil {
			fmt.Println(err)
			return
		}

		report.ACTIVITY = append(report.ACTIVITY, struct {
			ID          int    "json:\"activity_id\""
			DESCRIPTION string "json:\"activity_description\""
		}(descriptions.ACTIVITY))

		report.OBSERVATION = append(report.OBSERVATION, struct {
			ID          int    "json:\"observation_id\""
			DESCRIPTION string "json:\"observation_description\""
		}(descriptions.OBSERVATION))

		report.PENDENCY = append(report.PENDENCY, struct {
			ID          int    "json:\"pendency_id\""
			DESCRIPTION string "json:\"pendency_description\""
		}(descriptions.PENDENCY))
	}

	database.Commit()

	c.JSON(http.StatusOK, report)

}

func addReport(c *gin.Context) {
	var report Report
	var visit_id int

	if err := c.BindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		fmt.Println(err)
		return
	}

	db := databaseConnection()

	// Change variable name
	database, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return
	}

	if err := database.QueryRow(`SELECT id FROM Visit ORDER BY id DESC`).Scan(&visit_id); err != nil {
		if visit_id != 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "internal server error",
			})
			fmt.Println(err)
			return
		}
	}
	defer database.Rollback()

	report.VISIT_ID = visit_id + 1

	fmt.Println(visit_id)
	fmt.Println(report.VISIT_ID)

	preparedStatement, err := database.Prepare(`INSERT INTO Visit (date, car) VALUES (?, ?)`)

	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := preparedStatement.Exec(&report.DATE, &report.CAR); err != nil {
		fmt.Println(err)
		return
	}

	database.Commit()

	database, err = db.BeginTx(context.Background(), nil)

	if err != nil {
		return
	}

	preparedStatement, err = database.Prepare(`INSERT INTO Makes_visit (visit_id, employee_id) VALUES (?, ?)`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal server error",
		})
		fmt.Println(err)
		return
	}

	for i := 0; i < len(report.TEAM_IDS); i++ {
		_, err = preparedStatement.Exec(&report.VISIT_ID, &report.TEAM_IDS[i])

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			fmt.Println(err)
			return
		}
	}

	preparedStatement, err = database.Prepare(`INSERT INTO Recieves_visit(visit_id, building_id) VALUES (?, ?)`)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		fmt.Println(err)
		return
	}

	if _, err := preparedStatement.Exec(&report.VISIT_ID, &report.BUILDING_ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
		})
		fmt.Println(err)
		return
	}

	fmt.Println(report)

	preparedStatement, err = database.Prepare(`INSERT INTO Activity(visit_id, description) VALUES (?, ?)`)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		fmt.Println(err)
		return
	}

	for i := 0; i < len(report.ACTIVITY); i++ {
		fmt.Println(&report.ACTIVITY[i])
		if _, err := preparedStatement.Exec(&report.VISIT_ID, &report.ACTIVITY[i].DESCRIPTION); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			fmt.Println(err)
			return
		}
	}

	preparedStatement, err = database.Prepare(`INSERT INTO Observation(visit_id, description) VALUES (?, ?)`)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		fmt.Println(err)
		return
	}

	for i := 0; i < len(report.OBSERVATION); i++ {
		fmt.Println(report.OBSERVATION[i])
		if _, err := preparedStatement.Exec(&report.VISIT_ID, &report.OBSERVATION[i].DESCRIPTION); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			fmt.Println(err)
			return
		}
	}

	preparedStatement, err = database.Prepare(`INSERT INTO Pendency(visit_id, description) VALUES (?, ?)`)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
		})
		fmt.Println(err)
		return
	}

	for i := 0; i < len(report.PENDENCY); i++ {
		fmt.Println(report.PENDENCY[i])
		if _, err := preparedStatement.Exec(&report.VISIT_ID, &report.PENDENCY[i].DESCRIPTION); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
			})
			fmt.Println(err)
			return
		}
	}

	database.Commit()

	db.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func main() {
	config = parseDatbaseConfig("./configs/config.toml")

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
	router.GET("/visit/report/get", getReport)
	// router.GET("/visit/getByEmployee", getVisitByEmployee)

	/* Post endpoints */

	// employees
	router.POST("/employee/add", addEmployee)

	//clients
	router.POST("/client/add", addClient)

	// building
	router.POST("/building/add", addBuilding)

	// visit
	router.POST("/visit/report/add", addReport)

	// report
	router.POST("/report/test", addReport)

	router.Run(config.Network.ADDRESS + ":" + config.Network.PORT)
}
