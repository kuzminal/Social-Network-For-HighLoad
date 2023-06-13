package handler

import (
	"SocialNetHL/internal/store"
)

type Instance struct {
	store        store.Store
	readStorages *store.ReadNodes
}

func NewInstance(writeStorage store.Store, readStorages *store.ReadNodes) *Instance {
	return &Instance{
		store:        writeStorage,
		readStorages: readStorages,
	}
}
