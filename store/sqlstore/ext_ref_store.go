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

type SqlExtRefStore struct {
	SqlStore
	metrics einterfaces.MetricsInterface
}

func newSqlExtRefStore(sqlStore SqlStore, metrics einterfaces.MetricsInterface) store.ExtRefStore {
	s := &SqlExtRefStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.ExtRef{}, "ExtRef").SetKeys(false, "ExternalId", "ExternalPlatform")
		table.ColMap("ExternalId").SetMaxSize(64)
		table.ColMap("ExternalPlatform").SetMaxSize(15)
		table.ColMap("RealUserId").SetMaxSize(64)
		table.ColMap("AliasUserId").SetMaxSize(64)
	}

	return s
}

func (es SqlExtRefStore) createIndexesIfNotExists() {
	es.CreateIndexIfNotExists("idx_ext_ref_real_user", "ExtRef", "RealUserId")
	es.CreateIndexIfNotExists("idx_ext_ref_alias", "ExtRef", "AliasUserId")
	es.CreateIndexIfNotExists("idx_ext_ref_external", "ExtRef", "ExternalId")
}

func (es SqlExtRefStore) GetByRealUserIdAndPlatform(realUserId string, externalPlatform string) (*model.ExtRef, error) {
	var ext_ref *model.ExtRef

	err := es.GetReplica().SelectOne(&ext_ref,
		`SELECT
			*
		FROM
			ExtRef
		WHERE
			RealUserId = :Key1
			AND ExternalPlatform = :Key2`, map[string]string{"Key1": realUserId, "Key2": externalPlatform})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("ExtRef", fmt.Sprintf("realUserId=%s, platform=%s", realUserId, externalPlatform))
		}

		return nil, errors.Wrapf(err, "could not get ext ref with realUserId=%s, platform=%s", realUserId, externalPlatform)
	}

	return ext_ref, nil
}

func (es SqlExtRefStore) Save(ext_ref *model.ExtRef) (*model.ExtRef, error) {
	// if err := ext_ref.IsValid(); err != nil {
	// 	return nil, err
	// }

	if err := es.GetMaster().Insert(ext_ref); err != nil {
		return nil, errors.Wrap(err, "error saving ext ref")
	}

	return ext_ref, nil
}

// func (es SqlEmojiAccessStore) GetMultipleByUserId(ids []string) ([]*model.EmojiAccess, error) {
// 	keys, params := MapStringsToQueryParams(ids, "EmojiAccess")

// 	var emojiAccesses []*model.EmojiAccess

// 	if _, err := es.GetReplica().Select(&emojiAccesses,
// 		`SELECT
// 			*
// 		FROM
// 			EmojiAccess
// 		WHERE
// 			UserId IN `+keys+`
// 			`, params); err != nil {
// 		return nil, errors.Wrapf(err, "error getting emoji access by user ids %v", ids)
// 	}
// 	return emojiAccesses, nil
// }

// func (es SqlEmojiAccessStore) DeleteAccessByUserIdAndEmojiId(userId string, emojiId string) error {

// 	sql := `DELETE
// 		FROM EmojiAccess
// 	WHERE
// 	UserId = :UserId
// 	AND EmojiId = :EmojiId`

// 	queryParams := map[string]string{
// 		"UserId":  userId,
// 		"EmojiId": emojiId,
// 	}
// 	_, err := es.GetMaster().Exec(sql, queryParams)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (es SqlEmojiAccessStore) DeleteAccessByEmojiId(emojiId string) error {
// 	sql := `DELETE
// 		FROM EmojiAccess
// 	WHERE
// 		EmojiId = :EmojiId`

// 	queryParams := map[string]string{
// 		"EmojiId": emojiId,
// 	}

// 	_, err := es.GetMaster().Exec(sql, queryParams)
// 	if err != nil {
// 		//mlog.Warn("Failed to delete access", mlog.Err(err))
// 		return err
// 	}
// 	return nil
// }
