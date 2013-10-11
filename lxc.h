/*
 * lxc.h: Go bindings for lxc
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

extern bool lxc_container_clear_config_item(struct lxc_container *, char *);
extern bool lxc_container_clone(struct lxc_container *, const char *, int, const char *);
extern bool lxc_container_console(struct lxc_container *, int, int, int, int, int);
extern bool lxc_container_create(struct lxc_container *, char *, int, char **);
extern bool lxc_container_defined(struct lxc_container *);
extern bool lxc_container_destroy(struct lxc_container *);
extern bool lxc_container_freeze(struct lxc_container *);
extern bool lxc_container_load_config(struct lxc_container *, char *);
extern bool lxc_container_reboot(struct lxc_container *);
extern bool lxc_container_running(struct lxc_container *);
extern bool lxc_container_save_config(struct lxc_container *, char *);
extern bool lxc_container_set_cgroup_item(struct lxc_container *, char *key, char *);
extern bool lxc_container_set_config_item(struct lxc_container *, char *, char *);
extern bool lxc_container_set_config_path(struct lxc_container *, char *);
extern bool lxc_container_shutdown(struct lxc_container *, int);
extern bool lxc_container_start(struct lxc_container *, int, char **);
extern bool lxc_container_stop(struct lxc_container *);
extern bool lxc_container_unfreeze(struct lxc_container *);
extern bool lxc_container_wait(struct lxc_container *, char *, int);
extern char* lxc_container_config_file_name(struct lxc_container *);
extern char* lxc_container_get_cgroup_item(struct lxc_container *, char *);
extern char* lxc_container_get_config_item(struct lxc_container *, char *);
extern char** lxc_container_get_interfaces(struct lxc_container *);
extern char** lxc_container_get_ips(struct lxc_container *, char *, char *, int);
extern char* lxc_container_get_keys(struct lxc_container *, char *);
extern const char* lxc_container_get_config_path(struct lxc_container *);
extern const char* lxc_container_state(struct lxc_container *);
extern int lxc_container_attach_run_wait(struct lxc_container *, char **);
extern int lxc_container_attach(struct lxc_container *);
extern int lxc_container_console_getfd(struct lxc_container *, int);
extern pid_t lxc_container_init_pid(struct lxc_container *);
extern void lxc_container_want_daemonize(struct lxc_container *);
extern bool lxc_container_want_close_all_fds(struct lxc_container *);
extern bool lxc_container_may_control(struct lxc_container *);

extern int lxc_container_snapshot(struct lxc_container *);
extern bool lxc_container_snapshot_restore(struct lxc_container *, char *, char *);

extern struct lxc_snapshot* lxc_container_snapshot_list(struct lxc_container *);
extern int lxc_container_snapshot_list_size(struct lxc_container *);
