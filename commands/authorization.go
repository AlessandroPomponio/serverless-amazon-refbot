// Copyright 2019 Alessandro Pomponio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commands

import (
	"github.com/pkg/errors"

	"github.com/AlessandroPomponio/serverless-amazon-refbot/persistence"
	"github.com/AlessandroPomponio/serverless-amazon-refbot/repository"
)

func authorizeUser(userID int) error {

	// We need to make sure the user is an admin.
	isAdmin, err := persistence.IsUserAdmin(userID, repository.DynamoDBClient)
	if err != nil {
		err = errors.Errorf("retrieveLatestRequest: error while checking user status: %s", err)
		return err
	}

	if !isAdmin {
		err = errors.Errorf("retrieveLatestRequest: user %d is not authorized to perform this action", userID)
		return err
	}

	return nil

}
