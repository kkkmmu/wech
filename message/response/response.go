package message

type MessageSyncResponse struct {
	BaseResponse           BaseResponse `json:"BaseResponse"`
	AddMsgCount            int          `json:"AddMsgCount"`
	AddMsgList             []Message    `json:"AddMsgList"`
	ModContactCount        int          `json:"ModContactCount"`
	ModContactList         []Member     `json:"ModContactList"`
	DelContactCount        int          `json:"DelContactCount"`
	DelContactList         []Member     `json:"DelContactList"`
	ModChatRoomMemberCount int          `json:"ModChatRoomMemberCount"`
	ModChatRoomMemberList  []Member     `json:"ModChatRoomMemberList"`
	Profile                Profile      `json:"Profile"`
	ContinueFlag           int          `json:"ContinueFlag"`
	SyncKey                SyncKey      `json:"SyncKey"`
	Skey                   string       `json:"Skey"`
	SyncCheckKey           SyncKey      `json:"SyncCheckKey"`
}

type GetContactListResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	MemberCount  int          `json:"MemberCount"`
	MemberList   []Member     `json:"MemberList"`
	Seq          float64      `json:"Seq"`
}

type UploadMediaResponse struct {
	BaseResponse      BaseResponse `json:"BaseResponse"`
	MediaID           string       `json:"MediaID"`
	StartPos          int          `json:"StartPos"`
	CDNThumbImgHeight int          `json:"CDNThumbImgHeight"`
	CDNThumbImgWidth  int          `json:"CDNThumbImgWidth"`
}

type Response struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
}

type BaseResponse struct {
	Ret    int
	ErrMsg string
}

type InitResponse struct {
	BaseResponse        BaseResponse         `json:"BaseResponse"`
	Count               int                  `json:"Count"`
	User                User                 `json:"User"` //This is ourself
	ContactList         []Member             `json:"ContactList"`
	SyncKey             SyncKey              `json:"SyncKey"`
	ChatSet             string               `json:"ChatSet"`
	SKey                string               `json:"SKey"`
	ClientVersion       int                  `json:"ClientVersion"`
	SystemTime          int                  `json:"SystemTime"`
	GrayScale           int                  `json:"GrayScale"`
	InviteStartCount    int                  `json:"InviteStartCount"`
	MPSubscribeMsgCount int                  `json:"MPSubscribeMsgCount"`
	MPSubscribeMsgList  []MPSubscribeMsgList `json:"MPSubscribeMsgList"`
	ClickReportInterval int                  `json:"ClickReportInterval"`
}
