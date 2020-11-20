// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"io"
	"net/http"
)

type PublicEmoji struct {
	EmojiId string `json:"emoji_id"`
}

func (public_emoji *PublicEmoji) IsValid() *AppError {
	if !IsValidId(public_emoji.EmojiId) {
		return NewAppError("Emoji.IsValid", "model.emoji.id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func (public_emoji *PublicEmoji) ToJson() string {
	b, _ := json.Marshal(public_emoji)
	return string(b)
}

func PublicEmojiFromJson(data io.Reader) *PublicEmoji {
	var public_emoji *PublicEmoji
	json.NewDecoder(data).Decode(&public_emoji)
	return public_emoji
}
