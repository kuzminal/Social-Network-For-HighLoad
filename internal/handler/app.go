package handler

import "SocialNetHL/internal/store"

type Instance struct {
	store store.Store
}

func NewInstance(storage store.Store) *Instance {
	return &Instance{
		store: storage,
	}
}
