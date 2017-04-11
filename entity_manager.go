package gizmo

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/FoxComm/gizmo/models"
	_ "github.com/lib/pq" // Needed to allow database/sql to use Postgres.
	log "github.com/sirupsen/logrus"
)

const (
	gizmoTag = "gizmo"
	jsonTag  = "json"
)

// EntityManager is the interface for creating, managing, and deleting Entity.
type EntityManager interface {
	// Find retrieves the most recent version of a Entity object within a View.
	// None of the parameters are modified, including the type hint.
	Find(id int64, viewID int64, typeHint Entity) (Entity, error)

	// FindByCommit retrieves a Entity object at a specific commit. This will
	// retrieve the entire object, including all associated objects, as of that
	// commit. None of the parameters are modified, including the type hint.
	FindByCommit(commitID int64, typeHint Entity) (Entity, error)

	// Create saves a new Entity object as a new entity and returns the created
	// version of the object back. If the ID, Entity ID, or Commit ID of the
	// Entity object have previously been set, they will be ignored.
	Create(toCreate Entity, viewID int64) (Entity, error)

	// Update modifies a previously saved Entity object. The new version will be
	// branched from the Entity's Commit ID, and will use the ID and View ID to
	// save in the appropriate format. If the object has not previously been saved
	// the method will error.
	Update(toUpdate Entity) (Entity, error)

	// Delete performs a soft-delete on a Entity object. This must occur at the
	// most recent commit, so the Entity is identified by the ID and View ID.
	Delete(id int64, viewID int64) error
}

// NewEntityManager connects a PostgreSQL database with the supplied connection
// parameters and returns the created EntityManager.
func NewEntityManager(db *sql.DB) EntityManager {
	return &defaultEntityManager{db: db}
}

type defaultEntityManager struct {
	db *sql.DB
}

func (d *defaultEntityManager) Find(id int64, viewID int64, typeHint Entity) (Entity, error) {
	return nil, errors.New("Not implemented")
}

func (d *defaultEntityManager) FindByCommit(commitID int64, typeHint Entity) (Entity, error) {
	return nil, errors.New("Not implemented")
}

