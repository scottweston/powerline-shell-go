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

// Helpers

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

// Segment generators

func addHgInfo(conf config.Configuration, p powerline.Powerline) *powerline.Segment {
	var fmt_str string

	segment := powerline.Segment{}

	branch_colour := conf.Colours.Hg.BackgroundDefault
	text_colour := conf.Colours.Hg.Text

	hg, err := exec.Command("hg", "sum", "--color=never", "-y").Output()

	if err == nil {
		// branch:
		reBranch := regexp.MustCompile(`(?m)^branch: (.*)$`)
		matchBranch := reBranch.FindStringSubmatch(string(hg))

		// commit:
		// %d modified, %d added, %d removed, %d renamed, %d copied
		// %d deleted, %d unknown, %d unresolved, %d subrepos
		reModifed := regexp.MustCompile(`(?m)^commit:.* (.*) modified`)
		res_mod := reModifed.FindStringSubmatch(string(hg))
		reUntracked := regexp.MustCompile(`(?m)^commit:.* (.*) unknown`)
		res_untrk := reUntracked.FindStringSubmatch(string(hg))
		reAdded := regexp.MustCompile(`(?m)^commit:.* (.*) added`)
		res_added := reAdded.FindStringSubmatch(string(hg))
		reRemoved := regexp.MustCompile(`(?m)^commit:.* (.*) removed`)
		res_remove := reRemoved.FindStringSubmatch(string(hg))
		reClean := regexp.MustCompile(`(?m)^commit:.*\(clean\)`)
		res_clean := reClean.FindStringSubmatch(string(hg))

		// update:
		reUpdate := regexp.MustCompile(`(?m)^update: (.*) new`)
		res_update := reUpdate.FindStringSubmatch(string(hg))

		// phases:
		rePublic := regexp.MustCompile(`(?m)^phases:.* (.*) public`)
		res_public := rePublic.FindStringSubmatch(string(hg))
		reDraft := regexp.MustCompile(`(?m)^phases:.* (.*) draft`)
		res_draft := reDraft.FindStringSubmatch(string(hg))
		reSecret := regexp.MustCompile(`(?m)^phases:.* (.*) secret`)
		res_secret := reSecret.FindStringSubmatch(string(hg))

		if len(res_clean) == 0 {
			branch_colour = conf.Colours.Hg.BackgroundChanges
		}

		segment.Background = branch_colour
		segment.Foreground = text_colour
		segment.Weight = conf.Weights.Segments.Hg

		// branch name
		if len(matchBranch) > 0 {
			if matchBranch[1] != "default" {
				fmt_str = p.Branch + " " + matchBranch[1]
			} else {
				fmt_str = matchBranch[1]
			}
			segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Branch})
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
			segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Phases})
		}

		// updated files
		if len(res_update) > 0 {
			if res_update[1] != "1" {
				fmt_str = fmt.Sprintf("%s%s", res_update[1], p.Behind)
			} else {
				fmt_str = p.Behind
			}
			segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Sync})
		}

		// added files
		if len(res_added) > 0 {
			if res_added[1] != "1" {
				fmt_str = fmt.Sprintf("%s%s", res_added[1], p.Added)
			} else {
				fmt_str = p.Added
			}
			segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Added})
		}

		// modified files
		if len(res_mod) > 0 {
			if res_mod[1] != "1" {
				fmt_str = fmt.Sprintf("%s%s", res_mod[1], p.Modified)
			} else {
				fmt_str = p.Modified
			}
			segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Modified})
		}

		// untracked files
		if len(res_untrk) > 0 {
			if res_untrk[1] != "1" {
				fmt_str = fmt.Sprintf("%s%s", res_untrk[1], p.Untracked)
			} else {
				fmt_str = p.Untracked
			}
			segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Untracked})
		}

		// removed files
		if len(res_remove) > 0 {
			if res_remove[1] != "1" {
				fmt_str = fmt.Sprintf("%s%s", res_remove[1], p.Removed)
			} else {
				fmt_str = p.Removed
			}
			segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Removed})
		}

		return &segment
	} else {
		return nil
	}
}

