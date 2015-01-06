package main

import (
	"os"
	"os/user"
	"reflect"
	"strings"
	"testing"
)

func Test_addHostname_with_username(t *testing.T) {
	hostname, _ := os.Hostname()
	user, _ := user.Current()

	rootSegment := addHostname(true)
	want := []string{"015", "161", user.Username + "@" + hostname}

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addHostname_without_username(t *testing.T) {
	hostname, _ := os.Hostname()

	rootSegment := addHostname(false)
	want := []string{"015", "161", hostname}

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addVirtualEnvName_empty(t *testing.T) {
	var want []string
	rootSegment := addVirtulEnvName("")

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addVirtualEnvName_present(t *testing.T) {
	rootSegment := addVirtulEnvName("MyVirtEnv")
	want := []string{"000", "035", "MyVirtEnv"}

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addGitInfo_no_status(t *testing.T) {
	var want []string
	rootSegment := addGitInfo("", false)

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addGitInfo_not_staged(t *testing.T) {
	var want = []string{"000", "148", "master"}
	rootSegment := addGitInfo("master", false)

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addGitInfo_staged(t *testing.T) {
	var want = []string{"015", "161", "master"}
	rootSegment := addGitInfo("master", true)

	if !reflect.DeepEqual(rootSegment, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegment, want)
	}
}

func Test_addCwd_root(t *testing.T) {
	segments := [][]string{}

	dir := "/"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(segments, []string{"250", "237", "/"})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_root_one(t *testing.T) {
	segments := [][]string{}

	dir := "/gocode"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"250", "237", "gocode"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_root_two(t *testing.T) {
	segments := [][]string{}

	dir := "/gocode/src"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"250", "237", "gocode", ">", "244"},
		[]string{"250", "237", "src"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_root_three(t *testing.T) {
	segments := [][]string{}

	dir := "/gocode/src/github.com"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"250", "237", "gocode", ">", "244"},
		[]string{"250", "237", "...", ">", "244"},
		[]string{"250", "237", "github.com"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home(t *testing.T) {
	segments := [][]string{}

	dir := "~"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(segments, []string{"015", "031", "~"})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_one(t *testing.T) {
	segments := [][]string{}

	dir := "~/gocode"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"015", "031", "~"},
		[]string{"250", "237", "gocode"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_two(t *testing.T) {
	segments := [][]string{}

	dir := "~/gocode/src"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"015", "031", "~"},
		[]string{"250", "237", "gocode", ">", "244"},
		[]string{"250", "237", "src"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_three(t *testing.T) {
	segments := [][]string{}

	dir := "~/gocode/src/github.com"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"015", "031", "~"},
		[]string{"250", "237", "gocode", ">", "244"},
		[]string{"250", "237", "...", ">", "244"},
		[]string{"250", "237", "github.com"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_five(t *testing.T) {
	segments := [][]string{}

	dir := "~/gocode/src/github.com/wm/powerline-shell-go"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"015", "031", "~"},
		[]string{"250", "237", "gocode", ">", "244"},
		[]string{"250", "237", "...", ">", "244"},
		[]string{"250", "237", "powerline-shell-go"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}
