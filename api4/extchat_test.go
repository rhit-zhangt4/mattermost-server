// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	_ "image/gif"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinkAccount(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	Client := th.Client

	var externalPlatform = "telegram"
	var externalId = "11234567890"
	realUserId := th.BasicUser.Id
	// ext_ref := &model.ExtRef{
	// 	RealUserId:       th.BasicUser.Id,
	// 	ExternalId:       externalId,
	// 	ExternalPlatform: externalPlatform,
	// 	AliasUserId:      "",
	// }

	ok, resp := Client.IsLinked(realUserId, externalPlatform)
	CheckNoError(t, resp)
	require.Equal(t, ok, []byte("false"), "did not return false")

	_, resp = Client.LinkAccount(externalId, externalPlatform)
	CheckNoError(t, resp)

	ok, resp = Client.IsLinked(realUserId, externalPlatform)
	CheckNoError(t, resp)
	require.Equal(t, ok, []byte("true"), "did not return true")

}
