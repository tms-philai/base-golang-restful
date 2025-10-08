package email

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestTemplates(t *testing.T) string {
	tmpDir := t.TempDir()

	htmlContent := `<html><body><h1>Hello, {{.Name}}!</h1></body></html>`
	textContent := `Hello, {{.Name}}!`

	err := os.WriteFile(filepath.Join(tmpDir, "welcome.html"), []byte(htmlContent), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, "welcome.txt"), []byte(textContent), 0644)
	assert.NoError(t, err)

	return tmpDir
}

func TestTemplateManager_LoadTemplate(t *testing.T) {
	tmpDir := setupTestTemplates(t)

	tm := NewTemplateManager(tmpDir)

	err := tm.LoadTemplate("welcome")
	assert.NoError(t, err)

	assert.True(t, tm.HasTemplate("welcome"))
}

func TestTemplateManager_LoadAllTemplates(t *testing.T) {
	tmpDir := setupTestTemplates(t)

	tm := NewTemplateManager(tmpDir)

	err := tm.LoadAllTemplates()
	assert.NoError(t, err)

	templates := tm.ListTemplates()
	assert.Contains(t, templates, "welcome")
}

func TestTemplateManager_Render(t *testing.T) {
	tmpDir := setupTestTemplates(t)

	tm := NewTemplateManager(tmpDir)
	err := tm.LoadTemplate("welcome")
	assert.NoError(t, err)

	data := map[string]interface{}{
		"Name":    "John",
		"Subject": "Welcome Email",
	}

	rendered, err := tm.Render("welcome", data)
	assert.NoError(t, err)
	assert.NotNil(t, rendered)
	assert.Contains(t, rendered.HTMLBody, "Hello, John!")
	assert.Contains(t, rendered.TextBody, "Hello, John!")
	assert.Equal(t, "Welcome Email", rendered.Subject)
}

func TestTemplateManager_RenderNonExistent(t *testing.T) {
	tmpDir := t.TempDir()

	tm := NewTemplateManager(tmpDir)

	_, err := tm.Render("nonexistent", nil)
	assert.Error(t, err)
}

func TestTemplateManager_HasTemplate(t *testing.T) {
	tmpDir := setupTestTemplates(t)

	tm := NewTemplateManager(tmpDir)
	err := tm.LoadTemplate("welcome")
	assert.NoError(t, err)

	assert.True(t, tm.HasTemplate("welcome"))
	assert.False(t, tm.HasTemplate("nonexistent"))
}

func TestTemplateManager_ListTemplates(t *testing.T) {
	tmpDir := setupTestTemplates(t)

	tm := NewTemplateManager(tmpDir)
	err := tm.LoadAllTemplates()
	assert.NoError(t, err)

	templates := tm.ListTemplates()
	assert.Len(t, templates, 1)
	assert.Contains(t, templates, "welcome")
}

func TestTemplateManager_RegisterTemplate(t *testing.T) {
	tm := NewTemplateManager("")

	tmpl := EmailTemplate{
		Name:        "test",
		Subject:     "Test Subject",
		HTMLContent: "<html><body>Hello, {{.Name}}!</body></html>",
		TextContent: "Hello, {{.Name}}!",
	}

	err := tm.RegisterTemplate(tmpl)
	assert.NoError(t, err)

	assert.True(t, tm.HasTemplate("test"))

	data := map[string]interface{}{
		"Name": "Alice",
	}

	rendered, err := tm.Render("test", data)
	assert.NoError(t, err)
	assert.Contains(t, rendered.HTMLBody, "Hello, Alice!")
	assert.Contains(t, rendered.TextBody, "Hello, Alice!")
}

func TestTemplateManager_RegisterInvalidTemplate(t *testing.T) {
	tm := NewTemplateManager("")

	tmpl := EmailTemplate{
		Name:        "invalid",
		HTMLContent: "<html><body>{{.InvalidSyntax</body></html>",
	}

	err := tm.RegisterTemplate(tmpl)
	assert.Error(t, err)
}
