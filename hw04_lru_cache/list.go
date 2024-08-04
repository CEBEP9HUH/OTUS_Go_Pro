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
	len   int
}

func NewList() List {
	return new(list)
}

// mutable

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}
	if l.front == nil {
		l.front = item
		l.back = item
	} else {
		l.front.Prev = item
		item.Next = l.front
		l.front = l.front.Prev
	}
	l.len++
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}
	if l.front == nil {
		l.front = item
		l.back = item
	} else {
		l.back.Next = item
		item.Prev = l.back
		l.back = l.back.Next
	}
	l.len++
	return l.back
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	switch {
	case i == l.front && l.front.Next == nil:
		l.front = nil
		l.back = nil
	case i == l.front:
		l.front = l.front.Next
		l.front.Prev = nil
	case i == l.back:
		l.back = l.back.Prev
		l.back.Next = nil
	default:
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.front {
		return
	}
	if i == l.back {
		l.back = l.back.Prev
		l.back.Next = nil
	} else {
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = l.front.Prev
}

// immutable

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}
