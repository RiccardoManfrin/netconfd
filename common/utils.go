package comm

import (
	"encoding/json"
	"reflect"
)

//ListToMap converts a List into a Map by the provided key
func ListToMap(slice interface{}, key string) map[string]interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	if s.IsNil() {
		return nil
	}
	trueSlice := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		trueSlice[i] = s.Index(i).Interface()
	}

	mappedList := make(map[string]interface{})

	for _, l := range trueSlice {
		val := reflect.ValueOf(l)
		kval := val.FieldByName(key).String()
		mappedList[kval] = l
	}
	return mappedList
}

// JSONBytesEqual compares the JSON in two byte slices.
func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}
