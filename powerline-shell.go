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
	"github.com/scottweston/powerline-shell-go/powerline-config"
)

var build string

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

func addHgInfo(conf config.Configuration, p powerline.Powerline) [][]interface{} {
	var fmt_str string
	segments := [][]interface{}{}
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

		rePublic := regexp.MustCompile(`(?m)^phases:.* (.*) public`)
		res_public := rePublic.FindStringSubmatch(string(hg))
		reDraft := regexp.MustCompile(`(?m)^phases:.* (.*) draft`)
		res_draft := reDraft.FindStringSubmatch(string(hg))
		reSecret := regexp.MustCompile(`(?m)^phases:.* (.*) secret`)
		res_secret := reSecret.FindStringSubmatch(string(hg))

		if len(res_clean) == 0 {
			branch_colour = conf.Colours.Hg.BackgroundChanges
		}

		// branch name
		if len(matchBranch) > 0 {
			if matchBranch[1] != "default" {
				fmt_str = p.Branch + " " + matchBranch[1]
			} else {
				fmt_str = matchBranch[1]
			}
			if len(res_added) > 0 || len(res_mod) > 0 || len(res_untrk) > 0 || len(res_remove) > 0 || len(res_public) > 0 || len(res_draft) > 0 || len(res_secret) > 0 {
				segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str, p.SeparatorThin, text_colour})
			} else {
				segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
			}
		}
		// phases
		if len(res_public) > 0 || len(res_draft) > 0 || len(res_secret) > 0 {
			var public int = 0
			var draft int = 0
			var secret int = 0
			if len(res_public) > 0 {
				public, _ = strconv.Atoi(res_public[1])
			}
			if len(res_draft) > 0 {
				draft, _ = strconv.Atoi(res_draft[1])
			}
			if len(res_secret) > 0 {
				secret, _ = strconv.Atoi(res_secret[1])
			}
			total := public + draft + secret
			if total == 1 {
				fmt_str = p.Phases
			} else {
				fmt_str = fmt.Sprintf("%d%s", total, p.Phases)
			}
			if len(res_added) > 0 || len(res_mod) > 0 || len(res_untrk) > 0 || len(res_remove) > 0 {
				segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str, p.SeparatorThin, text_colour})
			} else {
				segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
			}

		}
		if len(res_added) > 0 {
			if res_added[1] != "1" {
				fmt_str = fmt.Sprintf("%s%s", res_added[1], p.Added)
			} else {
				fmt_str = p.Added
			}
			if len(res_mod) > 0 || len(res_untrk) > 0 || len(res_remove) > 0 {
				segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str, p.SeparatorThin, text_colour})
			} else {
				segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
			}
		}
		if len(res_mod) > 0 {
			if res_mod[1] != "1" {
				fmt_str = fmt.Sprintf("%s%s", res_mod[1], p.Modified)
			} else {
				fmt_str = p.Modified
			}
			if len(res_untrk) > 0 || len(res_remove) > 0 {
				segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str, p.SeparatorThin, text_colour})
			} else {
				segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
			}
		}
		if len(res_untrk) > 0 {
			if res_untrk[1] != "1" {
				fmt_str = fmt.Sprintf("%s%s", res_untrk[1], p.Untracked)
			} else {
				fmt_str = p.Untracked
			}
			if len(res_remove) > 0 {
				segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str, p.SeparatorThin, text_colour})
			} else {
				segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
			}
		}
		if len(res_remove) > 0 {
			if res_remove[1] != "1" {
				fmt_str = fmt.Sprintf("%s%s", res_remove[1], p.Removed)
			} else {
				fmt_str = p.Removed
			}
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
		}

		return segments
	} else {
		return nil
	}
}

