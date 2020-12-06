// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/pkg/errors"
)

type SqlEmojiAccessStore struct {
	SqlStore
	metrics einterfaces.MetricsInterface
}

func newSqlEmojiAccessStore(sqlStore SqlStore, metrics einterfaces.MetricsInterface) store.EmojiAccessStore {
	s := &SqlEmojiAccessStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.EmojiAccess{}, "EmojiAccess").SetKeys(false, "EmojiId", "UserId")
		table.ColMap("EmojiId").SetMaxSize(26)
		table.ColMap("UserId").SetMaxSize(64)
	}

	return s
}

func (es SqlEmojiAccessStore) createIndexesIfNotExists() {
	es.CreateIndexIfNotExists("idx_emoji_access_user", "EmojiAccess", "UserId")
	es.CreateIndexIfNotExists("idx_emoji_access_emoji", "EmojiAccess", "EmojiId")
}

func (es SqlEmojiAccessStore) Save(emoji_access *model.EmojiAccess) (*model.EmojiAccess, error) {
	if err := emoji_access.IsValid(); err != nil {
		return nil, err
	}

	if err := es.GetMaster().Insert(emoji_access); err != nil {
		return nil, errors.Wrap(err, "error saving emoji access")
	}

	return emoji_access, nil
}

func (es SqlEmojiAccessStore) GetByUserIdAndEmojiId(userId string, emojiId string) (*model.EmojiAccess, error) {
	var emoji_access *model.EmojiAccess

	err := es.GetReplica().SelectOne(&emoji_access,
		`SELECT
			*
		FROM
			EmojiAccess
		WHERE
			UserId = :Key1
			AND EmojiId = :Key2`, map[string]string{"Key1": userId, "Key2": emojiId})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("EmojiAccess", fmt.Sprintf("%s=%s", userId, emojiId))
		}

		return nil, errors.Wrapf(err, "could not get emoji access by %s with value %s", userId, emojiId)
	}

	return emoji_access, nil
}

func (es SqlEmojiAccessStore) GetMultipleByUserId(ids []string) ([]*model.EmojiAccess, error) {
	keys, params := MapStringsToQueryParams(ids, "EmojiAccess")

	var emojiAccesses []*model.EmojiAccess

	if _, err := es.GetReplica().Select(&emojiAccesses,
		`SELECT
			*
		FROM
			EmojiAccess
		WHERE
			UserId IN `+keys+`
			`, params); err != nil {
		return nil, errors.Wrapf(err, "error getting emoji access by user ids %v", ids)
	}
	return emojiAccesses, nil
}

func (es SqlEmojiAccessStore) DeleteAccessByUserIdAndEmojiId(userId string, emojiId string) error {

	sql := `DELETE
		FROM EmojiAccess
	WHERE
	UserId = :UserId
	AND EmojiId = :EmojiId`

	queryParams := map[string]string{
		"UserId":  userId,
		"EmojiId": emojiId,
	}
	_, err := es.GetMaster().Exec(sql, queryParams)
	if err != nil {
		return err
	}
	return nil
}

func (es SqlEmojiAccessStore) DeleteAccessByEmojiId(emojiId string) error {
	sql := `DELETE
		FROM EmojiAccess
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
