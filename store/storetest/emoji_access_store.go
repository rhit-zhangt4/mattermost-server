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

func TestEmojiAccessStore(t *testing.T, ss store.Store) {
	t.Run("EmojiAccessSave", func(t *testing.T) { testEmojiAccessSave(t, ss) })
	t.Run("SavePrivateEmoji", func(t *testing.T) { testSavePrivateEmoji(t, ss) })
	t.Run("EmojiAccessGetByUserIdAndEmojiId", func(t *testing.T) { testEmojiAccessGetByUserIdAndEmojiId(t, ss) })
	t.Run("EmojiAccessGetMultipleByUserId", func(t *testing.T) { testGetMultipleByUserId(t, ss) })
	t.Run("DeleteEmojiAccessGetByUserIdAndEmojiId", func(t *testing.T) { testDeleteEmojiAccessGetByUserIdAndEmojiId(t, ss) })
	t.Run("SavePrivateEmoji", func(t *testing.T) { testSavePrivateEmoji(t, ss) })
	t.Run("DeleteEmojiAccessGetByEmojiId", func(t *testing.T) { testDeleteEmojiAccessGetByEmojiId(t, ss) })

}

var testEmojiId1 = model.NewId()
var testUserId1 = model.NewId()

// var testEmojiId2 = model.NewId()
var testUserId2 = model.NewId()

func testDeleteEmojiAccessGetByUserIdAndEmojiId(t *testing.T, ss store.Store) {
	err := ss.EmojiAccess().DeleteAccessByUserIdAndEmojiId(testUserId2, testEmojiId1)
	require.Nil(t, err)
	_, err = ss.EmojiAccess().GetByUserIdAndEmojiId(testUserId2, testEmojiId1)
	require.NotNil(t, err)
}

func testDeleteEmojiAccessGetByEmojiId(t *testing.T, ss store.Store) {
	err := ss.EmojiAccess().DeleteAccessByEmojiId(testEmojiId1)
	require.Nil(t, err)
	_, err = ss.EmojiAccess().GetByUserIdAndEmojiId(testUserId1, testEmojiId1)
	require.NotNil(t, err)
	_, err = ss.EmojiAccess().GetByUserIdAndEmojiId(testUserId2, testEmojiId1)
	require.NotNil(t, err)
}

func testEmojiAccessGetByUserIdAndEmojiId(t *testing.T, ss store.Store) {
	result, err := ss.EmojiAccess().GetByUserIdAndEmojiId(testUserId1, testEmojiId1)
	require.Nil(t, err)
	assert.Equal(t, result.UserId, testUserId1)
	assert.Equal(t, result.EmojiId, testEmojiId1)
}

func testSavePrivateEmoji(t *testing.T, ss store.Store) {
	emoji_access3 := &model.EmojiAccess{
		EmojiId: testEmojiId1,
		UserId:  testUserId2,
	}
	result, err := ss.EmojiAccess().Save(emoji_access3)
	require.Nil(t, err)
	assert.Equal(t, result.UserId, testUserId2)
	assert.Equal(t, result.EmojiId, testEmojiId1)
}

func testEmojiAccessSave(t *testing.T, ss store.Store) {

	emoji_access1 := &model.EmojiAccess{
		EmojiId: testEmojiId1,
		UserId:  testUserId1,
	}
	// emoji_access2 := &model.EmojiAccess{
	// 	EmojiId: testEmojiId2,
	// 	UserId:  testUserId2,
	// }

	_, err := ss.EmojiAccess().Save(emoji_access1)
	require.Nil(t, err)
	// _, err = ss.EmojiAccess().Save(emoji_access2)
	// require.Nil(t, err)

	// assert.Len(t, emoji1.Id, 26, "should've set id for emoji")

}

func testGetMultipleByUserId(t *testing.T, ss store.Store) {
	result, err := ss.EmojiAccess().GetMultipleByUserId([]string{testUserId1})
	require.Nil(t, err)
	assert.Equal(t, result[0].EmojiId, testEmojiId1)

	result, err = ss.EmojiAccess().GetMultipleByUserId([]string{testUserId2})
	require.Nil(t, err)
	assert.Equal(t, result[0].EmojiId, testEmojiId1)
	// assert.Equal(t, result[1], testEmojiId2)
}
