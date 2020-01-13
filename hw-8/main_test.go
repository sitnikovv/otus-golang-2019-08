package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func Init() {
	numThreads   = 5
	maxError     = 7
}

func success() error {
	return nil
}

func fail() error {
	return errors.New("FAIL")
}

func TestRunWithIncorrectMaxError(t *testing.T) {
	var err error

	//	Максимальное количество ошибок меньше, чем количество потоков
	//	Вариант, когда максимальное количество ошибок больше чем количество заданий ошибкой не считаю
	//	Вариант, когда количество потоков больше чем тасков также ошибкой не считаю
	tasks := []func() error{success}
	err = Run(tasks, 5, 4)
	require.NotEqual(t, err, nil, "Incorrect errors: waiting error, getting nil")
}

func TestRunWithEmptyTasks(t *testing.T) {
	var err error

	//	Таски отсутсвуют
	tasks := make([]func() error, 0)
	err = Run(tasks, numThreads, maxError)
	require.NotEqual(t, err, nil, "Incorrect errors: waiting error, getting nil")
}

func TestRunWithNilTasks(t *testing.T) {
	var err error

	//	Таски не то что отсутствуют, они просто не существуют
	var tasks []func() error
	err = Run(tasks, numThreads, maxError)
	require.NotEqual(t, err, nil, "Incorrect errors: waiting error, getting nil")
}

func TestRunAllSuccessful(t *testing.T) {
	var err error

	//	Все успешные
	tasks := []func() error{success, success, success, success, success, success, success, success, success, success}
	err = Run(tasks, numThreads, maxError)
	msg := "All successful:"
	require.Equal(t, err, nil, "Incorrect errors: waiting nil, getting %v", err)
	require.Equal(t, countSuccess, numTasks, "%v Incorrect successful: waiting %v, getting %v", msg, numTasks, countSuccess)
	require.Equal(t, countFail, 0, "%v Incorrect fail: waiting %v, getting %v", msg, 0, countFail)
}

//	Здесь и далее смысле проверять error нет
func TestRunAllFail(t *testing.T) {
	//	Все неуспешные
	tasks := []func() error{fail, fail, fail, fail, fail, fail, fail, fail, fail, fail}
	msg := "All fail:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 0, "%v Incorrect successful: waiting %v, getting %v", msg, 0, countSuccess)
	require.Equal(t, countFail, maxError, "%v Incorrect fail: waiting %v, getting %v", msg, maxError, countFail)

}

//	Проверяем граничные значения
func TestRunWithBorder(t *testing.T) {
	//	Успешные в начале
	tasks := []func() error{success, success, success, fail, fail, fail, fail, fail, fail, fail}
	msg := "Successful begin:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 3, "%v Incorrect successful: waiting %v, getting %v", msg, 3, countSuccess)
	require.Equal(t, countFail, 7, "%v Incorrect fail: waiting %v, getting %v", msg, 7, countFail)

	//	Успешные в конце
	tasks = []func() error{fail, fail, fail, fail, fail, fail, fail, success, success, success}
	msg = "Successful end:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 0, "%v Incorrect successful: waiting %v, getting %v", msg, 0, countSuccess)
	require.Equal(t, countFail, 7, "%v Incorrect fail: waiting %v, getting %v", msg, 7, countFail)

	//	Успешные внутри
	tasks = []func() error{fail, fail, fail, fail, success, success, success, fail, fail, fail}
	msg = "Successful inner:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 3, "%v Incorrect successful: waiting %v, getting %v", msg, 3, countSuccess)
	require.Equal(t, countFail, 7, "%v Incorrect fail: waiting %v, getting %v", msg, 7, countFail)

	//	Успешные снаружи
	tasks = []func() error{success, success, fail, fail, fail, fail, fail, fail, fail, success}
	msg = "Successful outer:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 2, "%v Incorrect successful: waiting %v, getting %v", msg, 2, countSuccess)
	require.Equal(t, countFail, 7, "%v Incorrect fail: waiting %v, getting %v", msg, 7, countFail)

}

