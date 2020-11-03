// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package storetest

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmojiAccessStore(t *testing.T, ss store.Store) {
	t.Run("EmojiAccessSave", func(t *testing.T) { testEmojiAccessSave(t, ss) })
	t.Run("EmojiAccessGetByUserIdAndEmojiId", func(t *testing.T) { testEmojiAccessGetByUserIdAndEmojiId(t, ss) })
}

var testEmojiId = model.NewId()
var testUserId = model.NewId()

func testEmojiAccessGetByUserIdAndEmojiId(t *testing.T, ss store.Store) {
	_, err := ss.EmojiAccess().GetByUserIdAndEmojiId(testUserId, testEmojiId)
	require.Nil(t, err)
}

func testEmojiAccessSave(t *testing.T, ss store.Store) {

	emoji_access1 := &model.EmojiAccess{
		EmojiId: testEmojiId,
		UserId:  testUserId,
	}

	_, err := ss.EmojiAccess().Save(emoji_access1)
	require.Nil(t, err)

	// assert.Len(t, emoji1.Id, 26, "should've set id for emoji")

}
