// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (a *App) IsLinked(username string, platform string) bool {
	_, err := a.Srv().Store.ExtRef().GetByRealUserIdAndPlatform(username, platform)
	if err != nil {
		return false
	}
	return true

}

func (a *App) LinkAccount(extRef *model.ExtRef) *model.AppError {
	_, err := a.Srv().Store.ExtRef().Save(extRef)
	if err != nil {
		return model.NewAppError("LinkAccount", "app.ext_ref.link_account.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
