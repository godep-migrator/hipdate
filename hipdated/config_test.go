package main

import (
	"testing"
)

func TestConfigMergeBackendOverwrite(t *testing.T) {
	cfg1, cfg2 := NewConfig(), NewConfig()

	cfg1.Backend = "foo"
	cfg2.Backend = "bar"

	cfg1.Merge(cfg2)
	if cfg1.Backend != "bar" {
		t.Logf("Backend was not overwritten (Value: %s)\n", cfg1.Backend)
		t.Fail()
	}
}

func TestConfigMergeBackendDontOverwriteWithEmpty(t *testing.T) {
	cfg1, cfg2 := NewConfig(), NewConfig()

	cfg1.Backend = "foo"
	cfg2.Backend = ""

	cfg1.Merge(cfg2)
	if cfg1.Backend != "foo" {
		t.Logf("Backend was overwritten (Value: '%s')\n", cfg1.Backend)
		t.Fail()
	}
}

func TestConfigMergeSources(t *testing.T) {
	cfg1, cfg2 := NewConfig(), NewConfig()

	cfg1.Sources.Set([]string{"a", "b", "c"})
	cfg2.Sources.Set([]string{"c", "d", "e"})

	cfg1.Merge(cfg2)
	if len(cfg1.Sources) != 5 {
		t.Logf("Sources has wrong length (length '%d')\n", len(cfg1.Sources))
		t.Log("Which indicates duplicates")
		t.Fail()
	}
}

func TestConfigMergeOptions(t *testing.T) {
	cfg1, cfg2 := NewConfig(), NewConfig()

	cfg1.Options["foo"] = "bar"
	cfg2.Options["foo"] = "qux"
	cfg2.Options["bar"] = "baz"

	cfg1.Merge(cfg2)
	if cfg1.Options["foo"] != "qux" {
		t.Logf("Options[\"foo\"] doesn't match \"%s\" != \"qux\"", cfg1.Options["foo"])
		t.Fail()
	}
	if cfg1.Options["bar"] != "baz" {
		t.Logf("Options[\"bar\"] doesn't match \"%s\" != \"baz\"", cfg1.Options["bar"])
		t.Fail()
	}
}
