// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

type Handler struct {
	db              *db.WrapDb
	slashStartEpoch uint64
	isDev           bool
}

func NewHandler(db *db.WrapDb, isDev bool, slashStartEpoch uint64) *Handler {
	return &Handler{db: db, isDev: isDev, slashStartEpoch: slashStartEpoch}
}
