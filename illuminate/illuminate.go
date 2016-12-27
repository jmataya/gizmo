package illuminate

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/FoxComm/gizmo/models"
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

	fmt.Printf("Name: %s\n", st.Name())

	fmt.Println("Public Fields")
	fmt.Println("-------------")

	for i := 0; i < val.NumField(); i++ {
		field := st.Field(i)
		sv := val.Field(i)

		if isPublic(field.Name) {
			fmt.Printf("Name: %s\n", fieldName(field))
			fmt.Printf("Type: %s\n", typeName(field))
			fmt.Printf("Value: %v\n", sv.Interface())

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
