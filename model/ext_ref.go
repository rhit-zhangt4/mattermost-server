// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"io"
)

type ExtRef struct {
	RealUserId       string `json:"real_user_id"`
	AliasUserId      string `json:"alias_user_id"`
	ExternalId       string `json:"external_id"`
	ExternalPlatform string `json:"external_platform"`
}

// func (emoji_access *EmojiAccess) IsValid() *AppError {
// 	if !IsValidId(emoji_access.EmojiId) {
// 		return NewAppError("Emoji.IsValid", "model.emoji.id.app_error", nil, "", http.StatusBadRequest)
// 	}

// 	if len(emoji_access.UserId) > 26 {
// 		return NewAppError("Emoji.IsValid", "model.emoji.user_id.app_error", nil, "", http.StatusBadRequest)
// 	}

// 	return nil
// }

func (ext_ref *ExtRef) ToJson() string {
	b, _ := json.Marshal(ext_ref)
	return string(b)
}

func ExtRefFromJson(data io.Reader) *ExtRef {
	var ext_ref *ExtRef
	json.NewDecoder(data).Decode(&ext_ref)
	return ext_ref
}
