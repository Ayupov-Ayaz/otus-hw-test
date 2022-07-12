package hw04lrucache

import "sync"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Clear()
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

func NewListItem(value interface{}, next, prev *ListItem) *ListItem {
	return &ListItem{
		Value: value,
		Next:  next,
		Prev:  prev,
	}
}

type list struct {
	mu    *sync.RWMutex
	count int
	tail  *ListItem
	head  *ListItem
}

func newList() *list {
	return &list{
		mu: &sync.RWMutex{},
	}
}

func (l list) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.count
}

func (l list) Front() *ListItem {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.head
}

func (l list) Back() *ListItem {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.tail
}

func (l *list) pushFront(item *ListItem) {
	if l.head != nil {
		l.head.Prev = item
		item.Next = l.head
	}

	l.head = item

	if l.tail == nil {
		l.tail = item
	}

	item.Prev = nil

	l.increment()
}

func (l *list) increment() {
	l.count++
}

func (l *list) decrement() {
	l.count--
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.mu.Lock()
	item := NewListItem(v, l.head, nil)
	l.pushFront(item)
	l.mu.Unlock()

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.mu.Lock()
	item := NewListItem(v, nil, l.tail)
	if l.tail != nil {
		l.tail.Next = item
	}

	if l.head == nil {
		l.head = item
	}

	l.tail = item
	l.increment()
	l.mu.Unlock()

	return item
}

func (l *list) removeFromMiddle(item *ListItem) {
	if item.Prev != nil {
		item.Prev.Next = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	}

	l.decrement()
}

func (l *list) removeTail(item *ListItem) {
	if item.Prev != nil {
		item.Prev.Next = nil
		l.tail = item.Prev
	}

	l.decrement()
}

func (l *list) removeHead(item *ListItem) {
	if item.Next != nil {
		item.Next.Prev = nil
	}

	l.head = item.Next
	l.decrement()
}

func (l *list) remove(item *ListItem) {
	if item.Prev != nil && item.Next != nil {
		l.removeFromMiddle(item)
		return
	}

	if item.Prev != nil {
		l.removeTail(item)
	}

	if item.Next != nil {
		l.removeHead(item)
	}
}

func (l *list) Remove(item *ListItem) {
	l.mu.Lock()
	l.remove(item)
	l.mu.Unlock()
}

func (l *list) MoveToFront(item *ListItem) {
	l.mu.Lock()
	l.remove(item)
	l.pushFront(item)
	l.mu.Unlock()
}

func (l *list) Clear() {
	l.mu.Lock()
	l.tail = nil
	l.head = nil
	l.count = 0
	l.mu.Unlock()
}