func addGitInfo(conf config.Configuration, porcelain string, p powerline.Powerline) [][]interface{} {
	var fmt_str string
	segments := [][]interface{}{}
	branch_colour := conf.Colours.Git.BackgroundDefault
	text_colour := conf.Colours.Git.Text

	// what branch
	reBranch := regexp.MustCompile(`(?m)^## (([^ \.\n]+).*|.* on (\S+))$`)
	matchBranch := reBranch.FindStringSubmatch(porcelain)

	// detached?
	reDetached := regexp.MustCompile(`(?m)^## .* \(no branch\)`)
	matchDetached := reDetached.FindStringSubmatch(porcelain)

	// are we ahead/behind
	reStatus := regexp.MustCompile(`(?m)^## .* \[(ahead|behind) ([0-9]+)\]`)
	matchStatus := reStatus.FindStringSubmatch(porcelain)

	// added files
	add, _ := regexp.Compile(`(?m)^A. `)
	add_res := add.FindAllString(porcelain, -1)

	// modified files
	mod, _ := regexp.Compile(`(?m)^.M `)
	mod_res := mod.FindAllString(porcelain, -1)

	// uncommitted files
	uncom, _ := regexp.Compile(`(?m)^\?\? `)
	uncom_res := uncom.FindAllString(porcelain, -1)

	// removed files
	del, _ := regexp.Compile(`(?m)^(D.|.D) `)
	del_res := del.FindAllString(porcelain, -1)

	// conflicted files
	cfd, _ := regexp.Compile(`(?m)^DD|AU|UD|UA|DU|AA|UU .*$`)
	cfd_res := cfd.FindAllString(porcelain, -1)

	// any changes at all?
	if len(add_res) > 0 || len(mod_res) > 0 || len(uncom_res) > 0 || len(del_res) > 0 || len(cfd_res) > 0 {
		branch_colour = conf.Colours.Git.BackgroundChanges
	}

	// branch name
	if len(matchBranch) > 0 {
		if len(matchDetached) > 0 {
			fmt_str = p.Detached + " "
		} else {
			fmt_str = ""
		}
		if matchBranch[2] != "master" {
			fmt_str = fmt.Sprintf("%s%s ", fmt_str, p.Branch)
		}
		fmt_str = fmt.Sprintf("%s%s", fmt_str, matchBranch[2])

		if len(matchStatus) > 0 || len(add_res) > 0 || len(mod_res) > 0 || len(uncom_res) > 0 || len(del_res) > 0 || len(cfd_res) > 0 {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str, p.SeparatorThin, text_colour})
		} else {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
		}
	}

	// ahead/behind
	if len(matchStatus) > 0 {
		num, _ := strconv.Atoi(matchStatus[2])

		if matchStatus[1] == "behind" {
			if num > 1 {
				fmt_str = fmt.Sprintf("%s%s", matchStatus[2], p.Behind)
			} else {
				fmt_str = p.Behind
			}
		} else if matchStatus[1] == "ahead" {
			if num > 1 {
				fmt_str = fmt.Sprintf("%s%s", matchStatus[2], p.Ahead)
			} else {
				fmt_str = p.Ahead
			}
		} else {
			fmt_str = "unk"
		}

		if len(add_res) > 0 || len(mod_res) > 0 || len(uncom_res) > 0 || len(del_res) > 0 || len(cfd_res) > 0 {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str, p.SeparatorThin, text_colour})
		} else {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
		}
	}

	// added files
	if len(add_res) > 0 {
		if (len(add_res)) > 1 {
			fmt_str = fmt.Sprintf("%d%s", len(add_res), p.Added)
		} else {
			fmt_str = p.Added
		}

		if len(mod_res) > 0 || len(uncom_res) > 0 || len(del_res) > 0 || len(cfd_res) > 0 {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str, p.SeparatorThin, text_colour})
		} else {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
		}
	}

	// modified files
	if len(mod_res) > 0 {
		if (len(mod_res)) > 1 {
			fmt_str = fmt.Sprintf("%d%s", len(mod_res), p.Modified)
		} else {
			fmt_str = p.Modified
		}

		if len(uncom_res) > 0 || len(del_res) > 0 || len(cfd_res) > 0 {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str, p.SeparatorThin, text_colour})
		} else {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
		}
	}

	// untracked files
	if len(uncom_res) > 0 {
		if (len(uncom_res)) > 1 {
			fmt_str = fmt.Sprintf("%d%s", len(uncom_res), p.Untracked)
		} else {
			fmt_str = p.Untracked
		}

		if len(del_res) > 0 || len(cfd_res) > 0 {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str, p.SeparatorThin, text_colour})
		} else {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
		}
	}

	// deleted files
	if len(del_res) > 0 {
		if (len(del_res)) > 1 {
			fmt_str = fmt.Sprintf("%d%s", len(del_res), p.Removed)
		} else {
			fmt_str = p.Removed
		}

		if len(cfd_res) > 0 {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str, p.SeparatorThin, text_colour})
		} else {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt_str})
		}
	}

	// conflicted files
	if len(cfd_res) > 0 {
		if (len(cfd_res)) > 1 {
			segments = append(segments, []interface{}{text_colour, branch_colour, fmt.Sprintf("%d%s", len(cfd_res), p.Conflicted)})
		} else {
			segments = append(segments, []interface{}{text_colour, branch_colour, p.Conflicted})
		}
	}

	return segments
}

