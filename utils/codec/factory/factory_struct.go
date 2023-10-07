package factory

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unsafe"

	"github.com/liwei1dao/lego/utils/codec/codecore"
	"github.com/liwei1dao/lego/utils/codec/utils"

	"github.com/modern-go/reflect2"
)

type StructDescriptor struct {
	Type   reflect2.Type
	Fields []*Binding
}

type Binding struct {
	levels    []int
	Field     reflect2.StructField
	FromNames []string
	ToNames   []string
	Encoder   codecore.IEncoder
	Decoder   codecore.IDecoder
}

func encoderOfStruct(ctx codecore.ICtx, typ reflect2.Type) codecore.IEncoder {
	type bindingTo struct {
		binding *Binding
		toName  string
		ignored bool
	}
	orderedBindings := []*bindingTo{}
	structDescriptor := describeStruct(ctx, typ)
	for _, binding := range structDescriptor.Fields {
		for _, toName := range binding.ToNames {
			new := &bindingTo{
				binding: binding,
				toName:  toName,
			}
			for _, old := range orderedBindings {
				if old.toName != toName {
					continue
				}
				old.ignored, new.ignored = resolveConflictBinding(ctx.Config(), old.binding, new.binding)
			}
			orderedBindings = append(orderedBindings, new)
		}
	}
	if len(orderedBindings) == 0 {
		return &emptyStructEncoder{}
	}
	finalOrderedFields := []structFieldTo{}
	for _, bindingTo := range orderedBindings {
		if !bindingTo.ignored {
			finalOrderedFields = append(finalOrderedFields, structFieldTo{
				encoder: bindingTo.binding.Encoder.(*structFieldEncoder),
				toName:  bindingTo.toName,
			})
		}
	}
	return &structEncoder{typ, finalOrderedFields}
}

func decoderOfStruct(ctx codecore.ICtx, typ reflect2.Type) codecore.IDecoder {
	bindings := map[string]*Binding{}
	structDescriptor := describeStruct(ctx, typ)
	for _, binding := range structDescriptor.Fields {
		for _, fromName := range binding.FromNames {
			old := bindings[fromName]
			if old == nil {
				bindings[fromName] = binding
				continue
			}
			ignoreOld, ignoreNew := resolveConflictBinding(ctx.Config(), old, binding)
			if ignoreOld {
				delete(bindings, fromName)
			}
			if !ignoreNew {
				bindings[fromName] = binding
			}
		}
	}
	fields := map[string]*structFieldDecoder{}
	for k, binding := range bindings {
		fields[k] = binding.Decoder.(*structFieldDecoder)
	}

	if !ctx.Config().CaseSensitive {
		for k, binding := range bindings {
			if _, found := fields[strings.ToLower(k)]; !found {
				fields[strings.ToLower(k)] = binding.Decoder.(*structFieldDecoder)
			}
		}
	}
	return &structDecoder{ctx.Config(), typ, fields}
}

//结构第编辑码构建
func describeStruct(ctx codecore.ICtx, typ reflect2.Type) *StructDescriptor {
	structType := typ.(*reflect2.UnsafeStructType)
	embeddedBindings := []*Binding{}
	bindings := []*Binding{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if !utils.IsExported(field.Name()) { //内部字段不处理
			continue
		}
		tag, hastag := field.Tag().Lookup(ctx.Config().TagKey)
		if ctx.Config().OnlyTaggedField && !hastag && !field.Anonymous() {
			continue
		}
		if tag == "-" || field.Name() == "_" {
			continue
		}
		tagParts := strings.Split(tag, ",")
		if field.Anonymous() && (tag == "" || tagParts[0] == "") {
			if field.Type().Kind() == reflect.Struct {
				structDescriptor := describeStruct(ctx, field.Type())
				for _, binding := range structDescriptor.Fields {
					binding.levels = append([]int{i}, binding.levels...)
					omitempty := binding.Encoder.(*structFieldEncoder).omitempty
					binding.Encoder = &structFieldEncoder{field, binding.Encoder, omitempty}
					binding.Decoder = &structFieldDecoder{field, binding.Decoder}
					embeddedBindings = append(embeddedBindings, binding)
				}
				continue
			} else if field.Type().Kind() == reflect.Ptr {
				ptrType := field.Type().(*reflect2.UnsafePtrType)
				if ptrType.Elem().Kind() == reflect.Struct {
					structDescriptor := describeStruct(ctx, ptrType.Elem())
					for _, binding := range structDescriptor.Fields {
						binding.levels = append([]int{i}, binding.levels...)
						omitempty := binding.Encoder.(*structFieldEncoder).omitempty
						binding.Encoder = &dereferenceEncoder{binding.Encoder}
						binding.Encoder = &structFieldEncoder{field, binding.Encoder, omitempty}
						binding.Decoder = &dereferenceDecoder{ptrType.Elem(), binding.Decoder}
						binding.Decoder = &structFieldDecoder{field, binding.Decoder}
						embeddedBindings = append(embeddedBindings, binding)
					}
					continue
				}
			}
		}
		fieldNames := calcFieldNames(field.Name(), tagParts[0], tag)
		decoder := DecoderOfType(ctx.Append(field.Name()), field.Type())
		encoder := EncoderOfType(ctx.Append(field.Name()), field.Type())
		binding := &Binding{
			Field:     field,
			FromNames: fieldNames,
			ToNames:   fieldNames,
			Decoder:   decoder,
			Encoder:   encoder,
		}
		binding.levels = []int{i}
		bindings = append(bindings, binding)
	}
	return createStructDescriptor(ctx, typ, bindings, embeddedBindings)
}

