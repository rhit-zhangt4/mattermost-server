// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package storetest

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtRefStore(t *testing.T, ss store.Store) {
	t.Run("SaveExtRef", func(t *testing.T) { testSaveExtRef(t, ss) })
	t.Run("GetByAliasUserId", func(t *testing.T) { testGetByAliasUserId(t, ss) })
	t.Run("UpdateRealId", func(t *testing.T) { testUpdateRealId(t, ss) })
	t.Run("GetByExtIdAndPlatform", func(t *testing.T) { testGetByExtIdAndPlatform(t, ss) })
	t.Run("GetByRealUserIdAndPlatform", func(t *testing.T) { testGetByRealUserIdAndPlatform(t, ss) })
	t.Run("Unlink", func(t *testing.T) { testUnlink(t, ss) })
}

var testRealUserId1 = model.NewId()
var testAliasUserId1 = model.NewId()
var testExternalId1 = model.NewId()
var testExternalPlatform1 = "Telegram"

func testSaveExtRef(t *testing.T, ss store.Store) {
	extRef := &model.ExtRef{
		ExternalId:       testExternalId1,
		AliasUserId:      testAliasUserId1,
		ExternalPlatform: testExternalPlatform1,
	}
	result, err := ss.ExtRef().Save(extRef)
	require.Nil(t, err)
	assert.Equal(t, result.ExternalId, testExternalId1)
	assert.Equal(t, result.ExternalPlatform, testExternalPlatform1)
}

func testGetByExtIdAndPlatform(t *testing.T, ss store.Store) {
	result, err := ss.ExtRef().GetByExtIdAndPlatform(testExternalId1, testExternalPlatform1)
	require.Nil(t, err)
	assert.Equal(t, result.RealUserId, testRealUserId1)
}

func testGetByRealUserIdAndPlatform(t *testing.T, ss store.Store) {
	result, err := ss.ExtRef().GetByRealUserIdAndPlatform(testRealUserId1, testExternalPlatform1)
	require.Nil(t, err)
	assert.Equal(t, result.ExternalId, testExternalId1)
}

func testGetByAliasUserId(t *testing.T, ss store.Store) {
	result, err := ss.ExtRef().GetByAliasUserId(testAliasUserId1)
	require.Nil(t, err)
	assert.Equal(t, result.ExternalId, testExternalId1)
	assert.Equal(t, result.ExternalPlatform, testExternalPlatform1)
}

func testUpdateRealId(t *testing.T, ss store.Store) {
	err := ss.ExtRef().UpdateRealId(testRealUserId1, testExternalId1, testExternalPlatform1)
	require.Nil(t, err)
}

func testUnlink(t *testing.T, ss store.Store) {
	err := ss.ExtRef().Unlink(testRealUserId1, testExternalPlatform1)
	require.Nil(t, err)
	_, err = ss.ExtRef().GetByExtIdAndPlatform(testExternalId1, testExternalPlatform1)
	require.NotNil(t, err)
}
