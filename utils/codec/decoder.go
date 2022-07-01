package codec

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

//解码器
type Decoder struct {
	DefDecoder func(buf []byte, v interface{}) error //默认编码 扩展自定义编码方式
}

func (this *Decoder) Decoder(buf []byte, v interface{}) error {
	switch v := v.(type) {
	case nil:
		return fmt.Errorf("decoder: Decoder(nil)")
	case *string:
		*v = BytesToString(buf)
		return nil
	case *[]byte:
		*v = buf
		return nil
	case *int:
		*v = BytesToInt(buf)
		return nil
	case *int8:
		*v = BytesToInt8(buf)
		return nil
	case *int16:
		*v = BytesToInt16(buf)
		return nil
	case *int32:
		*v = BytesToInt32(buf)
		return nil
	case *int64:
		*v = BytesToInt64(buf)
		return nil
	case *uint:
		*v = BytesToUInt(buf)
		return nil
	case *uint8:
		*v = BytesToUInt8(buf)
		return nil
	case *uint16:
		*v = BytesToUInt16(buf)
		return nil
	case *uint32:
		*v = BytesToUInt32(buf)
		return nil
	case *uint64:
		*v = BytesToUInt64(buf)
		return nil
	case *float32:
		*v = BytesToFloat32(buf)
		return nil
	case *float64:
		*v = BytesToFloat64(buf)
		return nil
	case *bool:
		*v = BytesToBool(buf)
		return nil
	case *time.Time:
		var err error
		*v, err = time.Parse(time.RFC3339Nano, BytesToString(buf))
		return err
	case *time.Duration:
		*v = time.Duration(BytesToInt64(buf))
		return nil
	default:
		if this.DefDecoder != nil {
			return this.DefDecoder(buf, v)
		} else {
			return fmt.Errorf(
				"decoder: can't marshal %T (implement decoder.DefDecoder)", v)
		}
	}
}

func (this *Decoder) DecoderMap(data map[string][]byte, v interface{}) error {
	vof := reflect.ValueOf(v)
	if !vof.IsValid() {
		return fmt.Errorf("Decoder: DecoderMap(nil)")
	}
	if vof.Kind() != reflect.Ptr && vof.Kind() != reflect.Map && vof.Kind() != reflect.Slice {
		return fmt.Errorf("Decoder: DecoderMap(non-pointer %T)", v)
	}
	if vof.Kind() == reflect.Map {
		kt, vt := vof.Type().Key(), vof.Type().Elem()
		for k, v := range data {
			key := reflect.New(kt).Elem()
			if key.Kind() != reflect.Ptr {
				if err := this.Decoder(StringToBytes(k), key.Addr().Interface()); err != nil {
					return fmt.Errorf("Decoder: Decoder(non-pointer %T) err:%v", key, err)
				}
			} else {
				if err := this.Decoder(StringToBytes(k), key.Addr()); err != nil {
					return fmt.Errorf("Decoder: Decoder(non-pointer %T) err:%v", key, err)
				}
			}
			value := reflect.New(vt).Elem()
			if value.Kind() != reflect.Ptr {
				if err := this.Decoder(v, value.Addr().Interface()); err != nil {
					return fmt.Errorf("Decoder: Decoder(non-pointer %T) err:%v", value, err)
				}
			} else {
				value.Interface()
				if err := this.Decoder(v, value.Addr().Interface()); err != nil {
					return fmt.Errorf("Decoder: Decoder(non-pointer %T) err:%v", value, err)
				}
			}
			vof.SetMapIndex(key, value)
		}
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
			if value, ok := data[name]; ok {
				v := reflect.New(fieldInfo.Type).Elem()
				if fieldInfo.Type.Kind() != reflect.Ptr {
					if err := this.Decoder(value, v.Addr().Interface()); err != nil {
						return fmt.Errorf("Decoder: Decoder(non-pointer %T) err:%v", value, err)
					}
					elem.FieldByName(fieldInfo.Name).Set(v)
				} else {
					if err := this.Decoder(value, v); err != nil {
						return fmt.Errorf("Decoder: Decoder(non-pointer %T) err:%v", value, err)
					}
					elem.FieldByName(fieldInfo.Name).Set(v)
				}
			}
		}
	}
	return nil
}

