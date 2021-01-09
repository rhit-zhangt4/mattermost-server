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

type SqlSecretStore struct {
	SqlStore
	metrics einterfaces.MetricsInterface
}

func newSqlSecretStore(sqlStore SqlStore, metrics einterfaces.MetricsInterface) store.SecretStore {
	s := &SqlSecretStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Secret{}, "Secret").SetKeys(false, "SecretName")
		table.ColMap("SecretName").SetMaxSize(26)
		table.ColMap("SecretValue").SetMaxSize(100)
	}

	return s
}

func (es SqlSecretStore) createIndexesIfNotExists() {
	es.CreateIndexIfNotExists("idx_secret", "Secret", "SecretName")
}

// func (es SqlSecretStore) Save(public_emoji *model.Secret) (*model.Secret, error) {
// 	// if err := public_emoji.IsValid(); err != nil {
// 	// 	return nil, err
// 	// }

// 	if err := es.GetMaster().Insert(secret); err != nil {
// 		return nil, errors.Wrap(err, "error saving secret")
// 	}

// 	return secret, nil
// }

func (es SqlSecretStore) GetBySecretName(secretName string) (*model.Secret, error) {
	var secret *model.Secret

	err := es.GetReplica().SelectOne(&secret,
		`SELECT *
			FROM Secret
		WHERE
			SecretName = :secretName`, map[string]interface{}{"secretName": secretName})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Secret", fmt.Sprintf("%s", secretName))
		}

		return nil, errors.Wrapf(err, "could not get secret by name %s", secretName)
	}

	return secret, nil
}
