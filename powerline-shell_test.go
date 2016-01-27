package main

import (
        "os"
	"os/user"
	"reflect"
	"strings"
	"testing"
        "github.com/scottweston/powerline-shell-go/powerline-config"
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

        conf.SetDefaults()
	rootSegment := addGitInfo(conf, ">")

        want := append(segments,
          []interface{}{conf.Colours.Git.Text, conf.Colours.Git.BackgroundDefault, "\u2693 master"})

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addGitInfo_not_staged(t *testing.T) {
        var conf config.Configuration
	segments := [][]interface{}{}

        conf.SetDefaults()
	rootSegment := addGitInfo(conf, ">")

        want := append(segments,
	  []interface{}{conf.Colours.Git.Text, conf.Colours.Git.BackgroundDefault, "\u2693 master"})

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addCwd_root(t *testing.T) {
        var conf config.Configuration
        conf.SetDefaults()
	segments := [][]interface{}{}

	dir := "/"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, "...", ">")
	want := append(segments, []interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "/"})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_root_one(t *testing.T) {
        var conf config.Configuration
        conf.SetDefaults()
	segments := [][]interface{}{}

	dir := "/gocode"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, "...", ">")
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

	dir := "/gocode/src"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, "...", ">")
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode", ">", conf.Colours.Cwd.Text},
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

	dir := "/gocode/src/github.com"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf,parts, "...", ">")
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode", ">", conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "...", ">", conf.Colours.Cwd.Text},
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

	dir := "~"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, "...", ">")
	want := append(segments, []interface{}{conf.Colours.Cwd.HomeText, conf.Colours.Cwd.HomeBackground, "~"})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_one(t *testing.T) {
        var conf config.Configuration
        conf.SetDefaults()
	segments := [][]interface{}{}

	dir := "~/gocode"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, "...", ">")
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

	dir := "~/gocode/src"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, "...", ">")
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.HomeText, conf.Colours.Cwd.HomeBackground, "~"},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode", ">", conf.Colours.Cwd.Text},
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

	dir := "~/gocode/src/github.com"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, "...", ">")
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.HomeText, conf.Colours.Cwd.HomeBackground, "~"},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode", ">", conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "...", ">", conf.Colours.Cwd.Text},
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

	dir := "~/gocode/src/github.com/wm/powerline-shell-go"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(conf, parts, "...", ">")
	want := append(
		segments,
		[]interface{}{conf.Colours.Cwd.HomeText, conf.Colours.Cwd.HomeBackground, "~"},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "gocode", ">", conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "...", ">", conf.Colours.Cwd.Text},
		[]interface{}{conf.Colours.Cwd.Text, conf.Colours.Cwd.Background, "powerline-shell-go"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}
