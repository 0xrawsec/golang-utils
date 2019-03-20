package config

import (
	"testing"
)

var (
	configpath = "./test/config.json"
	conf       = Config{
		"test":   "foo",
		"foobar": 64,
		"array":  []string{"this", "is", "an", "array"}}
)

func TestDumps(t *testing.T) {
	dumps, err := conf.Dumps()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(dumps))
}

func TestLoads(t *testing.T) {
	dumps, err := conf.Dumps()
	if err != nil {
		t.Error(err)
	}
	loaded, err := Loads(dumps)
	if err != nil {
		t.Error(err)
	}
	t.Log(loaded)
}

func TestSetGet(t *testing.T) {
	conf.Set("Foo", map[string]string{"Foo": "bar"})
	val, err := conf.Get("Foo")
	if err == nil {
		t.Logf("%T, %[1]v", val)
	}
}

func TestAll(t *testing.T) {
	conf.Set("Foo", map[string]string{"foo": "bar"})
	dumps, err := conf.Dumps()
	if err != nil {
		t.Error(err)
	}
	loaded, err := Loads(dumps)
	if err != nil {
		t.Error(err)
	}
	t.Log(loaded)
	val, err := loaded.Get("Foo")
	if err == nil {
		t.Logf("%T, %[1]v", val)
		foo := val.(map[string]interface{})["foo"]
		t.Logf("%T, %[1]v", foo)
	}
}

func TestLoadJson(t *testing.T) {
	c, err := Load(configpath)
	if err != nil {
		panic(err)
	}
	m := c.GetRequiredSubConfig("misp")
	t.Log(m)
	l := c.GetRequiredSubConfig("log-search")
	t.Log(l)
	s := c.GetRequiredStringSlice("notification-recipients")
	t.Log(s)

}
