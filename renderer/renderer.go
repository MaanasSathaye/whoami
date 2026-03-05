package renderer

import (
	"fmt"
	"strings"

	"github.com/MaanasSathaye/whoami/model"
	"github.com/jung-kurt/gofpdf"
)

const (
	pageWidth  = 210.0
	pageHeight = 297.0
)

type RenderConfig struct {
	MarginTop             float64
	MarginBottom          float64
	MarginLeft            float64
	MarginRight           float64
	FontSizeName          float64
	FontSizeContact       float64
	FontSizeSectionHeader float64
	FontSizeEntryTitle    float64
	FontSizeSubtitle      float64
	FontSizeBullet        float64
	SpacingAfterName      float64
	SpacingAfterContact   float64
	SpacingAfterSection   float64
	SpacingAfterEntry     float64
	SpacingAfterSubtitle  float64
	SpacingBetweenBullets float64
	LineHeight            float64
	BulletIndent          float64
}

// calculateConfig calculates rendering configuration based on max pages
func calculateConfig(maxPages int) RenderConfig {
	var config RenderConfig

	switch maxPages {
	case 1:
		config.MarginTop = 15.0
		config.MarginBottom = 15.0
		config.MarginLeft = 15.0
		config.MarginRight = 15.0
		config.FontSizeName = 16.0
		config.FontSizeContact = 9.0
		config.FontSizeSectionHeader = 10.0
		config.FontSizeEntryTitle = 9.0
		config.FontSizeSubtitle = 9.0
		config.FontSizeBullet = 9.0
		config.SpacingAfterName = 1.0
		config.SpacingAfterContact = 2.0
		config.SpacingAfterSection = 2.0
		config.SpacingAfterEntry = 1.5
		config.SpacingAfterSubtitle = 0.5
		config.SpacingBetweenBullets = 0.3
		config.LineHeight = 4.0
		config.BulletIndent = 4.0
	case 2:
		config.MarginTop = 20.0
		config.MarginBottom = 20.0
		config.MarginLeft = 20.0
		config.MarginRight = 20.0
		config.FontSizeName = 18.0
		config.FontSizeContact = 10.0
		config.FontSizeSectionHeader = 11.0
		config.FontSizeEntryTitle = 10.0
		config.FontSizeSubtitle = 10.0
		config.FontSizeBullet = 10.0
		config.SpacingAfterName = 2.0
		config.SpacingAfterContact = 4.0
		config.SpacingAfterSection = 4.0
		config.SpacingAfterEntry = 3.0
		config.SpacingAfterSubtitle = 1.0
		config.SpacingBetweenBullets = 0.5
		config.LineHeight = 5.0
		config.BulletIndent = 5.0
	default:
		config.MarginTop = 25.0
		config.MarginBottom = 25.0
		config.MarginLeft = 25.0
		config.MarginRight = 25.0
		config.FontSizeName = 20.0
		config.FontSizeContact = 11.0
		config.FontSizeSectionHeader = 12.0
		config.FontSizeEntryTitle = 11.0
		config.FontSizeSubtitle = 11.0
		config.FontSizeBullet = 11.0
		config.SpacingAfterName = 3.0
		config.SpacingAfterContact = 5.0
		config.SpacingAfterSection = 5.0
		config.SpacingAfterEntry = 4.0
		config.SpacingAfterSubtitle = 2.0
		config.SpacingBetweenBullets = 1.0
		config.LineHeight = 6.0
		config.BulletIndent = 6.0
	}

	return config
}

// RenderPDF renders a Resume to a PDF file with dynamic spacing
func RenderPDF(resume *model.Resume, outputPath string, maxPages int) error {
	var err error
	var baseConfig RenderConfig
	var adjustedConfig RenderConfig
	var pdf *gofpdf.Fpdf
	var tr func(string) string

	baseConfig = calculateConfig(maxPages)
	adjustedConfig = calculateDynamicSpacing(resume, baseConfig, maxPages)

	pdf = gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(adjustedConfig.MarginLeft, adjustedConfig.MarginTop, adjustedConfig.MarginRight)
	pdf.SetAutoPageBreak(true, adjustedConfig.MarginBottom)
	pdf.AddPage()

	tr = pdf.UnicodeTranslatorFromDescriptor("")

	renderHeader(pdf, resume.Contact, adjustedConfig)
	renderSections(pdf, resume.Sections, adjustedConfig, tr)

	if err = pdf.OutputFileAndClose(outputPath); err != nil {
		return fmt.Errorf("failed to save PDF: %w", err)
	}

	return nil
}

