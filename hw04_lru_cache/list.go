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

type list struct {
	size int
	head *ListItem
	tail *ListItem
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	if l.head == nil && l.tail != nil {
		return l.tail
	}

	return l.head
}

func (l *list) Back() *ListItem {
	if l.tail == nil && l.head != nil {
		return l.head
	}

	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	listItem := &ListItem{
		Value: v,
	}

	if l.size == 1 && l.tail == nil {
		l.tail = l.head
		l.head = nil
	}

	if l.head == nil {
		l.head = listItem
	} else {
		listItem.Next = l.head
		l.head.Prev = listItem
		l.head = listItem
	}

	l.size++

	if l.size == 2 {
		l.fixRefs()
	}

	return l.head
}

// PushBack добавляем элемент в конец.
func (l *list) PushBack(v interface{}) *ListItem {
	listItem := &ListItem{
		Value: v,
	}

	if l.size == 1 && l.head == nil {
		l.head = l.tail
		l.tail = nil
	}

	if l.tail == nil {
		l.tail = listItem
	} else {
		listItem.Prev = l.tail
		l.tail.Next = listItem
		l.tail = listItem
	}

	l.size++

	if l.size == 2 {
		l.fixRefs()
	}

	return l.tail
}

func (l *list) Remove(i *ListItem) {
	if l.isTailEqualHead() && l.head == i {
		l.head = nil
		l.tail = nil
		i.Next = nil
		i.Prev = nil
		l.size = 0

		return
	}

	switch i {
	case l.head:
		l.head = i.Next
		if l.head != nil {
			l.head.Prev = nil
		}
	case l.tail:
		l.tail = i.Prev

		if l.tail != nil {
			l.tail.Next = nil
		}
	default:
		prevItem := i.Prev
		nextItem := i.Next

		prevItem.Next = nextItem
		nextItem.Prev = prevItem
	}

	l.size--
}

func (l *list) isTailEqualHead() bool {
	return l.head == l.tail
}

// MoveToFront переместит элемент в начало.
func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}

	prevItem := i.Prev
	nextItem := i.Next

	if prevItem != nil {
		i.Prev.Next = nextItem
	}

	if nextItem != nil {
		nextItem.Prev = prevItem
	}

	i.Prev = nil
	i.Next = l.head
	l.head = i
}

// fixRefs проставляет ссылки головы и хвоста друг на друга.
func (l *list) fixRefs() {
	if l.size != 2 {
		return
	}

	l.head.Next = l.tail
	l.tail.Prev = l.head
}

func NewList() List {
	return new(list)
}
