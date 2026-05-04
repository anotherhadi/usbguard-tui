package guard

import "errors"

var (
	ErrNotFound   = errors.New("usbguard not found in PATH")
	ErrPermission = errors.New("insufficient permissions to manage devices")
	ErrReadOnly   = errors.New("rules file is not writable")
)
