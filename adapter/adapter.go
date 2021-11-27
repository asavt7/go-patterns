package adapter

import "strconv"

type Pet struct {
	id   int
	name string
}

type Service interface {
	GetById(id string) interface{}
}

type AdapterGetById func(id int) (Pet, error)

func NewAdapterGetById(s Service) AdapterGetById {
	return func(id int) (Pet, error) {
		// convert input params
		idStr := strconv.Itoa(id)
		// call service
		res := s.GetById(idStr)
		//convert response
		return res.(Pet), nil
	}
}
