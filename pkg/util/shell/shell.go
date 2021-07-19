package shell

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/rs/zerolog/log"
)

//DefaultFailedCode constant
const DefaultFailedCode = 1

// RedirectToConsole Show command execution through the standard output or
// standard error Exit if error
func RedirectToConsole(out []byte, err error) {
	if err != nil {
		if out != nil {
			log.Error().Msg(string(out))
		}
		log.Fatal().Err(err).Msg(err.Error())
	}
	log.Info().Msg(string(out))
}

// RunCommand Run a shell command
func RunCommand(name string, args ...string) (stdout string, stderr string, exitCode int) {
	log.Debug().
		Msg("Run command: " + name + strings.Join(args, ""))
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Warn().
				Msg("Could not get exit code for failed program: " + name + strings.Join(args, ""))
			exitCode = DefaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	log.Debug().
		Str("stdout", stdout).
		Str("stderr", stderr).
		Int("exitCode", exitCode).
		Msg("command result")
	return
}

// GetFullPath Get Full Path
func GetFullPath(relativePath string) string {
	currentPath, _ := os.Getwd()
	fullPathTemplate := path.Join(currentPath, relativePath)
	return fullPathTemplate
}

// ExistsFileOrDirectory Exists File or Directory
func ExistsFileOrDirectory(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// FindFilesByPattern Find Files in a directory by pattern
func FindFilesByPattern(targetDir string, pattern []string) ([]string, error) {

	for _, v := range pattern {
		matches, err := filepath.Glob(path.Join(targetDir, v))

		if err != nil {
			return nil, err
		}

		if len(matches) != 0 {
			return matches, nil
		}
	}

	return nil, nil
}

// RemoveDirectory Remove a directory
func RemoveDirectory(directory string) error {
	d, err := os.Open(directory)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(directory, name))
		if err != nil {
			return err
		}
	}

	return os.Remove(directory)
}
