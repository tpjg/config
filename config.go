package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var (
	// ErrNotPointerToStruct is an error returned when the supplied parameter
	// is not a pointer to a struct
	ErrNotPointerToStruct = errors.New("Expected pointer to struct parameter")
)

// flagset is a local struct so the standard library flag.FlagSet can be extended
// with some functions (see functions below)
type flagset struct {
	*flag.FlagSet
}

// ReadStructFromEnv reads the fields of a supplied pointer to a struct v from
// the environment variables. This function is typically used to get the
// configuration information for a (12 factor) application from environment
// variables.
func ReadStructFromEnv(v interface{}) error {
	if !isPointerToStruct(v) {
		return ErrNotPointerToStruct
	}
	tmpfs := flagset{flag.NewFlagSet("environment", flag.ExitOnError)}
	tmpfs.setupFlagsForStruct(v)
	args := tmpfs.getArgsFromEnv()
	// Now parse the prepared flags with the arguments taken from the environment
	tmpfs.Parse(args)
	return nil
}

// ReadStructFromEnvOverrideWithArgs reads the fields of a supplied pointer to a
// struct v from the environment variables, and if any command line arguments
// are given they will be read and override the environment flags.
func ReadStructFromEnvOverrideWithArgs(v interface{}) error {
	if !isPointerToStruct(v) {
		return ErrNotPointerToStruct
	}
	tmpfs := flagset{flag.NewFlagSet(os.Args[0], flag.ExitOnError)}
	tmpfs.setupFlagsForStruct(v)
	args := tmpfs.getArgsFromEnv()
	// Append the flags that do not start with "test." as used go "go test" tools
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-test.") {
			args = append(args, arg)
		}
	}
	// Now parse the prepared flags with the arguments taken from the environment
	tmpfs.Parse(args)
	return nil
}

// isPointerToStruct only returns true if the supplied interface{} is a pointer to a struct
func isPointerToStruct(v interface{}) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false
	}
	if !rv.IsValid() {
		return false
	}
	rvp := reflect.ValueOf(rv.Elem())
	return rvp.Kind() == reflect.Struct
}

func (fs flagset) setupFlagsForStruct(v interface{}) {
	// Get the value and type of underlying struct
	rv := reflect.ValueOf(v).Elem()
	t := rv.Type()
	// Prepare to read each field from the environment
	for i := 0; i < rv.NumField(); i++ {
		fs.setupFlagForField(rv.Field(i), t.Field(i))
	}
}

func (fs flagset) getArgsFromEnv() []string {
	// This is a copy of the https://github.com/ianschenck/envflag Parse() function
	// that takes the environment variables and creates "fake" flag K/V pairs.
	env := os.Environ()
	// Clean up and "fake" some flag k/v pairs.
	args := make([]string, 0, len(env))
	for _, value := range env {
		if fs.Lookup(value[:strings.Index(value, "=")]) == nil {
			continue
		}
		args = append(args, fmt.Sprintf("-%s", value))
	}
	return args
}

// setupFlagForField analyses the supplied field (reflect.Value) and sets up to
// read that field from the environment with envflag
func (fs flagset) setupFlagForField(rv reflect.Value, sf reflect.StructField) {
	if !rv.CanSet() {
		return
	}
	name, val := getFlagNameAndDefault(sf)
	// Get pointer to the underlying object so it can be set, similar to
	// how the reflect.SetString() like routines work.
	p := unsafe.Pointer(rv.Addr().Pointer())
	// Check for type supported by flag package only
	switch (rv.Interface()).(type) {
	case bool:
		fs.setupFlagForBool((*bool)(p), name, val)
	case time.Duration:
		fs.setupFlagForDuration((*time.Duration)(p), name, val)
	case float64:
		fs.setupFlagForFloat64((*float64)(p), name, val)
	case int:
		fs.setupFlagForInt((*int)(p), name, val)
	case int64:
		fs.setupFlagForInt64((*int64)(p), name, val)
	case string:
		fs.setupFlagForString((*string)(p), name, val)
	case uint:
		fs.setupFlagForUint((*uint)(p), name, val)
	case uint64:
		fs.setupFlagForUint64((*uint64)(p), name, val)
	}
}

// getFlagNameAndDefault determines the name and default for a field based on
// either the struct field Name or Tag. Supported tags are:
// flag - to set the name instead of using the field Name
// default - to set the default value for the parameter
func getFlagNameAndDefault(sf reflect.StructField) (name string, val string) {
	name = sf.Tag.Get(`flag`)
	if name == "" {
		name = sf.Name
	}
	val = sf.Tag.Get(`default`)
	return name, val
}

func (fs flagset) setupFlagForBool(p *bool, name string, val string) {
	value, _ := strconv.ParseBool(val)
	fs.BoolVar(p, name, value, "")
}

func (fs flagset) setupFlagForDuration(p *time.Duration, name string, val string) {
	value, _ := time.ParseDuration(val)
	fs.DurationVar(p, name, value, "")
}

func (fs flagset) setupFlagForFloat64(p *float64, name string, val string) {
	value, _ := strconv.ParseFloat(val, 64)
	fs.Float64Var(p, name, value, "")
}

func (fs flagset) setupFlagForInt(p *int, name string, val string) {
	value, _ := strconv.ParseInt(val, 0, 64)
	fs.IntVar(p, name, int(value), "")
}

func (fs flagset) setupFlagForInt64(p *int64, name string, val string) {
	value, _ := strconv.ParseInt(val, 0, 64)
	fs.Int64Var(p, name, value, "")
}

func (fs flagset) setupFlagForString(p *string, name string, val string) {
	fs.StringVar(p, name, val, "")
}

func (fs flagset) setupFlagForUint(p *uint, name string, val string) {
	value, _ := strconv.ParseUint(val, 0, 64)
	fs.UintVar(p, name, uint(value), "")
}

func (fs flagset) setupFlagForUint64(p *uint64, name string, val string) {
	value, _ := strconv.ParseUint(val, 0, 64)
	fs.Uint64Var(p, name, value, "")
}
