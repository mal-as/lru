package mem

import "fmt"

// структура для хранения пар ключ-значение
type pair struct {
	key  string
	data []byte
}

// структура элемента в связном списке
type node struct {
	value pair
	next  *node
}

// структура - связный список
type list struct {
	head   *node
	length int
}

// конструктор для list
func newList() *list {
	return &list{}
}

// вставляет элементы в начало связного списка
func (l *list) insesrtFront(n *node) {
	if l.head == nil {
		l.head = n
		l.length++
	} else {
		n.next = l.head
		l.head = n
		l.length++
	}
}

// перемещает элемент в начало связного списка
func (l *list) moveToFront(n *node) {
	if l.head == nil {
		l.head = n
		l.length++
	} else {
		if n == l.head {
			return
		}
		elem := l.head
		for i := 0; i < l.length-1; i++ {
			if elem.next == n {
				elem.next = n.next
				n.next = l.head
				l.head = n
				break
			}
			elem = elem.next
		}
	}
}

// удаляет последний элемент
func (l *list) pop(idx int) {
	if l.length >= idx {
		elem := l.head

		for i := 0; i < idx-1; i++ {
			elem = elem.next
		}

		if elem != nil {
			elem.next = nil
		}
		l.length--
	}
}

// возвращает последний элемент
func (l *list) retriveLastElement() *node {
	elem := l.head

	for i := 0; i < l.length-1; i++ {
		elem = elem.next
	}

	return elem
}

// печатает связный список
func (l *list) print() {
	for elem := l.head; elem != nil; elem = elem.next {
		fmt.Printf("{%s: %s} -> ", elem.value.key, string(elem.value.data))
	}
}
