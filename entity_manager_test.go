package gizmo

import (
	"testing"

	"github.com/FoxComm/gizmo/models"
	"github.com/FoxComm/gizmo/testutils"

	log "github.com/sirupsen/logrus"
)

type SKU struct {
	EntityObject
	Price float64
}

type Variant struct {
	EntityObject
	Title string
	SKUs  []SKU
}

type Product struct {
	EntityObject
	Title string
}

func TestCreate(t *testing.T) {
	assert := testutils.NewAssert(t)
	log.SetLevel(log.DebugLevel)

	db := testutils.InitDB(t)
	defer db.Close()

	context := models.CreateObjectContext(t, db)
	product := Product{Title: "Fox Socks"}

	mgr := NewEntityManager(db)
	newProduct, err := mgr.Create(&product, context.ID)
	if err != nil {
		t.Error(err)
	}

	actualTitle := newProduct.(*Product).Title
	assert.Equal(product.Title, actualTitle)

	if newProduct.Identifier() == 0 {
		t.Error("Created ID should be greater than 0")
	}
	if newProduct.CommitID() == 0 {
		t.Error("CommitID should be greater than 0")
	}
	if newProduct.ViewID() == 0 {
		t.Error("ViewID should be greater than 0")
	}
}

func TestCreate_CustomAttributes(t *testing.T) {
	assert := testutils.NewAssert(t)
	log.SetLevel(log.DebugLevel)

	db := testutils.InitDB(t)
	defer db.Close()

	context := models.CreateObjectContext(t, db)
	product := Product{Title: "Fox Socks"}
	if err := product.SetAttribute("description", "A nice pair of socks"); err != nil {
		t.Fatal(err)
	}

	mgr := NewEntityManager(db)
	newProduct, err := mgr.Create(&product, context.ID)
	if err != nil {
		t.Error(err)
		return
	}

	actualDescription, _ := newProduct.Attribute("description")
	assert.Equal("A nice pair of socks", actualDescription)
}

func TestCreate_SimpleAssociation(t *testing.T) {
	assert := testutils.NewAssert(t)
	log.SetLevel(log.DebugLevel)

	db := testutils.InitDB(t)
	defer db.Close()

	context := models.CreateObjectContext(t, db)

	sku := SKU{Price: 999.0}
	mgr := NewEntityManager(db)
	newSKU, err := mgr.Create(&sku, context.ID)
	if err != nil {
		t.Error(err)
		return
	}

	variant := Variant{Title: "Fox Socks"}
	castedSKU := newSKU.(*SKU)
	variant.SKUs = []SKU{*castedSKU}
	newVariant, err := mgr.Create(&variant, context.ID)
	if err != nil {
		t.Error(err)
		return
	}

	actualTitle := newVariant.(*Variant).Title
	assert.Equal(variant.Title, actualTitle)
}
