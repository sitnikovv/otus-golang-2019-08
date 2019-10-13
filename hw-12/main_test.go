package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestSetEnvDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatalf("Fail to create emporary dir, reason: %s", err)
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	tmpfile, err := os.Create(dir + "/my_env")
	if err != nil {
		t.Fatalf("Fail to create file in temporary dir, reason: %s", err)
	}
	defer func() {
		_ = tmpfile.Close()
	}()

	if length, err := tmpfile.Write([]byte("my_value")); err != nil || length != 8 {
		t.Fatalf("Fail to write to file in temporary dir, reason: %s", err)
	}
	result, err := SetEnvDir(dir)
	if err != nil {
		t.Fatalf("Fail SetEnvDir func, reason: %s", err)
	}
	require.Equal(t, len(result), 1, "Incorrect get environment count: waiting 1, getting %v", len(result))
	require.Equal(t, result["my_env"], "my_value", "Incorrect environment value: waiting 'my_value', getting '%v'", result["my_env"])
}

func TestSetEnvironments(t *testing.T) {

	if err := os.Unsetenv("MY_TEST_VAR"); err != nil {
		t.Fatalf("Unable to delete environment variable, reason: %s", err)
	}
	defer func() {
		_ = os.Unsetenv("MY_TEST_VAR")
	}()

	if err := SetEnvironments(map[string]string{"MY_TEST_VAR": "test_value"}); err != nil {
		t.Fatalf("Unable to set environment variable, reason: %s", err)
	}

	if value, ok := os.LookupEnv("MY_TEST_VAR"); !ok {
		t.Fatalf("Variable is not set")
	} else {
		require.Equal(t, value, "test_value", "Incorrect environment value: waiting 'test_value', getting '%v'", value)
	}
}

func TestUsage(t *testing.T) {
	text := Usage()
	require.NotEqual(t, len(text), 0, "WTF???")
}

func TestRunFromOS(t *testing.T) {
	var (
		cmd *exec.Cmd
		err error
		script string
		output []byte
	)

	//	Компилим
	cmd = exec.Command("go", "build", "main.go")
	err = cmd.Run()
	require.Nil(t, err, "Unable to compile, reason '%v'", err)

	//	Переименовываем
	if runtime.GOOS == "windows" {
		script = "@echo %test_env% passed"
		err = os.Rename("main.exe", "envdir.exe")
	} else {
		script = "#!/bin/sh\necho $test_env passed"
		err = os.Rename("main", "envdir")
		if err != nil {
			err = os.Chmod("envdir", 0755)
		}
	}
	require.Nil(t, err, "Unable to rename executable file, reason '%v'", err)

	//	Созлаём каталог test, если его нет
	if _, err := os.Stat("test"); os.IsNotExist(err) {
		err = os.Mkdir("test", 0755)
		require.Nil(t, err, "Unable to create test dir, reason '%v'", err)
	}
	defer func() {
		_ = os.RemoveAll("test")
	}()

	//	Создаём файл проверочных данных и скрипт, задача которого - вывести переменную окружения
	err = ioutil.WriteFile("script.cmd", []byte(script), 0755)
	if err == nil {
		err = ioutil.WriteFile("test/test_env", []byte("test"), 0755)
	}
	require.Nil(t, err, "Unable to create test data for test, reason '%v'", err)
	defer func() {
		_ = os.Remove("script.cmd")
	}()

	//	Запускаем скрипт средствами операционки, проверяем вывод
	if runtime.GOOS == "windows" {
		output, err = exec.Command("envdir", "test", "script.cmd").Output()
	} else {
		output, err = exec.Command("./envdir", "test", "./script.cmd").Output()
	}
	require.Nil(t, err, "Unable to run test reason '%s'", err)
	script = strings.Trim(string(output), "\n\r")
	require.Equal(t, script, "test passed", "Incorrect result: waiting 'test passed', getting '%s'", script)
}
