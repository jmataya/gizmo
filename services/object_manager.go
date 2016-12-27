package services

import (
	"fmt"

	"github.com/FoxComm/gizmo/models"
)

type ObjectManager struct{}

func (om ObjectManager) Create(illuminated *models.IlluminatedObject) error {
	form := models.NewObjectForm(illuminated.Kind)
	shadow := models.NewObjectShadow()

	for name, attribute := range illuminated.Attributes {
		ref, err := form.AddAttribute(attribute.Value)
		if err != nil {
			return err
		}

		shadow.AddAttribute(name, attribute.Type, ref)
	}

	fmt.Printf("Form: %v\n", form)
	fmt.Printf("Shadow: %v\n", shadow)

	return nil
}
