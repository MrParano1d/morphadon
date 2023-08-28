package web

import "github.com/marlaone/engine/core"

const (
	OpHttpAny core.Operation = 1 << iota
	OpHttpGet
	OpHttpPost
	OpHttpPut
	OpHttpDelete
	OpHttpPatch
	OpHttpHead
	OpHttpOptions
)
