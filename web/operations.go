package web

import "github.com/mrparano1d/morphadon"

const (
	OpHttpAny morphadon.Operation = 1 << iota
	OpHttpGet
	OpHttpPost
	OpHttpPut
	OpHttpDelete
	OpHttpPatch
	OpHttpHead
	OpHttpOptions
)