// calculateDynamicSpacing performs a measurement pass and adjusts spacing to fill the page
func calculateDynamicSpacing(resume *model.Resume, baseConfig RenderConfig, maxPages int) RenderConfig {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(baseConfig.MarginLeft, baseConfig.MarginTop, baseConfig.MarginRight)
	pdf.SetAutoPageBreak(false, baseConfig.MarginBottom)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	renderHeader(pdf, resume.Contact, baseConfig)
	renderSections(pdf, resume.Sections, baseConfig, tr)

	_, currentY := pdf.GetXY()
	contentHeight := currentY

	availableHeight := float64(maxPages) * (pageHeight - baseConfig.MarginTop - baseConfig.MarginBottom)
	usedHeight := contentHeight - baseConfig.MarginTop

	adjustedConfig := baseConfig

	if usedHeight > availableHeight {
		excessHeight := usedHeight - availableHeight

		if maxPages == 1 {
			marginReduction := 8.0
			adjustedConfig.MarginTop -= marginReduction
			adjustedConfig.MarginBottom -= marginReduction
			availableHeight += marginReduction * 2
			excessHeight = usedHeight - availableHeight
		}

		if excessHeight > 0 {
			totalSpacing := baseConfig.SpacingAfterName +
				baseConfig.SpacingAfterContact +
				baseConfig.SpacingAfterSection*float64(len(resume.Sections)) +
				baseConfig.SpacingAfterEntry*float64(countEntries(resume)) +
				baseConfig.SpacingAfterSubtitle*float64(countSubtitles(resume)) +
				baseConfig.SpacingBetweenBullets*float64(countBullets(resume))

			if totalSpacing > 0 {
				compressionFactor := (totalSpacing - excessHeight) / totalSpacing
				if compressionFactor < 0.25 {
					compressionFactor = 0.25
				}

				adjustedConfig.SpacingAfterName *= compressionFactor
				adjustedConfig.SpacingAfterContact *= compressionFactor
				adjustedConfig.SpacingAfterSection *= compressionFactor
				adjustedConfig.SpacingAfterEntry *= compressionFactor
				adjustedConfig.SpacingAfterSubtitle *= compressionFactor
				adjustedConfig.SpacingBetweenBullets *= compressionFactor
				adjustedConfig.LineHeight *= compressionFactor * 1.1
			}
		}

		return adjustedConfig
	}

	extraSpace := availableHeight - usedHeight
	numSpacingGaps := countSpacingGaps(resume)

	if numSpacingGaps == 0 {
		return baseConfig
	}

	extraSpacingPerGap := extraSpace / float64(numSpacingGaps)

	adjustedConfig.SpacingAfterName += extraSpacingPerGap * 0.5
	adjustedConfig.SpacingAfterContact += extraSpacingPerGap * 1.0
	adjustedConfig.SpacingAfterSection += extraSpacingPerGap * 1.5
	adjustedConfig.SpacingAfterEntry += extraSpacingPerGap * 1.0
	adjustedConfig.SpacingAfterSubtitle += extraSpacingPerGap * 0.3
	adjustedConfig.SpacingBetweenBullets += extraSpacingPerGap * 0.2

	return adjustedConfig
}

// countSpacingGaps counts the number of spacing gaps in the resume for distribution
func countSpacingGaps(resume *model.Resume) int {
	count := 2

	for _, section := range resume.Sections {
		count++
		if section.Type == model.SectionTypeList {
			count += len(section.Items)
		} else {
			for _, entry := range section.Entries {
				count++
				if entry.Subtitle != "" {
					count++
				}
				count += len(entry.Bullets)
			}
		}
	}

	return count
}

// countEntries counts the total number of entries in the resume
func countEntries(resume *model.Resume) int {
	count := 0
	for _, section := range resume.Sections {
		if section.Type == model.SectionTypeStandard {
			count += len(section.Entries)
		}
	}
	return count
}

