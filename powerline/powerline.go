package powerline

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
)

type Part struct {
	Text   string
	Weight int
	Dirty  bool
}
type Parts []Part

func (slice Parts) Len() int {
	return len(slice)
}

func (slice Parts) Less(i, j int) bool {
	return slice[i].Weight > slice[j].Weight
}

func (slice Parts) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Segment struct {
	Foreground int
	Background int
	Weight     int
	Parts      Parts
}
type Segments []Segment

func (slice Segments) Len() int {
	return len(slice)
}

func (slice Segments) Less(i, j int) bool {
	return slice[i].Weight > slice[j].Weight
}

func (slice Segments) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Powerline struct {
	ShTemplate    string
	BashTemplate  string
	ColorTemplate string
	Reset         string
	Separator     string
	SeparatorThin string
	Ellipsis      string
	ReadOnly      string
	Phases        string
	Added         string
	Modified      string
	Untracked     string
	Removed       string
	Renamed       string
	Detached      string
	Attached      string
	Branch        string
	Ahead         string
	Behind        string
	Conflicted    string
	Dollar        string
	SetTitle      string
	Bold          string
	Segments      Segments
}

func (p *Powerline) Color(fore int, back int) string {
	return fmt.Sprintf(
		p.ShTemplate,
		fmt.Sprintf(p.ColorTemplate, fore, back),
	)
}

func (p *Powerline) ForegroundColor(fore int) string {
	return p.Color(38, fore)
}

func (p *Powerline) BackgroundColor(back int) string {
	return p.Color(48, back)
}

func (p *Powerline) AppendSegment(segment *Segment) {
	if segment != nil {
		p.Segments = append(p.Segments, *segment)
	}
}

func (p *Powerline) PrintSegments() string {
	var buffer bytes.Buffer
	var nextBackground string
	var text string

	// sort segments
	sort.Sort(p.Segments)

	for i, Seg := range p.Segments {

		// What color do we need to end the segment, this last background is
		// the next segments background
		if (i + 1) == len(p.Segments) {
			nextBackground = p.Reset
		} else {
			nextBackground = p.BackgroundColor(p.Segments[i+1].Background)
		}

		// sort parts
		sort.Sort(Seg.Parts)

		re := regexp.MustCompile("([$&\\\\`!])")
		for j, Part := range Seg.Parts {
			// escape dodgy shell injection characters
			text = Part.Text
			if Part.Dirty {
				text = re.ReplaceAllString(Part.Text, "\\$1")
			}
			// are we on the last part?
			if (j + 1) == len(Seg.Parts) {
				buffer.WriteString(fmt.Sprintf("%s%s %s %s%s%s",
					p.ForegroundColor(Seg.Foreground), p.BackgroundColor(Seg.Background),
					text,
					nextBackground, p.ForegroundColor(Seg.Background),
					p.Separator))
			} else {
				buffer.WriteString(fmt.Sprintf("%s%s %s %s%s%s",
					p.ForegroundColor(Seg.Foreground), p.BackgroundColor(Seg.Background),
					text,
					p.BackgroundColor(Seg.Background), p.ForegroundColor(Seg.Foreground), p.SeparatorThin))
			}
		}
	}

	buffer.WriteString(p.Reset)

	return buffer.String()
}

func NewPowerline(shell string, fancy bool) Powerline {
	p := Powerline{
		ReadOnly:      "\u2297",
		Separator:     "",
		SeparatorThin: "/",
		Ellipsis:      "\u2026",
		Branch:        "\u2607",
		Phases:        "+",
		Added:         "\u2714",
		Modified:      "\u270e",
		Untracked:     "\u2690",
		Removed:       "\u2716",
		Renamed:       "\u2608",
		Detached:      "\u2702",
		Ahead:         "\u21d1",
		Behind:        "\u21d3",
		Conflicted:    "\u203c",
	}

	if fancy {
		p.Separator = "\ue0b0"
		p.SeparatorThin = "\ue0b1"
		p.Branch = "\ue0a0"
	}

	switch shell {
	case "bash":
		p.ShTemplate = "\\[\\e%s\\]"
		p.ColorTemplate = "[%03d;5;%03dm"
		p.Reset = "\\[\\e[0m\\]"
		p.Bold = "\\[\\e[1m\\]"
		p.Dollar = "\\$"
		p.SetTitle = "\\[\\e]0;\\u@\\h: \\w\\a\\]"

	case "zsh":
		p.ShTemplate = "%s"
		// escape literal %'s (%%) as this gets passed through ShTemplate afterwards
		p.ColorTemplate = "%%{[%d;5;%dm%%}"
		// p.ColorTemplate = "%%{%%k{%d}%%f{%d}%%}"
		p.Reset = "%{%k%f%}"
		p.Bold = "%{[1m%}"
		p.Dollar = "%#"
		p.SetTitle = "%{\033]0;%n@%m: %~\007%}"
	}
	return p
}

// vim: ts=8 sw=8 noexpandtab:
