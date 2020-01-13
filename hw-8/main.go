//	Для проверяющего. Примитивы синхронизации мы изучали на следующем уроке

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//	Может возникнуть ситуация, когда обработчиков 10, а остановиться нужно после первой ошибки

var (
	//	Поведение проверки:
	//		one  - за проверкой ошибок следит одна горутина, нет необходимости в синхронизации, но это может стать узким местом
	//		self - за проверкой ошибок горутины следят самостоятельно
	behaviorRule = "one"
	numTasks     = 10
	numThreads   = 5
	maxError     = 7
	percentError = 60
	forceStop    bool

	countSuccess   int
	countFail      int

	//	Канал с тасками
	processChan chan func() error

	//	Канал с результатом(ами) выполнения
	resultChan chan error

	//	Канал для блокировки обработчиков для правила поведения "one", напомню, на этом этапе обучения мы не знаем про примитивы синхронизации
	blockChan chan bool
)

func Run(tasks []func() error, N int, M int) error {
	if M < N {
		return errors.New("count error don't be less when count threads")
	}
	if tasks == nil {
		return errors.New("task cannot be nil")
	}
	if len(tasks) == 0 {
		return errors.New("task cannot be empty")
	}
	processChan = make(chan func() error, len(tasks))
	resultChan = make(chan error)
	blockChan = make(chan bool)
	countSuccess = 0
	countFail    = 0
	forceStop    = false

	defer close(blockChan)

	//	Запускаем необходимое количество горутин, пусть висят себе, обрабатывают таски потихоньку
	var runningThreads sync.WaitGroup
	for i := 0; i < N; i++ {
		if behaviorRule == "one" {
			runningThreads.Add(1)
			go ProcessingOne(&runningThreads)
		}
	}

	//	Напихиваем полный канал тасков
	for _, task := range tasks {
		processChan <- task
	}

	//	Будем считывать результаты обработки бесконечно
	for result := range resultChan {

		//	В зависимости от полученного результата, увеличиваем соответствующий счётчик
		if result != nil {
			countFail++
			if countFail >= M {
				forceStop = true
			}
		} else {
			countSuccess++
		}
		blockChan <- true
		if countFail + countSuccess == len(tasks) || forceStop{
			close(processChan)
			close(resultChan)
		}
	}
	runningThreads.Wait()
	return nil
}

//	Обработчики для правила "one", логика следующая:
//	Читаем функцию из канала и если нет команды на выход, то выполняем функцию, записывая её результат в канал результата (простите за тафталогию)
//	Если же выставлен флаг срочного завершения, сворачиваем всё
//	Также сворачиваемся, если очередь забита выполняющимися задачами, но флаг срочного завершения всё-таки выставлен (проверяем по таймауту)
func ProcessingOne(runningThreads *sync.WaitGroup) {
	for ; ; {
		select {
		case data, ok := <-processChan:

			//	Проверяем,
			if !ok || forceStop {
				runningThreads.Done()
				return
			}
			resultChan <- data()
			_ = <-blockChan
		case <-time.After(time.Millisecond * 100):
			if forceStop {
				runningThreads.Done()
				return
			}
		}
	}
}

func Work() error {
	var err error

	//	Делаем рандомную задержку, вплоть до секунды
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

	//	И, с указанной вероятностью, генерируем ошибку
	if !(rand.Intn(99) >= percentError) {
		err = errors.New("FAIL")
	}
	return err
}

func main() {
	//	Инициализируем рандом
	rand.Seed(time.Now().UnixNano())

	//	Заполняем тасками слайс
	tasks := make([]func() error, numTasks)
	for i := 0; i < numTasks; i++ {
		tasks[i] = Work
	}

	//	Сообщаем, если не удалось обработать все потоки
	if err := Run(tasks, numThreads, maxError); err != nil {
		fmt.Println(err)
	}

	//	Результат
	fmt.Printf("Success: %d, Error: %d\n", countSuccess, countFail)
}