func createStructDescriptor(ctx codecore.ICtx, typ reflect2.Type, bindings []*Binding, embeddedBindings []*Binding) *StructDescriptor {
	structDescriptor := &StructDescriptor{
		Type:   typ,
		Fields: bindings,
	}
	processTags(structDescriptor, ctx.Config())
	allBindings := sortableBindings(append(embeddedBindings, structDescriptor.Fields...))
	sort.Sort(allBindings)
	structDescriptor.Fields = allBindings
	return structDescriptor
}

func processTags(structDescriptor *StructDescriptor, config *codecore.Config) {
	for _, binding := range structDescriptor.Fields {
		shouldOmitEmpty := false
		tagParts := strings.Split(binding.Field.Tag().Get(config.TagKey), ",")
		for _, tagPart := range tagParts[1:] {
			if tagPart == "omitempty" {
				shouldOmitEmpty = true
			}
		}
		binding.Decoder = &structFieldDecoder{binding.Field, binding.Decoder}
		binding.Encoder = &structFieldEncoder{binding.Field, binding.Encoder, shouldOmitEmpty}
	}
}

func calcFieldNames(originalFieldName string, tagProvidedFieldName string, wholeTag string) []string {
	// ignore?
	if wholeTag == "-" {
		return []string{}
	}
	// rename?
	var fieldNames []string
	if tagProvidedFieldName == "" {
		fieldNames = []string{originalFieldName}
	} else {
		fieldNames = []string{tagProvidedFieldName}
	}
	// private?
	isNotExported := unicode.IsLower(rune(originalFieldName[0])) || originalFieldName[0] == '_'
	if isNotExported {
		fieldNames = []string{}
	}
	return fieldNames
}

func resolveConflictBinding(config *codecore.Config, old, new *Binding) (ignoreOld, ignoreNew bool) {
	newTagged := new.Field.Tag().Get(config.TagKey) != ""
	oldTagged := old.Field.Tag().Get(config.TagKey) != ""
	if newTagged {
		if oldTagged {
			if len(old.levels) > len(new.levels) {
				return true, false
			} else if len(new.levels) > len(old.levels) {
				return false, true
			} else {
				return true, true
			}
		} else {
			return true, false
		}
	} else {
		if oldTagged {
			return true, false
		}
		if len(old.levels) > len(new.levels) {
			return true, false
		} else if len(new.levels) > len(old.levels) {
			return false, true
		} else {
			return true, true
		}
	}
}

//结构对象 编解码-----------------------------------------------------------------------------------------------------------------------
type structEncoder struct {
	typ    reflect2.Type
	fields []structFieldTo
}

func (codec *structEncoder) GetType() reflect.Kind {
	return reflect.Struct
}
func (this *structEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteObjectStart()
	isNotFirst := false
	for _, field := range this.fields {
		if field.encoder.omitempty && field.encoder.IsEmpty(ptr) {
			continue
		}
		if field.encoder.IsEmbeddedPtrNil(ptr) {
			continue
		}
		if isNotFirst {
			stream.WriteMemberSplit()
		}
		stream.WriteObjectFieldName(field.toName)
		field.encoder.Encode(ptr, stream)
		isNotFirst = true
	}
	stream.WriteObjectEnd()
	if stream.Error() != nil && stream.Error() != io.EOF {
		stream.SetErr(fmt.Errorf("%v.%s", this.typ, stream.Error().Error()))
	}
}

func (this *structEncoder) EncodeToMapJson(ptr unsafe.Pointer, w codecore.IWriter) (ret map[string]string, err error) {
	ret = make(map[string]string)
	for _, field := range this.fields {
		if field.encoder.omitempty && field.encoder.IsEmpty(ptr) {
			continue
		}
		if field.encoder.IsEmbeddedPtrNil(ptr) {
			continue
		}
		w.Reset(nil)
		field.encoder.Encode(ptr, w)
		if w.Error() != nil && w.Error() != io.EOF {
			err = w.Error()
			return
		}
		ret[field.toName] = string(w.Buffer())
	}
	return
}

