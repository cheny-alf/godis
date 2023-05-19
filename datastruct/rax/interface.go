package rax

type Rax interface {
	Insert(key, value string)
	Find(key string) (string, bool)
	Delete(key string) bool
	Size() int
}
