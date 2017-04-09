package models

import (
	"database/sql"
	"testing"
)

func createObjectForm(t *testing.T, db *sql.DB) ObjectForm {
	form := ObjectForm{
		Kind: "product",
		Attributes: map[string]interface{}{
			"abcdef": "a product",
		},
	}

	inserted, err := form.Insert(db)
	if err != nil {
		t.Error(err)
	}

	return inserted
}

func createObjectShadow(t *testing.T, db *sql.DB, form ObjectForm) ObjectShadow {
	var ref string
	for key := range form.Attributes {
		ref = key
		break
	}

	shadow := ObjectShadow{
		FormID: form.ID,
		Attributes: map[string]attribute{
			"title": attribute{
				Type: "string",
				Ref:  ref,
			},
		},
	}

	inserted, err := shadow.Insert(db)
	if err != nil {
		t.Error(err)
	}

	return inserted
}

func createObjectCommit(t *testing.T, db *sql.DB, form ObjectForm, shadow ObjectShadow) ObjectCommit {
	commit := ObjectCommit{
		FormID:   form.ID,
		ShadowID: shadow.ID,
	}

	inserted, err := commit.Insert(db)
	if err != nil {
		t.Error(err)
	}

	return inserted
}

func createObjectContext(t *testing.T, db *sql.DB) ObjectContext {
	context := ObjectContext{Name: "Default"}

	inserted, err := context.Insert(db)
	if err != nil {
		t.Error(err)
	}

	return inserted
}

func createObjectHead(t *testing.T, db *sql.DB, context ObjectContext, commit ObjectCommit) ObjectHead {
	head := ObjectHead{
		ContextID: context.ID,
		CommitID:  commit.ID,
	}

	inserted, err := head.Insert(db)
	if err != nil {
		t.Error(err)
	}

	return inserted
}

func createFullObject(t *testing.T, db *sql.DB, context ObjectContext) FullObject {
	form := createObjectForm(t, db)
	shadow := createObjectShadow(t, db, form)
	commit := createObjectCommit(t, db, form, shadow)
	return FullObject{Form: form, Shadow: shadow, Commit: commit}
}