// countSubtitles counts the total number of subtitles in the resume
func countSubtitles(resume *model.Resume) int {
	count := 0
	for _, section := range resume.Sections {
		if section.Type == model.SectionTypeStandard {
			for _, entry := range section.Entries {
				if entry.Subtitle != "" {
					count++
				}
			}
		}
	}
	return count
}

// countBullets counts the total number of bullet points in the resume
func countBullets(resume *model.Resume) int {
	count := 0
	for _, section := range resume.Sections {
		if section.Type == model.SectionTypeStandard {
			for _, entry := range section.Entries {
				count += len(entry.Bullets)
			}
		}
	}
	return count
}

// renderHeader renders the contact information header
func renderHeader(pdf *gofpdf.Fpdf, contact model.ContactInfo, config RenderConfig) {
	pageW, _ := pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	contentWidth := pageW - leftMargin - rightMargin

	pdf.SetFont("Times", "B", config.FontSizeName)
	pdf.CellFormat(contentWidth, 10, contact.Name, "", 1, "C", false, 0, "")
	pdf.Ln(config.SpacingAfterName)

	renderContactLine(pdf, contact, config, contentWidth, leftMargin)
	pdf.Ln(config.SpacingAfterContact)
}

type contactPart struct {
	text string
	url  string
}

// renderContactLine renders the contact information line with hyperlinks for socials
func renderContactLine(pdf *gofpdf.Fpdf, contact model.ContactInfo, config RenderConfig, contentWidth float64, leftMargin float64) {
	var parts []contactPart

	if contact.Location != "" {
		parts = append(parts, contactPart{text: contact.Location, url: ""})
	}
	if contact.Phone != "" {
		parts = append(parts, contactPart{text: contact.Phone, url: ""})
	}
	if contact.Email != "" {
		parts = append(parts, contactPart{text: contact.Email, url: ""})
	}
	if contact.LinkedIn != "" {
		parts = append(parts, contactPart{text: "LinkedIn", url: contact.LinkedIn})
	}
	if contact.GitHub != "" {
		parts = append(parts, contactPart{text: "GitHub", url: contact.GitHub})
	}
	if contact.Website != "" {
		parts = append(parts, contactPart{text: contact.Website, url: contact.Website})
	}

	pdf.SetFont("Times", "", config.FontSizeContact)

	var textParts []string
	for _, part := range parts {
		textParts = append(textParts, part.text)
	}
	fullText := strings.Join(textParts, " | ")
	totalWidth := pdf.GetStringWidth(fullText)

	startX := (contentWidth - totalWidth) / 2
	currentX := startX
	currentY, _ := pdf.GetY(), 0.0

	pdf.SetXY(leftMargin+currentX, currentY)

	for i, part := range parts {
		if i > 0 {
			separatorWidth := pdf.GetStringWidth(" | ")
			pdf.CellFormat(separatorWidth, 5, " | ", "", 0, "L", false, 0, "")
			currentX += separatorWidth
		}

		partWidth := pdf.GetStringWidth(part.text)

		if part.url != "" {
			pdf.SetTextColor(0, 0, 255)
			pdf.CellFormat(partWidth, 5, part.text, "", 0, "L", false, 0, part.url)
			pdf.SetTextColor(0, 0, 0)
		} else {
			pdf.CellFormat(partWidth, 5, part.text, "", 0, "L", false, 0, "")
		}

		currentX += partWidth
	}

	pdf.Ln(5)
}

// renderSections renders all sections
func renderSections(pdf *gofpdf.Fpdf, sections []model.Section, config RenderConfig, tr func(string) string) {
	for i, section := range sections {
		renderSectionHeader(pdf, section.Title, config)

		if section.Type == model.SectionTypeList {
			renderListSection(pdf, section.Items, config)
		} else {
			renderStandardSection(pdf, section.Entries, config, tr)
		}

		if i < len(sections)-1 {
			pdf.Ln(config.SpacingAfterSection)
		}
	}
}

