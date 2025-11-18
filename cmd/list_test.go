package cmd

import (
	"reflect"
	"testing"
)

func TestListCmdStructure(t *testing.T) {
	cmd := &ListCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check Specs field exists
	specsField := val.FieldByName("Specs")
	if !specsField.IsValid() {
		t.Error("ListCmd does not have Specs field")
	}

	// Check Long field exists
	longField := val.FieldByName("Long")
	if !longField.IsValid() {
		t.Error("ListCmd does not have Long field")
	}

	// Check JSON field exists
	jsonField := val.FieldByName("JSON")
	if !jsonField.IsValid() {
		t.Error("ListCmd does not have JSON field")
	}
}

func TestCLIHasListCommand(t *testing.T) {
	cli := &CLI{}
	val := reflect.ValueOf(cli).Elem()
	listField := val.FieldByName("List")

	if !listField.IsValid() {
		t.Fatal("CLI struct does not have List field")
	}

	// Check the type
	if listField.Type().Name() != "ListCmd" {
		t.Errorf("List field type: got %s, want ListCmd", listField.Type().Name())
	}
}
