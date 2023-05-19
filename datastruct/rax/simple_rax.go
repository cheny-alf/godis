package rax

type SimpleRax struct {
	children map[string]*SimpleRax
	value    string
	isLeaf   bool
	queue    []string
}

func NewRax() *SimpleRax {
	return &SimpleRax{children: make(map[string]*SimpleRax)}
}
func (s *SimpleRax) Insert(key, value string) {
	if s.children == nil {
		s.children = make(map[string]*SimpleRax)
	}
	if s.queue == nil {
		s.queue = make([]string, 0)
	}
	if _, ok := s.children[key]; !ok {
		s.queue = append(s.queue, key) // 在队列中添加新节点的key
		s.children[key] = &SimpleRax{
			value:  value,
			isLeaf: true,
		}
	}
}
func (s *SimpleRax) Find(key string) (string, bool) {
	if len(key) == 0 {
		return s.value, s.isLeaf
	}
	c, rest := key[:1], key[1:]
	if child, ok := s.children[c]; ok {
		return child.Find(rest)
	}
	return "", false
}
func (s *SimpleRax) Delete(key string) bool {
	if len(key) == 0 {
		if s.isLeaf {
			s.isLeaf = false
			s.value = ""
			return true
		}
		return false
	}
	c, rest := key[:1], key[1:]
	if child, ok := s.children[c]; ok {
		if child.Delete(rest) {
			if len(child.children) == 0 && !child.isLeaf {
				delete(s.children, c)
			}
			return true
		}
	}
	return false
}
func (s *SimpleRax) Size() int {
	size := 0
	if s.isLeaf {
		size++
	}
	for _, child := range s.children {
		size += child.Size()
	}
	return size
}

func (s *SimpleRax) RemoveFirstNode() {
	if len(s.queue) == 0 { // 如果队列为空，则直接返回
		return
	}
	firstKey := s.queue[0]
	firstNode := s.children[firstKey]
	if len(firstNode.queue) > 0 { // 如果 firstNode 还有子节点，则递归删除其子节点中插入时间最早的节点
		firstNode.RemoveFirstNode()
	} else { // 否则，删除 firstNode
		delete(s.children, firstKey)
		s.queue = s.queue[1:]
	}
}
