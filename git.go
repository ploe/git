// The package 'git' is a loose wrapper around git's collection of command
// line apps. It's so loose that the parameters most of the methods take
// are just passed in as args to the apps.
package git

import (
	"os"
	"os/exec"
	"bytes"
	"strings"
//	"fmt"
)

// Our type Repo is our interface to a git repository. The reason for this
// being that way is so we don't have to remember our working directory if
// for some mental reason we wind up working on many at once.
type Repo struct {
	dir string
	wd string
}

// Returns a Repo with its dir set. Does no error handling whatsoever.
func GetRepo(dir string) (* Repo) {
	return &Repo{dir: dir}
}

// Repo method that clones the Repo to the specified dir. You only need to
// specify where you want the Repo cloning to or other command line params
func (r *Repo) Clone(arg ...string) ([]byte, error) {
	arg = prependArg(r.dir, arg)
	return run("clone", arg...)
}

// Repo method that adds listed files to staging
func (r *Repo) Add(arg ...string) ([]byte, error) {
	return r.chdirRun("add", arg...)
}

// Repo method that commits. Only takes the commit message as an arg.
func (r *Repo) Commit(arg string) ([]byte, error) {
	return r.chdirRun("commit", "-m", arg)
}

// Repo method that adds all untracked files.
func (r *Repo) AddUntracked() ([]byte, error) {
	return r.addStdout("ls-files", "--others", "--exclude-standard")
}

// Repo method that adds all modified files.
// You're welcome...
func (r *Repo) AddModified() ([]byte, error) {
	return r.addStdout("ls-files", "-m")
}

// Temporarily switches the working directory before we fire off git
func (r *Repo) chdirRun(subcmd string, arg ...string) ([]byte, error) {
	r.flameOn()
	arg = prependArg(subcmd, arg)
	output, err := exec.Command("git", arg...).CombinedOutput()
	r.flameOff()
	return output, err
}

// Generic function that adds a list of files to the repo
func (r *Repo) addStdout(subcmd string, arg ...string) ([]byte, error) {
	stdout, err := r.chdirRun(subcmd, arg...)
	if (err != nil) {
		return stdout, err
	}

	n := bytes.Index(stdout, nil)
	files := strings.Split(string(stdout[:n]), "\n")
	tmp, err := r.Add(files...)
	return append(stdout, tmp...), err
}

// Changes working directory to Repo's dir but caches the working
// directory, for when we come back
func (r *Repo) flameOn() error {
	wd, err := os.Getwd()
	r.wd = wd
	if (err == nil) {
		err = os.Chdir(r.dir)
	}
	return err
}

func (r *Repo) flameOff() error {
	return os.Chdir(r.wd)
}

/// Dirty workaround for the fact that a slice can only be exploded from 
// the first element.
func prependArg(pre string, arg []string) ([]string) {
	buffer := make([]string, 1)
	buffer[0] = pre
	return append(buffer, arg...)
}

// Wrapper around exec calls, builds up the command
func run(subcmd string, arg ...string) ([]byte, error) {
	arg = prependArg(subcmd, arg)
	cmd := exec.Command("git", arg...)
	return cmd.CombinedOutput()
}