func addCwd(conf config.Configuration, cwdParts []string, p powerline.Powerline) [][]interface{} {
	segments := [][]interface{}{}
	back_col := conf.Colours.Cwd.Background
	fore_col := conf.Colours.Cwd.Text

	home := false
	if cwdParts[0] == "~" {
		cwdParts = cwdParts[1:len(cwdParts)]
		home = true
	}

	// limit part length, less than 3 makes no sense
	if conf.CwdMaxLength > 3 {
		for i, part := range cwdParts {
			if len(part) > conf.CwdMaxLength {
				sml := int(conf.CwdMaxLength/2 - 1)
				if sml > 0 {
					cwdParts[i] = part[0:sml] + p.Ellipsis + part[len(part)-sml:]
				}
			}
		}
	}

	if home {
		segments = append(segments, []interface{}{conf.Colours.Cwd.HomeText, conf.Colours.Cwd.HomeBackground, "~"})

		if len(cwdParts) > 2 {
			segments = append(segments, []interface{}{fore_col, back_col, cwdParts[0], p.SeparatorThin, fore_col})
			segments = append(segments, []interface{}{fore_col, back_col, p.Ellipsis, p.SeparatorThin, fore_col})
		} else if len(cwdParts) == 2 {
			segments = append(segments, []interface{}{fore_col, back_col, cwdParts[0], p.SeparatorThin, fore_col})
		}
	} else {
		if len(cwdParts[len(cwdParts)-1]) == 0 {
			segments = append(segments, []interface{}{fore_col, back_col, "/"})
		}

		if len(cwdParts) > 3 {
			segments = append(segments, []interface{}{fore_col, back_col, cwdParts[1], p.SeparatorThin, fore_col})
			segments = append(segments, []interface{}{fore_col, back_col, p.Ellipsis, p.SeparatorThin, fore_col})
		} else if len(cwdParts) > 2 {
			segments = append(segments, []interface{}{fore_col, back_col, cwdParts[1], p.SeparatorThin, fore_col})
		}
	}

	if len(cwdParts) != 0 && len(cwdParts[len(cwdParts)-1]) > 0 {
		segments = append(segments, []interface{}{fore_col, back_col, cwdParts[len(cwdParts)-1]})
	}

	return segments
}

func addVirtulEnvName(conf config.Configuration, virtualEnvName string) []interface{} {
	if virtualEnvName != "" {
		return []interface{}{conf.Colours.Virtualenv.Text, conf.Colours.Virtualenv.Background, virtualEnvName}
	}

	return nil
}

func addReturnCode(conf config.Configuration, ret_code int) []interface{} {
	if ret_code != 0 {
		return []interface{}{conf.Colours.Returncode.Text, conf.Colours.Returncode.Background, fmt.Sprintf("%d", ret_code)}
	}

	return nil
}

