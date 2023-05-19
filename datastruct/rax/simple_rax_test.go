package rax

import "testing"

func TestRax(t *testing.T) {
	rax := NewRax()
	// 测试插入操作
	rax.Insert("hello", "world")
	rax.Insert("hi", "there")
	rax.Insert("go", "lang")
	// 测试查找操作
	val1, ok := rax.Find("hello")
	if !ok || val1 != "world" {
		t.Errorf("Expected value of 'hello' is 'world', but got %v", val1)
	}
	val2, ok := rax.Find("hi")
	if !ok || val2 != "there" {
		t.Errorf("Expected value of 'hi' is 'there', but got %v", val2)
	}
	val3, ok := rax.Find("go")
	if !ok || val3 != "lang" {
		t.Errorf("Expected value of 'go' is 'lang', but got %v", val3)
	}
	// 测试删除操作
	rax.Delete("hello")
	_, ok = rax.Find("hello")
	if ok {
		t.Error("Expected 'hello' to have been deleted, but still exists")
	}
	// 测试计算树大小
	size := rax.Size()
	if size != 2 {
		t.Errorf("Expected size of Rax to be 2, but got %v", size)
	}
}
