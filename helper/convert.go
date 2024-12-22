package helper

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func StructToMap(input interface{}) (map[string]any, error) {
	// Validasi bahwa input adalah pointer ke struct atau struct
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference pointer jika diperlukan
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or a pointer to a struct, got %s", val.Kind())
	}

	// Konversi struct ke map[string]any
	result := make(map[string]any)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i) // Mendapatkan definisi field
		fieldValue := val.Field(i)

		// Lewati field yang tidak dapat diakses atau tidak diisi (nil/zero value)
		if !fieldValue.CanInterface() || isZeroValue(fieldValue) {
			continue
		}

		// Tambahkan field ke map
		fieldName := strings.ToLower(field.Name) // Bisa diganti sesuai kebutuhan
		result[fieldName] = fieldValue.Interface()
	}

	return result, nil
}

func StructToMapSlice(input interface{}) ([]map[string]any, error) {
	// Validasi bahwa input adalah slice pointer
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a pointer to a slice")
	}

	slice := val.Elem()
	result := make([]map[string]any, slice.Len())

	// Iterasi setiap elemen dalam slice
	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem() // Dereference pointer jika diperlukan
		}

		// Validasi bahwa item adalah struct
		if item.Kind() != reflect.Struct {
			return nil, fmt.Errorf("slice elements must be structs, got %s", item.Kind())
		}

		// Konversi struct ke map[string]any
		itemMap := make(map[string]any)
		for j := 0; j < item.NumField(); j++ {
			field := item.Type().Field(j) // Mendapatkan definisi field
			fieldValue := item.Field(j)

			// Lewati field yang tidak dapat diakses atau tidak diisi (nil/zero value)
			if !fieldValue.CanInterface() || isZeroValue(fieldValue) {
				continue
			}

			// Tambahkan field ke map
			fieldName := strings.ToLower(field.Name)
			itemMap[fieldName] = fieldValue.Interface()
		}
		result[i] = itemMap
	}

	return result, nil
}

// Fungsi untuk memeriksa apakah suatu nilai adalah nilai default (zero value)
func isZeroValue(value reflect.Value) bool {
	zeroValue := reflect.Zero(value.Type()).Interface()
	return reflect.DeepEqual(value.Interface(), zeroValue)
}

func ConvertFieldInData(data interface{}, fieldName string, targetType string) ([]map[string]interface{}, error) {
	// Pastikan data adalah slice
	dataSlice, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("data is not a slice of interface{}")
	}

	// Slice hasil konversi
	result := make([]map[string]interface{}, len(dataSlice))

	for i, item := range dataSlice {
		// Pastikan setiap item adalah map[string]interface{}
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("data item at index %d is not a map[string]interface{}", i)
		}

		// Cek apakah fieldName ada dan memiliki tipe yang sesuai
		if value, exists := itemMap[fieldName]; exists {
			switch targetType {
			case "int":
				if floatVal, ok := value.(float64); ok {
					itemMap[fieldName] = int(floatVal)
				} else {
					return nil, fmt.Errorf("field '%s' is not a float64", fieldName)
				}
			case "string":
				itemMap[fieldName] = fmt.Sprintf("%v", value)
			case "float64":
				if intVal, ok := value.(int); ok {
					itemMap[fieldName] = float64(intVal)
				} else if strVal, ok := value.(string); ok {
					floatVal, err := strconv.ParseFloat(strVal, 64)
					if err != nil {
						return nil, fmt.Errorf("failed to convert field '%s' to float64: %v", fieldName, err)
					}
					itemMap[fieldName] = floatVal
				} else {
					return nil, fmt.Errorf("field '%s' cannot be converted to float64", fieldName)
				}
			default:
				return nil, fmt.Errorf("unsupported target type: %s", targetType)
			}
		}

		// Tambahkan itemMap ke hasil
		result[i] = itemMap
	}

	return result, nil
}

