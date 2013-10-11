/*
 * list_snapshots.go
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

package main

import (
	"fmt"
	"github.com/caglar10ur/lxc"
)

func main() {
	for _, v := range lxc.Containers() {
		fmt.Printf("%s\n", v.Name())
		l, err := v.ListSnapshots()
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}

		for _, s := range l {
			fmt.Printf("Name: %s\n", s.Name)
			fmt.Printf("Comment path: %s\n", s.CommentPath)
			fmt.Printf("Timestamp: %s\n", s.Timestamp)
			fmt.Printf("LXC path: %s\n", s.Path)
            fmt.Println()
		}
	}
}
