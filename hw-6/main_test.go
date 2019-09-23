package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var data = Glist{}

func TestInit(t *testing.T) {
	if data.first != nil || data.last != nil || data.Len() != 0 {
		t.Fatalf("Fail initialization list")
	}
}

func TestGlist_PushFront(t *testing.T) {
	data.PushFront(1)
	if data.first.Value() != 1 || data.last.Value() != 1 || data.Len() != 1 {
		t.Fatalf("Fail insert first item at begin")
	}

	data.PushFront(2)
	if data.first.Value() != 2 || data.last.Value() != 1 || data.Len() != 2 {
		t.Fatalf("Fail insert new item at begin")
	}
}

func TestGlist_PushBack(t *testing.T) {
	data.PushBack(3)
	if data.first.Value() != 2 || data.last.Value() != 3 || data.Len() != 3 {
		t.Fatalf("Fail insert new item at end")
	}

	dataTest := Glist{}
	dataTest.PushBack("A")
	if dataTest.first.Value() != "A" || dataTest.last.Value() != "A" || dataTest.Len() != 1 {
		t.Fatalf("Fail insert new item at end")
	}
}

func TestGlist_First(t *testing.T) {
	result := data.First().Value()
	require.Equal(t, result, 2, "Incorrect get first value: waiting 2, getting %v", result)
}

func TestGlist_Last(t *testing.T) {
	result := data.Last().Value()
	require.Equal(t, result, 3, "Incorrect get last value: waiting 3, getting %v", result)
}

func TestGlist_Get(t *testing.T) {

	//	Тестируем выход за индекс
	if _, err := data.Get(10); err == nil {
		require.FailNow(t, "Fail parse error")
	}

	//	Заполняем данные
	data.PushBack(4)
	data.PushBack(5)
	data.PushBack(6)
	data.PushBack(7)
	data.PushBack(8)
	data.PushBack(9)

	//	Тест чтения нечётного количества элементов
	result, _ := data.Get(1)
	require.Equal(t, result.Value(), 1, "Incorrect get last value: waiting 1, getting %v", result.Value())
	result, _ = data.Get(0)
	require.Equal(t, result.Value(), 2, "Incorrect get last value: waiting 2, getting %v", result.Value())
	result, _ = data.Get(2)
	require.Equal(t, result.Value(), 3, "Incorrect get last value: waiting 3, getting %v", result.Value())
	result, _ = data.Get(3)
	require.Equal(t, result.Value(), 4, "Incorrect get last value: waiting 4, getting %v", result.Value())
	result, _ = data.Get(4)
	require.Equal(t, result.Value(), 5, "Incorrect get last value: waiting 5, getting %v", result.Value())
	result, _ = data.Get(5)
	require.Equal(t, result.Value(), 6, "Incorrect get last value: waiting 6, getting %v", result.Value())
	result, _ = data.Get(6)
	require.Equal(t, result.Value(), 7, "Incorrect get last value: waiting 7, getting %v", result.Value())
	result, _ = data.Get(7)
	require.Equal(t, result.Value(), 8, "Incorrect get last value: waiting 8, getting %v", result.Value())

	//	Тест чтения чётного количества элементов
	data.PushBack(10)
	result, _ = data.Get(1)
	require.Equal(t, result.Value(), 1, "Incorrect get last value: waiting 1, getting %v", result.Value())
	result, _ = data.Get(0)
	require.Equal(t, result.Value(), 2, "Incorrect get last value: waiting 2, getting %v", result.Value())
	result, _ = data.Get(2)
	require.Equal(t, result.Value(), 3, "Incorrect get last value: waiting 3, getting %v", result.Value())
	result, _ = data.Get(3)
	require.Equal(t, result.Value(), 4, "Incorrect get last value: waiting 4, getting %v", result.Value())
	result, _ = data.Get(4)
	require.Equal(t, result.Value(), 5, "Incorrect get last value: waiting 5, getting %v", result.Value())
	result, _ = data.Get(5)
	require.Equal(t, result.Value(), 6, "Incorrect get last value: waiting 6, getting %v", result.Value())
	result, _ = data.Get(6)
	require.Equal(t, result.Value(), 7, "Incorrect get last value: waiting 7, getting %v", result.Value())
	result, _ = data.Get(7)
	require.Equal(t, result.Value(), 8, "Incorrect get last value: waiting 8, getting %v", result.Value())
	result, _ = data.Get(8)
	require.Equal(t, result.Value(), 9, "Incorrect get last value: waiting 9, getting %v", result.Value())
	result, _ = data.Get(9)
	require.Equal(t, result.Value(), 10, "Incorrect get last value: waiting 10, getting %v", result.Value())
}

func TestGlist_Remove(t *testing.T) {

	//	Удаляем единственный элемент
	dataTest := Glist{}
	dataTest.PushFront("A")
	if result, err := dataTest.Remove(0); err != nil {
		require.FailNow(t, "Fail remove last element, reason: "+err.Error())
	} else {
		require.Equal(t, result.Value(), "A", "Remove incorrect element: waiting A, getting %v", result.Value())
	}

	//	Удаляем последний элемент
	if result, err := data.Remove(9); err != nil {
		require.FailNow(t, "Fail remove last element, reason: "+err.Error())
	} else {
		require.Equal(t, result.Value(), 10, "Remove incorrect element: waiting 10, getting %v", result.Value())
	}
	require.Equal(t, data.Len(), uint(9), "Remove incorrect resize after remove: waiting 9, getting %v", data.Len())

	//	Удаляем элемент в середине
	if result, err := data.Remove(5); err != nil {
		require.FailNow(t, "Fail remove last element, reason: "+err.Error())
	} else {
		require.Equal(t, result.Value(), 6, "Remove incorrect element: waiting 6, getting %v", result.Value())
	}
	require.Equal(t, data.Len(), uint(8), "Remove incorrect resize after remove: waiting 8, getting %v", data.Len())

	//	Удаляем элемент в начале
	if result, err := data.Remove(0); err != nil {
		require.FailNow(t, "Fail remove last element, reason: "+err.Error())
	} else {
		require.Equal(t, result.Value(), 2, "Remove incorrect element: waiting 2, getting %v", result.Value())
	}
	require.Equal(t, data.Len(), uint(7), "Remove incorrect resize after remove: waiting 7, getting %v", data.Len())

	//	Удаляем ещё элемент в середине (для проверки целостности)
	if result, err := data.Remove(5); err != nil {
		require.FailNow(t, "Fail remove last element, reason: "+err.Error())
	} else {
		require.Equal(t, result.Value(), 8, "Remove incorrect element: waiting 8, getting %v", result.Value())
	}
}

func TestGlist_Len(t *testing.T) {
	require.Equal(t, data.Len(), uint(6), "Remove incorrect resize after remove: waiting 8, getting %v", data.Len())
}

func TestList_Value(t *testing.T) {
	result,_ := data.Get(0)
	require.Equal(t, result.Value(), 1, "Remove incorrect resize after remove: waiting 1, getting %v", result.Value())
}

func TestList_Next(t *testing.T) {
	result,_ := data.Get(0)
	result, _ = result.Next()
	require.Equal(t, result.Value(), 3, "Remove incorrect resize after remove: waiting 3, getting %v", result.Value())
}

func TestList_Prev(t *testing.T) {
	result,_ := data.Get(3)
	result, _ = result.Prev()
	require.Equal(t, result.Value(), 4, "Remove incorrect resize after remove: waiting 4, getting %v", result.Value())
}

