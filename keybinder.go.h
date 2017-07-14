#ifndef __KEYBINDER_GO_H__
#define __KEYBINDER_GO_H__

#include <stdlib.h>
#include <keybinder.h>
#include <glib.h>

extern void goKeybinderHandler(char *keystring, void *user_data);

static inline void keybinder_handler(const char *keystring, void *user_data) {
	goKeybinderHandler((char *)keystring, user_data);
}

static inline gboolean _keybinder_bind(const char *keystring, void *user_data) {
	keybinder_bind(keystring, (KeybinderHandler)(keybinder_handler), user_data);
}

static inline void _keybinder_unbind(const char *keystring) {
	keybinder_unbind(keystring, (KeybinderHandler)(keybinder_handler));
}

#endif