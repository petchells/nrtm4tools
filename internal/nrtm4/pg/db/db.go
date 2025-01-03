/*
Package db is the database access tier.

Technical functions are in tx.go. Business functions are implemented in persist.go,
which needs a transaction to get the ball rolling. WithTransction should be
called from a top-level service tier, in this case the API.

Example:

var result Thing1

	err := db.WithTransaction(func(tx pgx.Tx) error {
		var err error
		persist := db.Dao{Tx: tx}
		// do things with the persist, set err if there is one...
		result = GetTheThing()
		// transaction will roll back if err is not nil or if a Panic happens
		return err
	})

	if err != nil {
		// transaction was rolled back
	}

// use the result...
*/
package db

import "github.com/petchells/nrtm4client/internal/nrtm4/util"

var logger = util.Logger
