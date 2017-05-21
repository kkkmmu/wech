package message

type RecommendInfo struct {
	UserName   string `json:"UserName"`
	NickName   string `json:"NickName"`
	QQNum      int    `json:"QQNum"`
	Province   string `json:"Province"`
	City       string `json:"City"`
	Content    string `json:"Content"`
	Signature  string `json:"Signature"`
	Alias      string `json:"Alias"`
	Scene      int    `json:"Scene"`
	VerifyFlag int    `json:"VerifyFlag"`
	AttrStatus int    `json:"AttrStatus"`
	Sex        int    `json:"Sex"`
	Ticket     string `json:"Ticket"`
	OpCode     int    `json:"OpCode"`
}
