// The package 'git' is a loose wrapper around git's collection of command
// line apps. It's so loose that the parameters most of the methods take
// are just passed in as args to the apps.
package git

import (
	"os"
	"os/exec"
	"bytes"
	"strings"
	//"fmt"
)

// Our type Repo is our interface to a git repository. The reason for this
// being that way is so we don't have to remember our working directory if
// for some mental reason we wind up working on many at once.
//
// The visible attributes Err and Output are the error and combined output
// (stdout, stderr) of the last executed shell command.
type Repo struct {
	Err error
	Output []byte

	dir string
	wd string
}

// Returns a Repo with its dir set. Does no error handling whatsoever.
func GetRepo(dir string) (* Repo) {
	return &Repo{dir: dir}
}

// Repo method that clones the Repo to the specified dir. Returns the 
// cloned repo as a Repo.
func (r *Repo) Clone(arg ...string) (* Repo) {
	arg = prependArg(r.dir, arg)
	r.run("clone", arg...)

	return GetRepo(arg[1])
}

// Repo method that adds listed files to staging. The args are variadic.
func (r *Repo) Add(arg ...string) (error) {
	return r.chdirRun("add", arg...)
}

// Repo method that commits. Only takes the commit message as an arg.
func (r *Repo) Commit(arg string) (error) {
	return r.chdirRun("commit", "-m", arg)
}

// Repo method that adds all untracked files.
func (r *Repo) AddUntracked() (error) {
	return r.addStdout("ls-files", "--others", "--exclude-standard")
}

// Repo method that adds all modified files.
// You're welcome...
func (r *Repo) AddModified() (error) {
	return r.addStdout("ls-files", "-m")
}

// Wrapper around exec calls, builds up the command
func (r* Repo) run(subcmd string, arg ...string) (error) {
	arg = prependArg(subcmd, arg)
	cmd := exec.Command("git", arg...)
	r.Output, r.Err = cmd.CombinedOutput()
	return r.Err
}

// Temporarily switches the working directory before we fire off git
func (r *Repo) chdirRun(subcmd string, arg ...string) (error) {
	r.flameOn()
	r.run(subcmd, arg...)
	r.flameOff()
	return r.Err
}

// Generic function that adds a list of files to the repo
func (r *Repo) addStdout(subcmd string, arg ...string) (error) {
	r.chdirRun(subcmd, arg...)
	if (r.Err != nil) {
		return r.Err
	}
	// Previous output is dropped since the last command was successful

	n := bytes.Index(r.Output, nil)
	files := strings.Split(string(r.Output[:n]), "\n")
	r.Add(files...)
	return r.Err
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
