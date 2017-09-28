package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"unsafe"

	"github.com/0xrawsec/golang-utils/log"
)

// Config : configuration structure definition
type Config map[string]Value

// Value : stored in the configuration
type Value interface{}

var (
	ErrNoSuchKey = errors.New("No such key")
)

//////////////////////////////// Utils /////////////////////////////////////////

func configErrorf(fmt string, i ...interface{}) {
	log.Errorf(fmt, i...)
	os.Exit(1)
}

func getRequiredError(key, ofType string, err error) {
	configErrorf("Cannot get mandatory parameter %s as %s: %s ", key, ofType, err)
}

////////////////////////////////////////////////////////////////////////////////

// Loads : loads a configuration structure from a data buffer
// @data : buffer containing the configuration object
// return (Config, error) : the Config struct filled from data, error code
func Loads(data []byte) (c Config, err error) {
	err = json.Unmarshal(data, &c)
	if err != nil {
		return
	}
	return
}

// Load : loads a configuration structure from a file
// @path : path where the configuration is stored as a json file
// return (Config, error) : the Config struct parsed, error code
func Load(path string) (c Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return Loads([]byte(data))
}

// Dumps : Dumps Config structure into a byte slice
// return ([]byte, error) : byte slice and error code
func (c *Config) Dumps() (dump []byte, err error) {
	dump, err = json.Marshal(c)
	if err != nil {
		return
	}
	return
}

// Debug : prints out the configuration in debug information
func (c *Config) Debug() {
	for key, val := range *c {
		log.Debugf("config[%s] = %v", key, val)
	}
}

// Get : get the Value associated to a key found in Config structure
// return (Value, error) : Value associated to key and error code
func (c *Config) Get(key string) (Value, error) {
	val, ok := (*c)[key]
	if !ok {
		return val, ErrNoSuchKey
	}
	return val, nil
}

// GetString gets the value associated to a key as string
// return (string, error)
func (c *Config) GetString(key string) (string, error) {
	val, ok := (*c)[key]
	if !ok {
		return "", ErrNoSuchKey
	}
	s, ok := val.(string)
	if !ok {
		return s, fmt.Errorf("Wrong type for %s (Type:%T Expecting:%T)", key, val, s)
	}
	return val.(string), nil
}

// GetInt64 gets the value associated to a key as int64
// return (int64, error)
func (c *Config) GetInt64(key string) (i int64, err error) {
	val, ok := (*c)[key]
	if !ok {
		return 0, ErrNoSuchKey
	}
	switch val.(type) {
	case int8:
		return int64(val.(int8)), nil
	case int16:
		return int64(val.(int16)), nil
	case int:
		return int64(val.(int)), nil
	case int32:
		return int64(val.(int32)), nil
	case int64:
		return val.(int64), nil
	case float64:
		// json loads float64 so handle that case
		return int64(val.(float64)), nil
	default:
		return 0, fmt.Errorf("Wrong type for %s (Type:%T Expecting:%T)", key, val, i)
	}
}

// GetUint64 gets the value associated to a key as uint64
// return (uint64, error)
func (c *Config) GetUint64(key string) (u uint64, err error) {
	val, ok := (*c)[key]
	if !ok {
		return 0, ErrNoSuchKey
	}
	switch val.(type) {
	case uint8:
		return uint64(val.(uint8)), nil
	case uint16:
		return uint64(val.(uint16)), nil
	case uint32:
		return uint64(val.(uint32)), nil
	case uint:
		return uint64(val.(uint)), nil
	case uint64:
		return val.(uint64), nil
	case float64:
		// json loads float64 so handle that case
		return uint64(val.(float64)), nil
	default:
		return 0, fmt.Errorf("Wrong type for %s (Type:%T Expecting:%T)", key, val, u)
	}
}

// GetSubConfig : get a subconfig referenced by key
// return (Config, error)
func (c *Config) GetSubConfig(key string) (Config, error) {
	val, err := c.Get(key)
	if err != nil {
		return Config{}, err
	}
	sc, ok := val.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Wrong type for %s (Type:%T Expecting:%T)", key, val, sc)
	}
	return *(*Config)(unsafe.Pointer(&(sc))), nil
}

// GetRequiredSubConfig : get a subconfig referenced by key
// return (Config)
func (c *Config) GetRequiredSubConfig(key string) Config {
	sc, err := c.GetSubConfig(key)
	if err != nil {
		getRequiredError(key, "map[string]interface{}", err)
	}
	return sc
}

// GetRequired : get the Value associated to a key found in Config structure and exit if
// not available
// return (Value) : Value associated to key if it exists
func (c *Config) GetRequired(key string) Value {
	val, err := c.Get(key)
	if err != nil {
		configErrorf("Configuration parameter %s is mandatory", key)
	}
	return val
}

func (c *Config) GetRequiredString(key string) string {
	s, err := c.GetString(key)
	if err != nil {
		getRequiredError(key, "string", err)
	}
	return s
}

func (c *Config) GetRequiredInt64(key string) int64 {
	val, err := c.GetInt64(key)
	if err != nil {
		getRequiredError(key, "int64", err)
	}
	return val
}

func (c *Config) GetRequiredUint64(key string) uint64 {
	val, err := c.GetUint64(key)
	if err != nil {
		getRequiredError(key, "uint64", err)
	}
	return val
}

func (c *Config) GetStringSlice(key string) (s []string, err error) {
	s = make([]string, 0)
	val, err := c.Get(key)
	if err != nil {
		return
	}
	ival, ok := val.([]interface{})
	if !ok {
		return s, fmt.Errorf("Wrong type for %s (Type:%T Expecting:%T)", key, val, []interface{}{})
	}
	for _, e := range ival {
		s = append(s, e.(string))
	}
	return
}

func (c *Config) GetRequiredStringSlice(key string) []string {
	ss, err := c.GetStringSlice(key)
	if err != nil {
		getRequiredError(key, "[]string", err)
	}
	return ss
}

func (c *Config) GetUint64Slice(key string) (u []uint64, err error) {
	u = make([]uint64, 0)
	val, err := c.Get(key)
	if err != nil {
		return
	}
	ival, ok := val.([]interface{})
	if !ok {
		return u, fmt.Errorf("Wrong type for %s (Type:%T Expecting:%T)", key, val, []interface{}{})
	}
	for _, e := range ival {
		u = append(u, e.(uint64))
	}
	return
}

func (c *Config) GetRequiredUint64Slice(key string) []uint64 {
	val, err := c.GetUint64Slice(key)
	if err != nil {
		getRequiredError(key, "[]uint64", err)
	}
	return val
}

func (c *Config) GetInt64Slice(key string) (i []int64, err error) {
	i = make([]int64, 0)
	val, err := c.Get(key)
	if err != nil {
		return
	}
	ival, ok := val.([]interface{})
	if !ok {
		return i, fmt.Errorf("Wrong type for %s (Type:%T Expecting:%T)", key, val, []interface{}{})
	}
	for _, e := range ival {
		i = append(i, e.(int64))
	}
	return
}

func (c *Config) GetRequiredInt64Slice(key string) []int64 {
	val, err := c.GetInt64Slice(key)
	if err != nil {
		getRequiredError(key, "[]int64", err)
	}
	return val
}

// Set : set parameter identified by key of the Config struct with a Value
func (c *Config) Set(key string, value interface{}) {
	(*c)[key] = value
}

// HasKey returns true if the configuration has the given key
func (c *Config) HasKey(key string) bool {
	_, ok := (*c)[key]
	return ok
}
