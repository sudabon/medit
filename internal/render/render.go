package render

import (
	"strings"
)

// Renderer wraps Glamour and re-renders only when content, width, or theme changes.
type Renderer struct {
	content  string
	width    int
	theme    Theme
	rendered string
	ready    bool
}

// New creates a Renderer using DetectTheme().
func New() *Renderer {
	return &Renderer{theme: DetectTheme()}
}

// Theme returns the current theme.
func (r *Renderer) Theme() Theme {
	return r.theme
}

// SetTheme switches the color scheme and invalidates the cache.
func (r *Renderer) SetTheme(theme Theme) {
	if r.theme == theme {
		return
	}
	r.theme = theme
	r.Invalidate()
}

// ToggleTheme flips between dark and light.
func (r *Renderer) ToggleTheme() Theme {
	if r.theme == ThemeLight {
		r.SetTheme(ThemeDark)
	} else {
		r.SetTheme(ThemeLight)
	}
	return r.theme
}

// Render returns ANSI-rendered Markdown for the given content and wrap width.
// Fenced code blocks are drawn with a rounded box border.
func (r *Renderer) Render(content string, width int) (string, error) {
	if width < 1 {
		width = 1
	}
	if r.ready && r.content == content && r.width == width {
		return r.rendered, nil
	}

	out, err := renderDocument(content, width, r.theme)
	if err != nil {
		return "", err
	}

	r.content = content
	r.width = width
	r.rendered = out
	r.ready = true
	return out, nil
}

func renderDocument(content string, width int, theme Theme) (string, error) {
	segs := splitFenced(content)
	if len(segs) == 0 {
		return renderMarkdown("", width, theme)
	}

	var b strings.Builder
	for i, seg := range segs {
		if seg.isCode {
			if i > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(renderBoxedCode(seg.code, seg.lang, width, theme))
			b.WriteByte('\n')
			continue
		}
		if strings.TrimSpace(seg.markdown) == "" {
			continue
		}
		part, err := renderMarkdown(seg.markdown, width, theme)
		if err != nil {
			return "", err
		}
		b.WriteString(part)
	}
	return b.String(), nil
}

// Invalidate forces the next Render call to recompute.
func (r *Renderer) Invalidate() {
	r.ready = false
}
