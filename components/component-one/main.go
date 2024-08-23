package main

import (
	"github.com/golemcloud/golem-go/std"

	"golem-go-project/components/component-one/binding"
)

func init() {
	binding.SetExportsGolemComponentOneComponentOneApi(&Impl{})
}

type Impl struct {
	counter uint64
}

func (i *Impl) Add(value uint64) {
	std.Init(std.Packages{Os: true, NetHttp: true})

	componentTwo := binding.NewComponentTwoApi(binding.GolemRpc0_1_0_TypesUri{Value: "uri"})
	defer componentTwo.Drop()
	componentTwo.BlockingAdd(value)

	componentThree := binding.NewComponentThreeApi(binding.GolemRpc0_1_0_TypesUri{Value: "uri"})
	defer componentThree.Drop()
	componentThree.BlockingAdd(value)

	i.counter += value
}

func (i *Impl) Get() uint64 {
	std.Init(std.Packages{Os: true, NetHttp: true})

	return i.counter
}

func main() {}
