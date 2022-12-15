package cobra

import "github.com/chsir-zy/anan/framework"

func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}
