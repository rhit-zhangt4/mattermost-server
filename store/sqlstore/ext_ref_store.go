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

func (es SqlExtRefStore) GetByAliasUserId(aliasUserId string) (*model.ExtRef, error) {
	var ext_ref *model.ExtRef

	err := es.GetReplica().SelectOne(&ext_ref,
		`SELECT
			*
		FROM
			ExtRef
		WHERE
			AliasUserId = :Key1
			`, map[string]string{"Key1": aliasUserId})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("ExtRef", fmt.Sprintf("aliasUserId=%s", aliasUserId))
		}

		return nil, errors.Wrapf(err, "could not get ext ref with aliasUserId=%s", aliasUserId)
	}

	return ext_ref, nil

}

func (es SqlExtRefStore) GetByExtIdAndPlatform(externalId string, externalPlatform string) (*model.ExtRef, error) {
	var ext_ref *model.ExtRef

	err := es.GetReplica().SelectOne(&ext_ref,
		`SELECT
			*
		FROM
			ExtRef
		WHERE
			ExternalId = :Key1
			AND ExternalPlatform = :Key2`, map[string]string{"Key1": externalId, "Key2": externalPlatform})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("ExtRef", fmt.Sprintf("externalId=%s, platform=%s", externalId, externalPlatform))
		}

		return nil, errors.Wrapf(err, "could not get ext ref with externalId=%s, platform=%s", externalId, externalPlatform)
	}

	return ext_ref, nil

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

func (es SqlExtRefStore) UpdateRealId(realUserId string, externalId string, externalPlatform string) error {
	if sqlResult, err := es.GetMaster().Exec(
		`UPDATE
			ExtRef
		SET
			RealUserId = :realUserId
		WHERE
			ExternalId = :externalId
			AND ExternalPlatform = :externalPlatform`, map[string]interface{}{"realUserId": realUserId, "externalId": externalId, "externalPlatform": externalPlatform}); err != nil {
		return errors.Wrap(err, "could not update realUserId")
	} else if rows, _ := sqlResult.RowsAffected(); rows == 0 {
		return store.NewErrNotFound("ExtRef", externalId)
	}
	return nil
}

func (es SqlExtRefStore) Unlink(realUserId string, externalPlatform string) error {
	sql := `DELETE
		FROM ExtRef
	WHERE
	RealUserId = :realUserId
	AND ExternalPlatform = :externalPlatform`

	queryParams := map[string]string{
		"realUserId":       realUserId,
		"externalPlatform": externalPlatform,
	}
	_, err := es.GetMaster().Exec(sql, queryParams)
	if err != nil {
		return err
	}
	return nil
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
