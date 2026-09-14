package bundb

import (
	"context"
	"errors"

	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"github.com/uptrace/bun"
)

type configDB struct {
	db    *bun.DB
	state *state.State
}

func (c *configDB) RuntimeConfig(ctx context.Context) *gtsmodel.RuntimeConfig {
	// Check if config stored in
	// the cache. Load it if not.
	config := c.state.Caches.DB.RuntimeConfig.Load()
	if config != nil {
		// All good!
		return config
	}

	// Not hydrated.
	//
	// Try to load existing from db.
	config = new(gtsmodel.RuntimeConfig)
	err := c.db.
		NewSelect().
		Table("runtime_config").
		Limit(1).
		Scan(ctx, config)

	if errors.Is(err, db.ErrNoEntries) {
		// No config entry yet,
		// create one + return it.
		return c.newConfig(ctx)
	}

	if err != nil {
		// Something wrong with
		// the db, this is fatal!
		panic(err)
	}

	// No problem.
	//
	// Store loaded config in
	// the cache and return it.
	c.state.Caches.DB.RuntimeConfig.Store(config)
	return config
}

func (c *configDB) newConfig(ctx context.Context) *gtsmodel.RuntimeConfig {
	// Config doesn't exist yet,
	// store new one with defaults.
	config := new(gtsmodel.RuntimeConfig)
	_, err := c.db.
		NewInsert().
		Model(config).
		Exec(ctx)
	if errors.Is(err, db.ErrAlreadyExists) {
		// Someone beat us to creation of
		// new config, return that instead.
		return c.RuntimeConfig(ctx)
	}

	if err != nil {
		// Something wrong with
		// the db, this is fatal!
		panic(err)
	}

	// No problem.
	//
	// Store created config in
	// the cache and return it.
	c.state.Caches.DB.RuntimeConfig.Store(config)
	return config
}

func (c *configDB) UpdateRuntimeConfig(ctx context.Context, runtimeConfig *gtsmodel.RuntimeConfig, columns ...string) error {
	// Update config in the db.
	if _, err := c.db.
		NewUpdate().
		Model(runtimeConfig).
		Column(columns...).
		WherePK().
		Exec(ctx); err != nil {
		return err
	}

	// Cache updated model.
	c.state.Caches.DB.RuntimeConfig.Store(runtimeConfig)
	return nil
}
