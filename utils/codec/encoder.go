package codec

import (
	"fmt"
	"reflect"
	"time"
)

//编码器
type Encoder struct {
	DefEncoder func(v interface{}) (data []byte, err error) //默认编码 扩展自定义编码方式
}

func (this *Encoder) Encoder(v interface{}) (data []byte, err error) {
	switch v := v.(type) {
	case nil:
		return StringToBytes(""), nil
	case string:
		return StringToBytes(v), nil
	case []byte:
		return v, nil
	case int:
		return IntToBytes(v), nil
	case int8:
		return Int8ToBytes(v), nil
	case int16:
		return Int16ToBytes(v), nil
	case int32:
		return Int32ToBytes(v), nil
	case int64:
		return Int64ToBytes(v), nil
	case uint:
		return UIntToBytes(v), nil
	case uint8:
		return UInt8ToBytes(v), nil
	case uint16:
		return UInt16ToBytes(v), nil
	case uint32:
		return UInt32ToBytes(v), nil
	case uint64:
		return UInt64ToBytes(v), nil
	case float32:
		return Float32ToBytes(v), nil
	case float64:
		return Float64ToBytes(v), nil
	case bool:
		return BoolToBytes(v), nil
	case time.Time:
		return v.AppendFormat([]byte{}, time.RFC3339Nano), nil
	case time.Duration:
		return Int64ToBytes(v.Nanoseconds()), nil
	default:
		if this.DefEncoder != nil {
			return this.DefEncoder(v)
		} else {
			return nil, fmt.Errorf(
				"encoder: can't marshal %T (implement encoder.DefEncoder)", v)
		}
	}
}

func (this *Encoder) EncoderToMap(v interface{}) (data map[string][]byte, err error) {
	vof := reflect.ValueOf(v)
	if !vof.IsValid() {
		return nil, fmt.Errorf("Encoder: EncoderToMap(nil)")
	}
	if vof.Kind() != reflect.Ptr && vof.Kind() != reflect.Map && vof.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Encoder: EncoderToMap(non-pointer %T)", v)
	}
	// vof = vof.Elem()
	data = make(map[string][]byte)
	if vof.Kind() == reflect.Map {
		keys := vof.MapKeys()
		for _, k := range keys {
			value := vof.MapIndex(k)
			var keydata []byte
			var valuedata []byte
			if keydata, err = this.Encoder(k.Interface()); err != nil {
				return nil, fmt.Errorf("Encoder: EncoderToMap(invalid type %T) err:%v", value.Interface(), err)
			}
			if valuedata, err = this.Encoder(value.Interface()); err != nil {
				return nil, fmt.Errorf("Encoder: EncoderToMap(invalid type %T) err:%v", value.Interface(), err)
			}
			data[BytesToString(keydata)] = valuedata
		}
		return
	} else if vof.Kind() == reflect.Slice {
		for i := 0; i < vof.Len(); i++ {
			value := vof.Index(i).Interface()
			var valuedata []byte
			if valuedata, err = this.Encoder(value); err != nil {
				return nil, fmt.Errorf("Encoder: EncoderToMap(invalid type %T) err:%v", value, err)
			}
			data[BytesToString(IntToBytes(i))] = valuedata
		}
		return
	} else if vof.Kind() == reflect.Ptr {
		elem := vof.Elem()
		relType := elem.Type()
		for i := 0; i < relType.NumField(); i++ {
			fieldInfo := relType.Field(i)
			tag := fieldInfo.Tag
			name := tag.Get("json")
			if len(name) == 0 {
				name = fieldInfo.Name
			}
			field := elem.Field(i).Interface()
			var valuedata []byte
			if valuedata, err = this.Encoder(field); err != nil {
				return nil, fmt.Errorf("Encoder: EncoderToMap(invalid type %T) err:%v", field, err)
			}
			data[name] = valuedata
		}
		return
	} else {
		return nil, fmt.Errorf("Encoder: EncoderToMap(invalid type %T)", v)
	}
}

