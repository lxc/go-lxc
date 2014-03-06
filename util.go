// Copyright © 2013, 2014, S.Çağlar Onur
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.
//
// Authors:
// S.Çağlar Onur <caglar@10ur.org>

// +build linux

package lxc

/*
#include <stdlib.h>
#include <lxc/lxccontainer.h>

static char** makeCharArray(int size) {
    return calloc(sizeof(char*), size);
}

static void setArrayString(char **array, char *string, int n) {
    array[n] = string;
}

static void freeCharArray(char **array, int size) {
    int i;
    for (i = 0; i < size; i++)
        free(array[i]);
    free(array);
}

static void freeSnapshotArray(struct lxc_snapshot *s, int size) {
    int i;
    for (i = 0; i < size; i++) {
        s[i].free(&s[i]);
    }
    free(s);
}
*/
import "C"

import (
	"unsafe"
)

func sptr(p uintptr) *C.char {
	return *(**C.char)(unsafe.Pointer(p))
}

func makeNullTerminatedArgs(args []string) **C.char {
	cparams := C.makeCharArray(C.int(len(args) + 1))
	for i, s := range args {
		C.setArrayString(cparams, C.CString(s), C.int(i))
	}
	C.setArrayString(cparams, nil, C.int(len(args)))
	return cparams
}

func freeNullTerminatedArgs(cArgs **C.char, length int) {
	C.freeCharArray(cArgs, C.int(length+1))
}

func convertArgs(cArgs **C.char) []string {
	if cArgs == nil {
		return nil
	}

	var s []string

	// duplicate
	for p := uintptr(unsafe.Pointer(cArgs)); sptr(p) != nil; p += unsafe.Sizeof(uintptr(0)) {
		s = append(s, C.GoString(sptr(p)))
	}

	// free the original
	C.freeCharArray(cArgs, C.int(len(s)))

	return s
}

func convertNArgs(cArgs **C.char, size int) []string {
	if cArgs == nil || size <= 0 {
		return nil
	}

	var s []string

	// duplicate
	p := uintptr(unsafe.Pointer(cArgs))
	for i := 0; i < size; i++ {
		s = append(s, C.GoString(sptr(p)))
		p += unsafe.Sizeof(uintptr(0))
	}

	// free the original
	C.freeCharArray(cArgs, C.int(size))

	return s
}

func freeSnapshots(snapshots *C.struct_lxc_snapshot, size int) {
	C.freeSnapshotArray(snapshots, C.int(size))
}
