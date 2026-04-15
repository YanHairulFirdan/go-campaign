package user

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestUserModuleDoesNotImportAppPackage(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "module.go", nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse module.go imports: %v", err)
	}

	for _, importSpec := range file.Imports {
		if importSpec.Path.Value == `"go-campaign.com/internal/app"` {
			t.Fatal("user module must not import internal/app; pass only the dependencies user needs")
		}
	}
}
