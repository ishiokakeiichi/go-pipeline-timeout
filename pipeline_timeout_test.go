package pipeline_timeout

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestExec_pipeline(t *testing.T){
	stout,sterr , code , err := Exec("cat  /etc/passwd|egrep ^root")
	if err != nil{t.Error(err.Error())}
	if len(sterr) > 0{t.Error(sterr)}
	if code != 0 {t.Error("not return 0")}
	s := strings.TrimRight(stout, "\n")
	var e string
	switch(runtime.GOOS){
	case "darwin":
		e = "root:*:0:0:System Administrator:/var/root:/bin/sh"
	case "linux":
		e = "root:x:0:0:root:/root:/bin/bash"
	}
	if s != e{
		t.Error(s)
	}
}

func TestExec_pipeline_ENV(t *testing.T){
	os.Setenv("FOO", "root")
	stout,sterr , code , err := Exec("cat /etc/passwd|egrep ^${FOO}", ParseEnv(true))
	if err != nil{t.Error(err.Error())}
	if len(sterr) > 0{t.Error(sterr)}
	if code != 0 {t.Error("not return 0")}
	s := strings.TrimRight(stout, "\n")

	var e string
	switch(runtime.GOOS){
	case "darwin":
		e = "root:*:0:0:System Administrator:/var/root:/bin/sh"
	case "linux":
		e = "root:x:0:0:root:/root:/bin/bash"
	}
	if s != e{
		t.Error(s)
	}
}

func TestExec_ENV(t *testing.T){
	os.Setenv("FOO", "bar")
	stout,sterr , code , err := Exec("echo $FOO", ParseEnv(true))
	if err != nil{t.Error(err.Error())}
	if len(sterr) > 0{t.Error(sterr)}
	if code != 0 {t.Error("not return 0")}
	s := strings.TrimRight(stout, "\n")
	if s != "bar"{
		t.Error(s)
	}
}

func TestExec_NO_ENV(t *testing.T){
	os.Setenv("FOO", "bar")
	stout,sterr , code , err := Exec("echo $FOO")
	if err != nil{t.Error(err.Error())}
	if len(sterr) > 0{t.Error(sterr)}
	if code != 0 {t.Error("not return 0")}
	s := strings.TrimRight(stout, "\n")
	if s != "$FOO"{
		t.Error(s)
	}
}

func TestExec_timeout(t *testing.T) {
	_, _, _, err := Exec("sleep 5;ls -l", Timeout(1))
	if err == nil{
		t.Error("no timeout error!")
	}else{
		if err.Error() != "timeout"{
			t.Error(err.Error())
		}
	}
}

func TestExec_timeout_2(t *testing.T) {
	_, _, _, err := Exec("ls -l|sleep 5;ls -l", Timeout(1))
	if err == nil{
		t.Error("no timeout error!")
	}else{
		if err.Error() != "timeout"{
			t.Error(err.Error())
		}
	}
}

func TestExec_pass_2(t *testing.T) {
	_, _, _, err := Exec("ls -l|sleep 2;ls -l", Timeout(3))
	if err != nil{
		t.Error(err.Error())
	}
}
