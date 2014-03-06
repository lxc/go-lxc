/*
 * destroy_snapshots.go
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
	for _, v := range lxc.Containers() {
		log.Printf("%s\n", v.Name())
		l, err := v.Snapshots()
		if err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}

		for _, s := range l {
			log.Printf("Destroying Snaphot: %s\n", s.Name)
			if err := v.DestroySnapshot(s); err != nil {
				log.Fatalf("ERROR: %s\n", err.Error())
			}
		}
	}
}
