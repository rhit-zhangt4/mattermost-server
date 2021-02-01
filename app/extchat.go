// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (a *App) IsLinked(username string, platform string) bool {
	_, err := a.Srv().Store.ExtRef().GetByRealUserIdAndPlatform(username, platform)
	return err != nil
}

func (a *App) LinkAccount(extRef *model.ExtRef) *model.AppError {
	_, err := a.Srv().Store.ExtRef().Save(extRef)
	if err != nil {
		return model.NewAppError("LinkAccount", "app.ext_ref.link_account.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (a *App) CreateAliasAccount(userName string, externalId string, platform string) *model.AppError {
	userModel := &model.User{Email: "",
		Nickname: userName,
		Password: "",
		Username: userName,
		IsAlias:  true,
	}
	user, err := a.Srv().Store.User().Save(userModel)
	if err != nil {
		return model.NewAppError("CreateAlias", "app.ext_ref.create_alias.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}
	ext_ref := &model.ExtRef{
		RealUserId:       "",
		ExternalId:       externalId,
		ExternalPlatform: platform,
		AliasUserId:      user.Id,
	}
	_, extRefErr := a.Srv().Store.ExtRef().Save(ext_ref)
	if extRefErr != nil {
		return model.NewAppError("CreateAlias", "app.ext_ref.save_ext_ref.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
