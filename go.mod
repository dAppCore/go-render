module dappco.re/go/html

go 1.26.0

require (
	dappco.re/go/core v0.8.0-alpha.1
	dappco.re/go/i18n v0.8.0-alpha.1
	dappco.re/go/io v0.8.0-alpha.1
	dappco.re/go/log v0.8.0-alpha.1
	dappco.re/go/process v0.8.0-alpha.1
	github.com/gin-gonic/gin v1.12.0
)

require (
	dappco.re/go/inference v0.8.0-alpha.1 // indirect
	golang.org/x/text v0.36.0 // indirect
)

replace (
	dappco.re/go/core => ../go
	dappco.re/go/i18n => ../go-i18n
	dappco.re/go/inference => ../go-inference
	dappco.re/go/io => ../go-io
	dappco.re/go/log => ../go-log
	dappco.re/go/process => ../go-process
)
