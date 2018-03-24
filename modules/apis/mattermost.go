package apis

import (
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

type MattermostIncomingHookRequest struct {
	Text string `json:"text"`
}

func PostMattermostEvent(url, text string) {
	ir := MattermostIncomingHookRequest{Text: text}
	util.TorJSONPOST(url, ir)
}
