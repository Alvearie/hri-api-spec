package main

import (
    "testing"
)

func TestValidateGood(t *testing.T) {
    response := Validate("example1.json")
    if response != "The document is valid\n" {
        t.Errorf("Validation of example1.json failed: %q", response)
    }
}

func TestValidateBad(t *testing.T) {
    response := Validate("example2.json")
    if response != "The document is not valid. see errors :\n- (root): status is required\n- recordCount: Invalid type. Expected: integer, given: string\n" {
        t.Errorf("Validation of example2.json failed: %q", response)
    }
}
