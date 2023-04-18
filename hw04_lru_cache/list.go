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
	Head *ListItem
	Tail *ListItem
	Size int
}

func (lst *list) Len() int {
	return lst.Size
}

func (lst *list) Front() *ListItem {
	return lst.Head
}

func (lst *list) Back() *ListItem {
	return lst.Tail
}

func (lst *list) PushFront(v interface{}) *ListItem {
	var newHead = new(ListItem)
	newHead.Next = lst.Head
	newHead.Prev = nil
	lst.Head.Prev = newHead
	lst.Head = newHead
	lst.Size++
	return lst.Head
}

func (lst *list) PushBack(v interface{}) *ListItem {
	var newTail = new(ListItem)
	newTail.Next = nil
	newTail.Prev = lst.Tail
	lst.Tail.Next = newTail
	lst.Tail = newTail
	lst.Size++
	return lst.Tail
}

func (lst *list) Remove(i *ListItem) {

	lst.Size--
}

func (lst *list) MoveToFront(i *ListItem) {

}

func NewList() List {
	return new(list)
}
