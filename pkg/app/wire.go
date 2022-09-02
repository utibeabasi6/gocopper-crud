package app

import (
	"github.com/google/wire"
	"github.com/utibeabasi6/gocopper-crud/pkg/todos"
)

var WireModule = wire.NewSet(
	todos.WireModule,
)
