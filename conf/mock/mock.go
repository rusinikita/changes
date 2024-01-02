package mock

import (
	"github.com/rusinikita/changes/conf"
)

//go:generate moq -rm -stub -out conf_mock.go . Conf
type Conf = conf.Conf
