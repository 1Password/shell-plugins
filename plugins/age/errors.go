package age

import "fmt"

var (
	ErrConflictingIdentityFlag = fmt.Errorf("conflict: the -i/--identity flag is automatically added by this plugin. Remove it from your command to continue")
	ErrUnknownCommand          = fmt.Errorf("unknown command")
)
