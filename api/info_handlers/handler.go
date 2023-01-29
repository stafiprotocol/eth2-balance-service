// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

type Handler struct {
	db                      *db.WrapDb
	unstakingStartTimestamp uint64
}

func NewHandler(db *db.WrapDb, unstakingStartTimestamp uint64) *Handler {
	return &Handler{db: db, unstakingStartTimestamp: unstakingStartTimestamp}
}
