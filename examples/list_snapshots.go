/*
 * list_snapshots.go
 *
 * Copyright © 2013, 2014, S.Çağlar Onur
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
	"log"

	"github.com/lxc/go-lxc"
)

func main() {
	c := lxc.Containers()
	for i := range c {
		log.Printf("%s\n", c[i].Name())
		l, err := c[i].Snapshots()
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
		}

		for _, s := range l {
			log.Printf("Name: %s\n", s.Name)
			log.Printf("Comment path: %s\n", s.CommentPath)
			log.Printf("Timestamp: %s\n", s.Timestamp)
			log.Printf("LXC path: %s\n", s.Path)
			log.Println()
		}
	}
}
