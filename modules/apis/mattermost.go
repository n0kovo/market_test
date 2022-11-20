package apis

import (
	"github.com/n0kovo/market_test/modules/util"
)

type MattermostIncomingHookRequest struct {
	Text string `json:"text"`
}

func PostMattermostEvent(url, text string) {
	ir := MattermostIncomingHookRequest{Text: text}
	util.TorJSONPOST(url, ir)
}
