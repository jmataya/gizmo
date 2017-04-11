package gizmo

import (
	"testing"

	"github.com/FoxComm/gizmo/models"
	"github.com/FoxComm/gizmo/testutils"

	log "github.com/sirupsen/logrus"
)

type Product struct {
	EntityObject
	Title string
}

func TestCreate(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	db := testutils.InitDB(t)
	defer db.Close()

	view := models.CreateView(t, db)
	product := Product{Title: "Fox Socks"}

	mgr := NewEntityManager(db)
	_, err := mgr.Create(&product, view.ID)
	if err != nil {
		t.Error(err)
	}
}
