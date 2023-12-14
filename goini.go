// Package ini provides methods to read ini data from files or from byte array
package goini

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"
)

type KeyValue struct {
	Item map[string]string
}

type Ini struct {
	fileName string
	buffer   []byte
	sections map[string]KeyValue
}

// NewIni creates a new Ini struct
func NewIni() *Ini {
	result := new(Ini)
	result.clear()
	return result
}

// LoadFromFile read the content (sections, keys and values) of an ini file and
// and keep all data in memory inside the ini struct.
// It returns error if the operation fails.
func (i *Ini) LoadFromFile(path string) error {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	i.LoadFromBytes(buffer)
	return nil
}

// LoadFromByte read the content (sections, keys and values) of an ini content
// and keep all data inside the ini struct.
func (i *Ini) LoadFromBytes(buffer []byte) {
	i.clear()
	i.buffer = buffer
	i.parseLines()
}

// KeyExists returns true if the key inside the section exists (even if has no value).
func (i *Ini) KeyExists(section string, key string) bool {
	_, exists := i.sections[section].Item[key]
	return exists
}

// SectionExists returns true if the section exists (even if is empty)
func (i *Ini) SectionExists(section string) bool {
	_, exists := i.sections[section]
	return exists
}

// GetSectionValues returns a slice with all the sections. Note: the empty
// section "" is always available.
func (i *Ini) GetSectionValues() []string {
	result := make([]string, 0, len(i.sections))
	for sec := range i.sections {
		result = append(result, sec)
	}
	return result
}

// GetInt gets the key value from section as Int.
// If error or value is empty or bad formed returns default value.
func (i *Ini) GetInt(section string, key string, def int) int {
	strValue := i.getString(section, key)
	result, err := strconv.ParseInt(strValue, 0, 32)
	if err != nil {
		return def
	}
	return int(result)
}

// GetUint gets the key value from section as Uint.
// If error or value is empty or bad formed returns default value.
func (i *Ini) GetUint(section string, key string, def uint) uint {
	strValue := i.getString(section, key)
	result, err := strconv.ParseUint(strValue, 0, 32)
	if err != nil {
		return def
	}
	return uint(result)
}

// GetFloat gets the key value from section as Float32.
// If error or value is empty or bad formed returns default value.
func (i *Ini) GetFloat(section string, key string, def float32) float32 {
	strValue := i.getString(section, key)
	result, err := strconv.ParseFloat(strValue, 10)
	if err != nil {
		return def
	}
	return float32(result)
}

// GetBool gets the key value from section as Bool.
// If error or value is empty or bad formed returns default value.
func (i *Ini) GetBool(section string, key string, def bool) bool {
	strValue := i.getString(section, key)
	result, err := strconv.ParseBool(strValue)
	if err != nil {
		return def
	}
	return bool(result)
}

// GetString gets the key value from section as String.
// If value is an empty string, returns default value.
func (i *Ini) GetString(section string, key string, def string) string {
	result := i.getString(section, key)
	if result == "" {
		result = def
	}
	return result
}

// GetStringSlice gets the key value from section as Slice of strings.
// The items must be separated by a character (string not rune).
// If an item is an empty string, it value is set to defaul value.
func (i *Ini) GetStringSlice(section string, key string, def string, sep string) []string {
	result := strings.Split(i.getString(section, key), sep)

	if sep == "" || len(result) == 0 {
		return result
	}
	// Set defaul value to empty strings
	for i := 0; i < len(result); i++ {
		if result[i] == "" {
			result[i] = def
		}
	}
	return result
}

// ==========================
// Internal methods
// ==========================

func (i *Ini) getString(section string, key string) string {
	return i.sections[section].Item[key]
}

func (i *Ini) clear() {
	i.sections = make(map[string]KeyValue)
	i.buffer = []byte{}
}

func (i *Ini) parseLines() {
	reader := bytes.NewReader(i.buffer)
	scanner := bufio.NewScanner(reader)
	currentSection := ""
	i.sections[currentSection] = KeyValue{Item: make(map[string]string)}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Empty lines or comments...
		if line == "" || line[0] == '#' || line[0] == ';' {
			continue
		}

		// Section?
		if line[0] == '[' {
			// Must end with "]" if not, ignore it
			if line[len(line)-1] != ']' {
				continue
			}
			currentSection = line[1 : len(line)-1]
			// Create keyvalue map if not exists
			if _, exists := i.sections[currentSection]; !exists {
				i.sections[currentSection] = KeyValue{Item: make(map[string]string)}
			}

			continue
		}

		pairs := strings.Split(line, "=")
		// Only lines with one '=' are allowed
		if len(pairs) != 2 {
			continue
		}
		key := strings.TrimSpace(pairs[0])
		value := strings.TrimSpace(pairs[1])
		i.sections[currentSection].Item[key] = value
	}
}

// ==========================
// WORK IN PROGRESS METHODS (not tested yet)....
// ==========================

// No es posible usar esta sintaxis en un método.
//      No compila ==> func (i *Ini) [T any]GetValue(section string, key string, def any) any
// Si se usa con la sintaxis normal, el método de llamada tiene una sintaxis horrible:
//      func (i *Ini) GetValue(section string, key string, def any) any
//      Ejemplo: ini.GetValue("8 bits colors", "red", int8(0)).int8()

// This method has an ugly syntax: GetValue("sect", "key", int64(0)).int64()
// A more specific method. If you wanna be more precise. IE, checks if a return
// type (int8) does not fit in the type from the ini. "width=325"

/* TODO:
func (i *Ini) GetValue(section string, key string, def any) any {
	defType := reflect.TypeOf(def)
	defValue := reflect.ValueOf(def)
	varSize := varSize(def) //int(unsafe.Sizeof(defValue.Interface()))

	result := reflect.New(defType).Elem()

	switch result.Kind() {
	// Int
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := strconv.ParseInt(i.getString(section, key), 10, varSize)
		if err != nil {
			result.SetInt(defValue.Int())
		} else {
			result.SetInt(x)
		}

		// Uint
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := strconv.ParseUint(i.getString(section, key), 10, varSize)
		if err != nil {
			result.SetUint(defValue.Uint())
		} else {
			result.SetUint(x)
		}

	// Bool
	case reflect.Bool:
		x, err := strconv.ParseBool(i.getString(section, key))
		if err != nil {
			result.SetBool(defValue.Bool())
		} else {
			result.SetBool(x)
		}

	// Float
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(i.getString(section, key), 10)
		if err != nil {
			result.SetFloat(defValue.Float())
		} else {
			result.SetFloat(x)
		}

		// String
	case reflect.String:
		x := i.getString(section, key)
		result.SetString(x)

	default:
		panic("Type not implemented")
		//return def
	}
	return result.Interface()
}

func varSize(x any) int {
    varType := reflect.TypeOf(x)
    varSize := varType.Size() * 8
    return int(varSize)
}
*/

/* TODO:
func GetSplitValues[T any](section string, key string, def T) []T {
	//func (i *Ini) GetSplitValues(section string, key string, def T) []T {
	panic("Not Implemented")
}
*/
