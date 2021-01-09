// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

type Secret struct {
	SecretName  string `json:"secret_name"`
	SecretValue string `json:"secret_value"`
}

// func (secret *Secret) IsValid() *AppError {
// 	if !IsValidId(secret.SecretName) {
// 		return NewAppError("Emoji.IsValid", "model.emoji.id.app_error", nil, "", http.StatusBadRequest)
// 	}
// 	return nil
// }

// func (secret *Secret) ToJson() string {
// 	b, _ := json.Marshal(secret)
// 	return string(b)
// }

// func SecretFromJson(data io.Reader) *PublicEmoji {
// 	var secret *Secret
// 	json.NewDecoder(data).Decode(&secret)
// 	return secret
// }
