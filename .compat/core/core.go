// SPDX-Licence-Identifier: EUPL-1.2

// Package core is a local compatibility bridge for sibling modules that have
// not yet moved their imports from dappco.re/go/core to dappco.re/go.
package core

import root "dappco.re/go"

type (
	Action                = root.Action
	ActionHandler         = root.ActionHandler
	AtomicPointer[T any]  = root.AtomicPointer[T]
	Context               = root.Context
	Core                  = root.Core
	CoreOption            = root.CoreOption
	Embed                 = root.Embed
	Fs                    = root.Fs
	Message               = root.Message
	Mutex                 = root.Mutex
	Once                  = root.Once
	Option                = root.Option
	Options               = root.Options
	Process               = root.Process
	Registry[T any]       = root.Registry[T]
	Result                = root.Result
	RWMutex               = root.RWMutex
	ServiceRuntime[T any] = root.ServiceRuntime[T]
	Startable             = root.Startable
	Stoppable             = root.Stoppable
	Translator            = root.Translator
)

var (
	As                  = root.As
	CleanPath           = root.CleanPath
	Concat              = root.Concat
	Contains            = root.Contains
	E                   = root.E
	Env                 = root.Env
	HasPrefix           = root.HasPrefix
	HasSuffix           = root.HasSuffix
	ID                  = root.ID
	Is                  = root.Is
	IsDigit             = root.IsDigit
	IsLetter            = root.IsLetter
	IsSpace             = root.IsSpace
	JSONMarshal         = root.JSONMarshal
	JSONMarshalString   = root.JSONMarshalString
	JSONUnmarshal       = root.JSONUnmarshal
	JSONUnmarshalString = root.JSONUnmarshalString
	Join                = root.Join
	Lower               = root.Lower
	New                 = root.New
	NewBuffer           = root.NewBuffer
	NewBuilder          = root.NewBuilder
	NewError            = root.NewError
	NewOptions          = root.NewOptions
	NewReader           = root.NewReader
	Path                = root.Path
	PathBase            = root.PathBase
	PathDir             = root.PathDir
	PathIsAbs           = root.PathIsAbs
	PathJoin            = root.PathJoin
	Print               = root.Print
	Println             = root.Println
	ReadAll             = root.ReadAll
	Replace             = root.Replace
	SHA256              = root.SHA256
	Security            = root.Security
	Split               = root.Split
	SplitN              = root.SplitN
	Sprintf             = root.Sprintf
	Trim                = root.Trim
	TrimPrefix          = root.TrimPrefix
	TrimSuffix          = root.TrimSuffix
	Upper               = root.Upper
	Warn                = root.Warn
	WithName            = root.WithName
)

func NewRegistry[T any]() *Registry[T] {
	return root.NewRegistry[T]()
}

func NewServiceRuntime[T any](c *Core, opts T) *ServiceRuntime[T] {
	return root.NewServiceRuntime(c, opts)
}
