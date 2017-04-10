package models

import (
	"testing"

	"github.com/FoxComm/gizmo/testutils"
)

func TestFindLatest_Basic(t *testing.T) {
	assert := testutils.NewAssert(t)

	db := testutils.InitDB(t)
	defer db.Close()

	context := CreateObjectContext(t, db)
	fullObject := createFullObject(t, db, context)
	head := createObjectHead(t, db, context, fullObject.Commit)

	latest, err := fullObject.FindLatest(db, fullObject.Form.ID, head.ContextID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(fullObject.Form.ID, latest.Form.ID)
	assert.Equal(fullObject.Shadow.ID, latest.Shadow.ID)
	assert.Equal(fullObject.Commit.ID, latest.Commit.ID)
}

func TestFindByCommit_Basic(t *testing.T) {
	assert := testutils.NewAssert(t)
	db := testutils.InitDB(t)
	defer db.Close()

	context := CreateObjectContext(t, db)
	fullObject := createFullObject(t, db, context)

	latest, err := fullObject.FindByCommit(db, fullObject.Commit.ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(fullObject.Form.ID, latest.Form.ID)
	assert.Equal(fullObject.Shadow.ID, latest.Shadow.ID)
	assert.Equal(fullObject.Commit.ID, latest.Commit.ID)
}
