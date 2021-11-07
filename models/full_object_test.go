package models

import (
	"testing"

	"github.com/jmataya/gizmo/testutils"
)

func TestFind_Basic(t *testing.T) {
	assert := testutils.NewAssert(t)
	db := testutils.InitDB(t)
	defer db.Close()

	fullObject := createFullObject(t, db)
	latest, err := fullObject.Find(db, fullObject.Commit.ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(fullObject.Form.ID, latest.Form.ID)
	assert.Equal(fullObject.Shadow.ID, latest.Shadow.ID)
	assert.Equal(fullObject.Commit.ID, latest.Commit.ID)
}
