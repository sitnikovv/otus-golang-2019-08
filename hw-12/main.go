package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	//	Проверяем окличество аргументов, их должно быть как минимум два - указание каталога и указание программы для запуска
	//	В случае конфуза выводим инструкцию
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println(Usage())
		return
	}

	//	Получаем мапу данных для переменных окружения, в случае ошибки сообщаем причину и выходим
	data, err := SetEnvDir(args[0])
	if err != nil {
		log.Fatalf("Error read directiry, reason: %s", err)
	}

	////	Выставляем переменные окружения, в случае ошибки сообщаем причину и выходим
	//if err = SetEnvironments(data); err != nil {
	//	log.Fatalf("Cannot set environments, reason: %s", err)
	//}

	//	Запускаем указанную программу, в случае ошибки запуска сообщаем причину и выходим. Дожидаться пока отработает - не будем, не царское это дело
	program := exec.Command(args[1], args[2:]...)
	program.Stdout = os.Stdout
	program.Env = SetEnvironmentsForCmd(data)
	if err = program.Start(); err != nil {
		log.Fatalf("Unable to start program %s with arguments %s, reason: %s", args[1], strings.Join(args[2:], " "), err)
	}
}

//	Читаем содержимой файлов из укащанного каталога в мапу, где имена файлов являются ключами, а содержимое файлов - значением
func SetEnvDir(dir string) (map[string]string, error) {
	env := make(map[string]string)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return env, err
	}
	for _, file := range files {
		if content, err := ioutil.ReadFile(dir + "/" + file.Name()); err != nil {
			return env, err
		} else {
			env[file.Name()] = string(content)
		}
	}
	return env, err
}

//	Выставляем переменные окружения из мапы с данными
func SetEnvironments(data map[string]string) error {
	for key, value := range data {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return nil
}

//	Выставляем переменные окружения из мапы с данными
func SetEnvironmentsForCmd(data map[string]string) []string {
	result := make([]string, 0, len(data))
	for key, value := range data {
		result = append(result, key + "=" + value)
	}
	return result
}

//	Возврат текста инструкции
func Usage() string {
	return `envdir  - runs another program with environment modified according to files in a specified directory.
Usage: envdir d child
       d is a single argument.  child consists of one or more arguments.

envdir sets various environment variables as specified by files in the directory named  d. It then runs child.`
}