func addLock(conf config.Configuration, cwd string, p powerline.Powerline) []interface{} {
	if !isWritableDir(cwd) {
		return []interface{}{conf.Colours.Lock.Text, conf.Colours.Lock.Background, p.ReadOnly}
	}

	return nil
}

func addHostname(conf config.Configuration, includeUsername bool, hostHash bool) []interface{} {
	hostname, err := os.Hostname()
	if err != nil {
		return nil
	}

	back := 12

	if hostHash {
		// create a colour hash for the hostname
		sum := 0
		for _, v := range hostname {
			sum += int(v)
		}
		back = sum % 15
	}

	if includeUsername {
		user, err := user.Current()
		if err != nil {
			return nil
		}
		hostname = user.Username + "@" + hostname
	}

	return []interface{}{16, back, hostname}
}

func addBatteryWarn(conf config.Configuration) []interface{} {
	battery, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/capacity")
	if err == nil {
		capacity, _ := strconv.Atoi(strings.Trim(string(battery), " \n"))
		if capacity <= conf.BatteryWarn {
			return []interface{}{conf.Colours.Battery.Text, conf.Colours.Battery.Background, fmt.Sprintf("%d%%", capacity)}
		}
	}
	return nil
}

func addDollarPrompt(conf config.Configuration, dollar string) []interface{} {
	return []interface{}{conf.Colours.Dollar.Text, conf.Colours.Dollar.Background, dollar}
}