// renderSectionHeader renders a section header with underline
func renderSectionHeader(pdf *gofpdf.Fpdf, title string, config RenderConfig) {
	pageW, _ := pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	contentWidth := pageW - leftMargin - rightMargin

	pdf.SetFont("Times", "B", config.FontSizeSectionHeader)
	pdf.CellFormat(contentWidth, 6, title, "", 1, "L", false, 0, "")

	x, y := pdf.GetXY()
	pdf.Line(x, y, x+contentWidth, y)
	pdf.Ln(2)
}

// renderStandardSection renders a standard section with entries
func renderStandardSection(pdf *gofpdf.Fpdf, entries []model.Entry, config RenderConfig, tr func(string) string) {
	for i, entry := range entries {
		renderEntry(pdf, entry, config, tr)
		if i < len(entries)-1 {
			pdf.Ln(config.SpacingAfterEntry)
		}
	}
}

// renderEntry renders a single entry
func renderEntry(pdf *gofpdf.Fpdf, entry model.Entry, config RenderConfig, tr func(string) string) {
	pageW, _ := pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	contentWidth := pageW - leftMargin - rightMargin

	if entry.Title != "" {
		renderEntryTitle(pdf, entry.Title, entry.Location, entry.Date, contentWidth, config)
	}

	if entry.Subtitle != "" {
		pdf.SetFont("Times", "I", config.FontSizeSubtitle)
		pdf.CellFormat(contentWidth, config.LineHeight, entry.Subtitle, "", 1, "L", false, 0, "")
		pdf.Ln(config.SpacingAfterSubtitle)
	}

	if len(entry.Bullets) > 0 {
		renderBullets(pdf, entry.Bullets, contentWidth, config, tr)
	}
}

// renderEntryTitle renders entry title with location and date
func renderEntryTitle(pdf *gofpdf.Fpdf, title, location, date string, contentWidth float64, config RenderConfig) {
	pdf.SetFont("Times", "B", config.FontSizeEntryTitle)

	titleText := title
	if location != "" {
		titleText = fmt.Sprintf("%s - %s", title, location)
	}

	dateWidth := 0.0
	if date != "" {
		pdf.SetFont("Times", "I", config.FontSizeEntryTitle)
		dateWidth = pdf.GetStringWidth(date)
	}

	pdf.SetFont("Times", "B", config.FontSizeEntryTitle)
	pdf.CellFormat(contentWidth-dateWidth, config.LineHeight, titleText, "", 0, "L", false, 0, "")

	if date != "" {
		pdf.SetFont("Times", "I", config.FontSizeEntryTitle)
		pdf.CellFormat(dateWidth, config.LineHeight, date, "", 1, "R", false, 0, "")
	} else {
		pdf.Ln(config.LineHeight)
	}
}

// renderBullets renders bullet points with styled text
func renderBullets(pdf *gofpdf.Fpdf, bullets [][]model.TextSegment, contentWidth float64, config RenderConfig, tr func(string) string) {
	bulletWidth := 3.0
	textWidth := contentWidth - config.BulletIndent - bulletWidth

	for _, segments := range bullets {
		x, y := pdf.GetXY()
		pdf.SetXY(x+config.BulletIndent, y)

		pdf.SetFont("Times", "", config.FontSizeBullet)
		pdf.CellFormat(bulletWidth, config.LineHeight, tr("•"), "", 0, "L", false, 0, "")

		renderStyledSegments(pdf, segments, textWidth, config)
		pdf.Ln(config.LineHeight + config.SpacingBetweenBullets)
	}
}

