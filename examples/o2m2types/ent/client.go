// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/examples/o2m2types/ent/migrate"

	"github.com/facebookincubator/ent/examples/o2m2types/ent/pet"
	"github.com/facebookincubator/ent/examples/o2m2types/ent/user"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Pet is the client for interacting with the Pet builders.
	Pet *PetClient
	// User is the client for interacting with the User builders.
	User *UserClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	c := config{log: log.Println}
	c.options(opts...)
	return &Client{
		config: c,
		Schema: migrate.NewSchema(c.driver),
		Pet:    NewPetClient(c),
		User:   NewUserClient(c),
	}
}

// Open opens a connection to the database specified by the driver name and a
// driver-specific data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %v", err)
	}
	cfg := config{driver: tx, log: c.log, debug: c.debug}
	return &Tx{
		config: cfg,
		Pet:    NewPetClient(cfg),
		User:   NewUserClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Pet.
//		Query().
//		Count(ctx)
//
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true}
	return &Client{
		config: cfg,
		Schema: migrate.NewSchema(cfg.driver),
		Pet:    NewPetClient(cfg),
		User:   NewUserClient(cfg),
	}
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...ent.Hook) {
	c.Pet.Use(hooks...)
	c.User.Use(hooks...)
}

// PetClient is a client for the Pet schema.
type PetClient struct {
	config
}

// NewPetClient returns a client for the Pet from the given config.
func NewPetClient(c config) *PetClient {
	return &PetClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `pet.Hooks(f(g(h())))`.
func (c *PetClient) Use(hooks ...ent.Hook) {
	c.hooks = append(c.hooks[:len(c.hooks):len(c.hooks)], hooks...)
}

// Create returns a create builder for Pet.
func (c *PetClient) Create() *PetCreate {
	hooks := c.hooks
	mutation := newPetMutation(OpCreate)
	return &PetCreate{config: c.config, hooks: hooks, mutation: mutation}
}

// Update returns an update builder for Pet.
func (c *PetClient) Update() *PetUpdate {
	hooks := c.hooks
	mutation := newPetMutation(OpUpdate)
	return &PetUpdate{config: c.config, hooks: hooks, mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PetClient) UpdateOne(pe *Pet) *PetUpdateOne {
	return c.UpdateOneID(pe.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *PetClient) UpdateOneID(id int) *PetUpdateOne {
	hooks := c.hooks
	mutation := newPetMutation(OpUpdateOne)
	mutation.id = &id
	return &PetUpdateOne{config: c.config, hooks: hooks, mutation: mutation}
}

// Delete returns a delete builder for Pet.
func (c *PetClient) Delete() *PetDelete {
	hooks := c.hooks
	mutation := newPetMutation(OpDelete)
	return &PetDelete{config: c.config, hooks: hooks, mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *PetClient) DeleteOne(pe *Pet) *PetDeleteOne {
	return c.DeleteOneID(pe.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *PetClient) DeleteOneID(id int) *PetDeleteOne {
	builder := c.Delete().Where(pet.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PetDeleteOne{builder}
}

// Create returns a query builder for Pet.
func (c *PetClient) Query() *PetQuery {
	return &PetQuery{config: c.config}
}

// Get returns a Pet entity by its id.
func (c *PetClient) Get(ctx context.Context, id int) (*Pet, error) {
	return c.Query().Where(pet.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PetClient) GetX(ctx context.Context, id int) *Pet {
	pe, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return pe
}

// QueryOwner queries the owner edge of a Pet.
func (c *PetClient) QueryOwner(pe *Pet) *UserQuery {
	query := &UserQuery{config: c.config}
	id := pe.ID
	step := sqlgraph.NewStep(
		sqlgraph.From(pet.Table, pet.FieldID, id),
		sqlgraph.To(user.Table, user.FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, pet.OwnerTable, pet.OwnerColumn),
	)
	query.sql = sqlgraph.Neighbors(pe.driver.Dialect(), step)

	return query
}

// UserClient is a client for the User schema.
type UserClient struct {
	config
}

// NewUserClient returns a client for the User from the given config.
func NewUserClient(c config) *UserClient {
	return &UserClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `user.Hooks(f(g(h())))`.
func (c *UserClient) Use(hooks ...ent.Hook) {
	c.hooks = append(c.hooks[:len(c.hooks):len(c.hooks)], hooks...)
}

// Create returns a create builder for User.
func (c *UserClient) Create() *UserCreate {
	hooks := c.hooks
	mutation := newUserMutation(OpCreate)
	return &UserCreate{config: c.config, hooks: hooks, mutation: mutation}
}

// Update returns an update builder for User.
func (c *UserClient) Update() *UserUpdate {
	hooks := c.hooks
	mutation := newUserMutation(OpUpdate)
	return &UserUpdate{config: c.config, hooks: hooks, mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserClient) UpdateOne(u *User) *UserUpdateOne {
	return c.UpdateOneID(u.ID)
}

// UpdateOneID returns an update builder for the given id.
func (c *UserClient) UpdateOneID(id int) *UserUpdateOne {
	hooks := c.hooks
	mutation := newUserMutation(OpUpdateOne)
	mutation.id = &id
	return &UserUpdateOne{config: c.config, hooks: hooks, mutation: mutation}
}

// Delete returns a delete builder for User.
func (c *UserClient) Delete() *UserDelete {
	hooks := c.hooks
	mutation := newUserMutation(OpDelete)
	return &UserDelete{config: c.config, hooks: hooks, mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *UserClient) DeleteOne(u *User) *UserDeleteOne {
	return c.DeleteOneID(u.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *UserClient) DeleteOneID(id int) *UserDeleteOne {
	builder := c.Delete().Where(user.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UserDeleteOne{builder}
}

// Create returns a query builder for User.
func (c *UserClient) Query() *UserQuery {
	return &UserQuery{config: c.config}
}

// Get returns a User entity by its id.
func (c *UserClient) Get(ctx context.Context, id int) (*User, error) {
	return c.Query().Where(user.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserClient) GetX(ctx context.Context, id int) *User {
	u, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return u
}

// QueryPets queries the pets edge of a User.
func (c *UserClient) QueryPets(u *User) *PetQuery {
	query := &PetQuery{config: c.config}
	id := u.ID
	step := sqlgraph.NewStep(
		sqlgraph.From(user.Table, user.FieldID, id),
		sqlgraph.To(pet.Table, pet.FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, user.PetsTable, user.PetsColumn),
	)
	query.sql = sqlgraph.Neighbors(u.driver.Dialect(), step)

	return query
}
