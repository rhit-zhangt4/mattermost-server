// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"io"
	"net/http"
)

type EmojiAccess struct {
	EmojiId string `json:"emoji_id"`
	UserId  string `json:"user_id"`
}

func (emoji_access *EmojiAccess) IsValid() *AppError {
	if !IsValidId(emoji_access.EmojiId) {
		return NewAppError("Emoji.IsValid", "model.emoji.id.app_error", nil, "", http.StatusBadRequest)
	}

	if len(emoji_access.UserId) > 26 {
		return NewAppError("Emoji.IsValid", "model.emoji.user_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (emoji_access *EmojiAccess) ToJson() string {
	b, _ := json.Marshal(emoji_access)
	return string(b)
}

func EmojiAccessFromJson(data io.Reader) *EmojiAccess {
	var emoji_access *EmojiAccess
	json.NewDecoder(data).Decode(&emoji_access)
	return emoji_access
}
