package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
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
	count int
	tail  *ListItem
	head  *ListItem
}

func NewList() *list {
	return &list{}
}

func (l list) Len() int {
	return l.count
}

func (l list) Front() *ListItem {
	return l.head
}

func (l list) Back() *ListItem {
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

	l.count++
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := NewListItem(v, l.head, nil)

	l.pushFront(item)

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := NewListItem(v, nil, l.tail)

	if l.tail != nil {
		l.tail.Next = item
	}

	if l.head == nil {
		l.head = item
	}

	l.tail = item
	l.count++

	return item
}

func (l *list) removeFromMiddle(item *ListItem) {
	if item.Prev != nil {
		item.Prev.Next = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	}

	l.count--
}

func (l *list) removeTail(item *ListItem) {
	if item.Prev != nil {
		item.Prev.Next = nil
	}

	l.count--
}

func (l *list) removeHead(item *ListItem) {
	if item.Next != nil {
		item.Next.Prev = nil
	}

	l.head = item.Next
	l.count--
}

func (l *list) Remove(item *ListItem) {
	if item.Prev != nil && item.Next != nil {
		l.removeFromMiddle(item)
	} else if item.Prev != nil {
		l.removeTail(item)
	} else if item.Next != nil {
		l.removeHead(item)
	}
}

func (l *list) MoveToFront(item *ListItem) {
	l.Remove(item)
	l.pushFront(item)
}
