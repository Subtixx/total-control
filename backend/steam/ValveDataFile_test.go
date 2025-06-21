package steam

import (
	"os"
	"strings"
	"testing"
)

func TestParseVDF_SimpleKeyValue(t *testing.T) {
	vdf := `"key"	"value"`
	root, err := ParseVDF(strings.NewReader(vdf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root.Key != "key" || root.Value != "value" {
		t.Errorf("expected key/value 'key'/'value', got '%s'/'%s'", root.Key, root.Value)
	}
}

func TestParseVDF_NestedStructure(t *testing.T) {
	vdf := `
"root"
{
	"child1"	"value1"
	"child2"
	{
		"grandchild"	"value2"
	}
}`
	root, err := ParseVDF(strings.NewReader(vdf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root.Key != "root" {
		t.Errorf("expected root key 'root', got '%s'", root.Key)
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(root.Children))
	}
	if root.Children[0].Key != "child1" || root.Children[0].Value != "value1" {
		t.Errorf("unexpected child1: %+v", root.Children[0])
	}
	child2 := root.Children[1]
	if child2.Key != "child2" || child2.Value != "" {
		t.Errorf("unexpected child2: %+v", child2)
	}
	if len(child2.Children) != 1 || child2.Children[0].Key != "grandchild" || child2.Children[0].Value != "value2" {
		t.Errorf("unexpected grandchild: %+v", child2.Children)
	}
}

func TestParseVDF_CommentsAndWhitespace(t *testing.T) {
	vdf := `
// this is a comment
"root"
{
	// another comment
	"key"	"value" // trailing comment
}`
	root, err := ParseVDF(strings.NewReader(vdf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root.Key != "root" {
		t.Errorf("expected root key 'root', got '%s'", root.Key)
	}
	if len(root.Children) != 1 || root.Children[0].Key != "key" || root.Children[0].Value != "value" {
		t.Errorf("unexpected child: %+v", root.Children)
	}
}

func TestParseVDF_InvalidLine(t *testing.T) {
	vdf := `"key" "value" "extra"`
	_, err := ParseVDF(strings.NewReader(vdf))
	if err == nil {
		t.Fatal("expected error for invalid VDF line, got nil")
	}
}

func TestParseVDF_LibraryFoldersFile(t *testing.T) {
	data, err := os.ReadFile("./data/libraryfolders.vdf")
	if err != nil {
		t.Fatalf("failed to read libraryfolders.vdf: %v", err)
	}

	root, err := ParseVDF(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root.Key != "libraryfolders" {
		t.Errorf("expected root key 'libraryfolders', got '%s'", root.Key)
	}
	if len(root.Children) < 1 {
		t.Errorf("expected at least 1 children, got %d", len(root.Children))
	}

	found := false
	for _, child := range root.Children {
		if child.Key != "0" {
			continue
		}
		if len(child.Children) < 1 {
			continue
		}
		if child.Children[0].Key != "path" {
			continue
		}

		for _, sub := range child.Children {
			if sub.Key == "path" && sub.Value == "/home/subtixx/.local/share/Steam" {
				found = true
			}
		}
	}

	if !found {
		t.Errorf("did not find expected library folder path")
	}
}
