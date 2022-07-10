package factory

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"
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
	Encoder   core.IEncoder
	Decoder   core.IDecoder
}

func encoderOfStruct(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
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
				old.ignored, new.ignored = resolveConflictBinding(ctx.Options(), old.binding, new.binding)
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

func decoderOfStruct(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	bindings := map[string]*Binding{}
	structDescriptor := describeStruct(ctx, typ)
	for _, binding := range structDescriptor.Fields {
		for _, fromName := range binding.FromNames {
			old := bindings[fromName]
			if old == nil {
				bindings[fromName] = binding
				continue
			}
			ignoreOld, ignoreNew := resolveConflictBinding(ctx.Options(), old, binding)
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

	if !ctx.Options().CaseSensitive {
		for k, binding := range bindings {
			if _, found := fields[strings.ToLower(k)]; !found {
				fields[strings.ToLower(k)] = binding.Decoder.(*structFieldDecoder)
			}
		}
	}
	return createStructDecoder(ctx.Options(), typ, fields)
}

//结构第编辑码构建
func describeStruct(ctx *core.Ctx, typ reflect2.Type) *StructDescriptor {
	structType := typ.(*reflect2.UnsafeStructType)
	embeddedBindings := []*Binding{}
	bindings := []*Binding{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag, hastag := field.Tag().Lookup(ctx.Options().TagKey)
		if ctx.Options().OnlyTaggedField && !hastag && !field.Anonymous() {
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

func createStructDescriptor(ctx *core.Ctx, typ reflect2.Type, bindings []*Binding, embeddedBindings []*Binding) *StructDescriptor {
	structDescriptor := &StructDescriptor{
		Type:   typ,
		Fields: bindings,
	}
	processTags(structDescriptor, ctx.ICodec)
	allBindings := sortableBindings(append(embeddedBindings, structDescriptor.Fields...))
	sort.Sort(allBindings)
	structDescriptor.Fields = allBindings
	return structDescriptor
}

func processTags(structDescriptor *StructDescriptor, codec core.ICodec) {
	for _, binding := range structDescriptor.Fields {
		shouldOmitEmpty := false
		tagParts := strings.Split(binding.Field.Tag().Get(codec.Options().TagKey), ",")
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

func createStructDecoder(opt *core.Options, typ reflect2.Type, fields map[string]*structFieldDecoder) core.IDecoder {
	if opt.DisallowUnknownFields {
		return &structDecoder{typ: typ, fields: fields, disallowUnknownFields: true}
	} else {
		return &structDecoder{opt, typ, fields, false}
	}
}

func resolveConflictBinding(opt *core.Options, old, new *Binding) (ignoreOld, ignoreNew bool) {
	newTagged := new.Field.Tag().Get(opt.TagKey) != ""
	oldTagged := old.Field.Tag().Get(opt.TagKey) != ""
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

func (this *structEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
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

func (this *structEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type structDecoder struct {
	opt                   *core.Options
	typ                   reflect2.Type
	fields                map[string]*structFieldDecoder
	disallowUnknownFields bool
}

func (this *structDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	if !extra.ReadObjectStart() {
		return
	}
	for extra.ReadMemberSplit() {
		this.decodeField(ptr, extra)
	}
	if extra.Error() != nil && extra.Error() != io.EOF && len(this.typ.Type1().Name()) != 0 {
		extra.SetErr(fmt.Errorf("%v.%s", this.typ, extra.Error().Error()))
	}
	extra.ReadObjectEnd()
}

func (this *structDecoder) decodeField(ptr unsafe.Pointer, extra core.IExtractor) {
	var field string
	var fieldDecoder *structFieldDecoder

	field = extra.ReadString()
	fieldDecoder = this.fields[field]
	if fieldDecoder == nil && !this.opt.CaseSensitive {
		fieldDecoder = this.fields[strings.ToLower(field)]
	}

	if fieldDecoder == nil {
		if this.disallowUnknownFields {
			msg := "found unknown field: " + field
			extra.SetErr(fmt.Errorf("decodeField %s", msg))
			return
		}
		if !extra.ReadKVSplit() {
			return
		}
		extra.Skip() //跳过一个数据单元
		return
	}
	if extra.ReadKVSplit() {
		return
	}
	fieldDecoder.Decode(ptr, extra)
}

//结构对象字段 编解码-----------------------------------------------------------------------------------------------------------------------
type structFieldTo struct {
	encoder *structFieldEncoder
	toName  string
}
type structFieldEncoder struct {
	field        reflect2.StructField
	fieldEncoder core.IEncoder
	omitempty    bool
}

func (encoder *structFieldEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	fieldPtr := encoder.field.UnsafeGet(ptr)
	encoder.fieldEncoder.Encode(fieldPtr, stream)
	if stream.Error() != nil && stream.Error() != io.EOF {
		stream.SetErr(fmt.Errorf("%s: %s", encoder.field.Name(), stream.Error().Error()))
	}
}

func (encoder *structFieldEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	fieldPtr := encoder.field.UnsafeGet(ptr)
	return encoder.fieldEncoder.IsEmpty(fieldPtr)
}

func (encoder *structFieldEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	isEmbeddedPtrNil, converted := encoder.fieldEncoder.(core.IsEmbeddedPtrNil)
	if !converted {
		return false
	}
	fieldPtr := encoder.field.UnsafeGet(ptr)
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(fieldPtr)
}

type structFieldDecoder struct {
	field        reflect2.StructField
	fieldDecoder core.IDecoder
}

func (decoder *structFieldDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	fieldPtr := decoder.field.UnsafeGet(ptr)
	decoder.fieldDecoder.Decode(fieldPtr, extra)
	if extra.Error() != nil && extra.Error() != io.EOF {
		extra.SetErr(fmt.Errorf("%s: %s", decoder.field.Name(), extra.Error().Error()))
	}
}

//Empty-----------------------------------------------------------------------------------------------------------------------
type emptyStructEncoder struct {
}

func (encoder *emptyStructEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	stream.WriteEmptyObject()
}

func (encoder *emptyStructEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}
