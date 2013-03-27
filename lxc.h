/*
 * lxc.h: Go bindings for lxc
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

extern bool container_create(struct lxc_container *, char *, char **);
extern bool container_defined(struct lxc_container *);
extern bool container_destroy(struct lxc_container *);
extern bool container_freeze(struct lxc_container *);
extern bool container_running(struct lxc_container *);
extern bool container_set_config_item(struct lxc_container *, char *, char *);
extern bool container_shutdown(struct lxc_container *, int);
extern bool container_start(struct lxc_container *, int, char **);
extern bool container_stop(struct lxc_container *);
extern bool container_unfreeze(struct lxc_container *);
extern bool container_wait(struct lxc_container *, char *, int);
extern char* container_config_file_name(struct lxc_container *);
extern char* container_get_config_item(struct lxc_container *, char *);
extern char* container_get_keys(struct lxc_container *, char *);
extern const char* container_state(struct lxc_container *);
extern pid_t container_init_pid(struct lxc_container *);
extern void container_want_daemonize(struct lxc_container *);
