package webview

import core "dappco.re/go"

func jsonMarshal(value any) ([]byte, resultFailure) {
	result := core.JSONMarshal(value)
	if !result.OK {
		if err, ok := result.Value.(resultFailure); ok {
			return nil, err
		}
		return nil, core.NewError(result.Error())
	}
	return result.Value.([]byte), nil
}
