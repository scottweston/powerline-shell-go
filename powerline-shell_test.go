package main

import (
	"github.com/scottweston/powerline-shell-go/powerline"
	"github.com/scottweston/powerline-shell-go/powerline-config"
	"os"
	"os/user"
	"reflect"
	"strings"
	"testing"
)

func Test_addHostname_with_username(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	hostname, _ := os.Hostname()
	user, _ := user.Current()

	p := powerline.NewPowerline("bash", false)

	rootSegment := addHostname(conf, true, false, p)
	var parts []powerline.Part
	parts = append(parts, powerline.Part{Text: user.Username + "@" + hostname, Dirty: true})
	want := powerline.Segment{Foreground: 16, Background: 12,
		Parts: parts}

	if !reflect.DeepEqual(rootSegment, &want) {
		t.Errorf("addHostname returned:\n  %+v\nnot:\n  %+v", rootSegment, &want)
	}
}

func Test_addVirtualEnvName_empty(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	var want *powerline.Segment
	rootSegment := addVirtulEnvName(conf, "")

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addVirtualEnvName returned:\n  %+v\nnot:\n  %+v", rootSegment, want)
	}
}

func Test_addVirtualEnvName_present(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	rootSegment := addVirtulEnvName(conf, "MyVirtEnv")
	var parts []powerline.Part
	parts = append(parts, powerline.Part{Text: "MyVirtEnv", Dirty: true})
	want := powerline.Segment{Foreground: conf.Colours.Virtualenv.Text, Background: conf.Colours.Virtualenv.Background, Parts: parts}

	if !reflect.DeepEqual(rootSegment, &want) {
		t.Errorf("addVirtualEnvName returned:\n  %+v\nnot:\n  %+v", rootSegment, &want)
	}
}

func Test_addGitInfo_no_status(t *testing.T) {
	var conf config.Configuration

	var porc string = `## master...origin/master
`

	p := powerline.NewPowerline("bash", false)

	conf.SetDefaults()
	rootSegment := addGitInfo(conf, porc, p)

	var parts []powerline.Part
	parts = append(parts, powerline.Part{Text: "master", Dirty: true})
	want := powerline.Segment{Foreground: conf.Colours.Git.Text,
		Background: conf.Colours.Git.BackgroundDefault,
		Parts:      parts}

	if !reflect.DeepEqual(rootSegment, &want) {
		t.Errorf("addGitInfo_no_status returned:\n  %+v\nnot:\n  %+v", rootSegment, &want)
	}
}

func Test_addGitInfo_not_staged(t *testing.T) {
	var conf config.Configuration

	var porc string = `## master...origin/master
 M modifed.go
A  added.go
D  deleted.go
DD conflicted.go
?? not_staged.go
`

	p := powerline.NewPowerline("bash", false)

	conf.SetDefaults()
	rootSegment := addGitInfo(conf, porc, p)

	var parts []powerline.Part
	parts = append(parts, powerline.Part{Text: "master", Dirty: true})
	parts = append(parts, powerline.Part{Text: p.Added, Dirty: true})
	parts = append(parts, powerline.Part{Text: p.Modified, Dirty: true})
	parts = append(parts, powerline.Part{Text: p.Untracked, Dirty: true})
	parts = append(parts, powerline.Part{Text: "2" + p.Removed, Dirty: true})
	parts = append(parts, powerline.Part{Text: p.Conflicted, Dirty: true})
	want := powerline.Segment{Foreground: conf.Colours.Git.Text,
		Background: conf.Colours.Git.BackgroundChanges,
		Parts:      parts}

	if !reflect.DeepEqual(rootSegment, &want) {
		t.Errorf("addGitInfo_not_staged returned:\n  %+v\nnot:\n  %+v", rootSegment, &want)
	}
}

func Test_addCwd_root(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()

	p := powerline.NewPowerline("bash", false)

	dir := "/"
	cwdparts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, cwdparts, p)

	var parts []powerline.Part
	var want []powerline.Segment
	parts = append(parts, powerline.Part{Text: "/", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.Text,
		Background: conf.Colours.Cwd.Background,
		Parts:      parts})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd_root returned:\n  %+v\nnot:\n %+v", rootSegments, want)
	}
}

func Test_addCwd_root_one(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()

	p := powerline.NewPowerline("bash", false)

	dir := "/gocode"
	cwdparts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, cwdparts, p)

	var parts []powerline.Part
	var want []powerline.Segment
	parts = append(parts, powerline.Part{Text: "/gocode", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.Text,
		Background: conf.Colours.Cwd.Background,
		Parts:      parts})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd_root_one returned:\n  %+v\nnot:\n  %+v", rootSegments, want)
	}
}

