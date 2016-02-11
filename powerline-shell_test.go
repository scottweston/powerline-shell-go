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

	rootSegment := addHostname(conf, true, false)
	want := []interface{}{16, 12, user.Username + "@" + hostname}

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addHostname_without_username(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	hostname, _ := os.Hostname()

	rootSegment := addHostname(conf, false, false)
	want := []interface{}{16, 12, hostname}

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addVirtualEnvName_empty(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	var want []interface{}
	rootSegment := addVirtulEnvName(conf, "")

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addVirtualEnvName_present(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	rootSegment := addVirtulEnvName(conf, "MyVirtEnv")
	want := []interface{}{0, 35, "MyVirtEnv"}

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addGitInfo_no_status(t *testing.T) {
	var conf config.Configuration
	segments := [][]interface{}{}

	var porc string = `## master
`

	p := powerline.NewPowerline("bash", false)

	conf.SetDefaults()
	rootSegment := addGitInfo(conf, porc, p)

	want := append(segments,
		[]interface{}{conf.Colours.Git.Text, conf.Colours.Git.BackgroundDefault, "master"})

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addGitInfo_not_staged(t *testing.T) {
	var conf config.Configuration
	segments := [][]interface{}{}

	var porc string = `## master
 M modifed.go
A  added.go
D  deleted.go
DD conflicted.go
?? not_staged.go
`

	p := powerline.NewPowerline("bash", false)

	conf.SetDefaults()
	rootSegment := addGitInfo(conf, porc, p)

	want := append(segments,
		[]interface{}{conf.Colours.Git.Text, conf.Colours.Git.BackgroundChanges, "master", p.SeparatorThin, conf.Colours.Git.Text},
		[]interface{}{conf.Colours.Git.Text, conf.Colours.Git.BackgroundChanges, p.Added, p.SeparatorThin, conf.Colours.Git.Text},
		[]interface{}{conf.Colours.Git.Text, conf.Colours.Git.BackgroundChanges, p.Modified, p.SeparatorThin, conf.Colours.Git.Text},
		[]interface{}{conf.Colours.Git.Text, conf.Colours.Git.BackgroundChanges, p.Untracked, p.SeparatorThin, conf.Colours.Git.Text},
		[]interface{}{conf.Colours.Git.Text, conf.Colours.Git.BackgroundChanges, "2" + p.Removed, p.SeparatorThin, conf.Colours.Git.Text},
		[]interface{}{conf.Colours.Git.Text, conf.Colours.Git.BackgroundChanges, p.Conflicted},
	)

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addCwd_root(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	segments := [][]interface{}{}

	p := powerline.NewPowerline("bash", false)

	dir := "/"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, p)
	want := append(segments, []interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "/"})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_root_one(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	segments := [][]interface{}{}

	p := powerline.NewPowerline("bash", false)

	dir := "/gocode"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, p)
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_root_two(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	segments := [][]interface{}{}

	p := powerline.NewPowerline("bash", false)

	dir := "/gocode/src"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, p)
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode", p.SeparatorThin, conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "src"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_root_three(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	segments := [][]interface{}{}

	p := powerline.NewPowerline("bash", false)

	dir := "/gocode/src/github.com"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, p)
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode", p.SeparatorThin, conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, p.Ellipsis, p.SeparatorThin, conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "github.com"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	segments := [][]interface{}{}

	p := powerline.NewPowerline("bash", false)

	dir := "~"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, p)
	want := append(segments, []interface{}{conf.Colours.Cwd.HomeText, conf.Colours.Cwd.HomeBackground, "~"})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_one(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	segments := [][]interface{}{}

	p := powerline.NewPowerline("bash", false)

	dir := "~/gocode"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, p)
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.HomeText, conf.Colours.Cwd.HomeBackground, "~"},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_two(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	segments := [][]interface{}{}

	p := powerline.NewPowerline("bash", false)

	dir := "~/gocode/src"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, p)
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.HomeText, conf.Colours.Cwd.HomeBackground, "~"},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode", p.SeparatorThin, conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "src"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_three(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	segments := [][]interface{}{}

	p := powerline.NewPowerline("bash", false)

	dir := "~/gocode/src/github.com"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, p)
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.HomeText, conf.Colours.Cwd.HomeBackground, "~"},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode", p.SeparatorThin, conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, p.Ellipsis, p.SeparatorThin, conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "github.com"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_five(t *testing.T) {
	var conf config.Configuration
	conf.SetDefaults()
	segments := [][]interface{}{}

	p := powerline.NewPowerline("bash", false)

	dir := "~/gocode/src/github.com/wm/powerline-shell-go"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, p)
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.HomeText, conf.Colours.Cwd.HomeBackground, "~"},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode", p.SeparatorThin, conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, p.Ellipsis, p.SeparatorThin, conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "powerline-shell-go"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}
