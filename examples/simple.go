package main

import (
	"database/sql"
	"encoding/json"
	
	"github.com/FoxComm/gizmo/models"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	db, err := sql.Open("postgres", "user=gizmo dbname=gizmo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	of := models.ObjectForm{
		Kind: "product",
		Attributes: map[string]interface{}{
			"abcdef": "a product",
		},
	}

	bytes, err := json.Marshal(of.Attributes)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO object_forms (kind, attributes) VALUES ($1, $2)", of.Kind, bytes)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT kind FROM object_forms")
	if err != nil {
		log.Fatal(err)
	}

	var kind string
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&kind)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("The kind is %s", kind)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
