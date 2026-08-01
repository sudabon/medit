package render

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	styles "github.com/charmbracelet/glamour/styles"
	"github.com/muesli/termenv"
)

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }
func uintPtr(u uint) *uint    { return &u }

// highContrastDark is Glamour's dark style with brighter body text and
// rendered headings (no raw "#" markdown prefixes).
var highContrastDark = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockPrefix: "\n",
			BlockSuffix: "\n",
			Color:       strPtr("255"),
		},
		Margin: uintPtr(0),
	},
	BlockQuote: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: strPtr("250"),
		},
		Indent:      uintPtr(1),
		IndentToken: strPtr("│ "),
	},
	List: ansi.StyleList{
		LevelIndent: 2,
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockSuffix: "\n",
			Color:       strPtr("81"),
			Bold:        boolPtr(true),
		},
	},
	H1: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          " ",
			Suffix:          " ",
			Color:           strPtr("231"),
			BackgroundColor: strPtr("63"),
			Bold:            boolPtr(true),
		},
	},
	H2: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "▍ ",
			Color:  strPtr("159"),
			Bold:   boolPtr(true),
		},
	},
	H3: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "• ",
			Color:  strPtr("117"),
			Bold:   boolPtr(true),
		},
	},
	H4: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: strPtr("111"),
			Bold:  boolPtr(true),
		},
	},
	H5: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: strPtr("110"),
			Bold:  boolPtr(true),
		},
	},
	H6: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: strPtr("109"),
			Bold:  boolPtr(true),
		},
	},
	Strikethrough: ansi.StylePrimitive{
		CrossedOut: boolPtr(true),
	},
	Emph: ansi.StylePrimitive{
		Italic: boolPtr(true),
		Color:  strPtr("255"),
	},
	Strong: ansi.StylePrimitive{
		Bold:  boolPtr(true),
		Color: strPtr("231"),
	},
	HorizontalRule: ansi.StylePrimitive{
		Color:  strPtr("245"),
		Format: "\n────────\n",
	},
	Item: ansi.StylePrimitive{
		BlockPrefix: "• ",
		Color:       strPtr("255"),
	},
	Enumeration: ansi.StylePrimitive{
		BlockPrefix: ". ",
		Color:       strPtr("159"),
	},
	Task: ansi.StyleTask{
		Ticked:   "[✓] ",
		Unticked: "[ ] ",
	},
	Link: ansi.StylePrimitive{
		Color:     strPtr("81"),
		Underline: boolPtr(true),
	},
	LinkText: ansi.StylePrimitive{
		Color: strPtr("159"),
		Bold:  boolPtr(true),
	},
	Image: ansi.StylePrimitive{
		Color:     strPtr("218"),
		Underline: boolPtr(true),
	},
	ImageText: ansi.StylePrimitive{
		Color:  strPtr("250"),
		Format: "Image: {{.text}} →",
	},
	Code: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          " ",
			Suffix:          " ",
			Color:           strPtr("231"),
			BackgroundColor: strPtr("238"),
		},
	},
	CodeBlock: ansi.StyleCodeBlock{
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: strPtr("252"),
			},
			Margin: uintPtr(0),
		},
		Chroma: styles.DarkStyleConfig.CodeBlock.Chroma,
	},
	Table: ansi.StyleTable{
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: strPtr("255"),
			},
		},
	},
	DefinitionDescription: ansi.StylePrimitive{
		BlockPrefix: "\n🠶 ",
	},
}

// highContrastLight is for light terminal / system Light appearance.
var highContrastLight = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockPrefix: "\n",
			BlockSuffix: "\n",
			Color:       strPtr("#1a1a1a"),
		},
		Margin: uintPtr(0),
	},
	BlockQuote: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: strPtr("#333333"),
		},
		Indent:      uintPtr(1),
		IndentToken: strPtr("│ "),
	},
	List: ansi.StyleList{
		LevelIndent: 2,
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockSuffix: "\n",
			Color:       strPtr("#003d99"),
			Bold:        boolPtr(true),
		},
	},
	H1: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          " ",
			Suffix:          " ",
			Color:           strPtr("#ffffff"),
			BackgroundColor: strPtr("#003d99"),
			Bold:            boolPtr(true),
		},
	},
	H2: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "▍ ",
			Color:  strPtr("#003d99"),
			Bold:   boolPtr(true),
		},
	},
	H3: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "• ",
			Color:  strPtr("#0052cc"),
			Bold:   boolPtr(true),
		},
	},
	H4: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: strPtr("#1a1a1a"),
			Bold:  boolPtr(true),
		},
	},
	H5: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: strPtr("#1a1a1a"),
			Bold:  boolPtr(true),
		},
	},
	H6: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: strPtr("#333333"),
			Bold:  boolPtr(true),
		},
	},
	Strikethrough: ansi.StylePrimitive{
		CrossedOut: boolPtr(true),
	},
	Emph: ansi.StylePrimitive{
		Italic: boolPtr(true),
		Color:  strPtr("#1a1a1a"),
	},
	Strong: ansi.StylePrimitive{
		Bold:  boolPtr(true),
		Color: strPtr("#000000"),
	},
	HorizontalRule: ansi.StylePrimitive{
		Color:  strPtr("#666666"),
		Format: "\n────────\n",
	},
	Item: ansi.StylePrimitive{
		BlockPrefix: "• ",
		Color:       strPtr("#1a1a1a"),
	},
	Enumeration: ansi.StylePrimitive{
		BlockPrefix: ". ",
		Color:       strPtr("#003d99"),
	},
	Task: ansi.StyleTask{
		Ticked:   "[✓] ",
		Unticked: "[ ] ",
	},
	Link: ansi.StylePrimitive{
		Color:     strPtr("#0052cc"),
		Underline: boolPtr(true),
	},
	LinkText: ansi.StylePrimitive{
		Color: strPtr("#003d99"),
		Bold:  boolPtr(true),
	},
	Image: ansi.StylePrimitive{
		Color:     strPtr("#a6005c"),
		Underline: boolPtr(true),
	},
	ImageText: ansi.StylePrimitive{
		Color:  strPtr("#333333"),
		Format: "Image: {{.text}} →",
	},
	Code: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          " ",
			Suffix:          " ",
			Color:           strPtr("#000000"),
			BackgroundColor: strPtr("#e8e8e8"),
		},
	},
	CodeBlock: ansi.StyleCodeBlock{
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: strPtr("#1a1a1a"),
			},
			Margin: uintPtr(0),
		},
		Chroma: styles.LightStyleConfig.CodeBlock.Chroma,
	},
	Table: ansi.StyleTable{
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: strPtr("#1a1a1a"),
			},
		},
	},
	DefinitionDescription: ansi.StylePrimitive{
		BlockPrefix: "\n🠶 ",
	},
}

// selectStyle returns the style for theme.
func selectStyle(theme Theme) ansi.StyleConfig {
	if theme == ThemeLight {
		return highContrastLight
	}
	return highContrastDark
}

func renderMarkdown(content string, width int, theme Theme) (string, error) {
	tr, err := glamour.NewTermRenderer(
		glamour.WithStyles(selectStyle(theme)),
		glamour.WithWordWrap(width),
		glamour.WithColorProfile(termenv.TrueColor),
	)
	if err != nil {
		return "", err
	}
	return tr.Render(content)
}
