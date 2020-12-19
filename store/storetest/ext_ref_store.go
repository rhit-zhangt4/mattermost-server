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

// // var testEmojiId2 = model.NewId()
// var testUserId2 = model.NewId()
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
	err := ss.ExtRef().Unlink(testExternalId1, testExternalPlatform1)
	require.Nil(t, err)
	_, err = ss.ExtRef().GetByExtIdAndPlatform(testExternalId1, testExternalPlatform1)
	require.NotNil(t, err)
}

// func testDeleteEmojiAccessGetByEmojiId(t *testing.T, ss store.Store) {
// 	err := ss.EmojiAccess().DeleteAccessByEmojiId(testEmojiId1)
// 	require.Nil(t, err)
// 	_, err = ss.EmojiAccess().GetByUserIdAndEmojiId(testUserId1, testEmojiId1)
// 	require.NotNil(t, err)
// 	_, err = ss.EmojiAccess().GetByUserIdAndEmojiId(testUserId2, testEmojiId1)
// 	require.NotNil(t, err)
// }

// func testEmojiAccessGetByUserIdAndEmojiId(t *testing.T, ss store.Store) {
// 	result, err := ss.ExtRef().GetByUserIdAndEmojiId(testRealUserId1, testExternalPlatform1)
// 	require.Nil(t, err)
// 	assert.Equal(t, result.UserId, testRealUserId1)
// 	assert.Equal(t, result.EmojiId, testEmojiId1)
// }

// func testSavePrivateEmoji(t *testing.T, ss store.Store) {
// 	emoji_access3 := &model.EmojiAccess{
// 		EmojiId: testEmojiId1,
// 		UserId:  testUserId2,
// 	}
// 	result, err := ss.EmojiAccess().Save(emoji_access3)
// 	require.Nil(t, err)
// 	assert.Equal(t, result.UserId, testUserId2)
// 	assert.Equal(t, result.EmojiId, testEmojiId1)
// }

// func testEmojiAccessSave(t *testing.T, ss store.Store) {

// 	emoji_access1 := &model.EmojiAccess{
// 		EmojiId: testEmojiId1,
// 		UserId:  testUserId1,
// 	}
// 	// emoji_access2 := &model.EmojiAccess{
// 	// 	EmojiId: testEmojiId2,
// 	// 	UserId:  testUserId2,
// 	// }

// 	_, err := ss.EmojiAccess().Save(emoji_access1)
// 	require.Nil(t, err)
// 	// _, err = ss.EmojiAccess().Save(emoji_access2)
// 	// require.Nil(t, err)

// 	// assert.Len(t, emoji1.Id, 26, "should've set id for emoji")

// }

// func testGetMultipleByUserId(t *testing.T, ss store.Store) {
// 	result, err := ss.EmojiAccess().GetMultipleByUserId([]string{testUserId1})
// 	require.Nil(t, err)
// 	assert.Equal(t, result[0].EmojiId, testEmojiId1)

// 	result, err = ss.EmojiAccess().GetMultipleByUserId([]string{testUserId2})
// 	require.Nil(t, err)
// 	assert.Equal(t, result[0].EmojiId, testEmojiId1)
// 	// assert.Equal(t, result[1], testEmojiId2)
// }
