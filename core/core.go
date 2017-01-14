package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/lxc/lxd"
)

func ContainerExec(name string, cmd []string, stdin_text string) string {
	client, err := lxd.NewClient(&lxd.DefaultConfig, "local")
	if err != nil {
		fmt.Println(err)
		return "Connect Error: "
	}
	// Exec
	env := map[string]string{}
	stdout := os.Stdout
	stderr := os.Stderr

	rescueStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = w
	io.WriteString(os.Stdin, stdin_text)
	status_code, _ := client.Exec(name, cmd, env, r, stdout, stderr, nil, 0, 0)
	w.Close()
	os.Stdin = rescueStdin
	return strconv.Itoa(status_code)
}

func Compile(language string, stdin_text string) map[string]string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Compile
	var status_code string
	if language == "clang" {
		cmd := []string{"gcc", "-o", "/tmp/a.out", "/tmp/code.c"}
		status_code = ContainerExec("jessie2", cmd, "")
		if status_code == "0" {
			cmd := []string{"/tmp/a.out"}
			status_code = ContainerExec("jessie2", cmd, stdin_text)
		}
	}

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	//fmt.Printf("Captured: %s", out)
	return map[string]string{"stdout": string(out), "status_code": status_code}
}

func CodePush(code string, extension string) error {
	client, err := lxd.NewClient(&lxd.DefaultConfig, "local")
	if err != nil {
		return err
	}

	extensions := map[string]string{"clang": ".c", "gcc": ".c", "python": ".py"}
	f := bytes.NewReader([]byte(code))
	err = client.PushFile("jessie2", "/tmp/code"+extensions[extension], -1, -1, "", f)
	return err
}
