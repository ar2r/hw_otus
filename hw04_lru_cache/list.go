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
	First *ListItem
	Last  *ListItem
	// init empty map
	Items map[interface{}]*ListItem
}

func (l *list) Len() int {
	return len(l.Items)
}

func (l *list) Front() *ListItem {
	return l.First
}

func (l *list) Back() *ListItem {
	return l.Last
}

func (l *list) PushFront(v interface{}) *ListItem {
	// todo: Проверить на дубли
	OldFirst := l.First
	NewFirst := &ListItem{Value: v, Next: OldFirst}

	if OldFirst != nil {
		NewFirst.Next = OldFirst
		OldFirst.Prev = NewFirst
	}

	if len(l.Items) == 0 {
		l.Last = NewFirst
		l.First = NewFirst
	} else {
		l.First = NewFirst
	}

	l.Items[v] = NewFirst

	return NewFirst
}

func (l *list) PushBack(v interface{}) *ListItem {
	// todo: Проверить на дубли
	OldLast := l.Last
	NewLast := &ListItem{Value: v, Prev: OldLast}

	if OldLast != nil {
		OldLast.Next = NewLast
		NewLast.Prev = OldLast
	}

	if len(l.Items) == 0 {
		l.First = NewLast
		l.Last = NewLast
	} else {
		l.Last = NewLast
	}

	l.Items[v] = NewLast

	return NewLast
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil && i.Next != nil {
		// Удаляем ссылку на элемент
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	if i.Prev == nil && i.Next != nil {
		// Удаляем ссылку на левый элемент
		i.Next.Prev = nil
	}

	if i.Next == nil && i.Prev != nil {
		// Удаляем ссылку на правый элемент
		i.Prev.Next = nil
	}

	delete(l.Items, i.Value)
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		// Элемент уже в начале списка
		return
	}

	// Удаляем ссылку на элемент
	i.Prev.Next = i.Next
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	// Добавляем элемент в начало списка
	OldFirst := l.First
	l.First = i
	l.First.Next = OldFirst
	OldFirst.Prev = l.First
	l.First.Prev = nil
}

func NewList() List {
	return &list{
		Items: make(map[interface{}]*ListItem),
	}
}
