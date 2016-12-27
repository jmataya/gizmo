package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/FoxComm/gizmo/illuminate"
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
	fmt.Println("More butter, more cream")

	m := Money{199, "USD"}
	s := SKU{1, "Test SKU", 14, m}

	illuminated, err := illuminate.EncodeSimple(&s)
	if err != nil {
		log.Fatal(err)
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(illuminated); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", body)
}
