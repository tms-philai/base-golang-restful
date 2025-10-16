package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sync"
	textTemplate "text/template"
)

type TemplateManager struct {
	htmlTemplates map[string]*template.Template
	textTemplates map[string]*textTemplate.Template
	templatesDir  string
	mu            sync.RWMutex
}

type RenderedTemplate struct {
	Subject  string
	HTMLBody string
	TextBody string
}

func NewTemplateManager(templatesDir string) *TemplateManager {
	return &TemplateManager{
		htmlTemplates: make(map[string]*template.Template),
		textTemplates: make(map[string]*textTemplate.Template),
		templatesDir:  templatesDir,
	}
}

func (tm *TemplateManager) LoadTemplate(name string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	htmlPath := filepath.Join(tm.templatesDir, name+".html")
	textPath := filepath.Join(tm.templatesDir, name+".txt")

	if _, err := os.Stat(htmlPath); err == nil {
		tmpl, err := template.ParseFiles(htmlPath)
		if err != nil {
			return fmt.Errorf("failed to parse HTML template: %w", err)
		}
		tm.htmlTemplates[name] = tmpl
	}

	if _, err := os.Stat(textPath); err == nil {
		tmpl, err := textTemplate.ParseFiles(textPath)
		if err != nil {
			return fmt.Errorf("failed to parse text template: %w", err)
		}
		tm.textTemplates[name] = tmpl
	}

	if _, hasHTML := tm.htmlTemplates[name]; !hasHTML {
		if _, hasText := tm.textTemplates[name]; !hasText {
			return fmt.Errorf("template %s not found", name)
		}
	}

	return nil
}

func (tm *TemplateManager) LoadAllTemplates() error {
	files, err := os.ReadDir(tm.templatesDir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	templateNames := make(map[string]bool)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		ext := filepath.Ext(name)
		if ext == ".html" || ext == ".txt" {
			baseName := name[:len(name)-len(ext)]
			templateNames[baseName] = true
		}
	}

	for name := range templateNames {
		if err := tm.LoadTemplate(name); err != nil {
			return err
		}
	}

	return nil
}

func (tm *TemplateManager) Render(name string, data interface{}) (*RenderedTemplate, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	result := &RenderedTemplate{}

	if htmlTmpl, exists := tm.htmlTemplates[name]; exists {
		var buf bytes.Buffer
		if err := htmlTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("failed to execute HTML template: %w", err)
		}
		result.HTMLBody = buf.String()
	}

	if textTmpl, exists := tm.textTemplates[name]; exists {
		var buf bytes.Buffer
		if err := textTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("failed to execute text template: %w", err)
		}
		result.TextBody = buf.String()
	}

	if result.HTMLBody == "" && result.TextBody == "" {
		return nil, fmt.Errorf("template %s not found", name)
	}

	if dataMap, ok := data.(map[string]interface{}); ok {
		if subject, ok := dataMap["Subject"].(string); ok {
			result.Subject = subject
		}
	}

	return result, nil
}

func (tm *TemplateManager) HasTemplate(name string) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	_, hasHTML := tm.htmlTemplates[name]
	_, hasText := tm.textTemplates[name]

	return hasHTML || hasText
}

func (tm *TemplateManager) ListTemplates() []string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	templates := make(map[string]bool)

	for name := range tm.htmlTemplates {
		templates[name] = true
	}
	for name := range tm.textTemplates {
		templates[name] = true
	}

	result := make([]string, 0, len(templates))
	for name := range templates {
		result = append(result, name)
	}

	return result
}

type EmailTemplate struct {
	Name        string
	Subject     string
	HTMLContent string
	TextContent string
}

func (tm *TemplateManager) RegisterTemplate(tmpl EmailTemplate) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tmpl.HTMLContent != "" {
		htmlTmpl, err := template.New(tmpl.Name).Parse(tmpl.HTMLContent)
		if err != nil {
			return fmt.Errorf("failed to parse HTML template: %w", err)
		}
		tm.htmlTemplates[tmpl.Name] = htmlTmpl
	}

	if tmpl.TextContent != "" {
		textTmpl, err := textTemplate.New(tmpl.Name).Parse(tmpl.TextContent)
		if err != nil {
			return fmt.Errorf("failed to parse text template: %w", err)
		}
		tm.textTemplates[tmpl.Name] = textTmpl
	}

	return nil
}
