// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package extchat

import (
	"github.com/mattermost/mattermost-server/v5/app"
	"github.com/mattermost/mattermost-server/v5/model"
)

type ExtChatAdapter interface {
	StartAuthentication(a app.AppIface, username string) *model.AppError
	VerifyPasscode(a app.AppIface, username string, code string) (*model.ExtRef, *model.AppError)
}
