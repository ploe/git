package git

import (
	"os/exec"
	"bytes"
	"strings"
)

// wrapper around exec calls
func run(subcmd string, arg ...string) ([]byte, error) {
	// Dirty workaround for the fact that a slice can only be exploded from the first element
	buffer := make([]string, 1)
	buffer[0] = subcmd
	arg = append(buffer, arg...)

	cmd := exec.Command("git", arg...)
	return cmd.CombinedOutput()
}

func Clone(arg ...string) ([]byte, error) {
	return run("clone", arg...)
}

func Add(arg ...string) ([]byte, error) {
	return run("add", arg...)
}

// generic function that adds a list of files to the repo
func addStdout(stdout []byte) ([]byte, error) {
	n := bytes.Index(stdout, nil)
	files := strings.Split(string(stdout[:n]), "\n")
	tmp, err := Add(files...)
	return append(stdout, tmp...), err
}

// adds all untracked files to the git repo
func AddUntracked() ([]byte, error) {
	stdout, err := run("ls-files", "--others", "--exclude-standard")
	if(err != nil) {
		return stdout, err
	}
	return addStdout(stdout)
}

func AddModified() ([]byte, error) {
	stdout, err := run("ls-files", "-m")
	if(err != nil) {
		return stdout, err
	}
	return addStdout(stdout)
}

func Commit(arg string) ([]byte, error) {
	return run("commit", "-m", arg)
}
