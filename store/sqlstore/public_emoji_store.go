// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/pkg/errors"
)

type SqlPublicEmojiStore struct {
	SqlStore
	metrics einterfaces.MetricsInterface
}

func newSqlPublicEmojiStore(sqlStore SqlStore, metrics einterfaces.MetricsInterface) store.PublicEmojiStore {
	s := &SqlPublicEmojiStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.PublicEmoji{}, "PublicEmoji").SetKeys(false, "EmojiId")
		table.ColMap("EmojiId").SetMaxSize(26)
	}

	return s
}

func (es SqlPublicEmojiStore) createIndexesIfNotExists() {
	es.CreateIndexIfNotExists("idx_public__emoji", "PublicEmoji", "EmojiId")
}

func (es SqlPublicEmojiStore) Save(public_emoji *model.PublicEmoji) (*model.PublicEmoji, error) {
	if err := public_emoji.IsValid(); err != nil {
		return nil, err
	}

	if err := es.GetMaster().Insert(public_emoji); err != nil {
		return nil, errors.Wrap(err, "error saving public emoji")
	}

	return public_emoji, nil
}

func (es SqlPublicEmojiStore) GetAllPublicEmojis() ([]*model.PublicEmoji, error) {
	var publicEmojies []*model.PublicEmoji
	if _, err := es.GetReplica().Select(&publicEmojies,
		`SELECT
			*
		FROM
			PublicEmoji`); err != nil {
		return nil, errors.Wrapf(err, "error getting all public emoji")
	}
	return publicEmojies, nil
}

func (es SqlPublicEmojiStore) CheckIsPublicEmojis(emojiId string) bool {

	count, err := es.GetReplica().SelectInt(`
		SELECT count(*)
			FROM PublicEmoji
		WHERE
			EmojiId = :EmojiId
			`, map[string]interface{}{"EmojiId": emojiId})

	if err != nil || count == 0 {
		return false
	}

	return true
}

func (es SqlPublicEmojiStore) DeleteAccessByEmojiId(emojiId string) error {
	sql := `DELETE
		FROM PublicEmoji
	WHERE
		EmojiId = :EmojiId`

	queryParams := map[string]string{
		"EmojiId": emojiId,
	}

	_, err := es.GetMaster().Exec(sql, queryParams)
	if err != nil {
		//mlog.Warn("Failed to delete access", mlog.Err(err))
		return err
	}
	return nil
}
