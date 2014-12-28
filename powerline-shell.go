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
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wm/powerline-shell-go/powerline"
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

func addCwd(cwd string, cwdParts []string, p powerline.Powerline) [][]string {
	segments := [][]string{}
	home := false
	if cwdParts[0] == "~" {
		cwdParts = cwdParts[1:len(cwdParts)]
		home = true
	}

	if !home && len(cwdParts) != 0 && len(cwdParts[len(cwdParts)-1]) > 0 {
		segments = append(segments, []string{"250", "237", "/", p.SeparatorThin, "244"})
	} else if !home {
		segments = append(segments, []string{"250", "237", "/"})
	} else {
		segments = append(segments, []string{"015", "031", "~"})
	}

	if len(cwdParts) >= 4 {
		segments = append(segments, []string{"250", "237", p.Ellipsis, p.SeparatorThin, "244"})
	} else if len(cwdParts) > 2 {
		if home {
			segments = append(segments, []string{"250", "237", p.Ellipsis, p.SeparatorThin, "244"})
		} else {
			segments = append(segments, []string{"250", "237", cwdParts[1], p.SeparatorThin, "244"})
		}
	}

	if len(cwdParts) != 0 && len(cwdParts[len(cwdParts)-1]) > 0 {
		segments = append(segments, []string{"250", "237", cwdParts[len(cwdParts)-1]})
	}

	return segments
}

func addVirtulEnvName() []string {
	_, _, virtualEnvName := getVirtualEnv()
	if virtualEnvName != "" {
		return []string{"000", "035", virtualEnvName}
	}

	return nil
}

func addLock(cwd string, p powerline.Powerline) []string {
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

func main() {
	shell := "zsh"

	if len(os.Args) > 1 {
		shell = os.Args[1]
	}

	p := powerline.NewPowerline(shell)
	cwd, cwdParts := getCurrentWorkingDir()

	p.AppendSegments(addCwd(cwd, cwdParts, p))
	p.AppendSegment(addVirtulEnvName())
	p.AppendSegment(addLock(cwd, p))
	p.AppendSegment(addGitInfo())
	p.AppendSegment(addDollarPrompt())

	fmt.Print(p.PrintSegments())
}