func ConvertFieldInMap(data interface{}, fieldName string, targetType string) (map[string]interface{}, error) {
	// Pastikan data adalah map[string]interface{}
	itemMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data is not a map[string]interface{}")
	}

	// Cek apakah fieldName ada dan memiliki tipe yang sesuai
	if value, exists := itemMap[fieldName]; exists {
		switch targetType {
		case "int":
			if floatVal, ok := value.(float64); ok {
				itemMap[fieldName] = int(floatVal)
			} else {
				return nil, fmt.Errorf("field '%s' is not a float64", fieldName)
			}
		case "string":
			itemMap[fieldName] = fmt.Sprintf("%v", value)
		case "float64":
			if intVal, ok := value.(int); ok {
				itemMap[fieldName] = float64(intVal)
			} else if strVal, ok := value.(string); ok {
				floatVal, err := strconv.ParseFloat(strVal, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to convert field '%s' to float64: %v", fieldName, err)
				}
				itemMap[fieldName] = floatVal
			} else {
				return nil, fmt.Errorf("field '%s' cannot be converted to float64", fieldName)
			}
		default:
			return nil, fmt.Errorf("unsupported target type: %s", targetType)
		}
	}

	return itemMap, nil
}

func ConvertFieldInSlice(data interface{}, fieldName string, targetType string) (interface{}, error) {
	// Gunakan refleksi untuk memeriksa tipe data
	v := reflect.ValueOf(data)

	// Pastikan data adalah slice
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("data is not a slice")
	}

	// Iterasi elemen-elemen slice
	result := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)

		// Pastikan elemen adalah pointer ke struct atau struct
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		if item.Kind() != reflect.Struct {
			return nil, fmt.Errorf("slice elements must be structs or pointers to structs")
		}

		// Dapatkan field berdasarkan nama
		field := item.FieldByName(fieldName)
		if !field.IsValid() {
			return nil, fmt.Errorf("field '%s' not found in struct", fieldName)
		}
		if !field.CanSet() {
			return nil, fmt.Errorf("field '%s' cannot be set", fieldName)
		}

		// Konversi nilai field sesuai targetType
		switch targetType {
		case "int":
			if field.Kind() == reflect.Float64 {
				field.SetInt(int64(field.Float()))
			} else {
				return nil, fmt.Errorf("field '%s' is not a float64", fieldName)
			}
		case "string":
			field.SetString(fmt.Sprintf("%v", field.Interface()))
		case "float64":
			if field.Kind() == reflect.Int {
				field.SetFloat(float64(field.Int()))
			} else if field.Kind() == reflect.String {
				floatVal, err := strconv.ParseFloat(field.String(), 64)
				if err != nil {
					return nil, fmt.Errorf("failed to convert field '%s' to float64: %v", fieldName, err)
				}
				field.SetFloat(floatVal)
			} else {
				return nil, fmt.Errorf("field '%s' cannot be converted to float64", fieldName)
			}
		default:
			return nil, fmt.Errorf("unsupported target type: %s", targetType)
		}

		// Tambahkan elemen yang telah diubah ke slice hasil
		result.Index(i).Set(item)
	}

	return result.Interface(), nil
}

// ConvertToMap menerima parameter interface{} dan mengubahnya menjadi map[string]interface{}
func ConvertToMap(input interface{}) (map[string]interface{}, error) {
	// Buat map untuk menyimpan hasil konversi
	result := make(map[string]interface{})

	// Dapatkan nilai refleksi dari input
	val := reflect.ValueOf(input)

	// Pastikan input adalah struct atau map
	if val.Kind() == reflect.Struct {
		// Jika input adalah struct, iterasi melalui field-nya
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			result[field.Name] = val.Field(i).Interface()
		}
	} else if val.Kind() == reflect.Map {
		// Jika input adalah map, iterasi melalui elemen-elemennya
		for _, key := range val.MapKeys() {
			result[fmt.Sprintf("%v", key.Interface())] = val.MapIndex(key).Interface()
		}
	} else {
		return nil, fmt.Errorf("input must be a struct or map")
	}

	return result, nil
}