// renderStyledSegments renders text segments with styling
func renderStyledSegments(pdf *gofpdf.Fpdf, segments []model.TextSegment, maxWidth float64, config RenderConfig) {
	currentX, currentY := pdf.GetXY()
	lineStartX := currentX
	remainingWidth := maxWidth

	// Helper to check if text starts with sentence-ending punctuation (not decimal/number punctuation)
	startsWithPunct := func(text string) bool {
		trimmed := strings.TrimLeft(text, " \t\n\r")
		if len(trimmed) == 0 {
			return false
		}
		// Only treat as punctuation if it's a standalone punctuation character followed by space or is alone
		// This excludes things like ".8" or ",000" which are parts of numbers
		if len(trimmed) == 1 {
			firstChar := rune(trimmed[0])
			return firstChar == ',' || firstChar == '.' || firstChar == ';' || firstChar == ':' || firstChar == '!' || firstChar == '?' || firstChar == ')' || firstChar == ']' || firstChar == '}'
		}
		if len(trimmed) >= 2 {
			firstChar := rune(trimmed[0])
			secondChar := rune(trimmed[1])
			isPunct := firstChar == ',' || firstChar == '.' || firstChar == ';' || firstChar == ':' || firstChar == '!' || firstChar == '?' || firstChar == ')' || firstChar == ']' || firstChar == '}'
			// Only count as punctuation if followed by space or letter (not a digit)
			return isPunct && (secondChar == ' ' || (secondChar >= 'a' && secondChar <= 'z') || (secondChar >= 'A' && secondChar <= 'Z'))
		}
		return false
	}

	for i, segment := range segments {
		style := ""
		switch segment.Style {
		case model.StyleBold:
			style = "B"
		case model.StyleItalic:
			style = "I"
		case model.StyleUnderline:
			style = "U"
		default:
			style = ""
		}

		pdf.SetFont("Times", style, config.FontSizeBullet)

		text := segment.Text

		// Trim leading whitespace if this segment starts with punctuation
		if startsWithPunct(text) {
			text = strings.TrimLeft(text, " \t\n\r")
		}

		// Check if next segment starts with punctuation
		nextIsPunct := false
		if i+1 < len(segments) {
			nextIsPunct = startsWithPunct(segments[i+1].Text)
		}

		if segment.URL != "" {
			segmentWidth := pdf.GetStringWidth(text)

			if segmentWidth > remainingWidth && currentX > lineStartX {
				pdf.Ln(config.LineHeight)
				currentX = lineStartX
				currentY += config.LineHeight
				pdf.SetXY(currentX, currentY)
				remainingWidth = maxWidth
			}

			pdf.SetTextColor(0, 0, 255)
			pdf.CellFormat(segmentWidth, config.LineHeight, text, "", 0, "L", false, 0, segment.URL)
			pdf.SetTextColor(0, 0, 0)
			currentX += segmentWidth
			remainingWidth -= segmentWidth
		} else {
			words := strings.Fields(text)
			for j, word := range words {
				isLastWordInSegment := (j == len(words)-1)
				isLastSegment := (i == len(segments)-1)

				// Add space after word unless:
				// 1. It's the last word in the last segment, OR
				// 2. It's the last word in this segment and next segment starts with punctuation
				addSpace := true
				if isLastWordInSegment && isLastSegment {
					addSpace = false
				} else if isLastWordInSegment && nextIsPunct {
					addSpace = false
				}

				wordText := word
				if addSpace {
					wordText = word + " "
				}

				wordWidth := pdf.GetStringWidth(wordText)

				if wordWidth > remainingWidth && currentX > lineStartX {
					pdf.Ln(config.LineHeight)
					currentX = lineStartX
					currentY += config.LineHeight
					pdf.SetXY(currentX, currentY)
					remainingWidth = maxWidth
				}

				pdf.CellFormat(wordWidth, config.LineHeight, wordText, "", 0, "L", false, 0, "")
				currentX += wordWidth
				remainingWidth -= wordWidth
			}
		}
	}
}

// renderListSection renders a list section (like interests/certifications)
func renderListSection(pdf *gofpdf.Fpdf, items []string, config RenderConfig) {
	pageW, _ := pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	contentWidth := pageW - leftMargin - rightMargin

	for _, item := range items {
		parts := strings.Split(item, ":")
		if len(parts) == 2 {
			label := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			pdf.SetFont("Times", "B", config.FontSizeBullet)
			labelWidth := pdf.GetStringWidth(label + ": ")
			pdf.CellFormat(labelWidth, config.LineHeight, label+": ", "", 0, "L", false, 0, "")

			pdf.SetFont("Times", "", config.FontSizeBullet)
			pdf.MultiCell(contentWidth-labelWidth, config.LineHeight, value, "", "L", false)
		} else {
			pdf.SetFont("Times", "", config.FontSizeBullet)
			pdf.MultiCell(contentWidth, config.LineHeight, item, "", "L", false)
		}
	}
}
