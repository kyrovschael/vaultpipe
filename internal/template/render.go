package template

import (
	"fmt"
	"regexp"
	"strings"
)

// varPattern matches {{vault:secret/path#field}} style placeholders.
var varPattern = regexp.MustCompile(`\{\{vault:([^#}]+)#([^}]+)\}\}`)

// Renderer replaces vault placeholders in strings with resolved secret values.
type Renderer struct {
	secrets map[string]map[string]string // path -> field -> value
}

// NewRenderer creates a Renderer backed by the provided secrets map.
func NewRenderer(secrets map[string]map[string]string) *Renderer {
	return &Renderer{secrets: secrets}
}

// Render replaces all {{vault:path#field}} placeholders in input.
// Returns an error if any placeholder cannot be resolved.
func (r *Renderer) Render(input string) (string, error) {
	var renderErr error
	result := varPattern.ReplaceAllStringFunc(input, func(match string) string {
		if renderErr != nil {
			return match
		}
		parts := varPattern.FindStringSubmatch(match)
		if len(parts) != 3 {
			renderErr = fmt.Errorf("malformed placeholder: %s", match)
			return match
		}
		path, field := strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2])
		fields, ok := r.secrets[path]
		if !ok {
			renderErr = fmt.Errorf("secret path not found: %s", path)
			return match
		}
		val, ok := fields[field]
		if !ok {
			renderErr = fmt.Errorf("field %q not found in path %s", field, path)
			return match
		}
		return val
	})
	if renderErr != nil {
		return "", renderErr
	}
	return result, nil
}

// RenderMap applies Render to every value in the provided map.
func (r *Renderer) RenderMap(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		rendered, err := r.Render(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = rendered
	}
	return out, nil
}
