package contact

// Contact is wx Account struct
type Contact struct {
	GGID            string
	UserName        string
	NickName        string
	HeadImgURL      string `json:"HeadImgUrl"`
	HeadHash        string
	RemarkName      string
	DisplayName     string
	StarFriend      float64
	Sex             float64
	Signature       string
	VerifyFlag      float64
	ContactFlag     float64
	HeadImgFlag     float64
	Province        string
	City            string
	Alias           string
	EncryChatRoomID string `json:"EncryChatRoomId"`
	Type            int
	MemberList      []*Contact
}

type Member struct {
	Uin               int64  `json:"Uin"`
	UserName          string `json:"UserName"`
	NickName          string `json:"NickName"`
	HeadImgUrl        string `json:"HeadImgUrl"`
	ContactFlag       int    `json:"ContactFlag"`
	MemberCount       int    `json:"MemberCount"`
	MemberList        []User `json:"MemberList"`
	RemarkName        string `json:"RemarkName"`
	HideInputBarFlag  int    `json:"HideInputBarFlag"`
	Sex               int    `json:"Sex"`
	Signature         string `json:"Signature"`
	VerifyFlag        int    `json:"VerifyFlag"`
	OwnerUin          int    `json:"OwnerUin"`
	PYInitial         string `json:"PYInitial"`
	PYQuanPin         string `json:"PYQuanPin"`
	RemarkPYInitial   string `json:"RemarkPYInitial"`
	RemarkPYQuanPin   string `json:"RemarkPYQuanPin"`
	StarFriend        int    `json:"StarFriend"`
	AppAccountFlag    int    `json:"AppAccountFlag"`
	Statues           int    `json:"Statues"`
	AttrStatus        int    `json:"AttrStatus"`
	Province          string `json:"Province"`
	City              string `json:"City"`
	Alias             string `json:"Alias"`
	SnsFlag           int    `json:"SnsFlag"`
	UniFriend         int    `json:"UniFriend"`
	DisplayName       string `json:"DisplayName"`
	ChatRoomId        int    `json:"ChatRoomId"`
	KeyWord           string `json:"KeyWord"`
	EncryChatRoomId   string `json:"EncryChatRoomId"`
	IsOwner           int    `json:"IsOwner"`
	HeadImgUpdateFlag int    `json:"HeadImgUpdateFlag"`
	ContactType       int    `json:"ContactType"`
	ChatRoomOwner     string `json:"ChatRoomOwner"`
}
