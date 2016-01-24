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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/scottweston/powerline-shell-go/powerline"
)

func getCurrentWorkingDir() (string, []string) {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	userDir := strings.Replace(dir, os.Getenv("HOME"), "~", 1)
	userDir = strings.TrimSuffix(userDir, "/")
	parts := strings.Split(userDir, "/")
	return dir, parts
}

func getVirtualEnv() string {
	virtualEnv := os.Getenv("VIRTUAL_ENV")
	if virtualEnv == "" {
		return ""
	}

	virtualEnvName := path.Base(virtualEnv)
	return virtualEnvName
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

func addHgInfo(conf Configuration, separator string) [][]string {
	var fmt_str string
	segments := [][]string{}
	branch_colour := conf.Colours.Hg.BackgroundDefault
	text_colour := conf.Colours.Hg.Text

	hg, err := exec.Command("hg", "sum", "--color=never", "-y").Output()

	if err == nil {
		reBranch := regexp.MustCompile(`(?m)^branch: (.*)$`)
		matchBranch := reBranch.FindStringSubmatch(string(hg))

		reModifed := regexp.MustCompile(`(?m)^commit:.* (.*) modified`)
		res_mod := reModifed.FindStringSubmatch(string(hg))
		reUntracked := regexp.MustCompile(`(?m)^commit:.* (.*) unknown`)
		res_untrk := reUntracked.FindStringSubmatch(string(hg))
		reAdded := regexp.MustCompile(`(?m)^commit:.* (.*) added`)
		res_added := reAdded.FindStringSubmatch(string(hg))
		reRemoved := regexp.MustCompile(`(?m)^commit:.* (.*) removed`)
		res_remove := reRemoved.FindStringSubmatch(string(hg))
		reClean := regexp.MustCompile(`(?m)^commit:.*clean`)
		res_clean := reClean.FindStringSubmatch(string(hg))

		if len(res_clean) == 0 {
			branch_colour = conf.Colours.Hg.BackgroundChanges
		}

		// branch name
		if len(matchBranch) > 0 {
			if len(res_added) > 0 || len(res_mod) > 0 || len(res_untrk) > 0 || len(res_remove) > 0 {
				segments = append(segments, []string{text_colour, branch_colour, matchBranch[1], separator, text_colour})
			} else {
				segments = append(segments, []string{text_colour, branch_colour, matchBranch[1]})
			}
		}
		if len(res_added) > 0 {
			if res_added[1] != "1" {
				fmt_str = fmt.Sprintf("%s\u2714", res_added[1])
			} else {
				fmt_str = fmt.Sprintf("\u2714")
			}
			if len(res_mod) > 0 || len(res_untrk) > 0 || len(res_remove) > 0 {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str, separator, text_colour})
			} else {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str})
			}
		}
		if len(res_mod) > 0 {
			if res_mod[1] != "1" {
				fmt_str = fmt.Sprintf("%s\u270e", res_mod[1])
			} else {
				fmt_str = fmt.Sprintf("\u270e")
			}
			if len(res_untrk) > 0 || len(res_remove) > 0 {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str, separator, text_colour})
			} else {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str})
			}
		}
		if len(res_untrk) > 0 {
			if res_untrk[1] != "1" {
				fmt_str = fmt.Sprintf("%s\u272a", res_untrk[1])
			} else {
				fmt_str = fmt.Sprintf("\u272a")
			}
			if len(res_remove) > 0 {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str, separator, text_colour})
			} else {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str})
			}
		}
		if len(res_remove) > 0 {
			if res_remove[1] != "1" {
				fmt_str = fmt.Sprintf("%s\u2620", res_remove[1])
			} else {
				fmt_str = fmt.Sprintf("\u2620")
			}
			segments = append(segments, []string{text_colour, branch_colour, fmt_str})
		}

		return segments
	} else {
		return nil
	}
}

