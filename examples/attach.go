/*
 * attach.go
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
	"flag"
	"log"
	"os"

	"gopkg.in/lxc/go-lxc.v1"
)

var (
	lxcpath string
	name    string
	clear   bool
)

func init() {
	flag.StringVar(&lxcpath, "lxcpath", lxc.DefaultConfigPath(), "Use specified container path")
	flag.StringVar(&name, "name", "rubik", "Name of the original container")
	flag.BoolVar(&clear, "clear", false, "Attach with clear environment")
	flag.Parse()
}

func main() {
	c, err := lxc.NewContainer(name, lxcpath)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	defer lxc.PutContainer(c)

	if clear {
		log.Printf("AttachShellWithClearEnvironment\n")
		if err := c.AttachShellWithClearEnvironment(); err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}

		log.Printf("RunCommandWithClearEnvironment\n")
		if _, err := c.RunCommandWithClearEnvironment(os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), "uname", "-a"); err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}

	} else {
		log.Printf("AttachShell\n")
		if err := c.AttachShell(); err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}

		log.Printf("RunCommand\n")
		if _, err := c.RunCommand(os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), "uname", "-a"); err != nil {
			log.Fatalf("ERROR: %s\n", err.Error())
		}
	}
}
