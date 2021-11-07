package gizmo

import (
	"fmt"
	"testing"

	"github.com/jmataya/gizmo/models"
	"github.com/jmataya/gizmo/testutils"

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

	view := models.CreateView(t, db)
	product := Product{Title: "Fox Socks"}

	mgr := NewEntityManager(db)
	newProduct, err := mgr.Create(&product, view.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal("product", newProduct.Kind())

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

	view := models.CreateView(t, db)
	product := Product{Title: "Fox Socks"}
	if err := product.SetAttribute("description", "A nice pair of socks"); err != nil {
		t.Fatal(err)
	}

	mgr := NewEntityManager(db)
	newProduct, err := mgr.Create(&product, view.ID)
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

	view := models.CreateView(t, db)

	sku := SKU{Price: 999.0}
	mgr := NewEntityManager(db)
	newSKU, err := mgr.Create(&sku, view.ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal("sku", newSKU.Kind())

	variant := Variant{Title: "Fox Socks"}
	castedSKU := newSKU.(*SKU)
	variant.SKUs = []SKU{*castedSKU}
	newVariant, err := mgr.Create(&variant, view.ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal("variant", newVariant.Kind())

	actualTitle := newVariant.(*Variant).Title
	assert.Equal(variant.Title, actualTitle)

	fmt.Printf("Variant: %+v\n", newVariant)

	skus, err := newVariant.RelationsByEntity("sku")
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(1, len(skus))
}

func TestFind(t *testing.T) {
	assert := testutils.NewAssert(t)
	log.SetLevel(log.DebugLevel)

	db := testutils.InitDB(t)
	defer db.Close()

	view := models.CreateView(t, db)
	product := Product{Title: "Fox Socks"}

	mgr := NewEntityManager(db)
	created, err := mgr.Create(&product, view.ID)
	if err != nil {
		t.Fatal(err)
	}

	var findProduct Product
	if err := mgr.Find(created.Identifier(), view.ID, &findProduct); err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal("Fox Socks", findProduct.Title)
}
