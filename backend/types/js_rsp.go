package types

type Code int

const (
	CodeOK Code = iota
	CodeFailed
)

type JSRsp struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}
