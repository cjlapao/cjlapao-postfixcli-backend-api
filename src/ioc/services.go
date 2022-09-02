package ioc

import (
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/log"
)

var Log = log.Get()

var Config = execution_context.Get().Configuration
