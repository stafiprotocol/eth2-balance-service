// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

type Handler struct {
	db *db.WrapDb
}

func NewHandler(db *db.WrapDb) *Handler {
	return &Handler{db: db}
}
