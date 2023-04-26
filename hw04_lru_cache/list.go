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
	head *ListItem
	tail *ListItem
	size int
}

func (lst *list) Len() int {
	return lst.size
}

func (lst *list) Front() *ListItem {
	return lst.head
}

func (lst *list) Back() *ListItem {
	return lst.tail
}

func (lst *list) PushFront(v interface{}) *ListItem {
	if lst.head == nil {
		lst.head = new(ListItem)
		lst.tail = lst.head
	} else {
		lst.head.Prev = new(ListItem)
		lst.head.Prev.Next = lst.head
		lst.head = lst.head.Prev
	}
	lst.head.Value = v
	lst.size++
	return lst.head
}

func (lst *list) PushBack(v interface{}) *ListItem {
	if lst.tail == nil {
		lst.tail = new(ListItem)
		lst.head = lst.tail
	} else {
		lst.tail.Next = new(ListItem)
		lst.tail.Next.Prev = lst.tail
		lst.tail = lst.tail.Next
	}
	lst.tail.Value = v
	lst.size++
	return lst.tail
}

func (lst *list) Remove(i *ListItem) {
	switch {
	case lst.tail == i:
		lst.tail = i.Prev
		lst.tail.Next = nil
		i.Prev = nil
		i = nil
	case lst.head == i: // This situation is impossible in case of cache task
		lst.head = i.Next
		lst.head.Prev = nil
		i.Next = nil
		i = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		i.Next = nil
		i.Prev = nil
		i = nil
	}
	lst.size--
}

func (lst *list) MoveToFront(i *ListItem) {
	if i == lst.head {
		return
	}
	if i == lst.tail {
		lst.tail.Prev.Next = nil
		lst.tail.Next = lst.head
		lst.head.Prev = lst.tail
		lst.tail = lst.tail.Prev
		lst.head = lst.head.Prev
		lst.head.Prev = nil
		return
	}
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	i.Next = lst.head
	lst.head.Prev = i
	lst.head = i
}

func NewList() List {
	return new(list)
}
