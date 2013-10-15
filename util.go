/*
 * util.go: Go bindings for lxc
 *
 * Copyright © 2013, S.Çağlar Onur
 *
 * Authors:
 * S.Çağlar Onur <caglar@10ur.org>
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.

 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.

 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
 */

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

func makeArgs(args []string) **C.char {
	cparams := C.makeCharArray(C.int(len(args)))
	for i, s := range args {
		C.setArrayString(cparams, C.CString(s), C.int(i))
	}
	return cparams
}

func freeArgs(cArgs **C.char, length int) {
	C.freeCharArray(cArgs, C.int(length))
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
