package main

import (
	"reflect"
	"testing"
)

func TestGetWords(t *testing.T) {

	// Test 1
	text1 := `
		aaa aaa aaa aaa aaa aaa aaa aaa aaa aaa aaa aaa
		bbb bbb bbb bbb bbb bbb bbb bbb bbb bbb bbb
		ccc ccc ccc ccc ccc ccc ccc ccc ccc ccc
		ddd ddd ddd ddd ddd ddd ddd ddd ddd
		eee eee eee eee eee eee eee eee
		fff fff fff fff fff fff fff
		ggg ggg ggg ggg ggg ggg
		hhh hhh hhh hhh hhh
		iii iii iii iii
		jjj jjj jjj
		kkk kkk
		lll
		aaa bbb1 ccc2 d d d ee e f ff
	`
	compare1 := []Word{{"aaa", 13}, {"bbb", 12}, {"ccc", 11}, {"ddd", 9},{"eee", 8}, {"fff", 7}, {"ggg", 6}, {"hhh", 5}, {"iii", 4}, {"jjj", 3}}
	result1 := getWords(text1, 0)
	if !reflect.DeepEqual(result1, compare1) {
		t.Fatalf("test failure:\nwaiting: %v\n incoming: %v", compare1, result1)
	}

	// Test 2
	text2 := `a bb a ccc a bb a a`
	compare2 := []Word{{"bb", 2}, {"ccc", 1}}

	result2 := getWords(text2, 2)
	if !reflect.DeepEqual(result2, compare2) {
		t.Fatalf("test failure:\nwaiting: %v\n incoming: %v", compare2, result2)
	}
}
