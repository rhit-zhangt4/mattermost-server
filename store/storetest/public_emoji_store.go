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

func TestPublicEmojiStore(t *testing.T, ss store.Store) {
	t.Run("PublicEmojiGetAllPublicEmoji", func(t *testing.T) { testPublicEmojiGetAllPublicEmoji(t, ss) })
	// t.Run("EmojiAccessGetByUserIdAndEmojiId", func(t *testing.T) { testEmojiAccessGetByUserIdAndEmojiId(t, ss) })
}

var testEmojiId = model.NewId()

func testPublicEmojiGetAllPublicEmoji(t *testing.T, ss store.Store) {
	_, err := ss.PublicEmoji().GetAllPublicEmojis()
	require.Nil(t, err)
	public_emoji := &model.PublicEmoji{
		EmojiId: testEmojiId,
	}

	result := ss.PublicEmoji().CheckIsPublicEmojis(testEmojiId)
	require.False(t, result)

	_, err = ss.PublicEmoji().Save(public_emoji)
	require.Nil(t, err)

	r, err := ss.PublicEmoji().GetAllPublicEmojis()
	require.Nil(t, err)
	assert.Len(t, r, 1, "should return one element")

	result = ss.PublicEmoji().CheckIsPublicEmojis(testEmojiId)
	require.True(t, result)

	err = ss.PublicEmoji().DeleteAccessByEmojiId(testEmojiId)
	require.Nil(t, err)

	result = ss.PublicEmoji().CheckIsPublicEmojis(testEmojiId)
	require.False(t, result)

}
