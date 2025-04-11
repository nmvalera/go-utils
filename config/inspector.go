package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// EncodeHookFunc is a function that can be used to encode a value
type EncodeHookFunc func(reflect.Value) (reflect.Value, error)

// InspectorConfig is the configuration for the Inspector
type InspectorConfig struct {
	TagNames   []string       // List of tag names to inspect
	EncodeHook EncodeHookFunc // Hook to apply to the value before encoding
	IncludeNil bool           // Whether to perform inspection even if a nil field is encountered
}

// Inspector is the main struct for inspecting a value
type Inspector struct {
	cfg *InspectorConfig
}

// NewInspector creates a new Inspector
func NewInspector(cfg *InspectorConfig) *Inspector {
	return &Inspector{cfg: cfg}
}

// FieldInfo represents information about a field in a struct
type FieldInfo struct {
	FieldParts []string            // Field name parts
	TagParts   map[string][]string // Map of inspected tag name to the tag values (each slice has the same length as the FieldParts slice)
	Value      reflect.Value       // Actual value of the field in the struct, possibly zeroed if the original value was unset (nil or within a nested nil struct field)
	Processed  reflect.Value       // Value after processing (in particular after EncodeHook has been applied), if not processed, it is the same as Value
	Encoded    string              // Encoded value of the field value
	IsNil      bool                // Whether the orignal value was nil or nested within a nil struct field
}

// Inspect inspects the given value and all its fields recursively
// It returns a map of field information where keys are the nested field names joined by "."
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

// inspect inspects the given value and populates the given map with the field information
// this function is recursive and inspects the given value and all its fields recursively
// - fieldParts is the slice of fieldName already inspected
// - val is the value to inspect
// - m is the map of fieldInfo to populate
// - tagParts is the map of tagName to the slice of the tag values already inspected
// - isNil is true if while inspecting we encountered a nil value
// - isPtr is true if the current value inspected is a pointer
func (i *Inspector) inspect(fieldParts []string, val reflect.Value, m map[string]*FieldInfo, tagParts map[string][]string, isNil, isPtr bool) error {
	// we process the value so in case a EncodeHook has been set for the value, it is applied
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
		// we encode the value so we get the processed value and the encoded string
		processedVal, encoded, err := i.encodeValue(processedVal)
		if err != nil {
			return fmt.Errorf("%s: failed to encode value: %w", join(fieldParts), err)
		}

		// we create the fieldInfo to populate
		fieldInfo := &FieldInfo{
			FieldParts: fieldParts,
			TagParts:   tagParts,
			Processed:  processedVal,
			Encoded:    encoded,
			IsNil:      isNil,
		}

		// we set the fieldInfo value
		valTyp := val.Type()
		if isNil {
			// if we encountered a nil value, we set the fieldInfo value to the zero value
			if isPtr {
				// if the value is a pointer, we set the fieldInfo value to a nil pointer to the type
				fieldInfo.Value = reflect.Zero(reflect.PointerTo(valTyp))
			} else {
				// if the value is not a pointer, we set the fieldInfo value to the zero value of the type
				fieldInfo.Value = reflect.Zero(valTyp)
			}
		} else {
			if isPtr {
				// if the value is a pointer, we set the fieldInfo value to the address of the value
				if val.CanAddr() {
					// if the value is addressable, we set the fieldInfo value to the address of the value
					fieldInfo.Value = val.Addr()
				} else {
					// if the value is not addressable, we create a new pointer to the value and set the fieldInfo value to it
					ptrElem := reflect.New(valTyp)
					ptrElem.Elem().Set(val)
					fieldInfo.Value = ptrElem
				}
			} else {
				// if the value is not a pointer, we set the fieldInfo value to the value itself
				fieldInfo.Value = val
			}
		}

		// we add the fieldInfo to the map
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
				// if we do not include nil values, we stop the inspection here
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

// inspectStruct inspects the given struct value and recursively inspects its fields
func (i *Inspector) inspectStruct(fieldParts []string, val reflect.Value, m map[string]*FieldInfo, tagParts map[string][]string, isNil bool) error {
	structType := val.Type()
	// inspect each field of the struct
	for idx := range val.NumField() {
		field := structType.Field(idx)
		fieldVal := val.Field(idx)

		fieldName := field.Name

		// we create a copy of the field parts so field's inspections do not concurrently modify the same fieldParts slice
		fieldPartsCopy := make([]string, len(fieldParts))
		copy(fieldPartsCopy, fieldParts)

		// we create a copy of the tag parts so field's inspections do not concurrently modify the same tagParts map
		tagPartsCopy := make(map[string][]string)
		for k, v := range tagParts {
			tagPartsCopy[k] = make([]string, len(v))
			copy(tagPartsCopy[k], v)
		}

		// we add the field name to the field parts
		fieldPartsCopy = append(fieldPartsCopy, fieldName)

		// we add the field tag parts to the tag parts
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

// join joins the given parts with a "." separator
func join(parts []string) string {
	return strings.Join(parts, ".")
}

// getKind returns a sub-kind of the given kind
// - for all reflect.Int* kinds, it returns reflect.Int
// - for all reflect.Uint* kinds, it returns reflect.Uint
// - for all reflect.Float* kinds, it returns reflect.Float32
// - for all other kinds, it returns the kind itself
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

// encodeValue encodes the given value and returns the processed value and the encoded string
// before encoding, it applies the EncodeHook to compute the processed value
// It supports a processed value of
// - base Kinds: reflect.Bool, reflect.Int, reflect.Uint, reflect.Float32, reflect.String
// - reflect.Pointer, in which case it recursively encode the pointer element and return a pointer to the processed element
// - reflect.Array and reflect.Slice, in which case it recursively encodes the elements and return a new array or slice of the processed elements
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
			// If the pointer is nil, we need to encode the zero value of the element type to know the possible processed type
			// then we return the processed zero value
			// the encoded string is arbitrarily set to "<nil>"
			processedZero, _, _ := i.encodeValue(reflect.Zero(processedVal.Type().Elem()))

			if processedZero.Kind() == reflect.String {
				return reflect.ValueOf(nilStr), nilStr, nil
			}

			return processedZero, nilStr, nil
		}

		// If the pointer is not nil, we need to encode the value it points to
		processedElem, encodedElem, err := i.encodeValue(processedVal.Elem())
		if err != nil {
			return reflect.Value{}, "", err
		}

		return processedElem, encodedElem, nil
	case reflect.Array, reflect.Slice:
		return i.encodeArrayValue(processedVal)
	}

	return reflect.Value{}, "", fmt.Errorf("unsupported type: %v", kind)
}

// encodeArrayValue encodes the given array or slice value and returns the processed value and the encoded string
// The array or slice is encoded by encoding each of its elements and joining the encoded elements with the " " separator
func (i *Inspector) encodeArrayValue(val reflect.Value) (processedVal reflect.Value, encoded string, err error) {
	// we encode the zero value of the element type to know the possible processed type
	// We do not need to check for errors, because we are just looking to know the returned type,
	// We make the assumption that the returned type is consistent across a given input type, this is actually not guaranteed, but is largely
	// acceptable for the vast majority of reasonable use cases and we are so far okay not supporting edge cases
	processedZero, _, _ := i.encodeValue(reflect.Zero(val.Type().Elem()))

	encodedElems := make([]string, val.Len())

	// we create a new array or slice of the processed element type
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
