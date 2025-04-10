package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type EncodeHookFunc func(reflect.Value) (reflect.Value, error)

type InspectorConfig struct {
	TagNames   []string
	EncodeHook EncodeHookFunc
	IncludeNil bool
}

type Inspector struct {
	cfg *InspectorConfig
}

func NewInspector(cfg *InspectorConfig) *Inspector {
	return &Inspector{cfg: cfg}
}

// FieldInfo represents information about a field in a struct
type FieldInfo struct {
	FieldParts []string            // Field name parts
	TagParts   map[string][]string // Tag parts
	Value      reflect.Value       // Actual value of the field in the struct (possibly zeroed)
	Processed  reflect.Value       // Value after processing
	Encoded    string              // Encoded value
	IsNil      bool                // Whether the orignal value was nil
}

func (i *Inspector) Inspect(c any) (map[string]*FieldInfo, error) {
	m := make(map[string]*FieldInfo)

	tagParts := make(map[string][]string)
	for _, tagName := range i.cfg.TagNames {
		tagParts[tagName] = []string{}
	}

	if err := i.inspect(nil, reflect.ValueOf(c), m, tagParts, false, false); err != nil {
		return nil, err
	}
	return m, nil
}

func (i *Inspector) inspect(fieldParts []string, val reflect.Value, m map[string]*FieldInfo, tagParts map[string][]string, isNil, isPtr bool) error {
	var err error
	processedVal := val
	if !isNil && i.cfg.EncodeHook != nil {
		processedVal, err = i.cfg.EncodeHook(val)
		if err != nil {
			return fmt.Errorf("%s: failed to process value: %w", join(fieldParts), err)
		}
	}

	processedKind := getKind(processedVal.Kind())

	switch processedKind {
	case reflect.Bool, reflect.Int, reflect.Uint, reflect.Float32, reflect.String, reflect.Array, reflect.Slice:
		valTyp := val.Type()

		processedVal, encoded, err := i.encodeValue(processedVal)
		if err != nil {
			return fmt.Errorf("%s: failed to encode value: %w", join(fieldParts), err)
		}

		fieldInfo := &FieldInfo{
			FieldParts: fieldParts,
			TagParts:   tagParts,
			Processed:  processedVal,
			Encoded:    encoded,
			IsNil:      isNil,
		}

		if isNil {
			if isPtr {
				fieldInfo.Value = reflect.Zero(reflect.PointerTo(valTyp))
			} else {
				fieldInfo.Value = reflect.Zero(valTyp)
			}
		} else {
			if isPtr {
				if val.CanAddr() {
					fieldInfo.Value = val.Addr()
				} else {
					ptrElem := reflect.New(valTyp)
					ptrElem.Elem().Set(val)
					fieldInfo.Value = ptrElem
				}
			} else {
				fieldInfo.Value = val
			}
		}

		m[join(fieldParts)] = fieldInfo

		return nil
	}

	kind := getKind(processedVal.Kind())
	switch kind {
	case reflect.Struct:
		return i.inspectStruct(fieldParts, val, m, tagParts, isNil)
	case reflect.Ptr:
		if val.IsNil() {
			if !i.cfg.IncludeNil {
				return nil
			}
			val = reflect.New(val.Type().Elem())
			isNil = true
		}
		return i.inspect(fieldParts, val.Elem(), m, tagParts, isNil, true)
	default:
		return fmt.Errorf("%s: unsupported type: %s", join(fieldParts), processedKind)
	}
}

func (i *Inspector) inspectStruct(fieldParts []string, val reflect.Value, m map[string]*FieldInfo, tagParts map[string][]string, isNil bool) error {
	structType := val.Type()
	for idx := range val.NumField() {
		field := structType.Field(idx)
		fieldVal := val.Field(idx)

		fieldName := field.Name

		fieldPartsCopy := make([]string, len(fieldParts))
		copy(fieldPartsCopy, fieldParts)

		tagPartsCopy := make(map[string][]string)
		for k, v := range tagParts {
			tagPartsCopy[k] = make([]string, len(v))
			copy(tagPartsCopy[k], v)
		}

		fieldPartsCopy = append(fieldPartsCopy, fieldName)
		for tagName := range tagParts {
			tagPart := field.Tag.Get(tagName)
			tagPart = strings.SplitN(tagPart, ",", 2)[0]
			tagPartsCopy[tagName] = append(tagPartsCopy[tagName], tagPart)
		}

		if err := i.inspect(fieldPartsCopy, fieldVal, m, tagPartsCopy, isNil, false); err != nil {
			return err
		}
	}

	return nil
}

func join(parts []string) string {
	return strings.Join(parts, ".")
}

func getKind(kind reflect.Kind) reflect.Kind {
	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}

func (i *Inspector) encodeValue(val reflect.Value) (reflect.Value, string, error) {
	processedVal := val
	if i.cfg.EncodeHook != nil {
		var err error
		processedVal, err = i.cfg.EncodeHook(val)
		if err != nil {
			return reflect.Value{}, "", err
		}
	}

	kind := getKind(processedVal.Kind())
	switch kind {
	case reflect.Bool:
		return processedVal, strconv.FormatBool(processedVal.Bool()), nil
	case reflect.Int:
		return processedVal, strconv.FormatInt(processedVal.Int(), 10), nil
	case reflect.Uint:
		return processedVal, strconv.FormatUint(processedVal.Uint(), 10), nil
	case reflect.Float32:
		return processedVal, strconv.FormatFloat(processedVal.Float(), 'f', -1, 32), nil
	case reflect.String:
		return processedVal, processedVal.String(), nil
	case reflect.Ptr:
		if processedVal.IsNil() {
			// If the pointer is nil, we need to encode the zero value of the element type, to get the possible processed type
			processedZero, _, _ := i.encodeValue(reflect.Zero(processedVal.Type().Elem()))
			return reflect.Zero(reflect.PointerTo(processedZero.Type())), nilStr, nil
		}

		// If the pointer is not nil, we need to encode the value it points to
		processedElem, encodedElem, err := i.encodeValue(processedVal.Elem())
		if err != nil {
			return reflect.Value{}, "", err
		}

		//
		if processedElem.CanAddr() {
			return processedElem.Addr(), encodedElem, nil
		}

		ptrElem := reflect.New(processedElem.Type())
		ptrElem.Elem().Set(processedElem)

		return ptrElem, encodedElem, nil
	case reflect.Array, reflect.Slice:
		return i.encodeArrayValue(processedVal)
	}

	return reflect.Value{}, "", fmt.Errorf("unsupported type: %v", kind)
}

func (i *Inspector) encodeArrayValue(val reflect.Value) (processedVal reflect.Value, encoded string, err error) {
	processedZero, _, _ := i.encodeValue(reflect.Zero(val.Type().Elem()))
	encodedElems := make([]string, val.Len())

	var processedElems reflect.Value
	switch val.Kind() {
	case reflect.Array:
		processedElems = reflect.New(reflect.ArrayOf(val.Len(), processedZero.Type())).Elem()
	case reflect.Slice:
		processedElems = reflect.MakeSlice(reflect.SliceOf(processedZero.Type()), val.Len(), val.Len())
	}

	for idx := range val.Len() {
		processedElem, encodedElem, err := i.encodeValue(val.Index(idx))
		if err != nil {
			return reflect.Value{}, "", err
		}

		processedElems.Index(idx).Set(processedElem)
		encodedElems[idx] = encodedElem
	}

	return processedElems, strings.Join(encodedElems, envSliceSep), nil
}