func (this *structEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type structDecoder struct {
	config *codecore.Config
	typ    reflect2.Type
	fields map[string]*structFieldDecoder
}

func (codec *structDecoder) GetType() reflect.Kind {
	return reflect.Struct
}
func (this *structDecoder) Decode(ptr unsafe.Pointer, r codecore.IReader) {
	if !r.ReadObjectStart() {
		return
	}
	if r.CheckNextIsObjectEnd() { //空对象直接跳出
		// r.ReadObjectEnd()
		return
	}

	this.decodeField(ptr, r)
	for r.ReadMemberSplit() {
		this.decodeField(ptr, r)
	}
	if r.Error() != nil && r.Error() != io.EOF && len(this.typ.Type1().Name()) != 0 {
		r.SetErr(fmt.Errorf("%v.%s", this.typ, r.Error().Error()))
	}
	r.ReadObjectEnd()
}

func (this *structDecoder) decodeField(ptr unsafe.Pointer, r codecore.IReader) {
	var field string
	var fieldDecoder *structFieldDecoder

	field = r.ReadString()
	fieldDecoder = this.fields[field]
	if fieldDecoder == nil && !this.config.CaseSensitive {
		fieldDecoder = this.fields[strings.ToLower(field)]
	}

	if fieldDecoder == nil {
		if this.config.DisallowUnknownFields {
			msg := "found unknown field: " + field
			r.SetErr(fmt.Errorf("decodeField %s", msg))
			return
		}
		if !r.ReadKVSplit() {
			return
		}
		r.Skip() //跳过一个数据单元
		return
	}
	if !r.ReadKVSplit() {
		return
	}
	fieldDecoder.Decode(ptr, r)
}

//解码对象从MapJson 中
func (this *structDecoder) DecodeForMapJson(ptr unsafe.Pointer, r codecore.IReader, extra map[string]string) (err error) {
	var fieldDecoder *structFieldDecoder
	ext := r.GetReader([]byte{}, nil)
	for k, v := range extra {
		fieldDecoder = this.fields[k]
		if fieldDecoder == nil && !this.config.CaseSensitive {
			fieldDecoder = this.fields[strings.ToLower(k)]
		}
		if fieldDecoder == nil {
			if this.config.DisallowUnknownFields {
				err = errors.New("found unknown field: " + k)
				return
			}
			continue
		}
		ext.ResetBytes(StringToBytes(v), nil)
		fieldDecoder.Decode(ptr, ext)
		if ext.Error() != nil && ext.Error() != io.EOF {
			err = ext.Error()
			return
		}
	}
	return
}

//结构对象字段 编解码-----------------------------------------------------------------------------------------------------------------------
type structFieldTo struct {
	encoder *structFieldEncoder
	toName  string
}
type structFieldEncoder struct {
	field        reflect2.StructField
	fieldEncoder codecore.IEncoder
	omitempty    bool
}

func (this *structFieldEncoder) GetType() reflect.Kind {
	return this.fieldEncoder.GetType()
}
func (encoder *structFieldEncoder) Encode(ptr unsafe.Pointer, w codecore.IWriter) {
	fieldPtr := encoder.field.UnsafeGet(ptr)
	encoder.fieldEncoder.Encode(fieldPtr, w)
	if w.Error() != nil && w.Error() != io.EOF {
		w.SetErr(fmt.Errorf("%s: %s", encoder.field.Name(), w.Error().Error()))
	}
}

func (encoder *structFieldEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	fieldPtr := encoder.field.UnsafeGet(ptr)
	return encoder.fieldEncoder.IsEmpty(fieldPtr)
}

func (encoder *structFieldEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	isEmbeddedPtrNil, converted := encoder.fieldEncoder.(codecore.IsEmbeddedPtrNil)
	if !converted {
		return false
	}
	fieldPtr := encoder.field.UnsafeGet(ptr)
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(fieldPtr)
}

type structFieldDecoder struct {
	field        reflect2.StructField
	fieldDecoder codecore.IDecoder
}

func (this *structFieldDecoder) GetType() reflect.Kind {
	return this.fieldDecoder.GetType()
}
func (decoder *structFieldDecoder) Decode(ptr unsafe.Pointer, r codecore.IReader) {
	fieldPtr := decoder.field.UnsafeGet(ptr)
	decoder.fieldDecoder.Decode(fieldPtr, r)
	if r.Error() != nil && r.Error() != io.EOF {
		r.SetErr(fmt.Errorf("%s: %s", decoder.field.Name(), r.Error().Error()))
	}
}

//Empty-----------------------------------------------------------------------------------------------------------------------
type emptyStructEncoder struct {
}

func (codec *emptyStructEncoder) GetType() reflect.Kind {
	return reflect.Struct
}
func (encoder *emptyStructEncoder) Encode(ptr unsafe.Pointer, w codecore.IWriter) {
	w.WriteEmptyObject()
}

func (encoder *emptyStructEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}
