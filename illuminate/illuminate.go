package illuminate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/FoxComm/gizmo/models"
	log "github.com/sirupsen/logrus"
)

const (
	gizmoTag = "gizmo"
	jsonTag  = "json"
)

// EncodeSimple converts a SimpleObject to an IlluminatedObject.
func EncodeSimple(simple models.SimpleObject) (*models.IlluminatedObject, error) {
	st := reflect.TypeOf(simple)
	val := reflect.ValueOf(simple)

	if st.Kind() == reflect.Ptr {
		st = st.Elem()
		val = val.Elem()
	} else if st.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Expected kind of parameter to be a struct or a pointer, not %v", st)
	}

	illuminated := &models.IlluminatedObject{
		FormID:     simple.Identifier(),
		Kind:       strings.ToLower(st.Name()),
		Attributes: map[string]models.IlluminatedAttribute{},
	}

	log.Debugf("Name of the type to encode %s", st.Name())
	log.Debugln("Discovered public fields:")

	for i := 0; i < val.NumField(); i++ {
		field := st.Field(i)
		sv := val.Field(i)

		if isPublic(field.Name) {
			log.Debugf("Name: %s", fieldName(field))
			log.Debugf("Type: %s", typeName(field))
			log.Debugf("Value: %v", sv.Interface())

			attribute := models.IlluminatedAttribute{
				Type:  typeName(field),
				Value: sv.Interface(),
			}

			illuminated.Attributes[fieldName(field)] = attribute
		}

	}

	return illuminated, nil
}

// Decode converts an IlluminatedObject to a SimpleObject.
func Decode(illuminated *models.IlluminatedObject, simple models.SimpleObject) error {
	log.Debugf("Decoding illuminated object")

	log.Debugf("Cataloging available fields on the simple object")
	simpleFields := make(map[string]string)

	refSimple := reflect.TypeOf(simple).Elem()
	for i := 0; i < refSimple.NumField(); i++ {
		field := refSimple.Field(i)
		log.Debugf(
			"Field name=%s, type=%s, tag=%s, path=%s",
			field.Name,
			field.Type,
			field.Tag,
			field.PkgPath)

		// If PkgPath is not empty, it indicates that the field is private and
		// therefore one that we can't set. Ignore it.
		if field.PkgPath != "" {
			continue
		}

		fieldName, err := illuminatedFieldName(field)
		if err != nil {
			return err
		}

		simpleFields[fieldName] = field.Name
	}

	elem := reflect.ValueOf(simple).Elem()
	if elem.Kind() == reflect.Struct {
		log.Debugf("Decoding attributes on illuminated object")
		for name, attribute := range illuminated.Attributes {
			log.Debugf("Decoding %s, with type=%s and value=%v", name, attribute.Type, attribute.Value)

			realName, ok := simpleFields[name]
			if !ok {
				return fmt.Errorf("Can't match %s to simple object", name)
			}

			simpleField := elem.FieldByName(realName)
			if !simpleField.IsValid() {
				return fmt.Errorf("Can't decode field %s to simple object", name)
			} else if !simpleField.CanSet() {
				return fmt.Errorf("Can't set field %s to simple object", name)
			}

			simpleField.Set(reflect.ValueOf(attribute.Value))
		}

	} else {
		return fmt.Errorf("Invalid kind %v, expected struct", elem.Kind())
	}

	return nil
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
	if str == "" {
		return ""
	}

	bytes := []byte(str)
	first := bytes[0]
	bytes[0] = first | ('a' - 'A')

	return string(bytes)
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

func fieldName(field reflect.StructField) string {
	head := strings.ToLower(field.Name[:1])
	tail := field.Name[1:]
	return head + tail
}

func typeName(field reflect.StructField) string {
	lowered := strings.ToLower(field.Type.Name())
	if lowered == "int32" || lowered == "int64" {
		return "int"
	}
	if lowered == "float32" || lowered == "float64" {
		return "float"
	}
	return lowered
}

func isPublic(field string) bool {
	firstChar, _ := utf8.DecodeRuneInString(field[:1])
	return unicode.IsUpper(firstChar)
}
