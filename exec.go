package hancho

import (
  "bytes"
  "io"
  "log"
  "os"
  "os/exec"
)

// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
type ProgressCmd struct {
  cmd *exec.Cmd
}

func progressCommand(name string, arg ...string) *ProgressCmd {
  checkCommand(name)

  cmd := exec.Command(name, arg...)
  pCmd := ProgressCmd{cmd}

  return &pCmd
}

func (p ProgressCmd) EnvAppend(name, value string) {
  p.cmd.Env = append(os.Environ(), name + "=" + value)
}

func (p ProgressCmd) SetDir(dir string) {
  p.cmd.Dir = dir
}

func (p ProgressCmd) Wait() {
  var stdoutBuf, stderrBuf bytes.Buffer
  cmd := p.cmd

  stdoutIn, _ := cmd.StdoutPipe()
  stderrIn, _ := cmd.StderrPipe()

  var errStdout, errStderr error
  stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
  stderr := io.MultiWriter(os.Stderr, &stderrBuf)
  err := cmd.Start()
  if err != nil {
    log.Fatalf("cmd.Start() failed with '%s'\n", err)
  }

  go func() {
    _, errStdout = io.Copy(stdout, stdoutIn)
  }()

  go func() {
    _, errStderr = io.Copy(stderr, stderrIn)
  }()

  err = cmd.Wait()
  if err != nil {
    log.Fatalf("cmd.Run() failed with %s\n", err)
  }
  if errStdout != nil || errStderr != nil {
    log.Fatal("failed to capture stdout or stderr\n")
  }
}

func checkCommand(name string) {
  _, err := exec.LookPath(name)
  checkError(err)
}
