package file

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-yaml"
)

// Parse parses the YAML data into the destination structure.
// dest must be a pointer to the target structure.
//
// Example:
//
//	var data = `
//	a: test
//	b: 6
//	`
//
//	var schema struct {
//		A string
//		B int
//	}
//
//	err := Parse(data, &schema)
//	if err != nil {
//		// handle error
//	}
//
//	fmt.Println(schema) // Output: {A: "test", B: 6}
func Parse(data string, dest any) error {
	// Check dest type
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("yaml parser: dest must be a pointer")
	}

	if err := yaml.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("yaml parser: failed to parse YAML: %w", err)
	}
	return nil
}

// Mapping maps data from src to copy using YAML as an intermediary.
// copy must be a pointer to the target structure.
//
// Example:
//
//	var src = struct {
//		A string
//		B int
//	}{
//		A: "example",
//		B: 42,
//	}
//
//	var dest struct {
//		A string
//	}
//
//	err := Mapping(src, &dest)
//	if err != nil {
//		// handle error
//	}
//
//	fmt.Println(dest) // Output: {A: "example"}
func Mapping(src, copy any) error {
	// Marshal src to YAML
	data, err := yaml.Marshal(src)
	if err != nil {
		return fmt.Errorf("yaml mapper: failed to marshal source: %w", err)
	}

	// Check dest type
	if reflect.TypeOf(copy).Kind() != reflect.Ptr {
		return fmt.Errorf("yaml parser: copy must be a pointer")
	}

	// Unmarshal YAML to copy
	if err := yaml.Unmarshal(data, copy); err != nil {
		return fmt.Errorf("yaml mapper: failed to unmarshal to destination: %w", err)
	}
	return nil
}
