package main

import "fmt"

type Store interface {
	Save(i interface{}) (interface{}, error)
	Delete(id int) error
	GetById(id int) (interface{}, error)
}

type MockStore struct {
}

func (m MockStore) Save(i interface{}) (interface{}, error) {
	return i, nil
}

func (m MockStore) Delete(_ int) error {
	return nil
}

func (m MockStore) GetById(_ int) (interface{}, error) {
	return "test", nil
}

type MyDecorator struct {
	Store
}

func (m MyDecorator) Save(i interface{}) (interface{}, error) {
	fmt.Println("Before call")
	// warning : call method m.Save(i) will call this method recursively
	save, err := m.Store.Save(i)
	fmt.Println("After call")
	return save, err
}

func main() {
	store := MockStore{}
	decoratedStore := MyDecorator{store}

	_, _ = decoratedStore.Save("qwe")
	_ = decoratedStore.Delete(1)
}
