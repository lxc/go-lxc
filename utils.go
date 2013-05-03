/*
 * utils.go: Go bindings for lxc
 *
 * Copyright © 2013, S.Çağlar Onur
 *
 * Authors:
 * S.Çağlar Onur <caglar@10ur.org>
 *
 * This library is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2, as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along
 * with this program; if not, write to the Free Software Foundation, Inc.,
 * 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package lxc

// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

func makeArgs(args []string) []*C.char {
	ret := make([]*C.char, len(args)+1)
	for i, s := range args {
		ret[i] = C.CString(s)
	}
	ret[len(ret)-1] = nil
	return ret
}

func freeArgs(cArgs []*C.char) {
	for _, s := range cArgs {
		C.free(unsafe.Pointer(s))
	}
}
