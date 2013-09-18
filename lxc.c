/*
 * lxc.c: Go bindings for lxc
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

#include <stdbool.h>

#include <lxc/lxc.h>
#include <lxc/lxccontainer.h>
#include <lxc/attach_options.h>

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

bool lxc_container_create(struct lxc_container *c, char *t, int flags, char **argv) {
    return c->create(c, t, NULL, NULL, !!(flags & LXC_CREATE_QUIET), argv);
}

bool lxc_container_start(struct lxc_container *c, int useinit, char **argv) {
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

extern int lxc_container_console_getfd(struct lxc_container *c, int ttynum) {
    int masterfd;

    if (c->console_getfd(c, &ttynum, &masterfd) < 0) {
        return -1;
    }
    return masterfd;
}

extern bool lxc_container_console(struct lxc_container *c, int ttynum, int stdinfd, int stdoutfd, int stderrfd, int escape) {

    if (c->console(c, ttynum, stdinfd, stdoutfd, stderrfd, escape) == 0) {
        return true;
    }
    return false;
}

extern char** lxc_container_get_interfaces(struct lxc_container *c) {
    return c->get_interfaces(c);
}

extern char** lxc_container_get_ips(struct lxc_container *c, char *interface, char *family, int scope) {
    return c->get_ips(c, interface, family, scope);
}

extern int lxc_container_attach(struct lxc_container *c) {
    int ret;
    pid_t pid;
    lxc_attach_options_t default_options = LXC_ATTACH_OPTIONS_DEFAULT;

/*
    remount_sys_proc
    When using -s and the mount namespace is not included, this flag will cause lxc-attach to remount /proc and /sys to reflect the current other namespace contexts.
    default_options.attach_flags |= LXC_ATTACH_REMOUNT_PROC_SYS;

    elevated_privileges
    Do  not  drop privileges when running command inside the container. If this option is specified, the new process will not be added to the container's cgroup(s) and it will not drop its capabilities before executing.
    default_options.attach_flags &= ~(LXC_ATTACH_MOVE_TO_CGROUP | LXC_ATTACH_DROP_CAPABILITIES | LXC_ATTACH_APPARMOR);

    Specify the namespaces to attach to, as a pipe-separated list, e.g. NETWORK|IPC. Allowed values are MOUNT, PID, UTSNAME, IPC, USER and NETWORK.
    default_options.namespaces = namespace_flags; // lxc_fill_namespace_flags(arg, &namespace_flags);

    Specify the architecture which the kernel should appear to be running as to the command executed.
    default_options.personality = new_personality; // lxc_config_parse_arch(arg);

    Keep the current environment for attached programs.
    Clear the environment before attaching, so no undesired environment variables leak into the container.

    default_options.env_policy = env_policy; // LXC_ATTACH_KEEP_ENV or LXC_ATTACH_CLEAR_ENV

    default_options.extra_env_vars = extra_env;
    default_options.extra_keep_env = extra_keep;
*/

    ret = c->attach(c, lxc_attach_run_shell, NULL, &default_options, &pid);
    if (ret < 0)
        return -1;

    ret = lxc_wait_for_pid_status(pid);
    if (ret < 0)
        return -1;

    if (WIFEXITED(ret))
        return WEXITSTATUS(ret);

    return -1;
}

extern int lxc_container_attach_run_wait(struct lxc_container *c, char **argv) {
    int ret;
    lxc_attach_options_t default_options = LXC_ATTACH_OPTIONS_DEFAULT;

    ret = c->attach_run_wait(c, &default_options, argv[0], (const char * const*)argv);
    if (WIFEXITED(ret) && WEXITSTATUS(ret) == 255)
        return -1;
    return ret;
}