func addGitInfo(conf config.Configuration, porcelain string, p powerline.Powerline) *powerline.Segment {
	var fmt_str string

	segment := powerline.Segment{}

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

	// renamed files
	rename, _ := regexp.Compile(`(?m)^R. `)
	rename_res := rename.FindAllString(porcelain, -1)

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
	if len(rename_res) > 0 || len(add_res) > 0 || len(mod_res) > 0 || len(uncom_res) > 0 || len(del_res) > 0 || len(cfd_res) > 0 {
		branch_colour = conf.Colours.Git.BackgroundChanges
	}

	segment.Background = branch_colour
	segment.Foreground = text_colour
	segment.Weight = conf.Weights.Segments.Git

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
		segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Branch})
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
		segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Sync})
	}

	// renamed files
	if len(rename_res) > 0 {
		if (len(rename_res)) > 1 {
			fmt_str = fmt.Sprintf("%d%s", len(rename_res), p.Renamed)
		} else {
			fmt_str = p.Renamed
		}
		segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Renamed})
	}

	// added files
	if len(add_res) > 0 {
		if (len(add_res)) > 1 {
			fmt_str = fmt.Sprintf("%d%s", len(add_res), p.Added)
		} else {
			fmt_str = p.Added
		}
		segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Added})
	}

	// modified files
	if len(mod_res) > 0 {
		if (len(mod_res)) > 1 {
			fmt_str = fmt.Sprintf("%d%s", len(mod_res), p.Modified)
		} else {
			fmt_str = p.Modified
		}
		segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Modified})
	}

	// untracked files
	if len(uncom_res) > 0 {
		if (len(uncom_res)) > 1 {
			fmt_str = fmt.Sprintf("%d%s", len(uncom_res), p.Untracked)
		} else {
			fmt_str = p.Untracked
		}
		segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Untracked})
	}

	// deleted files
	if len(del_res) > 0 {
		if (len(del_res)) > 1 {
			fmt_str = fmt.Sprintf("%d%s", len(del_res), p.Removed)
		} else {
			fmt_str = p.Removed
		}
		segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Deleted})
	}

	// conflicted files
	if len(cfd_res) > 0 {
		if (len(cfd_res)) > 1 {
			fmt_str = fmt.Sprintf("%d%s", len(cfd_res), p.Conflicted)
		} else {
			fmt_str = p.Conflicted
		}
		segment.Parts = append(segment.Parts, powerline.Part{Text: fmt_str, Weight: conf.Weights.Parts.Conflicted})
	}

	return &segment
}

func addCwd(conf config.Configuration, cwdParts []string, p powerline.Powerline) []powerline.Segment {
	segment := []powerline.Segment{}

	back_col := conf.Colours.Cwd.Background
	fore_col := conf.Colours.Cwd.Text

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

	// are we under our home?
	if cwdParts[0] == "~" {
		segment = append(segment, powerline.Segment{Foreground: conf.Colours.Cwd.HomeText, Background: conf.Colours.Cwd.HomeBackground, Weight: conf.Weights.Segments.Cwd})
		segment[len(segment)-1].Parts = append(segment[len(segment)-1].Parts, powerline.Part{Text: cwdParts[0]})
		cwdParts = cwdParts[1:]
	}

	if len(cwdParts) == 0 {
		return segment
	}

	if cwdParts[0] == "" {
		if len(cwdParts) > 1 {
			cwdParts = cwdParts[1:]
		}
		cwdParts[0] = "/" + cwdParts[0]
	}

	segment = append(segment, powerline.Segment{Foreground: fore_col, Background: back_col, Weight: conf.Weights.Segments.Cwd})
	segment[len(segment)-1].Parts = append(segment[len(segment)-1].Parts, powerline.Part{Text: cwdParts[0]})
	cwdParts = cwdParts[1:]

	// if there's only one more we show it, otherwise it's an ellipsis then the last part
	if len(cwdParts) == 1 {
		segment[len(segment)-1].Parts = append(segment[len(segment)-1].Parts, powerline.Part{Text: cwdParts[0]})
	} else if len(cwdParts) > 1 {
		segment[len(segment)-1].Parts = append(segment[len(segment)-1].Parts, powerline.Part{Text: p.Ellipsis})
		segment[len(segment)-1].Parts = append(segment[len(segment)-1].Parts, powerline.Part{Text: cwdParts[len(cwdParts)-1]})
	}

	return segment
}

func addVirtulEnvName(conf config.Configuration, virtualEnvName string) *powerline.Segment {
	if virtualEnvName != "" {
		segment := powerline.Segment{Foreground: conf.Colours.Virtualenv.Text, Background: conf.Colours.Virtualenv.Background, Weight: conf.Weights.Segments.Virtualenv}
		segment.Parts = append(segment.Parts, powerline.Part{Text: virtualEnvName})
		return &segment
	}
	return nil
}

func addReturnCode(conf config.Configuration, ret_code int) *powerline.Segment {
	if ret_code != 0 {
		segment := powerline.Segment{Foreground: conf.Colours.Returncode.Text, Background: conf.Colours.Returncode.Background, Weight: conf.Weights.Segments.Returncode}
		segment.Parts = append(segment.Parts, powerline.Part{Text: fmt.Sprintf("%d", ret_code)})
		return &segment
	}
	return nil
}

