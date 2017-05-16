package gizmo

import (
	"testing"
)

func TestRelations(t *testing.T) {
	var entity Entity

	var tests = []struct {
		relations map[string][]int64
	}{
		{map[string][]int64{}},
		{map[string][]int64{
			"sku": []int64{1, 2, 3},
		}},
		{map[string][]int64{
			"sku":     []int64{1, 2, 3},
			"product": []int64{1, 2, 3},
		}},
	}

	for _, test := range tests {
		entity = &EntityObject{relations: test.relations}
		want := test.relations
		wantKeys := extractKeys(want)

		got := entity.Relations()
		gotKeys := extractKeys(got)

		if compareStrings(wantKeys, gotKeys) == false {
			t.Errorf("Relation Keys = %q, want %q", gotKeys, wantKeys)
			continue
		}

		for _, key := range wantKeys {
			wantIDs := want[key]
			gotIDs := got[key]

			if compareInts(wantIDs, gotIDs) == false {
				t.Errorf("Relation IDs %s = %q, want %q", key, gotIDs, wantIDs)
			}
		}
	}
}

func TestRelationByEntity(t *testing.T) {
	entity := createEntityObject()

	var tests = []struct {
		entityType string
		want       []int64
		wantErr    string
	}{
		{"product", []int64{4, 5}, ""},
		{"sku", []int64{1, 2, 3}, ""},
		{"", nil, "Entity type must be non-empty"},
		{"variant", []int64{}, ""},
		{"image", []int64{}, ""},
	}

	for _, test := range tests {
		got, err := entity.RelationsByEntity(test.entityType)
		if errorMsg(err) != test.wantErr {
			t.Errorf("RelationsByEntity(%s) = %s, want %s", test.entityType, errorMsg(err), test.wantErr)
		} else if test.wantErr == "" && compareInts(test.want, got) == false {
			t.Errorf("RelationsByEntity(%s) = %q, want %q", test.entityType, got, test.want)
		}
	}
}

func TestSetRelation(t *testing.T) {
	var tests = []struct {
		entityType string
		entityID   int64
		want       []int64
		wantErr    string
	}{
		{"product", 6, []int64{4, 5, 6}, ""},
		{"product", 4, []int64{4, 5}, ""},
		{"sku", 4, []int64{1, 2, 3, 4}, ""},
		{"", 3, nil, "Entity type must be non-empty"},
		{"sku", 0, nil, "Entity ID must be greater than 0"},
		{"variant", 10, []int64{10}, ""},
	}

	for _, test := range tests {
		entity := createEntityObject()
		// FIX ME: Wrap this in an object, don't do a crappy typecast.
		err := entity.(EntityUpdater).SetRelation(test.entityType, test.entityID)
		if errorMsg(err) != test.wantErr {
			t.Errorf(
				"SetRelation(%s, %d) = %s, want %s",
				test.entityType,
				test.entityID,
				errorMsg(err),
				test.wantErr,
			)
			continue
		}

		// If we expected an error, don't try to retrieve.
		if test.wantErr != "" {
			continue
		}

		got, err := entity.RelationsByEntity(test.entityType)
		if err != nil {
			t.Errorf("RelationsByEntity(%s) got error %s, want none", test.entityType, err.Error())
			continue
		}

		if compareInts(test.want, got) == false {
			t.Errorf("RelationsByEntity(%s) = %q, want %q", test.entityType, got, test.want)
		}
	}
}

