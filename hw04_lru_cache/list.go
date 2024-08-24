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
	head   *ListItem
	tail   *ListItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}

	if l.head == nil { // Если список пустой
		l.head = newItem
		l.tail = newItem
	} else {
		newItem.Next = l.head
		l.head.Prev = newItem
		l.head = newItem
	}

	l.length++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}

	if l.tail == nil { // Если список пустой
		l.head = newItem
		l.tail = newItem
	} else {
		newItem.Prev = l.tail
		l.tail.Next = newItem
		l.tail = newItem
	}

	l.length++
	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.head { // Если элемент уже находится в начале
		return
	}

	// Удаляем элемент из текущего положения
	l.Remove(i)

	// Вставляем его в начало
	i.Next = l.head
	i.Prev = nil

	if l.head != nil {
		l.head.Prev = i
	}

	l.head = i

	if l.tail == nil { // Если список был пуст
		l.tail = i
	}

	l.length++
}

func NewList() List {
	return new(list)
}
