package core

import (
	"testing"
)

func TestContainerPushFile(t *testing.T) {
	err := CodePush("hello", "clang")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestContainerExec(t *testing.T) {
	retcode := ContainerExec("jessie2", []string{"id"}, "")
	if retcode != "0" {
		t.Error("Return Code is not 0")
	}
}

func TestCompile(t *testing.T) {
	err := CodePush("#include <stdio.h>\nint main() {\n    printf(\"HELLO\\n\");\n    return 0;\n}\n", "clang")
	if err != nil {
		t.Error("Transfer Error on TestCompile")
	}
	result := Compile("clang", "")
	if result["stdout"] != "HELLO\n" {
		t.Error("Stdout is not excepted")
	}
}