func addLock(conf config.Configuration, cwd string, p powerline.Powerline) *powerline.Segment {
	if !IsWritableDir(cwd) {
		segment := powerline.Segment{Foreground: conf.Colours.Lock.Text, Background: conf.Colours.Lock.Background, Weight: conf.Weights.Segments.Lock}
		segment.Parts = append(segment.Parts, powerline.Part{Text: p.ReadOnly})
		return &segment
	}
	return nil
}

func addHostname(conf config.Configuration, includeUsername bool, hostHash bool, p powerline.Powerline) *powerline.Segment {
	hostname, err := os.Hostname()
	if err != nil {
		return nil
	}

	if len(hostname) > conf.HostnameMaxLength {
		sml := int(conf.HostnameMaxLength/2 - 1)
		if sml > 0 {
			hostname = hostname[0:sml] + p.Ellipsis + hostname[len(hostname)-sml:]
		}
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
		if conf.HostnameMaxLength > 0 {
			hostname = user.Username + "@" + hostname
		} else {
			hostname = user.Username
		}
	}

	segment := powerline.Segment{Foreground: 16, Background: back, Weight: conf.Weights.Segments.Hostname}
	segment.Parts = append(segment.Parts, powerline.Part{Text: hostname})
	return &segment
}

func addBatteryWarn(conf config.Configuration) *powerline.Segment {
	battery, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/capacity")
	if err == nil {
		capacity, _ := strconv.Atoi(strings.Trim(string(battery), " \n"))
		if capacity <= conf.BatteryWarn {
			segment := powerline.Segment{Foreground: conf.Colours.Battery.Text, Background: conf.Colours.Battery.Background, Weight: conf.Weights.Segments.Battery}
			segment.Parts = append(segment.Parts, powerline.Part{Text: fmt.Sprintf("%d%%", capacity)})
			return &segment
		}
	}
	return nil
}

func addDollarPrompt(conf config.Configuration, dollar string) *powerline.Segment {
	segment := powerline.Segment{Foreground: conf.Colours.Dollar.Text, Background: conf.Colours.Dollar.Background, Weight: -1000}
	segment.Parts = append(segment.Parts, powerline.Part{Text: dollar})
	return &segment
}

func main() {
	var configuration config.Configuration
	var set_title string = ""
	configuration.SetDefaults()
	shell := "bash"
	last_retcode := 0

	user, err := user.Current()
	var data []byte
	if err == nil {
		data, err = ioutil.ReadFile(user.HomeDir + "/.config/powerline-shell-go/config.json")
	} else if home, found := syscall.Getenv("HOME"); found {
		data, err = ioutil.ReadFile(home + "/.config/powerline-shell-go/config.json")
	}
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
		last_retcode, err = strconv.Atoi(os.Args[2])
		if err != nil {
			if os.Args[2] == "install" {
				if shell == "bash" {
					fmt.Println(`function _update_ps1() { export PS1="$(powerline-shell-go bash $? 2> /dev/null)"; };
export PROMPT_COMMAND="_update_ps1; $PROMPT_COMMAND";`)
				} else if shell == "zsh" {
					fmt.Println(`function powerline_precmd() { export PS1="$(powerline-shell-go zsh $? 2> /dev/null)"; };
function install_powerline_precmd() { for s in "${precmd_functions[@]}"; do; if [ "$s" = "powerline_precmd" ]; then; return; fi; done; precmd_functions+=(powerline_precmd); };
install_powerline_precmd;`)
				} else {
					fmt.Printf("echo Unsupported shell: %s;\n", shell)
				}
				os.Exit(0)
			}
		}
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
		if configuration.Icons.Powerline.Renamed != "" {
			p.Renamed = configuration.Icons.Powerline.Renamed
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
		if configuration.Icons.Plain.Renamed != "" {
			p.Renamed = configuration.Icons.Plain.Renamed
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
		p.AppendSegment(addHostname(configuration, true, true, p))
	}
	if configuration.ShowCwd {
		parts := addCwd(configuration, cwdParts, p)
		for _, element := range parts {
			p.AppendSegment(&element)
		}
	}
	if configuration.ShowWritable {
		p.AppendSegment(addLock(configuration, cwd, p))
	}
	if configuration.ShowGit {
		porcelain, err := exec.Command("git", "status", "--ignore-submodules", "-b", "--porcelain").Output()
		if err == nil {
			p.AppendSegment(addGitInfo(configuration, string(porcelain), p))
		}
	}
	if configuration.ShowHg {
		p.AppendSegment(addHgInfo(configuration, p))
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
