module dappco.re/go/core/html

go 1.26.0

require (
	dappco.re/go/core/i18n v0.1.8
	dappco.re/go/core/io v0.2.0
	dappco.re/go/core/log v0.1.0
	github.com/stretchr/testify v1.11.1
)

require (
	dappco.re/go/core v0.5.0 // indirect
	forge.lthn.ai/core/go-inference v0.1.4 // indirect
	forge.lthn.ai/core/go-log v0.0.4 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/text v0.35.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	dappco.re/go/core => ../../../../core/go
	dappco.re/go/core/i18n => ../../../../core/go-i18n
	dappco.re/go/core/io => ../../../../core/go-io
	dappco.re/go/core/log => ../../../../core/go-log
)
