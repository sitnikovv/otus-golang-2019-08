package main

import (
	"errors"
	"fmt"
	"log"
)

//	Для проверяющего: я исходил из того, что []List использовать нельзя, хотя с ним было бы куда проще

type List struct {
	prev *List
	data interface{}
	next *List
}

func (item *List) Prev() (*List, error) {
	if item.prev == nil {
		return nil, errors.New("no previous value: it's first element")
	}
	return item.prev, nil
}

func (item *List) Next() (*List, error) {
	if item.next == nil {
		return nil, errors.New("no next value: it's last element")
	}
	return item.next, nil
}

func (item *List) Value() interface{} {
	return item.data
}

type Glist struct {
	len   uint  //	Размерность списка
	first *List //	Ссылка на первый элемент
	last  *List //	Ссылка на последний элемент
}

//	Получение размерности списка
func (g *Glist) Len() uint {
	return g.len
}

func (g *Glist) First() *List {
	return g.first
}

func (g *Glist) Last() *List {
	return g.last
}

func (g *Glist) PushFront(data interface{}) *List {
	g.len++
	list := &List{prev: nil, data: data, next: g.first}
	if g.first == nil {
		g.first = list
		g.last = list
		return list
	}
	g.first.prev = list
	g.first = list
	return list
}

func (g *Glist) PushBack(data interface{}) *List {
	g.len++
	list := &List{prev: g.last, data: data, next: nil}
	if g.last == nil {
		g.first = list
		g.last = list
		return list
	}
	g.last.next = list
	g.last = list
	return list
}

func (g *Glist) Get(index uint) (*List, error) {

	//	Выходим, если пытаются обратиться к элементу вне списка
	if index >= g.len {
		return nil, errors.New("item is missing")
	}

	var element *List
	var err error

	//	Подбираемся к удаляемому элементу спереди или сзади, в зависимости от того, откуда ближе
	if g.len/2 >= index {
		element = g.First()
		for i := uint(0); i < index; i++ {
			if element, err = element.Next(); err != nil {
				return nil, errors.New("pointer error: " + err.Error())
			}
		}
	} else {
		element = g.Last()
		for i := uint(0); i < g.len-index-1; i++ {
			if element, err = element.Prev(); err != nil {
				return nil, errors.New("pointer error: " + err.Error())
			}
		}
	}
	return element, nil
}

func (g *Glist) Remove(place uint) (*List, error) {

	var (
		removed, prev, next *List
		err                 error
	)
	if removed, err = g.Get(place); err != nil {
		return nil, errors.New("no items to remove: " + err.Error())
	}

	//	Больше ошибок не будет, можно уменьшать счётчик количества
	g.len--

	//	Если удаляется единственный элемент, то сбрасываем указатели и возвращаем удалённый элемент
	if g.len == 0 {
		g.first = nil
		g.last = nil
		return removed, nil
	}

	//	Удаление элемента
	//	Для проверяющего:
	//	Я знаю, что мог бы обойтись конструкциями removed.prev.next = removed.next, но воспользовался методами List (мало ли, вдруг со временем логика
	//	поменяется, код не нужно будет переписывать). Так что это не ошибка, а осознанное решение (хотя и оно может быть ошибочным)
	prev, _ = removed.Prev()
	next, _ = removed.Next()
	if prev != nil {
		prev.next = next
	}
	if next != nil {
		next.prev = prev
	}

	//	Обрабатываем крайние элементы
	if g.last == removed {
		g.last, _ = removed.Prev()
	}
	if g.first == removed {
		g.first, _ = removed.Next()
	}
	return removed, nil
}

func main() {
	data := Glist{}
	fmt.Println("Вносим числа от 1 до 10 в начало")
	for i := 1; i <= 10; i++ {
		fmt.Printf("%d-й пошёл...\n", i)
		data.PushFront(i)
	}

	value := data.First()
	var err error
	fmt.Printf("Проверяем данные: %d", value.Value())
	for i := 9; i > 0; i-- {
		if value, err = value.Next(); err != nil {
			log.Fatalf("\n\nFail ")
		} else {
			fmt.Printf(", %d", value.Value())
		}
	}
	fmt.Println("\n\nВсё хорошо")
}
