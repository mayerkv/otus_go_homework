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
	front *ListItem
	back  *ListItem
	cnt   int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.cnt
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{v, nil, nil}

	if l.front == nil {
		l.front = item
		l.back = item
	} else {
		l.front.Prev = item
		item.Next = l.front
		l.front = item
	}

	l.cnt++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{v, nil, nil}

	if l.front == nil {
		l.front = item
		l.back = item
	} else {
		item.Prev = l.back
		l.back.Next = item
		l.back = item
	}

	l.cnt++

	return item
}

func (l *list) Remove(i *ListItem) {
	if l.front == nil || i == nil {
		return
	}

	if l.front == i {
		l.front = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	l.cnt--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	if l.front == i {
		return
	}

	if l.back == i {
		i.Prev.Next = nil
		l.back = i.Prev
	} else if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	l.front.Prev = i
	i.Next = l.front
	i.Prev = nil
	l.front = i
}
