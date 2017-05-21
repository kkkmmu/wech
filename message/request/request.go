package message

// BaseRequest is a base for all wx api request.
type BaseRequest struct {
	XMLName xml.Name `xml:"error" json:"-"`

	Ret        int    `xml:"ret" json:"-"`
	Message    string `xml:"message" json:"-"`
	Wxsid      string `xml:"wxsid" json:"Sid"`
	Skey       string `xml:"skey"`
	DeviceID   string `xml:"-"`
	Wxuin      string `xml:"wxuin" json:"Uin"`
	PassTicket string `xml:"pass_ticket" json:"-"`
}

type CommonReqBody struct {
	BaseRequest        *BaseRequest
	Msg                interface{}
	SyncKey            *SyncKey
	rr                 int
	Code               int
	FromUserName       string
	ToUserName         string
	ClientMsgId        int64
	ClientMediaId      int
	TotalLen           string
	StartPos           int
	DataLen            string
	MediaType          int
	Scene              int
	Count              int
	List               []Member
	Opcode             int
	SceneList          []int
	SceneListCount     int
	VerifyContent      string
	VerifyUserList     []*VerifyUser
	VerifyUserListSize int
	skey               string
	MemberCount        int
	MemberList         []*Member
	Topic              string
}
