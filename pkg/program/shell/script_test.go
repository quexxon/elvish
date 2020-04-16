package shell

import (
	"testing"
)

func TestScript_File(t *testing.T) {
	f := setup()
	defer f.cleanup()
	writeFile("a.elv", "echo hello")

	Script(f.fds(), []string{"a.elv"}, &ScriptConfig{})

	f.testOut(t, 1, "hello\n")
}

func TestScript_BadFile(t *testing.T) {
	f := setup()
	defer f.cleanup()

	ret := Script(f.fds(), []string{"a.elv"}, &ScriptConfig{})

	if ret != 2 {
		t.Errorf("got ret %v, want 2", ret)
	}
	f.testOutSnippet(t, 2, "cannot read script")
}

func TestScript_Cmd(t *testing.T) {
	f := setup()
	defer f.cleanup()

	Script(f.fds(), []string{"echo hello"}, &ScriptConfig{Cmd: true})

	f.testOut(t, 1, "hello\n")
}

func TestScript_DoesNotCompile(t *testing.T) {
	f := setup()
	defer f.cleanup()

	ret := Script(f.fds(), []string{"echo $a"}, &ScriptConfig{Cmd: true})

	if ret != 2 {
		t.Errorf("got ret %v, want 2", ret)
	}
	f.testOutSnippet(t, 2, "compilation error")
}

func TestScript_DoesNotCompile_JSON(t *testing.T) {
	f := setup()
	defer f.cleanup()

	ret := Script(f.fds(), []string{"echo $a"},
		&ScriptConfig{Cmd: true, CompileOnly: true, JSON: true})

	if ret != 2 {
		t.Errorf("got ret %v, want 2", ret)
	}
	f.testOutSnippet(t, 1, "variable $a not found")
}

func TestScript_Exception(t *testing.T) {
	f := setup()
	defer f.cleanup()

	ret := Script(f.fds(), []string{"fail failure"}, &ScriptConfig{Cmd: true})
	if ret != 2 {
		t.Errorf("got ret %v, want 2", ret)
	}
	f.testOutSnippet(t, 2, "fail failure")
}

func TestScript_Exception_CompileOnly(t *testing.T) {
	f := setup()
	defer f.cleanup()

	ret := Script(f.fds(), []string{"fail failure"}, &ScriptConfig{
		Cmd: true, CompileOnly: true})
	if ret != 0 {
		t.Errorf("got ret %v, want 0", ret)
	}
}