//redis 存储编码 针对 string 转换
func (this *Encoder) EncoderString(v interface{}) (data string, err error) {
	switch v := v.(type) {
	case nil:
		return "", nil
	case string:
		return v, nil
	case []byte:
		return BytesToString(v), nil
	case int:
		return IntToString(int64(v)), nil
	case int8:
		return IntToString(int64(v)), nil
	case int16:
		return IntToString(int64(v)), nil
	case int32:
		return IntToString(int64(v)), nil
	case int64:
		return IntToString(int64(v)), nil
	case uint:
		return UintToString(uint64(v)), nil
	case uint8:
		return UintToString(uint64(v)), nil
	case uint16:
		return UintToString(uint64(v)), nil
	case uint32:
		return UintToString(uint64(v)), nil
	case uint64:
		return UintToString(uint64(v)), nil
	case float32:
		return UintToString(uint64(v)), nil
	case float64:
		return UintToString(uint64(v)), nil
	case bool:
		if v {
			return IntToString(1), nil
		}
		return IntToString(0), nil
	case time.Time:
		return BytesToString(v.AppendFormat([]byte{}, time.RFC3339Nano)), nil
	case time.Duration:
		return IntToString(v.Nanoseconds()), nil
	default:
		if this.DefEncoder != nil {
			var b []byte
			b, err = this.DefEncoder(v)
			if err != nil {
				return "", err
			}
			return BytesToString(b), nil
		} else {
			return "", fmt.Errorf(
				"encoder: can't marshal %T (implement encoder.DefEncoder)", v)
		}
	}
}

//redis 存储编码 针对 string 转换
func (this *Encoder) EncoderToMapString(v interface{}) (data map[string]string, err error) {
	vof := reflect.ValueOf(v)
	if !vof.IsValid() {
		return nil, fmt.Errorf("Encoder: EncoderToMap(nil)")
	}
	if vof.Kind() != reflect.Ptr && vof.Kind() != reflect.Map && vof.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Encoder: EncoderToMap(non-pointer %T)", v)
	}
	// vof = vof.Elem()
	data = make(map[string]string)
	if vof.Kind() == reflect.Map {
		keys := vof.MapKeys()
		for _, k := range keys {
			value := vof.MapIndex(k)
			var keydata string
			var valuedata string
			if keydata, err = this.EncoderString(k.Interface()); err != nil {
				return nil, fmt.Errorf("Encoder: EncoderToMap(invalid type %T) err:%v", value.Interface(), err)
			}
			if valuedata, err = this.EncoderString(value.Interface()); err != nil {
				return nil, fmt.Errorf("Encoder: EncoderToMap(invalid type %T) err:%v", value.Interface(), err)
			}
			data[keydata] = valuedata
		}
		return
	} else if vof.Kind() == reflect.Slice {
		for i := 0; i < vof.Len(); i++ {
			value := vof.Index(i).Interface()
			var valuedata string
			if valuedata, err = this.EncoderString(value); err != nil {
				return nil, fmt.Errorf("Encoder: EncoderToMap(invalid type %T) err:%v", value, err)
			}
			data[IntToString(int64(i))] = valuedata
		}
		return
	} else if vof.Kind() == reflect.Ptr {
		elem := vof.Elem()
		relType := elem.Type()
		for i := 0; i < relType.NumField(); i++ {
			fieldInfo := relType.Field(i)
			tag := fieldInfo.Tag
			name := tag.Get("json")
			if len(name) == 0 {
				name = fieldInfo.Name
			}
			field := elem.Field(i).Interface()
			var valuedata string
			if valuedata, err = this.EncoderString(field); err != nil {
				return nil, fmt.Errorf("Encoder: EncoderToMap(invalid type %T) err:%v", field, err)
			}
			data[name] = valuedata
		}
		return
	} else {
		return nil, fmt.Errorf("Encoder: EncoderToMap(invalid type %T)", v)
	}
}

//redis 存储编码 针对 string 转换
func (this *Encoder) EncoderToSliceString(v interface{}) (data []string, err error) {
	vof := reflect.ValueOf(v)
	if !vof.IsValid() {
		return nil, fmt.Errorf("Encoder: EncoderToSliceString(nil)")
	}
	if vof.Kind() != reflect.Ptr && vof.Kind() != reflect.Map && vof.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Encoder: EncoderToSliceString(non-pointer %T)", v)
	}
	// vof = vof.Elem()
	if vof.Kind() == reflect.Slice {
		for i := 0; i < vof.Len(); i++ {
			value := vof.Index(i).Interface()
			var valuedata string
			if valuedata, err = this.EncoderString(value); err != nil {
				return nil, fmt.Errorf("Encoder: EncoderToSliceString(invalid type %T) err:%v", value, err)
			}
			data[i] = valuedata
		}
		return
	} else {
		return nil, fmt.Errorf("Encoder: EncoderToSliceString(invalid type %T)", v)
	}
}
