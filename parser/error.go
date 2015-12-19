package parser

import (
	"fmt"
)

type NoSuchPluginError struct {
	Lang string
}

func (err *NoSuchPluginError) Error() string {
	return fmt.Sprintf("Plugin for language %s does not exist", err.Lang)
}