//redis 存储解析 针对 string 转换
func (this *Decoder) DecoderString(buf string, v interface{}) (err error) {
	switch v := v.(type) {
	case nil:
		return fmt.Errorf("decoder: Decoder(nil)")
	case *string:
		*v = buf
		return nil
	case *[]byte:
		*v = StringToBytes(buf)
		return nil
	case *int:
		*v, err = Atoi(StringToBytes(buf))
		return
	case *int8:
		var n int64
		n, err = ParseInt(StringToBytes(buf), 10, 8)
		if err != nil {
			return err
		}
		*v = int8(n)
		return
	case *int16:
		var n int64
		n, err = ParseInt(StringToBytes(buf), 10, 16)
		if err != nil {
			return err
		}
		*v = int16(n)
		return
	case *int32:
		var n int64
		n, err = ParseInt(StringToBytes(buf), 10, 32)
		if err != nil {
			return err
		}
		*v = int32(n)
		return
	case *int64:
		var n int64
		n, err = ParseInt(StringToBytes(buf), 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return
	case *uint:
		var n uint64
		n, err = ParseUint(StringToBytes(buf), 10, 64)
		if err != nil {
			return err
		}
		*v = uint(n)
		return
	case *uint8:
		var n uint64
		n, err = ParseUint(StringToBytes(buf), 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(n)
		return
	case *uint16:
		var n uint64
		n, err = ParseUint(StringToBytes(buf), 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(n)
		return
	case *uint32:
		var n uint64
		n, err = ParseUint(StringToBytes(buf), 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(n)
		return
	case *uint64:
		var n uint64
		n, err = ParseUint(StringToBytes(buf), 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return
	case *float32:
		var f float64
		f, err = strconv.ParseFloat(buf, 32)
		*v = float32(f)
		return
	case *float64:
		*v, err = strconv.ParseFloat(buf, 64)
		return
	case *bool:
		*v, err = strconv.ParseBool(buf)
		return
	case *time.Time:
		var err error
		*v, err = time.Parse(time.RFC3339Nano, buf)
		return err
	case *time.Duration:
		var n int64
		n, err = ParseInt(StringToBytes(buf), 10, 64)
		if err != nil {
			return err
		}
		*v = time.Duration(n)
		return
	default:
		if this.DefDecoder != nil {
			return this.DefDecoder(StringToBytes(buf), v)
		} else {
			return fmt.Errorf(
				"decoder: can't marshal %T (implement decoder.DefDecoder)", v)
		}
	}
}

//redis 存储解析 针对 string 转换
func (this *Decoder) DecoderMapString(data map[string]string, v interface{}) error {
	vof := reflect.ValueOf(v)
	if !vof.IsValid() {
		return fmt.Errorf("Decoder: DecoderMap(nil)")
	}
	if vof.Kind() != reflect.Ptr && vof.Kind() != reflect.Map && vof.Kind() != reflect.Slice {
		return fmt.Errorf("Decoder: DecoderMap(non-pointer %T)", v)
	}
	if vof.Kind() == reflect.Map {
		kt, vt := vof.Type().Key(), vof.Type().Elem()
		for k, v := range data {
			key := reflect.New(kt).Elem()
			if key.Kind() != reflect.Ptr {
				if err := this.DecoderString(k, key.Addr().Interface()); err != nil {
					return fmt.Errorf("Decoder: Decoder(non-pointer %T) err:%v", key, err)
				}
			} else {
				if err := this.DecoderString(k, key.Addr()); err != nil {
					return fmt.Errorf("Decoder: Decoder(non-pointer %T) err:%v", key, err)
				}
			}
			value := reflect.New(vt).Elem()
			if value.Kind() != reflect.Ptr {
				if err := this.DecoderString(v, value.Addr().Interface()); err != nil {
					return fmt.Errorf("Decoder: Decoder(non-pointer %T) err:%v", value, err)
				}
			} else {
				value.Interface()
				if err := this.DecoderString(v, value.Addr().Interface()); err != nil {
					return fmt.Errorf("Decoder: Decoder(non-pointer %T) err:%v", value, err)
				}
			}
			vof.SetMapIndex(key, value)
		}
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
			if value, ok := data[name]; ok {
				v := reflect.New(fieldInfo.Type).Elem()
				if err := this.DecoderString(value, v.Addr().Interface()); err != nil {
					return fmt.Errorf("Decoder: Decoder(non-pointer %T) err:%v", value, err)
				}
				elem.FieldByName(fieldInfo.Name).Set(v)
			}
		}
	}
	return nil
}

//redis 存储解析 针对 string 转换
func (this *Decoder) DecoderSliceString(data []string, v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
	} else {
		panic("Input param is not a slice")
	}

	sl := reflect.ValueOf(v)

	if t.Kind() == reflect.Ptr {
		sl = sl.Elem()
	}

	st := sl.Type()
	fmt.Printf("Slice Type %s:\n", st)

	sliceType := st.Elem()
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
	}
	fmt.Printf("Slice Elem Type %v:\n", sliceType)
	for _, buf := range data {
		if sl.Len() < sl.Cap() {
			sl.Set(sl.Slice(0, sl.Len()+1))
			elem := sl.Index(sl.Len() - 1)
			if elem.IsNil() {
				elem.Set(reflect.New(sliceType))
			}
			this.DecoderString(buf, elem.Elem().Addr().Interface())
			continue
		}
		elem := reflect.New(sliceType)
		sl.Set(reflect.Append(sl, elem))
		this.DecoderString(buf, elem.Elem().Addr().Interface())
	}
	return nil
}
