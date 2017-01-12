package illuminate

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/FoxComm/gizmo/models"
	log "github.com/sirupsen/logrus"
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
