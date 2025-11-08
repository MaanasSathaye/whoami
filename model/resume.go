package model

type StyleType int

const (
	StyleNormal StyleType = iota
	StyleBold
	StyleItalic
	StyleUnderline
	StyleHighlight
)

type ContactInfo struct {
	Name     string
	Location string
	Phone    string
	Email    string
	LinkedIn string
	GitHub   string
	Website  string
}

type TextSegment struct {
	Text  string
	Style StyleType
	URL   string
}

type Entry struct {
	Title       string
	Subtitle    string
	Location    string
	Date        string
	Description string
	Bullets     [][]TextSegment
}

type SectionType int

const (
	SectionTypeStandard SectionType = iota
	SectionTypeList
)

type Section struct {
	Title   string
	Type    SectionType
	Entries []Entry
	Items   []string
}

type Resume struct {
	Contact  ContactInfo
	Sections []Section
}
