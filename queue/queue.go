package queue

import "errors"

type QueueItem struct {
	next  *QueueItem
	Value interface{}
}

type Queue struct {
	head *QueueItem
	last *QueueItem
}

func NewQueue() *Queue {
	q := Queue{}
	q.head = nil
	q.last = nil
	return &q
}

var (
	ErrQueueEmpty       = errors.New("queue empty")
	ErrQueueIsNil       = errors.New("pointer to queue must be not nil")
	ErrQueueItemIsNil   = errors.New("pointer to item must be not nil")
	ErrIteratorFinished = errors.New("iterator finished")
)

func (q *Queue) Push(qi *QueueItem) error {
	if q == nil {
		return ErrQueueIsNil
	}
	if qi == nil {
		return ErrQueueItemIsNil
	}
	if q.head == nil {
		q.head = qi
	} else {
		q.last.next = qi
	}
	qi.next = nil
	q.last = qi
	return nil
}

func (q *Queue) Pop() (*QueueItem, error) {
	if q.head == nil {
		return nil, ErrQueueEmpty
	}
	qi := q.head

	err := q.DeleteFirst()
	if err != nil {
		return nil, err
	}
	return qi, nil
}

func (q *Queue) Get() (*QueueItem, error) {
	if q.head == nil {
		return nil, ErrQueueEmpty
	}
	return q.head, nil
}

func (q *Queue) Next(qi *QueueItem) (*QueueItem, error) {
	if qi == nil {
		return nil, ErrIteratorFinished
	}
	return qi.next, nil
}

func (q *Queue) DeleteFirst() error {
	if q == nil {
		return ErrQueueIsNil
	}
	if q.head == nil {
		return ErrQueueEmpty
	}
	if q.head.next == nil {
		q.last = nil
	}
	q.head = q.head.next
	return nil
}

func (q *Queue) Delete(qi *QueueItem) error {
	if q == nil {
		return ErrQueueIsNil
	}
	if q.head == nil {
		return ErrQueueEmpty
	}
	if qi == nil {
		return ErrQueueItemIsNil
	}
	qii := q.head
	var qii_prev *QueueItem
	for {
		if qii == qi {
			if qii_prev == nil {
				q.head = qii.next
			} else {
				qii_prev.next = qii.next
			}
		}

		if qii.next == nil {
			break
		} else {
			qii = qii.next
		}
	}
	if q.head != nil {
		if q.head.next == nil {
			q.last = nil
		}
	}
	return nil
}
