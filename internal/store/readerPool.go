package store

import (
	"context"
	"log"
	"sync"
	"time"
)

type ReadNodes[T any] struct {
	mu      sync.Mutex
	Current int
	Nodes   []Backend[T]
}

type Backend[T any] struct {
	Id     string
	Store  T
	IsDead bool
}

func NewReadNode[T any]() ReadNodes[T] {
	return ReadNodes[T]{
		Current: 0,
		Nodes:   []Backend[T]{},
	}
}

func (r ReadNodes[T]) GetReadNode() T {
	r.mu.Lock()
	storage := r.Nodes[r.Current]
	if storage.IsDead {
		if r.Current == len(r.Nodes)-1 {
			r.Current = 0
		} else {
			r.Current++
		}
	}
	storage = r.Nodes[r.Current]
	r.mu.Unlock()
	return storage.Store
}

func (r ReadNodes[T]) SetDead(b bool, backend Backend[T]) {
	r.mu.Lock()
	for _, bckg := range r.Nodes {
		if bckg.Id == backend.Id {
			bckg.IsDead = b
		}
	}
	r.mu.Unlock()
}

func isAlive(T any) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := T.(Store).Ping(ctx)
	if err != nil {
		log.Printf("Unreachable host, error", err.Error())
		return false
	}
	return true
}

func (r ReadNodes[T]) HealthCheck() {
	t := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-t.C:
			for _, backend := range r.Nodes {
				isAlive := isAlive(backend.Store)
				r.SetDead(!isAlive, backend)
				msg := "ok"
				if !isAlive {
					msg = "dead"
				}
				log.Printf("%v checked by healthcheck", msg)
			}
		}
	}
}
