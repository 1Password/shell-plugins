//go:build !windows

package aws

import (
	"os"
	"syscall"
)

// ownedByCurrentUser reports whether the file is owned by the user running the plugin. It returns
// true when the owner cannot be determined.
func ownedByCurrentUser(fi os.FileInfo) bool {
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return true
	}
	return st.Uid == uint32(os.Geteuid())
}
