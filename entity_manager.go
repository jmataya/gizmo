package gizmo

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jmataya/gizmo/dal"
	"github.com/jmataya/gizmo/models"
	"github.com/gedex/inflector"
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
	Find(id int64, viewID int64, out Entity) error

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

func (d *defaultEntityManager) Find(id int64, viewID int64, out Entity) error {
	return errors.New("Not implemented")
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

	log.Debugln("Getting relations from Entity")
	relations, err := relationsFromEntity(toCreate)
	if err != nil {
		return nil, err
	}

	log.Debugln("Insert the EntityVersion")
	version := models.EntityVersion{
		ContentCommitID: newFullObject.Commit.ID,
		Kind:            newFullObject.Form.Kind,
		Relations:       relations,
	}
	newVersion, err := version.Insert(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	log.Debugf("Inserted EntityVersion with ID=%d", newVersion.ID)

	log.Debugln("Insert the EntityRoot")
	root := models.EntityRoot{Kind: fullObject.Form.Kind}
	newRoot, err := dal.InsertEntityRoot(tx, root)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	log.Debugf("Inserted EntityRoot with ID=%d", newRoot.ID)

	log.Debugln("Insert the EntityHead")
	head := models.EntityHead{
		RootID:    newRoot.ID,
		ViewID:    viewID,
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
	createdEntity := reflect.New(entityType).Interface().(Entity)
	if err = fullToEntity(newFullObject, createdEntity); err != nil {
		tx.Rollback()
		return nil, err
	}

	// FIX ME: Updater should be an object that wraps the entity, not a typecast.
	entityUpdater := createdEntity.(EntityUpdater)

	log.Debugln("Setting relations")

	entityUpdater.SetIdentifier(newRoot.ID)
	entityUpdater.SetCommitID(newVersion.ID)
	entityUpdater.SetViewID(viewID)
	entityUpdater.SetRelations(relations)

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
			// FIX ME: Wrap this in an object, don't do a crappy typecast.
			if err := entity.(EntityUpdater).SetAttribute(name, attrValue); err != nil {
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

	// FIX ME: Wrap this in an object, don't do a crappy typecast.
	return entity.(EntityUpdater).SetKind(full.Form.Kind)
}

func isEntity(fieldType reflect.Type) bool {
	entityInterface := reflect.TypeOf((*Entity)(nil)).Elem()

	switch fieldType.Kind() {
	case reflect.Slice:
		return isEntity(fieldType.Elem())
	case reflect.Struct:
		return reflect.PtrTo(fieldType).Implements(entityInterface)
	default:
		return fieldType.Implements(entityInterface)
	}
}

type reflectedFields struct {
	Info  reflect.StructField
	Value reflect.Value
}

func extractEntity(entity Entity) (name string, fields []reflectedFields, err error) {
	name = ""
	fields = []reflectedFields{}
	err = nil

	entityType := reflect.TypeOf(entity)
	entityVal := reflect.ValueOf(entity)

	if entityType.Kind() == reflect.Ptr {
		log.Debugln("Pointer type detected for entity, getting underlying element")
		entityType = entityType.Elem()
		entityVal = entityVal.Elem()
	} else if entityType.Kind() != reflect.Ptr {
		err = fmt.Errorf("Expected kind of parameter to be a struct or pointer, not %v", entityType)
		return
	}

	name = entityType.Name()
	log.Debugf("Name of the type whose fields are being extracted is %s", name)

	for i := 0; i < entityVal.NumField(); i++ {
		field := reflectedFields{
			Info:  entityType.Field(i),
			Value: entityVal.Field(i),
		}

		fields = append(fields, field)
	}

	return
}

func relationsFromEntity(entity Entity) (models.EntityRelations, error) {
	_, fields, err := extractEntity(entity)
	if err != nil {
		return nil, err
	}

	log.Debugln("Discovering relations")

	relations := map[string][]int64{}
	for _, field := range fields {
		if fieldIsPublic(field.Info) && !field.Info.Anonymous && isEntity(field.Info.Type) {
			fieldName := strings.ToLower(inflector.Singularize(field.Info.Name))
			fieldValue := field.Value.Interface()

			log.Debugf("Found relation %s with value %+v", fieldName, fieldValue)

			switch field.Value.Kind() {
			case reflect.Slice:
				entitySlice := reflect.ValueOf(fieldValue)
				for i := 0; i < entitySlice.Len(); i++ {
					relations[fieldName], err = appendToCommitList(relations[fieldName], entitySlice.Index(i))
					if err != nil {
						return nil, err
					}
				}
			case reflect.Struct:
				relations[fieldName], err = appendToCommitList(relations[fieldName], field.Value)
				if err != nil {
					return nil, err
				}
			case reflect.Ptr:
				entityValue := field.Value.Elem()
				relations[fieldName], err = appendToCommitList(relations[fieldName], entityValue)
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("Unexpected relation type %v", field.Value.Kind())
			}
		}
	}

	return relations, nil
}

func appendToCommitList(commits []int64, value reflect.Value) ([]int64, error) {
	entity, ok := value.Interface().(Entity)
	if !ok {
		return nil, errors.New("Cannot convert value to Entity")
	}

	if commits == nil {
		commits = []int64{}
	}

	return append(commits, entity.CommitID()), nil
}

func entityToFull(entity Entity) (*models.FullObject, error) {
	name, fields, err := extractEntity(entity)
	if err != nil {
		return nil, err
	}

	kind := strings.ToLower(name)
	form := models.NewObjectForm(kind)
	shadow := models.NewObjectShadow()

	log.Debugln("Discovering public fields")

	for _, field := range fields {
		if fieldIsPublic(field.Info) && !field.Info.Anonymous && !isEntity(field.Info.Type) {
			fName := fieldName(field.Info)
			fType := typeName(field.Info.Type)
			fVal := field.Value.Interface()

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