func (d *defaultEntityManager) Create(toCreate Entity, viewID int64) (Entity, error) {
	log.Debugln("Starting a transaction for creation")
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	log.Debugln("Converting Entity properties to FullObject")
	fullObject, err := entityToFull(toCreate)
	if err != nil {
		return nil, err
	}

	log.Debugln("Insert the FullObject")
	newFullObject, err := fullObject.Insert(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	log.Debugln("Insert the EntityVersion")
	version := models.EntityVersion{ContentCommitID: newFullObject.Commit.ID}
	newVersion, err := version.Insert(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	log.Debugf("Inserted EntityVersion with ID=%d", newVersion.ID)

	log.Debugln("Insert the EntityRoot")
	root := models.EntityRoot{Kind: fullObject.Form.Kind}
	newRoot, err := root.Insert(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	log.Debugf("Inserted EntityRoot with ID=%d", newRoot.ID)

	log.Debugln("Insert the EntityHead")
	head := models.EntityHead{
		RootID:    newRoot.ID,
		ContextID: viewID,
		VersionID: newVersion.ID,
	}

	newHead, err := head.Insert(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	log.Debugf("Inserted EntityHead with ID=%d", newHead.ID)

	log.Debugln("Convert content back to Entity")
	entityType := reflect.ValueOf(toCreate).Type().Elem()
	fmt.Printf("Entity type = %v\n", entityType)
	createdEntity := reflect.New(entityType).Interface().(Entity)
	if err = fullToEntity(newFullObject, createdEntity); err != nil {
		tx.Rollback()
		return nil, err
	}

	createdEntity.SetIdentifier(newRoot.ID)
	createdEntity.SetCommitID(newVersion.ID)
	createdEntity.SetViewID(viewID)

	return createdEntity, tx.Commit()
}

func (d *defaultEntityManager) Update(toUpdate Entity) (Entity, error) {
	return nil, errors.New("Not implemented")
}

func (d *defaultEntityManager) Delete(id int64, viewID int64) error {
	return errors.New("Not implemented")
}

func fullToEntity(full models.FullObject, entity Entity) error {
	log.Debugln("Converting FullObject to Entity")

	log.Debugln("Cataloging available fields on the Entity")
	entityFields := make(map[string]string)
	refEntity := reflect.TypeOf(entity).Elem()

	for i := 0; i < refEntity.NumField(); i++ {
		field := refEntity.Field(i)
		log.Debugf(
			"Field name=%s, type=%s, path=%s",
			field.Name,
			field.Type,
			field.Tag,
			field.PkgPath)

		// If PkgPath is not empty, it indicates that the field is private and
		// therefore one that we can't set. Ignore it. Also ignore embedded fields.
		// In the future, stop ignoring nested Entity objects.
		if field.PkgPath != "" || field.Anonymous || isEntity(field.Type) {
			log.Debugf("Ignoring %s", field.Name)
			continue
		}

		fieldName, err := illuminatedFieldName(field)
		if err != nil {
			return err
		}

		entityFields[fieldName] = field.Name
	}

	elem := reflect.ValueOf(entity).Elem()
	log.Debugln("Decoding attributes on FullObject")
	for name, attribute := range full.Shadow.Attributes {
		log.Debugf("Decoding %s", name)

		attrValue, ok := full.Form.Attributes[attribute.Ref]
		if !ok {
			return fmt.Errorf("Unable to find Form attribute for %s", name)
		}

		realName, ok := entityFields[name]
		if !ok {
			log.Debugf("Setting custom attribute %s", name)
			if err := entity.SetAttribute(name, attrValue); err != nil {
				return err
			}

			continue
		}

		entityField := elem.FieldByName(realName)
		if !entityField.IsValid() {
			return fmt.Errorf("Can't decode field %s to Entity", name)
		} else if !entityField.CanSet() {
			return fmt.Errorf("Can't set field %s to Entity", name)
		}

		entityField.Set(reflect.ValueOf(attrValue))
	}

	return nil
}

func isEntity(fieldType reflect.Type) bool {
	entityInterface := reflect.TypeOf((*Entity)(nil)).Elem()
	fieldTypePtr := reflect.PtrTo(fieldType)

	return fieldType.Implements(entityInterface) || fieldTypePtr.Implements(entityInterface)
}

func entityToFull(entity Entity) (*models.FullObject, error) {
	entityType := reflect.TypeOf(entity)
	entityVal := reflect.ValueOf(entity)

	if entityType.Kind() == reflect.Ptr {
		log.Debugln("Pointer type detected for entity, getting underlying element")
		entityType = entityType.Elem()
		entityVal = entityVal.Elem()
	} else if entityType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Expected kind of parameter to be a struct or pointer, not %v", entityType)
	}

	log.Debugf("Name of the type to convert is %s", entityType.Name())

	kind := strings.ToLower(entityType.Name())
	form := models.NewObjectForm(kind)
	shadow := models.NewObjectShadow()

	log.Debugln("Discovering public fields")

	for i := 0; i < entityVal.NumField(); i++ {
		fieldInfo := entityType.Field(i)
		fieldVal := entityVal.Field(i)

		if fieldIsPublic(fieldInfo) && !fieldInfo.Anonymous {
			fName := fieldName(fieldInfo)
			fType := typeName(fieldInfo.Type)
			fVal := fieldVal.Interface()

			log.Debugf("Converting Name=%s, Type=%s, Value=%v", fName, fType, fVal)

			log.Debugln("Adding value to form")
			ref, err := form.AddAttribute(fVal)
			if err != nil {
				return nil, err
			}

			log.Debugln("Adding ref to shadow")
			if err := shadow.AddAttribute(fName, fType, ref); err != nil {
				return nil, err
			}
		}
	}

	log.Debugln("Discovering custom attributes")

	for name, val := range entity.Attributes() {
		fType := typeName(reflect.TypeOf(val))

		log.Debugf("Converting Custom Name=%s, Type=%s, Value=%v", name, fType, val)

		log.Debugln("Adding value to form")
		ref, err := form.AddAttribute(val)
		if err != nil {
			return nil, err
		}

		log.Debugln("Adding ref to shadow")
		if err := shadow.AddAttribute(name, fType, ref); err != nil {
			return nil, err
		}
	}

	return &models.FullObject{
		Form:   *form,
		Shadow: *shadow,
	}, nil
}

func fieldIsPublic(field reflect.StructField) bool {
	firstChar, _ := utf8.DecodeRuneInString(field.Name[:1])
	return unicode.IsUpper(firstChar)
}

func fieldName(field reflect.StructField) string {
	return lowercaseFirst(field.Name)
}

func getTagValue(tag, key string) (string, error) {
	if tag == "" {
		return "", nil
	} else if key == "" {
		return "", errors.New("Must specify a non-empty key to retrieve from tag")
	}

	// Append the key with ':"', because they must exist and those characters can
	// only exist in a key.
	formattedKey := fmt.Sprintf("%s:\"", key)

	// Find the start of the value.
	keyStartIdx := strings.Index(tag, formattedKey)
	if keyStartIdx == -1 {
		return "", nil
	}

	valueStartIdx := keyStartIdx + len(formattedKey)

	// Find the next quotation mark (it marks the end of the value)
	valueLength := strings.Index(tag[valueStartIdx:], "\"")
	if valueLength == -1 {
		return "", fmt.Errorf("Malformed tag %s", tag)
	}

	valueEndIdx := valueStartIdx + valueLength
	return tag[valueStartIdx:valueEndIdx], nil
}

func illuminatedFieldName(field reflect.StructField) (string, error) {
	tag := string(field.Tag)

	// Check to see if there's a value 'gizmo'.
	if gizmoVal, err := getTagValue(tag, gizmoTag); err != nil {
		return "", err
	} else if gizmoVal != "" {
		return gizmoVal, nil
	}

	// Check to see if there's a value 'json'.
	if jsonVal, err := getTagValue(tag, jsonTag); err != nil {
		return "", err
	} else if jsonVal != "" {
		return jsonVal, nil
	}

	return lowercaseFirst(field.Name), nil
}

func lowercaseFirst(str string) string {
	head := strings.ToLower(str[:1])
	tail := str[1:]
	return head + tail
}

func typeName(tp reflect.Type) string {
	switch tp.Kind() {
	case reflect.Int32:
		return "int"
	case reflect.Int64:
		return "int"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "float"
	default:
		return strings.ToLower(tp.Name())
	}
}
