package rax

type SimpleRax struct {
	children map[string]*SimpleRax
	value    string
	isLeaf   bool
}

func NewRax() *SimpleRax {
	return &SimpleRax{children: make(map[string]*SimpleRax)}
}
func (r *SimpleRax) Insert(key, value string) {
	if len(key) == 0 {
		r.isLeaf = true
		r.value = value
	} else {
		c, rest := key[:1], key[1:]
		if child, ok := r.children[c]; ok {
			child.Insert(rest, value)
		} else {
			r.children[c] = NewRax()
			r.children[c].Insert(rest, value)
		}
	}
}
func (r *SimpleRax) Find(key string) (string, bool) {
	if len(key) == 0 {
		return r.value, r.isLeaf
	}
	c, rest := key[:1], key[1:]
	if child, ok := r.children[c]; ok {
		return child.Find(rest)
	}
	return "", false
}
func (r *SimpleRax) Delete(key string) bool {
	if len(key) == 0 {
		if r.isLeaf {
			r.isLeaf = false
			r.value = ""
			return true
		}
		return false
	}
	c, rest := key[:1], key[1:]
	if child, ok := r.children[c]; ok {
		if child.Delete(rest) {
			if len(child.children) == 0 && !child.isLeaf {
				delete(r.children, c)
			}
			return true
		}
	}
	return false
}
func (r *SimpleRax) Size() int {
	size := 0
	if r.isLeaf {
		size++
	}
	for _, child := range r.children {
		size += child.Size()
	}
	return size
}
