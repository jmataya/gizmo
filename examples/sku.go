package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/FoxComm/gizmo/illuminate"
	"github.com/FoxComm/gizmo/models"
	"github.com/FoxComm/gizmo/services"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Money struct {
	Value    uint   `json:"value"`
	Currency string `json:"currency"`
}

type SKU struct {
	id       uint
	Title    string `json:"title"`
	Number   int32  `json:"number"`
	UnitCost Money  `json:"unitCost"`
}

func (s *SKU) Identifier() uint {
	return s.id
}

func (s *SKU) SetIdentifier(id uint) {
	s.id = id
}

func main() {
	db, err := sql.Open("postgres", "user=gizmo dbname=gizmo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	om, _ := services.NewObjectManager(db)
	context, err := models.ObjectContext{Name: "default"}.Insert(db)
	if err != nil {
		log.Fatal(err)
	}

	m := Money{199, "USD"}
	s := SKU{1, "Test SKU", 14, m}

	illuminated, err := illuminate.EncodeSimple(&s)
	if err != nil {
		log.Fatal(err)
	}

	illuminated.ContextID = context.ID

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(illuminated); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", body)

	if _, err := om.Create(illuminated); err != nil {
		log.Fatal(err)
	}
}
