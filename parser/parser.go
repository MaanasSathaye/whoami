package parser

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/MaanasSathaye/whoami/model"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

var (
	datePattern     = regexp.MustCompile(`\{([^}]+)\}`)
	locationPattern = regexp.MustCompile(`^(.+?)\s*(?:-|–)\s*(.+?)(?:\s*\{|$)`)
)

// ParseMarkdown parses a markdown file into a Resume structure
func ParseMarkdown(content []byte) (*model.Resume, error) {
	var err error
	var frontmatterEnd int
	var contact model.ContactInfo
	var markdownContent []byte
	var doc ast.Node

	resume := model.Resume{}

	frontmatterEnd, contact, err = parseFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	resume.Contact = contact
	markdownContent = content[frontmatterEnd:]
	doc = goldmark.New().Parser().Parse(text.NewReader(markdownContent))
	resume.Sections = parseSections(doc, markdownContent)

	return &resume, nil
}

// parseFrontmatter extracts YAML frontmatter from the content
func parseFrontmatter(content []byte) (int, model.ContactInfo, error) {
	var err error
	var contact model.ContactInfo
	var start int
	var end int
	var endIndex int
	var yamlContent []byte
	var frontmatter map[string]string

	if !bytes.HasPrefix(content, []byte("---\n")) && !bytes.HasPrefix(content, []byte("---\r\n")) {
		return 0, contact, fmt.Errorf("no frontmatter found")
	}

	start = bytes.Index(content, []byte("---"))
	if start == -1 {
		return 0, contact, fmt.Errorf("no frontmatter start marker")
	}

	end = bytes.Index(content[start+3:], []byte("---"))
	if end == -1 {
		return 0, contact, fmt.Errorf("no frontmatter end marker")
	}

	endIndex = start + 3 + end + 3
	yamlContent = content[start+3 : start+3+end]

	if err = yaml.Unmarshal(yamlContent, &frontmatter); err != nil {
		return 0, contact, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	contact.Name = frontmatter["name"]
	contact.Location = frontmatter["location"]
	contact.Phone = frontmatter["phone"]
	contact.Email = frontmatter["email"]
	contact.LinkedIn = frontmatter["linkedin"]
	contact.GitHub = frontmatter["github"]
	contact.Website = frontmatter["website"]

	return endIndex, contact, nil
}

// parseSections extracts sections from the markdown AST
func parseSections(node ast.Node, source []byte) []model.Section {
	var sections []model.Section
	var currentSection *model.Section
	var currentEntry *model.Entry

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n.Kind() {
		case ast.KindHeading:
			heading := n.(*ast.Heading)
			headingText := extractText(n, source)

			if heading.Level == 1 {
				if currentSection != nil {
					if currentEntry != nil {
						currentSection.Entries = append(currentSection.Entries, *currentEntry)
						currentEntry = nil
					}
					sections = append(sections, *currentSection)
				}
				currentSection = &model.Section{
					Title:   strings.TrimSpace(headingText),
					Type:    model.SectionTypeStandard,
					Entries: []model.Entry{},
				}
			} else if heading.Level == 2 && currentSection != nil {
				if currentEntry != nil {
					currentSection.Entries = append(currentSection.Entries, *currentEntry)
				}
				currentEntry = &model.Entry{}
				parseH2Entry(headingText, currentEntry)
			} else if heading.Level == 3 && currentEntry != nil {
				currentEntry.Subtitle = strings.TrimSpace(headingText)
			}

		case ast.KindList:
			if currentEntry != nil {
				list := n.(*ast.List)
				bullets := parseList(list, source)
				currentEntry.Bullets = bullets
			} else if currentSection != nil {
				list := n.(*ast.List)
				items := parseSimpleList(list, source)
				if len(items) > 0 {
					currentSection.Type = model.SectionTypeList
					currentSection.Items = items
				}
			}
		}

		return ast.WalkContinue, nil
	})

	if currentEntry != nil && currentSection != nil {
		currentSection.Entries = append(currentSection.Entries, *currentEntry)
	}
	if currentSection != nil {
		sections = append(sections, *currentSection)
	}

	return sections
}

// parseH2Entry parses an H2 heading into entry title, location, and date
func parseH2Entry(text string, entry *model.Entry) {
	var cleanText string

	dateMatch := datePattern.FindStringSubmatch(text)
	if len(dateMatch) > 1 {
		entry.Date = strings.TrimSpace(dateMatch[1])
		cleanText = datePattern.ReplaceAllString(text, "")
	} else {
		cleanText = text
	}

	locationMatch := locationPattern.FindStringSubmatch(cleanText)
	if len(locationMatch) > 2 {
		entry.Title = strings.TrimSpace(locationMatch[1])
		entry.Location = strings.TrimSpace(locationMatch[2])
	} else {
		entry.Title = strings.TrimSpace(cleanText)
	}
}

// parseList parses a list into styled text segments
func parseList(list ast.Node, source []byte) [][]model.TextSegment {
	var bullets [][]model.TextSegment

	ast.Walk(list, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindListItem {
			segments := parseStyledText(n, source)
			if len(segments) > 0 {
				bullets = append(bullets, segments)
			}
		}
		return ast.WalkContinue, nil
	})

	return bullets
}

// parseSimpleList parses a list into plain text items
func parseSimpleList(list ast.Node, source []byte) []string {
	var items []string

	ast.Walk(list, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindListItem {
			text := extractText(n, source)
			items = append(items, strings.TrimSpace(text))
		}
		return ast.WalkContinue, nil
	})

	return items
}

// parseStyledText parses text with inline styling into TextSegments
func parseStyledText(node ast.Node, source []byte) []model.TextSegment {
	var segments []model.TextSegment

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n.Kind() {
		case ast.KindText:
			textNode := n.(*ast.Text)
			text := string(textNode.Segment.Value(source))
			segments = append(segments, model.TextSegment{
				Text:  text,
				Style: model.StyleNormal,
			})

		case ast.KindEmphasis:
			emphNode := n.(*ast.Emphasis)
			text := extractText(n, source)
			style := model.StyleItalic
			if emphNode.Level == 2 {
				style = model.StyleBold
			}
			segments = append(segments, model.TextSegment{
				Text:  text,
				Style: style,
			})
			return ast.WalkSkipChildren, nil

		case ast.KindLink:
			linkNode := n.(*ast.Link)
			text := extractText(n, source)
			url := string(linkNode.Destination)
			segments = append(segments, model.TextSegment{
				Text:  text,
				Style: model.StyleNormal,
				URL:   url,
			})
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})

	return segments
}

// extractText recursively extracts all text from a node
func extractText(node ast.Node, source []byte) string {
	var buf bytes.Buffer

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindText {
			textNode := n.(*ast.Text)
			buf.Write(textNode.Segment.Value(source))
		}
		return ast.WalkContinue, nil
	})

	return buf.String()
}
