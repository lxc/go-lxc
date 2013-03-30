/*
 * state.go: Go bindings for lxc
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

// #include <lxc/lxc.h>
// #include <lxc/lxccontainer.h>
import "C"

type State int

const (
	STOPPED  State = C.STOPPED
	STARTING State = C.STARTING
	RUNNING  State = C.RUNNING
	STOPPING State = C.STOPPING
	ABORTING State = C.ABORTING
	FREEZING State = C.FREEZING
	FROZEN   State = C.FROZEN
	THAWED   State = C.THAWED
)

var stateMap = map[string]State{
	"STOPPED":  STOPPED,
	"STARTING": STARTING,
	"RUNNING":  RUNNING,
	"STOPPING": STOPPING,
	"ABORTING": ABORTING,
	"FREEZING": FREEZING,
	"FROZEN":   FROZEN,
	"THAWED":   THAWED,
}

// State as string
func (t State) String() string {
	switch t {
	case STOPPED:
		return "STOPPED"
	case STARTING:
		return "STARTING"
	case RUNNING:
		return "RUNNING"
	case STOPPING:
		return "STOPPING"
	case ABORTING:
		return "ABORTING"
	case FREEZING:
		return "FREEZING"
	case FROZEN:
		return "FROZEN"
	case THAWED:
		return "THAWED"
	}
	return "<INVALID>"
}
