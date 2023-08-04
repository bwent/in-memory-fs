package util

import (
	"errors"
)

type queue []*File

func (q *queue) PushBack(element *File) {
	*q = append(*q, element)
}

func (q *queue) PopFront() (*File, error) {
	if len(*q) == 0 {
		return nil, errors.New("queue is empty")
	}

	head := (*q)[0]
	*q = (*q)[1:]
	return head, nil
}

func (q *queue) Empty() bool {
	return len(*q) == 0
}

func (q *queue) Size() int {
	return len(*q)
}

func (q *queue) Peek() (*File, error) {
	if len(*q) == 0 {
		return nil, errors.New("queue is empty")
	}
	return (*q)[0], nil
}