//	Проверяем на не граничные значения
func TestRunNoBorder(t *testing.T) {

	//	Успешные в начале
	tasks := []func() error{success, success, success, success, fail, fail, fail, fail, fail, fail}
	msg := "Successful begin:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 4, "%v Incorrect successful: waiting %v, getting %v", msg, 4, countSuccess)
	require.Equal(t, countFail, 6, "%v Incorrect fail: waiting %v, getting %v", msg, 6, countFail)

	//	Успешные в конце
	tasks = []func() error{fail, fail, fail, fail, fail, fail, success, success, success, success}
	msg = "Successful end:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 4, "%v Incorrect successful: waiting %v, getting %v", msg, 4, countSuccess)
	require.Equal(t, countFail, 6, "%v Incorrect fail: waiting %v, getting %v", msg, 6, countFail)

	//	Успешные внутри
	tasks = []func() error{fail, fail, fail, success, success, success, success, fail, fail, fail}
	msg = "Successful inner:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 4, "%v Incorrect successful: waiting %v, getting %v", msg, 4, countSuccess)
	require.Equal(t, countFail, 6, "%v Incorrect fail: waiting %v, getting %v", msg, 6, countFail)

	//	Успешные снаружи
	tasks = []func() error{success, success, fail, fail, fail, fail, fail, fail, success, success}
	msg = "Successful outer:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 4, "%v Incorrect successful: waiting %v, getting %v", msg, 4, countSuccess)
	require.Equal(t, countFail, 6, "%v Incorrect fail: waiting %v, getting %v", msg, 6, countFail)

	//	Один успешный в начале
	tasks = []func() error{success, fail, fail, fail, fail, fail, fail, fail, fail, fail}
	msg = "Successful first:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 1, "%v Incorrect successful: waiting %v, getting %v", msg, 1, countSuccess)
	require.Equal(t, countFail, 7, "%v Incorrect fail: waiting %v, getting %v", msg, 7, countFail)

	//	Один успешный в конце
	tasks = []func() error{fail, fail, fail, fail, fail, fail, fail, fail, fail, success}
	msg = "Successful last:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 0, "%v Incorrect successful: waiting %v, getting %v", msg, 0, countSuccess)
	require.Equal(t, countFail, 7, "%v Incorrect fail: waiting %v, getting %v", msg, 7, countFail)

	//	Один неуспешный в начале
	tasks = []func() error{fail, success, success, success, success, success, success, success, success, success}
	msg = "Fail first:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 9, "%v Incorrect successful: waiting %v, getting %v", msg, 9, countSuccess)
	require.Equal(t, countFail, 1, "%v Incorrect fail: waiting %v, getting %v", msg, 1, countFail)

	//	Один неуспешный в конце
	tasks = []func() error{success, success, success, success, success, success, success, success, success, fail}
	msg = "Fail last:"
	_ = Run(tasks, numThreads, maxError)
	require.Equal(t, countSuccess, 9, "%v Incorrect successful: waiting %v, getting %v", msg, 9, countSuccess)
	require.Equal(t, countFail, 1, "%v Incorrect fail: waiting %v, getting %v", msg, 1, countFail)

}

func processingOneSuccessfulHelper() error {
	forceStop = true
	return nil
}

func processingOneFailHelper() error {
	forceStop = true
	return errors.New("FAIL")
}

func TestProcessingOne(t *testing.T) {
	processChan = make(chan func() error, 1)
	resultChan = make(chan error, 1)
	blockChan = make(chan bool, 1)

	//	Проверяем таймаут
	forceStop    = true
	var runningThreads sync.WaitGroup
	runningThreads.Add(1)
	ProcessingOne(&runningThreads)
	runningThreads.Wait()
	chanLength := len(resultChan)
	require.Equal(t, chanLength, 0, "Incorrect timeout processing: waiting 0, getting %v", chanLength)

	//	Проверяем успешный вариант
	forceStop    = false
	processChan <- processingOneSuccessfulHelper
	blockChan <- true
	runningThreads.Add(1)
	ProcessingOne(&runningThreads)
	runningThreads.Wait()
	result := <-resultChan
	require.Equal(t, result, nil, "Incorrect successful processing: waiting nil, getting %v", result)

	//	Проверяем неуспешный вариант
	forceStop    = false
	processChan <- processingOneFailHelper
	blockChan <- true
	runningThreads.Add(1)
	ProcessingOne(&runningThreads)
	runningThreads.Wait()
	result = <-resultChan
	require.NotEqual(t, result, nil, "Incorrect fail processing: waiting error, getting nil")
}

func TestWork(t *testing.T) {
	percentError = 0
	result := Work()
	require.Equal(t, result, nil, "Incorrect work result: waiting nil, getting %v", result)

	percentError = 100
	result = Work()
	require.NotEqual(t, result, nil, "Incorrect work result: waiting error, getting nil")
}