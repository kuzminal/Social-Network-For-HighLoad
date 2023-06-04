package store

import (
	"context"
	"log"
	"sync"
	"time"
)

type ReadNodes struct {
	mu      sync.Mutex
	Current int
	Nodes   []Backend
}

type Backend struct {
	Store  Store
	IsDead bool
}

func GetReadNode(readStorage *ReadNodes) Store {
	readStorage.mu.Lock()
	storage := readStorage.Nodes[readStorage.Current]
	if storage.IsDead {
		if readStorage.Current == len(readStorage.Nodes)-1 {
			readStorage.Current = 0
		} else {
			readStorage.Current++
		}
	}
	storage = readStorage.Nodes[readStorage.Current]
	readStorage.mu.Unlock()
	return storage.Store
}

func SetDead(b bool, readStorage *ReadNodes, backend Backend) {
	readStorage.mu.Lock()
	for _, bckg := range readStorage.Nodes {
		if bckg == backend {
			bckg.IsDead = b
		}
	}
	readStorage.mu.Unlock()
}

func isAlive(store Store) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := store.Ping(ctx)
	if err != nil {
		log.Printf("Unreachable host, error", err.Error())
		return false
	}
	return true
}

func HealthCheck(readStorage *ReadNodes) {
	t := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-t.C:
			for _, backend := range readStorage.Nodes {
				isAlive := isAlive(backend.Store)
				SetDead(!isAlive, readStorage, backend)
				msg := "ok"
				if !isAlive {
					msg = "dead"
				}
				log.Printf("%v checked by healthcheck", msg)
			}
		}
	}
}
