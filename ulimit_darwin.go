// +build darwin

package main

import (
	"fmt"
	"syscall"
)

func init() {
	err := checkAndSetUlimit()
	if err != nil {
		log.Error(err)
	}
}

func checkAndSetUlimit() error {
	ipfsFileDescNum := uint64(1000000000)

	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("error getting rlimit: %s", err)
	}

	var setting bool
	if rLimit.Cur < ipfsFileDescNum {
		if rLimit.Max < ipfsFileDescNum {
			log.Error("adjusting max")
			rLimit.Max = ipfsFileDescNum
		}
		fmt.Printf("Adjusting current ulimit to %d...\n", ipfsFileDescNum)
		rLimit.Cur = ipfsFileDescNum
		setting = true
	}

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("error setting ulimit: %s", err)
	}

	if setting {
		fmt.Printf("Successfully raised file descriptor limit to %d.\n", ipfsFileDescNum)
	}

	return nil
}
