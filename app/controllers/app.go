package controllers

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) GetConfig() revel.Result {
	// TODO: Finalise this, what do we want to return here
	return c.Render()
}

func connectToDB() (*sql.DB, error) {
	return sql.Open("postgres", "user=chrisbattarbee dbname=climate_control sslmode=disable")
}

func (c App) SubmitData(id int) revel.Result {
	var jsonData map[string]interface{}
	c.Params.BindJSON(&jsonData)

	query := "INSERT INTO data(node_id, timestamp, temperature, humidity, air_quality) VALUES ("
	query += strconv.Itoa(id)
	query += ", '"
	query += jsonData["timestamp"].(string)
	query += "', "
	query += strconv.FormatFloat(jsonData["temperature"].(float64), 'f', 6, 64)
	query += ", "
	query += strconv.FormatFloat(jsonData["humidity"].(float64), 'f', 6, 64)
	query += ", "
	query += strconv.FormatFloat(jsonData["air_qual"].(float64), 'f', 6, 64)
	query += ");"

	db, err := connectToDB()
	if err != nil {
		log.Println("Could not connect to postgres database, error: " + err.Error())
	}

	_, err = db.Exec(query)

	if err != nil {
		log.Println("Could not execute insert query (add_data)")
		log.Println(err)
	} else {
		log.Println("Clime data entry successfully made from " + strconv.Itoa(id) + " at " + jsonData["timestamp"].(string))
	}
	db.Close()
	return c.Render()
}

func (c App) AddNode() revel.Result {
	var jsonData map[string]interface{}
	c.Params.BindJSON(&jsonData)

	query := "INSERT INTO nodes VALUES ("
	query += strconv.Itoa(int(jsonData["node_id"].(float64)))
	query += ",'"
	query += jsonData["node_location"].(string)
	query += "');"

	db, err := connectToDB()
	if err != nil {
		log.Println("Could not connect to postgres database, error: " + err.Error())
	}

	_, err = db.Exec(query)

	if err != nil {
		log.Println("Could not execute insert query (add node)")
		log.Println(err)
	} else {
		log.Println("Successfully added node with ID: " + strconv.Itoa(int(jsonData["node_id"].(float64))))
	}
	db.Close()
	return c.Render()
}

func (c App) DeleteNode(id int) revel.Result {
	query := "DELETE FROM nodes WHERE node_id=" + strconv.Itoa(id) + ";"
	db, err := connectToDB()
	if err != nil {
		log.Println("Could not connect to postgres database, error: " + err.Error())
	}

	_, err = db.Exec(query)
	if err != nil {
		log.Println("Could not delete node with ID: " + strconv.Itoa(id) + " from the database.")
		log.Println(err)
	} else {
		log.Println("Successfully deleted node with ID: " + strconv.Itoa(id) + " from the database.")
	}

	db.Close()
	return c.Render()
}

func (c App) UpdateLocation(id int) revel.Result {
	var jsonData map[string]interface{}
	c.Params.BindJSON(&jsonData)

	db, err := connectToDB()
	if err != nil {
		log.Println("Could not connect to postgres database, error: " + err.Error())
	}

	query := "UPDATE nodes SET node_location = '"
	query += jsonData["location"].(string) + "' WHERE node_id = " + strconv.Itoa(id) + ";"

	_, err = db.Exec(query)
	if err != nil {
		log.Println("Could not update location of node with ID: " + strconv.Itoa(id))
		log.Println(err)
	} else {
		log.Println("Succesfully update location of node with ID: " + strconv.Itoa(id) + " to " + jsonData["location"].(string))
	}
	db.Close()
	return c.Render()
}
