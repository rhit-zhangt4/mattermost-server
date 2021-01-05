// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"net/http"

	"github.com/mattermost/mattermost-server/v5/extchat"
)

func (api *API) InitExtChat() {
	api.BaseRoutes.ExtChat.Handle("/authenticate", api.ApiHandler(startAuthentication)).Methods("POST")
	api.BaseRoutes.ExtChat.Handle("/verify", api.ApiSessionRequired(verifyPasscode)).Methods("GET")
	// api.BaseRoutes.Emojis.Handle("/search", api.ApiSessionRequired(searchEmojis)).Methods("POST")
	// api.BaseRoutes.Emojis.Handle("/autocomplete", api.ApiSessionRequired(autocompleteEmojis)).Methods("GET")
	// api.BaseRoutes.Emoji.Handle("", api.ApiSessionRequired(deleteEmoji)).Methods("DELETE")
	// api.BaseRoutes.Emoji.Handle("", api.ApiSessionRequired(getEmoji)).Methods("GET")
	// api.BaseRoutes.EmojiByName.Handle("", api.ApiSessionRequired(getEmojiByName)).Methods("GET")
	// api.BaseRoutes.Emoji.Handle("/image", api.ApiSessionRequiredTrustRequester(getEmojiImage)).Methods("GET")
	// api.BaseRoutes.Emojis.Handle("/private", api.ApiSessionRequired(createPrivateEmoji)).Methods("POST")
	// api.BaseRoutes.Emojis.Handle("/private", api.ApiSessionRequired(getPrivateEmojiList)).Methods("GET")
	// api.BaseRoutes.Emoji.Handle("/privateimage", api.ApiSessionRequiredTrustRequester(getPrivateEmojiImage)).Methods("GET")
	// api.BaseRoutes.Emoji.Handle("/checkprivate", api.ApiSessionRequiredTrustRequester(getCanAccessPrivateEmojiImage)).Methods("GET")
	// api.BaseRoutes.Emoji.Handle("/save", api.ApiSessionRequiredTrustRequester(savePrivateEmoji)).Methods("POST")
	// api.BaseRoutes.Emojis.Handle("/public", api.ApiSessionRequiredTrustRequester(getPublicEmojiList)).Methods("GET")
	// api.BaseRoutes.Emoji.Handle("/access", api.ApiSessionRequired(deleteEmojiAccess)).Methods("DELETE")
	// api.BaseRoutes.Emoji.Handle("/withAccess", api.ApiSessionRequired(deleteEmojiWithAccess)).Methods("DELETE")
}

func getAdapterFromPlatform(platform string) (extchat.ExtChatAdapter, bool) {
	switch platform {
	case "telegram":
		return &extchat.TelegramAdapter{}, true
	default:
		return nil, false
	}
}

func startAuthentication(c *Context, w http.ResponseWriter, r *http.Request) {
	adapter, ok := getAdapterFromPlatform(c.Params.ExtChatPlatform)
	if !ok {
		//error
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		c.SetInvalidUrlParam("username")
		return
	}
	err := adapter.StartAuthentication(c.App, username)
	if err != nil {
		//error
	}
	ReturnStatusOK(w)
}

func verifyPasscode(c *Context, w http.ResponseWriter, r *http.Request) {

}
