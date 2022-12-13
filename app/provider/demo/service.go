package demo

import (
	"github.com/chsir-zy/anan/framework"
)

type Service struct {
	c framework.Container
}

func (s *Service) GetAllStudent() []Student {
	return []Student{
		{
			Id:   1,
			Name: "foo",
		},
		{
			Id:   2,
			Name: "bar",
		},
	}
}

func NewService(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)

	return &Service{c: c}, nil
}
