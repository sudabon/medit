package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
	"github.com/sudabon/medit/internal/loader"
	"github.com/sudabon/medit/internal/render"
)

const (
	statusHeight = 1
	keyHints     = "s:search e:edit r:reload t:theme q:quit"
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

	// lines is the rendered document before search highlighting.
	lines     []string
	searching bool
	input     textinput.Model
	query     string
	matches   []searchMatch
	matchIdx  int
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
	input := textinput.New()
	input.Prompt = "search: "
	input.Cursor.SetMode(cursor.CursorStatic)
	return Model{
		doc:      doc,
		renderer: render.New(),
		input:    input,
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
	if m.searching {
		return m.handleSearchKey(msg)
	}

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

	case "s":
		m.pendingG = false
		m.message = ""
		m.searching = true
		m.input.Width = max(1, m.width-runewidth.StringWidth(m.input.Prompt)-1)
		m.input.SetValue("")
		return m, m.input.Focus()

	case "n":
		m.pendingG = false
		m.message = ""
		m.stepMatch(1)
		return m, nil

	case "p":
		m.pendingG = false
		m.message = ""
		m.stepMatch(-1)
		return m, nil

	case "esc":
		m.pendingG = false
		m.message = ""
		m.clearSearch()
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

func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.searching = false
		m.input.Blur()
		if q := m.input.Value(); q != "" {
			m.runSearch(q)
		}
		return m, nil

	case "esc", "ctrl+c":
		m.searching = false
		m.input.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// runSearch finds query in the rendered document and jumps to the first
// match at or below the top of the current view, wrapping to the start.
func (m *Model) runSearch(query string) {
	m.query = query
	m.matches = findMatches(m.lines, query)
	if len(m.matches) == 0 {
		m.query = ""
		m.refreshContent()
		m.message = "not found: " + query
		return
	}
	m.matchIdx = 0
	for i, mt := range m.matches {
		if mt.line >= m.viewport.YOffset {
			m.matchIdx = i
			break
		}
	}
	m.refreshContent()
	m.revealMatch()
}

// stepMatch moves to the next (dir > 0) or previous match, wrapping around.
func (m *Model) stepMatch(dir int) {
	if len(m.matches) == 0 {
		return
	}
	n := len(m.matches)
	m.matchIdx = ((m.matchIdx+dir)%n + n) % n
	m.refreshContent()
	m.revealMatch()
}

func (m *Model) clearSearch() {
	if m.query == "" {
		return
	}
	m.query = ""
	m.matches = nil
	m.matchIdx = 0
	m.refreshContent()
}

// revealMatch scrolls so the current match is visible, placing it about a
// third of the way down the screen when it is off-screen.
func (m *Model) revealMatch() {
	line := m.matches[m.matchIdx].line
	top := m.viewport.YOffset
	if line >= top && line < top+m.viewport.Height {
		return
	}
	m.viewport.SetYOffset(max(0, line-m.viewport.Height/3))
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
	m.lines = strings.Split(out, "\n")
	m.matches = findMatches(m.lines, m.query)
	if m.matchIdx >= len(m.matches) {
		m.matchIdx = 0
	}
	m.refreshContent()
	return nil
}

// refreshContent pushes the rendered lines, with search highlights, into the viewport.
func (m *Model) refreshContent() {
	lines := m.lines
	if len(m.matches) > 0 {
		lines = highlightLines(m.lines, m.matches, m.matchIdx)
	}
	m.viewport.SetContent(strings.Join(lines, "\n"))
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
	if m.searching {
		line := m.input.View()
		if w := ansi.StringWidth(line); w < m.width {
			line += strings.Repeat(" ", m.width-w)
		}
		return statusBarStyle(m.renderer.Theme()).Width(max(0, m.width)).Render(line)
	}

	name := filepath.Base(m.doc.Path)
	pct := int(m.viewport.ScrollPercent() * 100)
	left := name
	if m.doc.IsNew {
		left += " [New]"
	}
	left += fmt.Sprintf("  %d%%", pct)
	if len(m.matches) > 0 {
		left += fmt.Sprintf("  /%s [%d/%d]", m.query, m.matchIdx+1, len(m.matches))
	}

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
