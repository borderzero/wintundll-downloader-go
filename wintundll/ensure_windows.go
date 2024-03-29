//go:build windows
// +build windows

package wintundll

import (
	"errors"
	"fmt"
	"net/http"
	"syscall"

	"golang.org/x/sys/windows"
)

// Ensure ensures the presence of the wintun dll in the system.
// It returns an incompatibility error on non-windows.
func Ensure(opts ...EnsureOption) error {
	config := getConfiguration(opts...)

	admin, err := isRunningAsAdministator()
	if err != nil {
		fmt.Errorf("failed to determine if executable is running as administrator: %v", err)
	}
	if !admin {
		return errors.New("executable is not running as administrator")
	}

	if _, err = syscall.LoadDLL(config.dllPathInToEnsure); err == nil {
		return nil
	}

	err = downloadAndMoveFromZip(
		http.Client{Timeout: config.downloadTimeout},
		config.downloadURL,
		config.dllPathInUnzippedDir,
		config.dllPathInToEnsure,
	)
	if err != nil {
		return fmt.Errorf("failed to get wintun.dll from remote: %v", err)
	}

	if _, err = syscall.LoadDLL(config.dllPathInToEnsure); err != nil {
		return fmt.Errorf("still failed to load wintun.dll after fresh download: %v", err)
	}

	return nil
}

// isRunningAsAdministator returns true if the executable is running
// as a system administrator. This function is copied directly from
// https://github.com/golang/go/issues/28804 with additional comments.
func isRunningAsAdministator() (bool, error) {
	var sid *windows.SID

	// Although this looks scary, it is directly copied from the
	// official windows documentation. The Go API for this is a
	// direct wrap around the official C++ API.
	// See https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,      // identAuth *windows.SidIdentifierAuthority
		2,                                   // subAuth byte
		windows.SECURITY_BUILTIN_DOMAIN_RID, // subAuth0 uint32
		windows.DOMAIN_ALIAS_RID_ADMINS,     // subAuth1 uint32
		0,                                   // subAuth2 uint32
		0,                                   // subAuth3 uint32
		0,                                   // subAuth4 uint32
		0,                                   // subAuth5 uint32
		0,                                   // subAuth6 uint32
		0,                                   // subAuth7 uint32
		&sid,                                // sid **windows.SID
	)
	if err != nil {
		return false, fmt.Errorf("failed to initialize admin SID: %v", err)
	}
	defer windows.FreeSid(sid)

	// This appears to cast a null pointer so I'm not sure why this
	// works, but this guy says it does and it Works for Me™:
	// https://github.com/golang/go/issues/28804#issuecomment-438838144
	token := windows.Token(0)
	defer token.Close()

	// check if the token is a member of the admin SID
	admin, err := token.IsMember(sid)
	if err != nil {
		return false, fmt.Errorf("token membership error: %v", err)
	}
	return admin, nil
}