func main() {
	var configuration config.Configuration
	var set_title string = ""
	configuration.SetDefaults()
	shell := "bash"
	last_retcode := 0

	user, _ := user.Current()
	data, err := ioutil.ReadFile(user.HomeDir + "/.config/powerline-shell-go/config.json")
	if err == nil {
		err = json.Unmarshal(data, &configuration)
		if err != nil {
			fmt.Printf("configuration error(%s)> ", err)
			os.Exit(1)
		}
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "version" || os.Args[1] == "build" {
			if build != "" {
				fmt.Println(build)
			} else {
				fmt.Println("unknown")
			}
			os.Exit(0)
		} else {
			shell = os.Args[1]
		}
	}

	if len(os.Args) > 2 {
		last_retcode, _ = strconv.Atoi(os.Args[2])
	}

	if shell != "bash" && shell != "zsh" {
		fmt.Printf("unsupported shell(%s)> ", shell)
		os.Exit(1)
	}

	var p powerline.Powerline
	if _, found := syscall.Getenv("LC_POWERLINE"); found {
		p = powerline.NewPowerline(shell, true)
		if configuration.Icons.Powerline.Added != "" {
			p.Added = configuration.Icons.Powerline.Added
		}
		if configuration.Icons.Powerline.Ahead != "" {
			p.Ahead = configuration.Icons.Powerline.Ahead
		}
		if configuration.Icons.Powerline.Behind != "" {
			p.Behind = configuration.Icons.Powerline.Behind
		}
		if configuration.Icons.Powerline.Branch != "" {
			p.Branch = configuration.Icons.Powerline.Branch
		}
		if configuration.Icons.Powerline.Conflicted != "" {
			p.Conflicted = configuration.Icons.Powerline.Conflicted
		}
		if configuration.Icons.Powerline.Detached != "" {
			p.Detached = configuration.Icons.Powerline.Detached
		}
		if configuration.Icons.Powerline.Ellipsis != "" {
			p.Ellipsis = configuration.Icons.Powerline.Ellipsis
		}
		if configuration.Icons.Powerline.Modified != "" {
			p.Modified = configuration.Icons.Powerline.Modified
		}
		if configuration.Icons.Powerline.Phases != "" {
			p.Phases = configuration.Icons.Powerline.Phases
		}
		if configuration.Icons.Powerline.ReadOnly != "" {
			p.ReadOnly = configuration.Icons.Powerline.ReadOnly
		}
		if configuration.Icons.Powerline.Removed != "" {
			p.Removed = configuration.Icons.Powerline.Removed
		}
		if configuration.Icons.Powerline.SeparatorThin != "" {
			p.SeparatorThin = configuration.Icons.Powerline.SeparatorThin
		}
		if configuration.Icons.Powerline.Separator != "" {
			p.Separator = configuration.Icons.Powerline.Separator
		}
		if configuration.Icons.Powerline.Untracked != "" {
			p.Untracked = configuration.Icons.Powerline.Untracked
		}
	} else {
		p = powerline.NewPowerline(shell, false)
		if configuration.Icons.Plain.Added != "" {
			p.Added = configuration.Icons.Plain.Added
		}
		if configuration.Icons.Plain.Ahead != "" {
			p.Ahead = configuration.Icons.Plain.Ahead
		}
		if configuration.Icons.Plain.Behind != "" {
			p.Behind = configuration.Icons.Plain.Behind
		}
		if configuration.Icons.Plain.Branch != "" {
			p.Branch = configuration.Icons.Plain.Branch
		}
		if configuration.Icons.Plain.Conflicted != "" {
			p.Conflicted = configuration.Icons.Plain.Conflicted
		}
		if configuration.Icons.Plain.Detached != "" {
			p.Detached = configuration.Icons.Plain.Detached
		}
		if configuration.Icons.Plain.Ellipsis != "" {
			p.Ellipsis = configuration.Icons.Plain.Ellipsis
		}
		if configuration.Icons.Plain.Modified != "" {
			p.Modified = configuration.Icons.Plain.Modified
		}
		if configuration.Icons.Plain.Phases != "" {
			p.Phases = configuration.Icons.Plain.Phases
		}
		if configuration.Icons.Plain.ReadOnly != "" {
			p.ReadOnly = configuration.Icons.Plain.ReadOnly
		}
		if configuration.Icons.Plain.Removed != "" {
			p.Removed = configuration.Icons.Plain.Removed
		}
		if configuration.Icons.Plain.SeparatorThin != "" {
			p.SeparatorThin = configuration.Icons.Plain.SeparatorThin
		}
		if configuration.Icons.Plain.Separator != "" {
			p.Separator = configuration.Icons.Plain.Separator
		}
		if configuration.Icons.Plain.Untracked != "" {
			p.Untracked = configuration.Icons.Plain.Untracked
		}
	}
	cwd, cwdParts := getCurrentWorkingDir()

	if term, found := syscall.Getenv("TERM"); found {
		if strings.Contains(term, "xterm") || strings.Contains(term, "rxvt") {
			set_title = p.SetTitle
		}
	}

	if configuration.ShowVirtualEnv {
		p.AppendSegment(addVirtulEnvName(configuration, getVirtualEnv()))
	}
	if _, found := syscall.Getenv("SSH_CLIENT"); found {
		p.AppendSegment(addHostname(configuration, true, true))
	}
	if configuration.ShowCwd {
		p.AppendSegments(addCwd(configuration, cwdParts, p))
	}
	if configuration.ShowWritable {
		p.AppendSegment(addLock(configuration, cwd, p))
	}
	if configuration.ShowGit {
		porcelain, err := exec.Command("git", "status", "--ignore-submodules", "-b", "--porcelain").Output()
		if err == nil {
			p.AppendSegments(addGitInfo(configuration, string(porcelain), p))
		}
	}
	if configuration.ShowHg {
		p.AppendSegments(addHgInfo(configuration, p))
	}
	if configuration.ShowReturnCode {
		p.AppendSegment(addReturnCode(configuration, last_retcode))
	}
	if configuration.BatteryWarn > 0 {
		p.AppendSegment(addBatteryWarn(configuration))
	}
	p.AppendSegment(addDollarPrompt(configuration, p.Dollar))

	fmt.Print(set_title, p.PrintSegments(), " ")
}
