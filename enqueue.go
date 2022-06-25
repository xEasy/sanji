package sanji

import (
	"strconv"
	"strings"
)

func (s *Luffy) Enqueue(queue string, data []byte, delay uint32) (jobID string, err error) {
	return s.EnqueueWithPriority(queue, data, PriorityLevelHight, delay)
}

func (s *Luffy) EnqueueWithPriority(queue string, data []byte, priority PriorityLevel, delay uint32) (jobID string, err error) {
	q := s.priorityQueue(queue, priority)
	return s.client.Publish(q, data, 0, 3, delay)
}

func (s *Luffy) BatchEnqueue(queue string, data []any, delay uint32) (jobID []string, err error) {
	return s.BatchEnqueueWithPriority(queue, data, PriorityLevelHight, delay)
}

func (s *Luffy) BatchEnqueueWithPriority(queue string, data []any, priority PriorityLevel, delay uint32) (jobID []string, err error) {
	q := s.priorityQueue(queue, priority)
	return s.client.BatchPublish(q, data, 0, 3, delay)
}

func (s *Luffy) priorityQueue(queue string, priority PriorityLevel) string {
	return strings.Join([]string{queue, strconv.Itoa(int(priority))}, "_")
}