func addGitInfo(conf Configuration, separator string) [][]string {
	var fmt_str string
	segments := [][]string{}
	branch_colour := conf.Colours.Git.BackgroundDefault
	text_colour := conf.Colours.Git.Text

	human, err := exec.Command("git", "status", "--ignore-submodules").Output()
	if err == nil {
		porcelain, _ := exec.Command("git", "status", "--ignore-submodules", "--porcelain").Output()

		// any changes at all?
		staged := !strings.Contains(string(human), "nothing to commit")
		if staged {
			branch_colour = conf.Colours.Git.BackgroundChanges
		}

		// what branch
		reBranch := regexp.MustCompile(`^(HEAD detached at|HEAD detached from|# On branch|On branch) (\S+)`)
		matchBranch := reBranch.FindStringSubmatch(string(human))

		// are we ahead/behind
		reStatus := regexp.MustCompile(`Your branch is (ahead|behind).*?([0-9]+) comm`)
		matchStatus := reStatus.FindStringSubmatch(string(human))

		// added files
		add, _ := regexp.Compile(`(?m)^A  .*$`)
		add_res := add.FindAllString(string(porcelain), -1)

		// modified files
		mod, _ := regexp.Compile(`(?m)^ M .*$`)
		mod_res := mod.FindAllString(string(porcelain), -1)

		// uncommitted files
		uncom, _ := regexp.Compile(`(?m)^\?\? .*$`)
		uncom_res := uncom.FindAllString(string(porcelain), -1)

		// removed files
		del, _ := regexp.Compile(`(?m)^D  .*$`)
		del_res := del.FindAllString(string(porcelain), -1)

		// conflicted files
		cfd, _ := regexp.Compile(`(?m)^(DD|AU|UD|UA|DU|AA|UU) .*$`)
		cfd_res := cfd.FindAllString(string(porcelain), -1)

		// branch name
		if len(matchBranch) > 0 {
			if strings.Contains(matchBranch[1], "detached") {
				fmt_str = fmt.Sprintf("\u2704 %s", matchBranch[2])
			} else {
				fmt_str = fmt.Sprintf("\u2693 %s", matchBranch[2])
			}

			if len(matchStatus) > 0 || len(mod_res) > 0 || len(uncom_res) > 0 || len(del_res) > 0 || len(cfd_res) > 0 {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str, separator, text_colour})
			} else {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str})
			}
		}

		// ahead/behind
		if len(matchStatus) > 0 {
			num, _ := strconv.Atoi(matchStatus[2])

			if matchStatus[1] == "behind" {
				if num > 1 {
					fmt_str = fmt.Sprintf("%s\u25bc", matchStatus[2])
				} else {
					fmt_str = fmt.Sprintf("\u25bc")
				}
			} else if matchStatus[1] == "ahead" {
				if num > 1 {
					fmt_str = fmt.Sprintf("%s\u25b2", matchStatus[2])
				} else {
					fmt_str = fmt.Sprintf("\u25b2")
				}
			} else {
				fmt_str = "unk"
			}

			if len(add_res) > 0 || len(mod_res) > 0 || len(uncom_res) > 0 || len(del_res) > 0 || len(cfd_res) > 0 {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str, separator, text_colour})
			} else {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str})
			}
		}

		// added files
		if len(add_res) > 0 {
			if (len(add_res)) > 1 {
				fmt_str = fmt.Sprintf("%d\u2714", len(add_res))
			} else {
				fmt_str = fmt.Sprintf("\u2714")
			}

			if len(mod_res) > 0 || len(uncom_res) > 0 || len(del_res) > 0 || len(cfd_res) > 0 {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str, separator, text_colour})
			} else {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str})
			}
		}

		// modified files
		if len(mod_res) > 0 {
			if (len(mod_res)) > 1 {
				fmt_str = fmt.Sprintf("%d\u270e", len(mod_res))
			} else {
				fmt_str = fmt.Sprintf("\u270e")
			}

			if len(uncom_res) > 0 || len(del_res) > 0 || len(cfd_res) > 0 {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str, separator, text_colour})
			} else {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str})
			}
		}

		// untracked files
		if len(uncom_res) > 0 {
			if (len(uncom_res)) > 1 {
				fmt_str = fmt.Sprintf("%d\u272a", len(uncom_res))
			} else {
				fmt_str = fmt.Sprintf("\u272a")
			}

			if len(del_res) > 0 || len(cfd_res) > 0 {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str, separator, text_colour})
			} else {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str})
			}
		}

		// deleted files
		if len(del_res) > 0 {
			if (len(del_res)) > 1 {
				fmt_str = fmt.Sprintf("%d\u2620", len(uncom_res))
			} else {
				fmt_str = fmt.Sprintf("\u2620")
			}

			if len(cfd_res) > 0 {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str})
			} else {
				segments = append(segments, []string{text_colour, branch_colour, fmt_str})
			}
		}

		// conflicted files
		if len(cfd_res) > 0 {
			if (len(cfd_res)) > 1 {
				segments = append(segments, []string{text_colour, branch_colour, fmt.Sprintf("%d\u273c", len(cfd_res))})
			} else {
				segments = append(segments, []string{text_colour, branch_colour, fmt.Sprintf("\u273c")})
			}
		}

		return segments
	} else {
		return nil
	}
}