func Test_addCwd_root_two(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()

	p := powerline.NewPowerline("bash", false)

	dir := "/gocode/src"
	cwdparts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, cwdparts, p)

	var parts []powerline.Part
	var want []powerline.Segment
	parts = append(parts, powerline.Part{Text: "/gocode", Dirty: true})
	parts = append(parts, powerline.Part{Text: "src", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.Text,
		Background: conf.Colours.Cwd.Background,
		Parts:      parts})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd_root_two returned:\n  %+v\nnot:\n  %+v", rootSegments, want)
	}
}

func Test_addCwd_root_three(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()

	p := powerline.NewPowerline("bash", false)

	dir := "/gocode/src/github.com"
	cwdparts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, cwdparts, p)

	var parts []powerline.Part
	var want []powerline.Segment
	parts = append(parts, powerline.Part{Text: "/gocode", Dirty: true})
	parts = append(parts, powerline.Part{Text: p.Ellipsis, Dirty: false})
	parts = append(parts, powerline.Part{Text: "github.com", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.Text,
		Background: conf.Colours.Cwd.Background,
		Parts:      parts})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd_root_three returned:\n  %+v\nnot:\n  %+v", rootSegments, want)
	}
}

func Test_addCwd_home(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()

	p := powerline.NewPowerline("bash", false)

	dir := "~"
	cwdparts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, cwdparts, p)

	var parts []powerline.Part
	var want []powerline.Segment
	parts = append(parts, powerline.Part{Text: "~", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.HomeText,
		Background: conf.Colours.Cwd.HomeBackground,
		Parts:      parts})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd_home returned:\n  %+v\nnot:\n  %+v", rootSegments, want)
	}
}

func Test_addCwd_home_one(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()

	p := powerline.NewPowerline("bash", false)

	dir := "~/gocode"
	cwdparts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, cwdparts, p)

	var parts []powerline.Part
	var want []powerline.Segment
	parts = append(parts, powerline.Part{Text: "~", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.HomeText,
		Background: conf.Colours.Cwd.HomeBackground,
		Parts:      parts})
	var subparts []powerline.Part
	subparts = append(subparts, powerline.Part{Text: "gocode", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.Text,
		Background: conf.Colours.Cwd.Background,
		Parts:      subparts})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd_home_one returned:\n  %+v\nnot:\n  %+v", rootSegments, want)
	}
}

func Test_addCwd_home_two(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()

	p := powerline.NewPowerline("bash", false)

	dir := "~/gocode/src"
	cwdparts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, cwdparts, p)

	var parts []powerline.Part
	var want []powerline.Segment
	parts = append(parts, powerline.Part{Text: "~", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.HomeText,
		Background: conf.Colours.Cwd.HomeBackground,
		Parts:      parts})
	var subparts []powerline.Part
	subparts = append(subparts, powerline.Part{Text: "gocode", Dirty: true})
	subparts = append(subparts, powerline.Part{Text: "src", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.Text,
		Background: conf.Colours.Cwd.Background,
		Parts:      subparts})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd_home_two returned:\n  %+v\nnot:\n  %+v", rootSegments, want)
	}
}

func Test_addCwd_home_three(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()

	p := powerline.NewPowerline("bash", false)

	dir := "~/gocode/src/github.com"
	cwdparts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, cwdparts, p)

	var parts []powerline.Part
	var want []powerline.Segment
	parts = append(parts, powerline.Part{Text: "~", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.HomeText,
		Background: conf.Colours.Cwd.HomeBackground,
		Parts:      parts})
	var subparts []powerline.Part
	subparts = append(subparts, powerline.Part{Text: "gocode", Dirty: true})
	subparts = append(subparts, powerline.Part{Text: p.Ellipsis})
	subparts = append(subparts, powerline.Part{Text: "github.com", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.Text,
		Background: conf.Colours.Cwd.Background,
		Parts:      subparts})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd_home_three returned:\n  %+v\nnot:\n  %+v", rootSegments, want)
	}
}

func Test_addCwd_home_five(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()

	p := powerline.NewPowerline("bash", false)

	dir := "~/gocode/src/github.com/wm/powerline-shell-go"
	cwdparts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, cwdparts, p)

	var parts []powerline.Part
	var want []powerline.Segment
	parts = append(parts, powerline.Part{Text: "~", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.HomeText,
		Background: conf.Colours.Cwd.HomeBackground,
		Parts:      parts})
	var subparts []powerline.Part
	subparts = append(subparts, powerline.Part{Text: "gocode", Dirty: true})
	subparts = append(subparts, powerline.Part{Text: p.Ellipsis})
	subparts = append(subparts, powerline.Part{Text: "power…ll-go", Dirty: true})
	want = append(want, powerline.Segment{Foreground: conf.Colours.Cwd.Text,
		Background: conf.Colours.Cwd.Background,
		Parts:      subparts})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd_home_five returned:\n  %+v\nnot:\n  %+v", rootSegments, want)
	}
}

// vim: ts=8 sw=8 smartindent noexpandtab:
