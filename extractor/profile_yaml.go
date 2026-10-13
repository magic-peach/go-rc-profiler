package extractor

import (
	"sort"

	"gopkg.in/yaml.v3"
)

var profileFieldOrder = []string{
	"apiVersion",
	"kind",
	"metadata",
	"workload",
	"extensions",
	"conditions",
}

// MarshalProfileYAML serializes a profile map using the field order defined
// by the Runtime Conditions Profile specification. Any additional root fields
// are emitted afterward in lexical order.
func MarshalProfileYAML(profile map[string]any) ([]byte, error) {
	root := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	seen := make(map[string]bool, len(profile))
	appendField := func(key string) error {
		value, ok := profile[key]
		if !ok {
			return nil
		}
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}
		valueNode := &yaml.Node{}
		if err := valueNode.Encode(value); err != nil {
			return err
		}
		root.Content = append(root.Content, keyNode, valueNode)
		seen[key] = true
		return nil
	}

	for _, key := range profileFieldOrder {
		if err := appendField(key); err != nil {
			return nil, err
		}
	}

	extraFields := make([]string, 0, len(profile)-len(seen))
	for key := range profile {
		if !seen[key] {
			extraFields = append(extraFields, key)
		}
	}
	sort.Strings(extraFields)
	for _, key := range extraFields {
		if err := appendField(key); err != nil {
			return nil, err
		}
	}

	document := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{root}}
	return yaml.Marshal(document)
}
