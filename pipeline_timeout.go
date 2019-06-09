package pipeline_timeout

import (
"bytes"
"errors"
"github.com/mattn/go-shellwords"
"os/exec"
"strings"
"syscall"
"time"
"unsafe"
)

// timeout 処理用
type Timeout time.Duration
//標準入力処理用
type Stdin string
type ParseEnv bool
type execOptions struct {
	timeout  time.Duration
	stdin    string
	parseEnv bool
}

//
//  コマンドのステータスコードのセット
//
func setExecCode(err error, status *int){
	if e2, ok := err.(*exec.ExitError); ok{
		if s, ok := e2.Sys().(syscall.WaitStatus); ok {
			tmp := s.ExitStatus()
			status = &tmp
		}
	}
}

//
//  コマンドの実行
//
func timeoutCommand(commands [][]string, timeout time.Duration, stdin string)(string, string, int, error){
	var status = -1 // デフォルトの返却ステータス
	var err error

	start := time.Now()
	tmpTimeout := timeout
	execCommand := make([]*exec.Cmd, len(commands))
	for i, c := range commands {
		execCommand[i] = exec.Command(c[0], c[1:]...)
		if i == 0 && len(stdin) > 0{
			execCommand[i].Stdin = strings.NewReader(stdin)
		}
		if i > 0 {
			if execCommand[i].Stdin, err = execCommand[i-1].StdoutPipe(); err != nil{
				return "", "", status, err
			}
		}
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	execCommand[len(execCommand)-1].Stdout = &stdout
	execCommand[len(execCommand)-1].Stderr = &stderr
	for _, c := range execCommand {
		if err := c.Start(); err != nil{
			return "", string(err.Error()), status,err
		}
	}

	if timeout == 0{
		for i, c := range execCommand {
			if err = c.Wait(); err != nil{
				setExecCode(err, &status)
				if i == len(execCommand) -1{
					out := stdout.Bytes()
					eout := stderr.Bytes()
					return *(*string)(unsafe.Pointer(&out)), *(*string)(unsafe.Pointer(&eout)), status, err
				}else {
					return "", string(err.Error()), status, err
				}
			}
		}
	}else{
		for i, c := range execCommand {
			if tmpTimeout <=0 {
				for _, c := range execCommand {
					_ = c.Process.Kill()
				}
				return "", "timeout", status, errors.New("timeout")
			}
			err2 := make(chan error)
			go func(){err2 <- c.Wait()}()
			select{
			case err := <-err2:
				timeout -= time.Duration(time.Since(start))
				if err != nil{
					setExecCode(err, &status)
					if i == len(execCommand) -1{
						out := stdout.Bytes()
						eout := stderr.Bytes()
						return *(*string)(unsafe.Pointer(&out)), *(*string)(unsafe.Pointer(&eout)), status, err
					}else {
						return "", string(err.Error()), status, err
					}
				}
			case <-time.After(timeout):
				for _, c := range execCommand {
					_ = c.Process.Kill()
				}
				return "", "timeout", status, errors.New("timeout")
			}
		}
	}

	out := stdout.Bytes()
	eout := stderr.Bytes()
	return *(*string)(unsafe.Pointer(&out)), *(*string)(unsafe.Pointer(&eout)), 0, err
}


func Exec(commandline string, options ...interface{})(string, string, int, error){
	// コマンドに "&" 禁止 TODO

	// 引数に shellPipe.Timeout(int)が指定された場合
	// execOpt.timeout に timeout時間を指定する
	execOpt := execOptions{}
	execOpt.timeout = time.Duration(0 * time.Second)
	execOpt.stdin   = ""
	execOpt.parseEnv = false
	for _,o := range options{
		switch v:= o.(type){
		case Timeout:
			execOpt.timeout = time.Duration(v) * time.Duration(time.Second)
		case Stdin:
			execOpt.stdin = string(v)
		case ParseEnv:
			execOpt.parseEnv = bool(v)
		}
	}

	var status = -1
	splitPipe := strings.Split(commandline,"|")
	commands := make([][]string,len(splitPipe))

	parser := shellwords.NewParser()
	parser.ParseBacktick = true
	parser.ParseEnv = execOpt.parseEnv
	for i, cmd := range splitPipe {
		args ,err := parser.Parse(cmd)
		if err != nil {
			return "", string(err.Error()), status, err
		}
		commands[i] = args
	}


	return timeoutCommand(commands, execOpt.timeout, execOpt.stdin)


}