func addCwd(conf Configuration, cwdParts []string, ellipsis string, separator string) [][]string {
	segments := [][]string{}
	back_col := conf.Colours.Cwd.Background
	fore_col := conf.Colours.Cwd.Text

	home := false
	if cwdParts[0] == "~" {
		cwdParts = cwdParts[1:len(cwdParts)]
		home = true
	}

	if home {
		segments = append(segments, []string{conf.Colours.Cwd.HomeText, conf.Colours.Cwd.HomeBackground, "~"})

		if len(cwdParts) > 2 {
			segments = append(segments, []string{fore_col, back_col, cwdParts[0], separator, fore_col})
			segments = append(segments, []string{fore_col, back_col, ellipsis, separator, fore_col})
		} else if len(cwdParts) == 2 {
			segments = append(segments, []string{fore_col, back_col, cwdParts[0], separator, fore_col})
		}
	} else {
		if len(cwdParts[len(cwdParts)-1]) == 0 {
			segments = append(segments, []string{fore_col, back_col, "/"})
		}

		if len(cwdParts) > 3 {
			segments = append(segments, []string{fore_col, back_col, cwdParts[1], separator, fore_col})
			segments = append(segments, []string{fore_col, back_col, ellipsis, separator, fore_col})
		} else if len(cwdParts) > 2 {
			segments = append(segments, []string{fore_col, back_col, cwdParts[1], separator, fore_col})
		}
	}

	if len(cwdParts) != 0 && len(cwdParts[len(cwdParts)-1]) > 0 {
		segments = append(segments, []string{fore_col, back_col, cwdParts[len(cwdParts)-1]})
	}

	return segments
}

func addVirtulEnvName(conf Configuration, virtualEnvName string) []string {
	if virtualEnvName != "" {
		return []string{conf.Colours.Virtualenv.Text, conf.Colours.Virtualenv.Background, virtualEnvName}
	}

	return nil
}

func addReturnCode(conf Configuration, ret_code int) []string {
	if ret_code != 0 {
		return []string{conf.Colours.Returncode.Text, conf.Colours.Returncode.Background, fmt.Sprintf("%d", ret_code)}
	}

	return nil
}

func addLock(conf Configuration, cwd string, lock string) []string {
	if !isWritableDir(cwd) {
		return []string{conf.Colours.Lock.Text, conf.Colours.Lock.Background, lock}
	}

	return nil
}

func addHostname(conf Configuration, includeUsername bool) []string {
	hostname, err := os.Hostname()
	if err != nil {
		return nil
	}

	// create a colour hash for the hostname
	sum := 0
	for _, v := range hostname {
		sum += int(v)
	}

	if includeUsername {
		user, err := user.Current()
		if err != nil {
			return nil
		}
		hostname = user.Username + "@" + hostname
	}

	return []string{"016", fmt.Sprintf("%03d", sum%15), hostname}
}

func addDollarPrompt(conf Configuration) []string {
	return []string{conf.Colours.Dollar.Text, conf.Colours.Dollar.Background, "\\$"}
}

