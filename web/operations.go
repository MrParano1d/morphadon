package web

import "github.com/mrparano1d/morphadon/core"

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
