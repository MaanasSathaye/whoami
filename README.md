# WhoAmI - Resume Builder

A lightweight, open-source resume builder written in Go that converts Markdown files into professional PDF resumes.

## Features

- Write your resume in simple Markdown syntax
- Generate professional PDF output with clickable hyperlinks for social profiles
- Dynamic vertical spacing that automatically fills the page based on content
- Control spacing and layout with max_pages parameter
- Support for text styling (bold, italic)
- Automatic formatting for sections, entries, and bullet points
- Clean, professional layout matching standard resume formats

## Installation

```bash
git clone https://github.com/MaanasSathaye/whoami.git
cd whoami
go build -o whoami
```

## Usage

```bash
./whoami <output.pdf> <input.md> [max_pages]
```

### Parameters

- `output.pdf`: Path to the output PDF file
- `input.md`: Path to the input Markdown file
- `max_pages` (optional): Target page count for layout optimization
  - `1`: Compact layout (smaller fonts, tighter spacing)
  - `2`: Standard layout (default, balanced spacing)
  - `3+`: Spacious layout (larger fonts, more spacing)

### Examples

```bash
# Generate resume with default (2-page) layout
./whoami resume.pdf resume.md

# Generate compact one-page resume
./whoami resume.pdf resume.md 1

# Generate spacious multi-page resume
./whoami resume.pdf resume.md 3
```

### Dynamic Spacing

The resume builder uses a two-pass rendering approach:

1. **Measurement Pass**: Calculates the total height of all content with base spacing
2. **Adjustment Pass**: Distributes extra vertical space proportionally across spacing gaps

This ensures that resumes fill the specified number of pages vertically without excessive whitespace at the bottom, creating a more professional and balanced appearance. The spacing adjustment is weighted to prioritize:
- Section gaps (highest weight)
- Entry gaps (medium weight)
- Inter-bullet spacing (lowest weight)

## Markdown Syntax

### Frontmatter (Required)

Start your resume with YAML frontmatter containing your contact information:

```markdown
---
name: John Doe
location: Austin, Texas
phone: 123-555-1234
email: john.doe@email.com
linkedin: https://linkedin.com/in/johndoe
github: https://github.com/johndoe
website: https://johndoe.com
---
```

**Note**: The `linkedin`, `github`, and `website` fields should contain full URLs and will be rendered as clickable blue hyperlinks in the PDF. The displayed text will be "LinkedIn", "GitHub", or the full website URL respectively. The `email` field is displayed as plain text without a hyperlink.

### Sections

Use H1 headers for section titles:

```markdown
# EXPERIENCE

# EDUCATION

# PROJECTS & OPEN-SOURCE

# CERTIFICATIONS & INTERESTS
```

### Entries

Use H2 headers for company/organization names with optional location and dates:

```markdown
## Company Name - Location {Date Range}
## University Name - City, State {Graduation Date}
```

The format is: `Title - Location {Date}`
- Location is separated by ` - `
- Date is enclosed in `{}`
- Both location and date are optional

### Position Titles

Use H3 headers for position titles or descriptions:

```markdown
### Senior Software Engineer
### B.S Computer Science
```

### Bullet Points

Use standard Markdown bullet points with optional text styling:

```markdown
- Built **scalable** systems with *high performance*
- Improved metrics by **50%** through optimization
- Implemented feature using *best practices*
```

### Text Styling

- `**text**` - Bold
- `*text*` - Italic

### List Sections

For sections like Certifications & Interests, use bullet points with `key: value` format:

```markdown
# CERTIFICATIONS & INTERESTS

- Certifications: GCP Professional Data Engineer (Feb 2022), AWS Solutions Architect (Jan 2023)
- Interests: Photography, Hiking, Chess, Guitar
```

## Example Resume

See `template_resume.md` for a complete template with example formatting and instructions.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions welcome! Please open an issue or submit a pull request.
