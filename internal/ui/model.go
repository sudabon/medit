package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/sudabon/medit/internal/loader"
	"github.com/sudabon/medit/internal/render"
)

const (
	statusHeight = 1
	keyHints     = "e:edit r:reload t:theme q:quit"
)

func statusBarStyle(theme render.Theme) lipgloss.Style {
	if theme == render.ThemeLight {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#003d99")).
			Bold(true)
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("231")).
		Background(lipgloss.Color("24")).
		Bold(true)
}

// editorFinishedMsg is sent after tea.ExecProcess returns.
type editorFinishedMsg struct {
	err error
}

// Model is the Bubble Tea viewer model.
type Model struct {
	doc      *loader.Document
	renderer *render.Renderer
	viewport viewport.Model
	ready    bool
	width    int
	height   int
	pendingG bool
	message  string
	errMsg   string
}

// Run starts the TUI for the given document.
func Run(doc *loader.Document) error {
	m := New(doc)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// New constructs a viewer model.
func New(doc *loader.Document) Model {
	return Model{
		doc:      doc,
		renderer: render.New(),
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		vpHeight := max(1, msg.Height-statusHeight)
		cw := render.ContentWidth(msg.Width)
		if !m.ready {
			m.viewport = viewport.New(cw, vpHeight)
			m.ready = true
			if err := m.applyContent(false); err != nil {
				m.errMsg = err.Error()
			}
			return m, nil
		}

		percent := m.viewport.ScrollPercent()
		m.viewport.Width = cw
		m.viewport.Height = vpHeight
		if err := m.applyContent(false); err != nil {
			m.errMsg = err.Error()
		}
		m.restoreScrollPercent(percent)
		return m, nil

	case editorFinishedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("editor exited: %v", msg.err)
		} else {
			m.message = ""
		}
		y := m.viewport.YOffset
		if err := m.doc.Reload(); err != nil {
			m.errMsg = err.Error()
		} else {
			m.errMsg = ""
			if err := m.applyContent(true); err != nil {
				m.errMsg = err.Error()
			} else {
				m.viewport.SetYOffset(y)
			}
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Clear transient messages on interaction (except after editor).
	clearMsg := true

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "e", "i":
		m.pendingG = false
		return m, m.openEditor()

	case "r":
		m.pendingG = false
		y := m.viewport.YOffset
		if err := m.doc.Reload(); err != nil {
			m.errMsg = err.Error()
			m.message = ""
		} else {
			m.errMsg = ""
			m.message = "reloaded"
			if err := m.applyContent(true); err != nil {
				m.errMsg = err.Error()
			} else {
				m.viewport.SetYOffset(y)
			}
		}
		return m, nil

	case "t":
		m.pendingG = false
		theme := m.renderer.ToggleTheme()
		m.message = "theme: " + theme.String()
		if err := m.applyContent(true); err != nil {
			m.errMsg = err.Error()
		}
		return m, nil

	case "j", "down":
		m.pendingG = false
		m.viewport.ScrollDown(1)

	case "k", "up":
		m.pendingG = false
		m.viewport.ScrollUp(1)

	case "ctrl+d":
		m.pendingG = false
		m.viewport.HalfPageDown()

	case "ctrl+u":
		m.pendingG = false
		m.viewport.HalfPageUp()

	case "f", " ", "pgdown":
		m.pendingG = false
		m.viewport.PageDown()

	case "b", "pgup":
		m.pendingG = false
		m.viewport.PageUp()

	case "G":
		m.pendingG = false
		m.viewport.GotoBottom()

	case "g":
		if m.pendingG {
			m.pendingG = false
			m.viewport.GotoTop()
		} else {
			m.pendingG = true
			clearMsg = false
		}

	default:
		m.pendingG = false
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	if clearMsg {
		m.message = ""
	}
	return m, nil
}

func (m *Model) applyContent(force bool) error {
	if force {
		m.renderer.Invalidate()
	}
	width := max(1, m.viewport.Width)
	out, err := m.renderer.Render(m.doc.Content, width)
	if err != nil {
		return err
	}
	m.viewport.SetContent(out)
	return nil
}

func (m *Model) restoreScrollPercent(percent float64) {
	total := m.viewport.TotalLineCount()
	h := m.viewport.Height
	if total <= h {
		m.viewport.GotoTop()
		return
	}
	maxOff := total - h
	offset := int(percent * float64(maxOff))
	m.viewport.SetYOffset(offset)
}

func (m Model) openEditor() tea.Cmd {
	cmdName, extra := resolveEditor()
	line := render.ApproximateSourceLine(
		m.viewport.ScrollPercent(),
		render.SourceLineCount(m.doc.Content),
	)
	args := append([]string{}, extra...)
	args = append(args, fmt.Sprintf("+%d", line), m.doc.Path)

	c := exec.Command(cmdName, args...) //nolint:gosec // editor from env / vi by design
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err: err}
	})
}

// View implements tea.Model.
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}
	body := m.viewport.View()
	if pad := render.LeftPad(m.width, m.viewport.Width); pad > 0 {
		body = lipgloss.NewStyle().MarginLeft(pad).Render(body)
	}
	return lipgloss.JoinVertical(lipgloss.Left, body, m.statusBar())
}

func (m Model) statusBar() string {
	name := filepath.Base(m.doc.Path)
	pct := int(m.viewport.ScrollPercent() * 100)
	left := name
	if m.doc.IsNew {
		left += " [New]"
	}
	left += fmt.Sprintf("  %d%%", pct)

	msg := m.message
	if m.errMsg != "" {
		msg = m.errMsg
	}

	right := keyHints
	if msg != "" {
		right = msg + "  |  " + keyHints
	}

	gap := 2
	avail := max(0, m.width)
	leftW := runewidth.StringWidth(left)
	rightW := runewidth.StringWidth(right)
	space := avail - leftW - rightW - gap
	if space < 1 {
		// Prefer showing left + truncated right.
		remain := max(0, avail-leftW-1)
		right = runewidth.Truncate(right, remain, "…")
		space = max(0, avail-leftW-runewidth.StringWidth(right))
	}
	line := left + strings.Repeat(" ", space) + right
	if runewidth.StringWidth(line) < avail {
		line += strings.Repeat(" ", avail-runewidth.StringWidth(line))
	}
	line = runewidth.Truncate(line, avail, "…")
	return statusBarStyle(m.renderer.Theme()).Width(avail).Render(line)
}

// resolveEditor returns the editor command and extra args from
// $VISUAL → $EDITOR → vi, splitting fields on whitespace.
func resolveEditor() (string, []string) {
	for _, env := range []string{"VISUAL", "EDITOR"} {
		if v := strings.TrimSpace(os.Getenv(env)); v != "" {
			fields := strings.Fields(v)
			if len(fields) == 0 {
				continue
			}
			return fields[0], fields[1:]
		}
	}
	return "vi", nil
}