type Configuration struct {
	ShowWritable   bool `json:"showWritable"`
	ShowVirtualEnv bool `json:"showVirtualEnv"`
	ShowCwd        bool `json:"showCwd"`
	ShowGit        bool `json:"showGit"`
	ShowHg         bool `json:"showHg"`
	ShowReturnCode bool `json:"showReturnCode"`
	Colours        struct {
		Hg struct {
			BackgroundDefault string `json:"backgroundDefault"`
			BackgroundChanges string `json:"backgroundChanges"`
			Text              string `json:"text"`
		} `json:"hg"`
		Git struct {
			BackgroundDefault string `json:"backgroundDefault"`
			BackgroundChanges string `json:"backgroundChanges"`
			Text              string `json:"text"`
		} `json:"git"`
		Cwd struct {
			Background     string `json:"background"`
			Text           string `json:"text"`
			HomeBackground string `json:"homeBackground"`
			HomeText       string `json:"homeText"`
		} `json:"cwd"`
		Virtualenv struct {
			Background string `json:"background"`
			Text       string `json:"text"`
		} `json:"virtualenv"`
		Returncode struct {
			Background string `json:"background"`
			Text       string `json:"text"`
		} `json:"returncode"`
		Lock struct {
			Background string `json:"background"`
			Text       string `json:"text"`
		} `json:"lock"`
		Dollar struct {
			Background string `json:"background"`
			Text       string `json:"text"`
		} `json:"dollar"`
	} `json:"colours"`
}

func (self *Configuration) SetDefaults() {
	self.ShowWritable = true
	self.ShowVirtualEnv = true
	self.ShowCwd = true
	self.ShowGit = true
	self.ShowHg = true
	self.ShowReturnCode = true
	self.Colours.Hg.BackgroundDefault = "022"
	self.Colours.Hg.BackgroundChanges = "064"
	self.Colours.Hg.Text = "251"
	self.Colours.Git.BackgroundDefault = "017"
	self.Colours.Git.BackgroundChanges = "021"
	self.Colours.Git.Text = "251"
	self.Colours.Cwd.Background = "040"
	self.Colours.Cwd.Text = "237"
	self.Colours.Cwd.HomeBackground = "031"
	self.Colours.Cwd.HomeText = "015"
	self.Colours.Virtualenv.Background = "035"
	self.Colours.Virtualenv.Text = "000"
	self.Colours.Returncode.Background = "196"
	self.Colours.Returncode.Text = "016"
	self.Colours.Lock.Background = "124"
	self.Colours.Lock.Text = "254"
	self.Colours.Dollar.Background = "240"
	self.Colours.Dollar.Text = "015"
}

func main() {
	var configuration Configuration
	configuration.SetDefaults()
	shell := "bash"
	last_retcode := 0

	user, _ := user.Current()
	data, err := ioutil.ReadFile(user.HomeDir + "/.config/powerline-shell-go/config.json")
	if err == nil {
		err = json.Unmarshal(data, &configuration)
		if err != nil {
			fmt.Println("configuration.error()$ ")
			os.Exit(1)
		}
	}

	if len(os.Args) > 1 {
		shell = os.Args[1]
	}

	if len(os.Args) > 2 {
		last_retcode, _ = strconv.Atoi(os.Args[2])
	}

	p := powerline.NewPowerline(shell)
	cwd, cwdParts := getCurrentWorkingDir()

	if configuration.ShowVirtualEnv {
		p.AppendSegment(addVirtulEnvName(configuration, getVirtualEnv()))
	}
	if _, found := syscall.Getenv("SSH_CLIENT"); found {
		p.AppendSegment(addHostname(configuration, true))
	}
	if configuration.ShowCwd {
		p.AppendSegments(addCwd(configuration, cwdParts, p.Ellipsis, p.SeparatorThin))
	}
	if configuration.ShowWritable {
		p.AppendSegment(addLock(configuration, cwd, p.Lock))
	}
	if configuration.ShowGit {
		p.AppendSegments(addGitInfo(configuration, p.SeparatorThin))
	}
	if configuration.ShowHg {
		p.AppendSegments(addHgInfo(configuration, p.SeparatorThin))
	}
	if configuration.ShowReturnCode {
		p.AppendSegment(addReturnCode(configuration, last_retcode))
	}
	p.AppendSegment(addDollarPrompt(configuration))

	fmt.Print(p.PrintSegments(), " ")
}
