package comm

import (
	"encoding/json"
	"reflect"
)

//MatchComparator takes a left and right element and compare them.
//In case of matching true is returned, false otherwise
type MatchComparator func(l interface{}, r interface{}) bool

//ListCompare compare two lists and returns the first non matching item of the subset first list argument or nil if
//All subset elements were matched
//In case the superset and subset have same length this corresponds to the equality condition
//(based on the match comparator)
func ListCompare(subsetSlice interface{}, supersetSlice interface{}, match MatchComparator) interface{} {
	subsetTrueSlice := SliceInterfaceToSlice(subsetSlice)
	supersetTrueSlice := SliceInterfaceToSlice(supersetSlice)
	matched := false
	for _, l := range subsetTrueSlice {
		for _, r := range supersetTrueSlice {
			if match(l, r) {
				matched = true
				break
			}
		}
		if matched == true {
			matched = false
			continue
		} else {
			return l
		}
	}
	return nil
}

//SliceInterfaceToSlice Convert a Slice interface to an actual explicit slide
func SliceInterfaceToSlice(slice interface{}) []interface{} {
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
	return trueSlice
}

//ListToMap converts a List into a Map by the provided key
func ListToMap(slice interface{}, key string) map[string]interface{} {
	trueSlice := SliceInterfaceToSlice(slice)

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
