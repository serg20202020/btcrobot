package util

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

func toString(s interface{}) string {
	if v, ok := s.(string); ok {
		return v
	}
	return fmt.Sprintf("%v", s)
}

func toFloat(s interface{}) float64 {
	var ret float64
	switch v := s.(type) {
	case float64:
		ret = v
	case int64:
		ret = float64(v)
	case string:
		ret, _ = strconv.ParseFloat(v, 64)
	}
	return ret
}

func float2str(i float64) string {
	return strconv.FormatFloat(i, 'f', -1, 64)
}

// 将url.Values（表单数据）转换为Model（struct）
func ConvertAssign(dest interface{}, form url.Values) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return fmt.Errorf("convertAssign(non-pointer %s)", destType)
	}
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	if destValue.Kind() != reflect.Struct {
		return fmt.Errorf("convertAssign(non-struct %s)", destType)
	}
	destType = destValue.Type()
	fieldNum := destType.NumField()
	for i := 0; i < fieldNum; i++ {
		// struct 字段的反射类型（StructField）
		fieldType := destType.Field(i)
		// 非导出字段不处理
		if fieldType.PkgPath != "" {
			continue
		}
		tag := fieldType.Tag.Get("json")
		fieldValue := destValue.Field(i)
		val := form.Get(tag)
		// 字段本身的反射类型（field type）
		fieldValType := fieldType.Type
		switch fieldValType.Kind() {
		case reflect.Int:
			if len(form[tag]) > 1 {
				// TODO:多个值如何处理？
			}
			if val == "" {
				continue
			}
			tmp, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			fieldValue.SetInt(int64(tmp))
		case reflect.String:
			if len(form[tag]) > 1 {
				// TODO:多个值如何处理？
			}
			fieldValue.SetString(val)
		default:

		}
	}
	return nil
}

func Struct2Map(dest map[string]interface{}, src interface{}) error {
	if dest == nil {
		return fmt.Errorf("Struct2Map(dest is %v)", dest)
	}
	srcType := reflect.TypeOf(src)
	srcValue := reflect.Indirect(reflect.ValueOf(src))
	if srcValue.Kind() != reflect.Struct {
		return fmt.Errorf("Struct2Map(non-struct %s)", srcType)
	}
	srcType = srcValue.Type()
	fieldNum := srcType.NumField()
	for i := 0; i < fieldNum; i++ {
		// struct 字段的反射类型（StructField）
		fieldType := srcType.Field(i)
		// 非导出字段不处理
		if fieldType.PkgPath != "" {
			continue
		}
		tag := fieldType.Tag.Get("json")
		fieldValue := srcValue.Field(i)
		dest[tag] = fieldValue.Interface()
	}
	return nil
}