func TestUpdateRelation(t *testing.T) {
	var tests = []struct {
		entityType string
		oldID      int64
		newID      int64
		want       []int64
		wantErr    string
	}{
		{"sku", 1, 12, []int64{12, 2, 3}, ""},
		{"sku", 3, 4, []int64{1, 2, 4}, ""},
		{"sku", 1, 2, nil, "New mapping to 2 of type sku is already an existing mapping"},
		{"product", 1, 12, nil, "Mapping to 1 of type product is not found"},
		{"variant", 3, 3, nil, "Mapping to 3 of type variant is not found"},
		{"", 1, 2, nil, "Entity type must be non-empty"},
		{"image", 1, 2, nil, "Mapping to 1 of type image is not found"},
		{"sku", 0, 2, nil, "Old ID must be greater than 0"},
		{"product", 1, 0, nil, "New ID must be greater than 0"},
	}

	for _, test := range tests {
		entity := createEntityObject()
		entityType := test.entityType
		oldID := test.oldID
		newID := test.newID

		// FIX ME: Wrap this in an object, don't do a crappy typecast.
		err := entity.(EntityUpdater).UpdateRelation(entityType, oldID, newID)
		if errorMsg(err) != test.wantErr {
			t.Errorf(
				"UpdateRelation(%s, %d, %d) = %s, want %s",
				entityType,
				oldID,
				newID,
				errorMsg(err),
				test.wantErr,
			)
			continue
		}

		// If we expected an error, don't do the next part.
		if test.wantErr != "" {
			continue
		}

		got, err := entity.RelationsByEntity(entityType)
		if err != nil {
			t.Errorf("RelationsByEntity(%s) got error %s, want none", entityType, err.Error())
			continue
		}

		if compareInts(test.want, got) == false {
			t.Errorf("RelationsByEntity(%s) = %q, want %q", entityType, got, test.want)
		}
	}
}

func TestRemoveRelation(t *testing.T) {
	var tests = []struct {
		entityType string
		entityID   int64
		want       []int64
		wantErr    string
	}{
		{"sku", 1, []int64{2, 3}, ""},
		{"product", 4, []int64{5}, ""},
		{"option", 17, []int64{}, ""},
		{"image", 12, nil, "Mapping to 12 of type image is not found"},
		{"product", 55, nil, "Mapping to 55 of type product is not found"},
		{"", 55, nil, "Entity type must be non-empty"},
		{"product", 0, nil, "Entity ID must be greater than 0"},
	}

	for _, test := range tests {
		entity := createEntityObject()
		entityType := test.entityType
		entityID := test.entityID

		err := entity.(EntityUpdater).RemoveRelation(entityType, entityID)
		if errorMsg(err) != test.wantErr {
			t.Errorf(
				"RemoveRelation(%s, %d) = %s, want %s",
				entityType,
				entityID,
				errorMsg(err),
				test.wantErr,
			)
			continue
		}

		// If we expected an error, don't do the next part.
		if test.wantErr != "" {
			continue
		}

		got, err := entity.RelationsByEntity(entityType)
		if err != nil {
			t.Errorf("RelationsByEntity(%s) got error %s, want none", entityType, err.Error())
			continue
		}

		if compareInts(test.want, got) == false {
			t.Errorf("RelationsByEntity(%s) = %q, want %q", entityType, got, test.want)
		}
	}
}

func createEntityObject() Entity {
	return &EntityObject{
		relations: map[string][]int64{
			"sku":     []int64{1, 2, 3},
			"product": []int64{4, 5},
			"variant": []int64{},
			"option":  []int64{17},
		},
	}
}

func errorMsg(err error) string {
	if err != nil {
		return err.Error()
	}

	return ""
}

func extractKeys(relations map[string][]int64) []string {
	i := 0
	keys := make([]string, len(relations))
	for key := range relations {
		keys[i] = key
		i++
	}

	return keys
}

func compareInts(a []int64, b []int64) bool {
	left := sortInts(a)
	right := sortInts(b)
	if len(left) != len(right) {
		return false
	}

	for i := 0; i < len(left); i++ {
		if left[i] != right[i] {
			return false
		}
	}

	return true
}

func sortInts(a []int64) []int64 {
	size := len(a)
	if size < 2 {
		return a
	}

	for i := 0; i < size; i++ {
		for j := size - 1; j >= i+1; j-- {
			if a[j] < a[j-1] {
				a[j], a[j-1] = a[j-1], a[j]
			}
		}
	}

	return a
}

func compareStrings(a []string, b []string) bool {
	left := sortStrings(a)
	right := sortStrings(b)
	if len(left) != len(right) {
		return false
	}

	for i := 0; i < len(left); i++ {
		if left[i] != right[i] {
			return false
		}
	}

	return true
}

func sortStrings(a []string) []string {
	size := len(a)
	if size < 2 {
		return a
	}

	for i := 0; i < size; i++ {
		for j := size - 1; j >= i+1; j-- {
			if a[j] < a[j-1] {
				a[j], a[j-1] = a[j-1], a[j]
			}
		}
	}

	return a
}
