package main

import "context"

type store struct {
	//TODO: add mongodb
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(context.Context) error {
	return nil
}
