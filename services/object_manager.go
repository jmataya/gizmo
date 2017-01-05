package services

import (
	"database/sql"
	"errors"

	"github.com/FoxComm/gizmo/models"
)

// ObjectManager manages the way that objects get stored in the database.
type ObjectManager struct {
	db *sql.DB
}

// NewObjectManager creates a new ObjectManager.
func NewObjectManager(db *sql.DB) (*ObjectManager, error) {
	if db == nil {
		return nil, errors.New("DB handle must be initialized")
	}

	return &ObjectManager{db: db}, nil
}

// Create saves a new IlluminatedObject.
func (om ObjectManager) Create(illuminated *models.IlluminatedObject) (*models.IlluminatedObject, error) {
	form := models.NewObjectForm(illuminated.Kind)
	shadow := models.NewObjectShadow()

	for name, attribute := range illuminated.Attributes {
		ref, err := form.AddAttribute(attribute.Value)
		if err != nil {
			return nil, err
		}

		shadow.AddAttribute(name, attribute.Type, ref)
	}

	newForm, err := form.Insert(om.db)
	if err != nil {
		return nil, err
	}

	shadow.FormID = newForm.ID
	newShadow, err := shadow.Insert(om.db)
	if err != nil {
		return nil, err
	}

	commit := models.ObjectCommit{FormID: newForm.ID, ShadowID: newShadow.ID}
	newCommit, err := commit.Insert(om.db)
	if err != nil {
		return nil, err
	}

	head := models.ObjectHead{ContextID: illuminated.ContextID, CommitID: newCommit.ID}
	if _, err := head.Insert(om.db); err != nil {
		return nil, err
	}

	return &models.IlluminatedObject{
		ContextID:  illuminated.ContextID,
		FormID:     newForm.ID,
		Kind:       illuminated.Kind,
		Attributes: illuminated.Attributes,
	}, nil
}
