package powerline

import (
	"bytes"
	"fmt"
)

type Powerline struct {
	ShTemplate    string
	BashTemplate  string
	ColorTemplate string
	Reset         string
	Lock          string
	Network       string
	Separator     string
	SeparatorThin string
	Ellipsis      string
	Dollar        string
	Bold          string
	Segments      [][]interface{}
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

func (p *Powerline) AppendSegment(segment []interface{}) {
	if segment != nil {
		p.Segments = append(p.Segments, segment)
	}
}

func (p *Powerline) AppendSegments(segments [][]interface{}) {
	for _, segment := range segments {
		p.AppendSegment(segment)
	}
}

func (p *Powerline) PrintSegments() string {
	var nextBackground string
	var buffer bytes.Buffer
	for i, Segment := range p.Segments {
		if (i + 1) == len(p.Segments) {
			nextBackground = p.Reset
		} else {
			nextBackground = p.BackgroundColor(p.Segments[i+1][1].(int))
		}
		if len(Segment) == 3 {
			buffer.WriteString(fmt.Sprintf("%s%s %s %s%s%s", p.ForegroundColor(Segment[0].(int)), p.BackgroundColor(Segment[1].(int)), Segment[2].(string), nextBackground, p.ForegroundColor(Segment[1].(int)), p.Separator))
		} else {
			buffer.WriteString(fmt.Sprintf("%s%s %s %s%s%s", p.ForegroundColor(Segment[0].(int)), p.BackgroundColor(Segment[1].(int)), Segment[2], nextBackground, p.ForegroundColor(Segment[4].(int)), Segment[3]))
		}
	}

	buffer.WriteString(p.Reset)

	return buffer.String()
}

func NewPowerline(shell string) Powerline {
	p := Powerline{
		Lock:          "\uE0A2",
		Network:       "\uE0A2",
		Separator:     "\uE0B0",
		SeparatorThin: "\uE0B1",
		Ellipsis:      "\u2026",
	}

	switch shell {
	case "bash":
		p.ShTemplate = "\\[\\e%s\\]"
		p.ColorTemplate = "[%03d;5;%03dm"
		p.Reset = "\\[\\e[0m\\]"
		p.Bold = "\\[\\e[1m\\]"
		p.Dollar = "\\$"

	case "zsh":
		p.ShTemplate = "%s"
                // escape literal %'s (%%) as this gets passed through ShTemplate afterwards
                p.ColorTemplate = "%%{[%d;5;%dm%%}"
		// p.ColorTemplate = "%%{%%k{%d}%%f{%d}%%}"
		p.Reset = "%{%k%f%}"
		p.Bold = "%{[1m%}"
		p.Dollar = "%#"
	}
	return p
}
