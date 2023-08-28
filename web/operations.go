package web

import "github.com/marlaone/morphadon/core"

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
