package envutil

import "reflect"

// GetStructFieldNames : 구조체의 field 이름 획득
func GetStructFieldNames(inputStructAddr interface{}, skipFields []string) []string {
	var result = []string{}
	val := reflect.ValueOf(inputStructAddr).Elem()
Exit:
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		for _, skipField := range skipFields {
			if skipField == fieldName {
				continue Exit
			}
		}
		result = append(result, fieldName)
	}

	return result
}
