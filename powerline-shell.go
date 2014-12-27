// Copyright 2014 Matt Martz <matt@sivel.net>
// All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func getCurrentWorkingDir() (string, []string) {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	userDir := strings.Replace(dir, os.Getenv("HOME"), "~", 1)
	parts := strings.Split(userDir, "/")
	return dir, parts
}

func getVirtualEnv() (string, []string, string) {
	var parts []string
	virtualEnv := os.Getenv("VIRTUAL_ENV")
	if virtualEnv == "" {
		return "", parts, ""
	}

	parts = strings.Split(virtualEnv, "/")

	virtualEnvName := path.Base(virtualEnv)
	return virtualEnv, parts, virtualEnvName
}

func isWritableDir(dir string) bool {
	tmpPath := path.Join(dir, ".powerline-write-test")
	_, err := os.Create(tmpPath)
	if err != nil {
		return false
	}
	os.Remove(tmpPath)
	return true
}

func getGitInformation() (string, bool) {
	var status string
	var staged bool
	stdout, _ := exec.Command("git", "status", "--ignore-submodules").Output()
	reBranch := regexp.MustCompile(`^(HEAD detached at|HEAD detached from|On branch) (\S+)`)
	matchBranch := reBranch.FindStringSubmatch(string(stdout))
	if len(matchBranch) > 0 {
		if matchBranch[2] == "detached" {
			status = matchBranch[2]
		} else {
			status = matchBranch[2]
		}
	}

	reStatus := regexp.MustCompile(`Your branch is (ahead|behind).*?([0-9]+) comm`)
	matchStatus := reStatus.FindStringSubmatch(string(stdout))
	if len(matchStatus) > 0 {
		status = fmt.Sprintf("%s %s", status, matchStatus[2])
		if matchStatus[1] == "behind" {
			status = fmt.Sprintf("%s\u21E3", status)
		} else if matchStatus[1] == "ahead" {
			status = fmt.Sprintf("%s\u21E1", status)
		}
	}

	staged = !strings.Contains(string(stdout), "nothing to commit")
	if strings.Contains(string(stdout), "Untracked files") {
		status = fmt.Sprintf("%s +", status)
	}

	return status, staged
}

type Powerline struct {
	ZshTemplate   string
	BashTemplate  string
	ColorTemplate string
	Reset         string
	Lock          string
	Network       string
	Separator     string
	SeparatorThin string
	Ellipsis      string
	Segments      [][]string
}

func (p *Powerline) Color(prefix string, code string) string {
	return fmt.Sprintf(p.ZshTemplate, fmt.Sprintf(p.ColorTemplate, prefix, code))
}

func (p *Powerline) ForegroundColor(code string) string {
	return p.Color("$FG", code)
}

func (p *Powerline) BackgroundColor(code string) string {
	return p.Color("$BG", code)
}

func (p *Powerline) AppendSegment(segment []string) {
	if segment != nil {
		p.Segments = append(p.Segments, segment)
	}
}

func (p *Powerline) PrintSegments() string {
	var nextBackground string
	var buffer bytes.Buffer
	for i, Segment := range p.Segments {
		if (i + 1) == len(p.Segments) {
			nextBackground = p.Reset
		} else {
			nextBackground = p.BackgroundColor(p.Segments[i+1][1])
		}
		if len(Segment) == 3 {
			buffer.WriteString(fmt.Sprintf("%s%s %s %s%s%s", p.ForegroundColor(Segment[0]), p.BackgroundColor(Segment[1]), Segment[2], nextBackground, p.ForegroundColor(Segment[1]), p.Separator))
		} else {
			buffer.WriteString(fmt.Sprintf("%s%s %s %s%s%s", p.ForegroundColor(Segment[0]), p.BackgroundColor(Segment[1]), Segment[2], nextBackground, p.ForegroundColor(Segment[4]), Segment[3]))
		}
	}

	buffer.WriteString(p.Reset)

	return buffer.String()
}

func main() {
	home := false
	p := Powerline{
		ZshTemplate:   "%s",
		ColorTemplate: "%%{%s[%s]%%}",
		Reset:         "%{$reset_color%}",
		Lock:          "\uE0A2",
		Network:       "\uE0A2",
		Separator:     "\uE0B0",
		SeparatorThin: "\uE0B1",
		Ellipsis:      "\u2026",
	}
	cwd, cwdParts := getCurrentWorkingDir()
	if cwdParts[0] == "~" {
		cwdParts = cwdParts[1:len(cwdParts)]
		home = true
	}

	if !home && len(cwdParts) != 0 && len(cwdParts[len(cwdParts)-1]) > 0 {
		p.Segments = append(p.Segments, []string{"250", "237", "/", p.SeparatorThin, "244"})
	} else if !home {
		p.Segments = append(p.Segments, []string{"250", "237", "/"})
	} else {
		p.Segments = append(p.Segments, []string{"015", "031", "~"})
	}

	if len(cwdParts) >= 4 {
		p.Segments = append(p.Segments, []string{"250", "237", p.Ellipsis, p.SeparatorThin, "244"})
	} else if len(cwdParts) > 2 {
		if home {
			p.Segments = append(p.Segments, []string{"250", "237", p.Ellipsis, p.SeparatorThin, "244"})
		} else {
			p.Segments = append(p.Segments, []string{"250", "237", cwdParts[1], p.SeparatorThin, "244"})
		}
	}

	if len(cwdParts) != 0 && len(cwdParts[len(cwdParts)-1]) > 0 {
		p.Segments = append(p.Segments, []string{"250", "237", cwdParts[len(cwdParts)-1]})
	}

	p.AppendSegment(addVirtulEnvName())
	p.AppendSegment(addLock(cwd, p))
	p.AppendSegment(addGitInfo())
	p.AppendSegment(addDollarPrompt())

	fmt.Print(p.PrintSegments())
}

func addVirtulEnvName() []string {
	_, _, virtualEnvName := getVirtualEnv()
	if virtualEnvName != "" {
		return []string{"000", "035", virtualEnvName}
	}

	return nil
}

func addLock(cwd string, p Powerline) []string {
	if !isWritableDir(cwd) {
		return []string{"254", "124", p.Lock}
	}

	return nil
}

func addGitInfo() []string {
	gitStatus, gitStaged := getGitInformation()
	if gitStatus != "" {
		if gitStaged {
			return []string{"015", "161", gitStatus}
		} else {
			return []string{"000", "148", gitStatus}
		}
	} else {
		return nil
	}
}

func addDollarPrompt() []string {
	return []string{"015", "236", "\\$"}
}
