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

#include <stdio.h>
#include <stdbool.h>

#include <lxc/lxc.h>
#include <lxc/lxccontainer.h>

bool lxc_container_defined(struct lxc_container *c) {
	return c->is_defined(c);
}

const char* lxc_container_state(struct lxc_container *c) {
	return c->state(c);
}

bool lxc_container_running(struct lxc_container *c) {
	return c->is_running(c);
}

bool lxc_container_freeze(struct lxc_container *c) {
	return c->freeze(c);
}

bool lxc_container_unfreeze(struct lxc_container *c) {
	return c->unfreeze(c);
}

pid_t lxc_container_init_pid(struct lxc_container *c) {
	return c->init_pid(c);
}

void lxc_container_want_daemonize(struct lxc_container *c) {
	c->want_daemonize(c);
}

bool lxc_container_create(struct lxc_container *c, char *t, char **argv) {
    return c->create(c, t, NULL, NULL, LXC_CREATE_QUIET, argv);
}

bool lxc_container_start(struct lxc_container *c, int useinit, char ** argv) {
	return c->start(c, useinit, argv);
}

bool lxc_container_stop(struct lxc_container *c) {
	return c->stop(c);
}

bool lxc_container_reboot(struct lxc_container *c) {
	return c->reboot(c);
}

bool lxc_container_shutdown(struct lxc_container *c, int timeout) {
	return c->shutdown(c, timeout);
}

char* lxc_container_config_file_name(struct lxc_container *c) {
	return c->config_file_name(c);
}

bool lxc_container_destroy(struct lxc_container *c) {
	return c->destroy(c);
}

bool lxc_container_wait(struct lxc_container *c, char *state, int timeout) {
	return c->wait(c, state, timeout);
}

char* lxc_container_get_config_item(struct lxc_container *c, char *key) {
	int len = c->get_config_item(c, key, NULL, 0);
	if (len <= 0) {
		return NULL;
	}

	char* value = (char*)malloc(sizeof(char)*len + 1);
	if (c->get_config_item(c, key, value, len + 1) != len) {
		return NULL;
	}
	return value;
}

bool lxc_container_set_config_item(struct lxc_container *c, char *key, char *value) {
	return c->set_config_item(c, key, value);
}

bool lxc_container_clear_config_item(struct lxc_container *c, char *key) {
	return c->clear_config_item(c, key);
}

char* lxc_container_get_keys(struct lxc_container *c, char *key) {
	int len = c->get_keys(c, key, NULL, 0);
	if (len <= 0) {
		return NULL;
	}

	char* value = (char*)malloc(sizeof(char)*len + 1);
	if (c->get_keys(c, key, value, len + 1) != len) {
		return NULL;
	}
	return value;
}

char* lxc_container_get_cgroup_item(struct lxc_container *c, char *key) {
	int len = c->get_cgroup_item(c, key, NULL, 0);
	if (len <= 0) {
		return NULL;
	}

	char* value = (char*)malloc(sizeof(char)*len + 1);
	if (c->get_cgroup_item(c, key, value, len + 1) != len) {
		return NULL;
	}
	return value;
}

bool lxc_container_set_cgroup_item(struct lxc_container *c, char *key, char *value) {
	return c->set_cgroup_item(c, key, value);
}

const char* lxc_container_get_config_path(struct lxc_container *c) {
	return c->get_config_path(c);
}

bool lxc_container_set_config_path(struct lxc_container *c, char *path) {
	return c->set_config_path(c, path);
}

bool lxc_container_load_config(struct lxc_container *c, char *alt_file) {
	return c->load_config(c, alt_file);
}

bool lxc_container_save_config(struct lxc_container *c, char *alt_file) {
	return c->save_config(c, alt_file);
}

bool lxc_container_clone(struct lxc_container *c, const char *newname, int flags, const char *bdevtype) {
    return c->clone(c, newname, NULL, flags, bdevtype, NULL, 0, NULL) != NULL;
}
