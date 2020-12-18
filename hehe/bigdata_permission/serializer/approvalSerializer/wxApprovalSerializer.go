package approvalSerializer

type WxMessageCallbackRequest struct {
	SelectedKey string `json:"selected_key"`
	Uid         int64  `json:"uid"`
	MsgID       int64  `json:"msg_id"`
	Extra       string `json:"extra"`
}