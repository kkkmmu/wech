package main

import (
	"crypto/tls"
	"errors"
	"github.com/mdp/qrterminal"
	//qrcode "github.com/skip2/go-qrcode"
	"encoding/json"
	"encoding/xml"
	//"io"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	//"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//- 微信目前分为2个版本，所以在获取接口时候请求的路径也不一样，
//	很早以前注册的用户请求地址一般为wx.qq.com，
//	新注册用户为wx2.qq.com，
//  导致很多开发者在开发微信网页版的时候返现有些用户能登录并获取到消息，有的只能登录不能获取到消息
//  微信返回码RetCode和相应的解决方案：0－正常；1－失败，refresh；1101/1100－登出/失败，refresh/重新登录；1203－恭喜您，几个小时后重试，没有解决方案；Selector：2－新消息，6/7－进入/离开聊天界面通常是在手机上进行操作，重新初始化即可，0－正常
type wechat struct {
	client      *http.Client
	cacheName   string
	userAgent   string
	uuid        string
	Redirect    string
	UUID        string
	Scan        string
	login       chan bool
	Ticket      string            `json:"Ticket"`
	Lang        string            `json:"Lang"`
	BaseURL     string            `json:"BaseURL"`
	BaseRequest *BaseRequest      `json:"BaseRequest"`
	DeviceID    string            `json:"DeviceID"`
	GroupList   []Member          `json:"GroupList"`
	SyncServer  string            `json:"SyncServer"`
	SyncKey     *SyncKey          `json:"SyncKey"`
	UserName    string            `json:"UserName"`
	NickName    string            `json:"NickName"`
	ContactDB   map[string]Member `json:"ContactDB"`
	Cookies     []*http.Cookie    `json:"Cookies"`
}

type Cache struct {
	BaseResponse *BaseResponse  `json:"BaseResponse"`
	DeviceID     string         `json:"DeviceID"`
	SyncKey      *SyncKey       `json:"SyncKey"`
	UserName     string         `json:"UserName"`
	NickName     string         `json:"NickName"`
	SyncServer   string         `json:"SyncServer"`
	Ticket       string         `json:"Ticket"`
	Lang         string         `json:"Lang"`
	Cookies      []*http.Cookie `json:"Cookies"`
}

type StatusNotifyResponse struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	MsgId        string        `json:"MsgId"`
}

type User struct {
	UserName          string `json:"UserName"`
	Uin               int64  `json:"Uin"`
	NickName          string `json:"NickName"`
	HeadImgUrl        string `json:"HeadImgUrl" xml:""`
	RemarkName        string `json:"RemarkName" xml:""`
	PYInitial         string `json:"PYInitial" xml:""`
	PYQuanPin         string `json:"PYQuanPin" xml:""`
	RemarkPYInitial   string `json:"RemarkPYInitial" xml:""`
	RemarkPYQuanPin   string `json:"RemarkPYQuanPin" xml:""`
	HideInputBarFlag  int    `json:"HideInputBarFlag" xml:""`
	StarFriend        int    `json:"StarFriend" xml:""`
	Sex               int    `json:"Sex" xml:""`
	Signature         string `json:"Signature" xml:""`
	AppAccountFlag    int    `json:"AppAccountFlag" xml:""`
	VerifyFlag        int    `json:"VerifyFlag" xml:""`
	ContactFlag       int    `json:"ContactFlag" xml:""`
	WebWxPluginSwitch int    `json:"WebWxPluginSwitch" xml:""`
	HeadImgFlag       int    `json:"HeadImgFlag" xml:""`
	SnsFlag           int    `json:"SnsFlag" xml:""`
}

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

type CommonReqBody struct {
	BaseRequest        *BaseRequest
	Msg                interface{}
	SyncKey            *SyncKey
	rr                 int
	Code               int
	FromUserName       string
	ToUserName         string
	ClientMsgId        int
	ClientMediaId      int
	TotalLen           int
	StartPos           int
	DataLen            int
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

/*
{
	"MsgId": "7318483579373924965",
	"FromUserName": "@e24439096308969756e667d06f33a50e",
	"ToUserName": "@b18d3f16a138505fd2ef663815925561948e2f970a910bba56396a5a62e7bf30",
	"MsgType": 1,
	"Content": "睡觉睡觉就是计算机计算机三级到你家都觉得那女的",
	"Status": 3,
	"ImgStatus": 1,
	"CreateTime": 1494261110,
	"VoiceLength": 0,
	"PlayLength": 0,
	"FileName": "",
	"FileSize": "",
	"MediaId": "",
	"Url": "",
	"AppMsgType": 0,
	"StatusNotifyCode": 0,
	"StatusNotifyUserName": "",
	"RecommendInfo": {
	    "UserName": "",
	    "NickName": "",
	    "QQNum": 0,
	    "Province": "",
	    "City": "",
	    "Content": "",
	    "Signature": "",
	    "Alias": "",
	    "Scene": 0,
	    "VerifyFlag": 0,
	    "AttrStatus": 0,
	    "Sex": 0,
	    "Ticket": "",
	    "OpCode": 0
	}
	,
	"ForwardFlag": 0,
	"AppInfo": {
	    "AppID": "",
	    "Type": 0
	}
	,
	"HasProductId": 0,
	"Ticket": "",
	"ImgHeight": 0,
	"ImgWidth": 0,
	"SubMsgType": 0,
	"NewMsgId": 7318483579373924965,
	"OriContent": ""
    }
*/
type Message struct {
	MsgId                string        `json:"MsgId"`
	FromUserName         string        `json:"FromUserName"`
	ToUserName           string        `json:"ToUserName"`
	MsgType              int           `json:"MsgType"`
	Content              string        `json:"Content"`
	Status               int           `json:"Status"`
	ImgStatus            int           `json:"ImgStatus"`
	CreateTime           int           `json:"CreateTime"`
	VoiceLength          int           `json:"VoiceLength"`
	PlayLength           int           `json:"PlayLength"`
	FileName             string        `json:"FileName"`
	FileSize             string        `json:"FileSize"`
	MediaId              string        `json:"MediaId"`
	Url                  string        `json:"Url"`
	AppMsgType           int           `json:"AppMsgType"`
	StatusNotifyCode     int           `json:"StatusNotifyCode"`
	StatusNotifyUserName string        `json:"StatusNotifyUserName"`
	RecommendInfo        RecommendInfo `json:"RecommendInfo"`
	ForwardFlag          int           `json:"ForwardFlag"`
	AppInfo              AppInfo       `json:"AppInfo"`
	HasProductId         int           `json:"HasProductId"`
	Ticket               string        `json:"Ticket"`
	ImgHeight            int           `json:"ImgHeight"`
	ImgWidth             int           `json:"ImgWidth"`
	SubMsgType           int           `json:"SubMsgType"`
	NewMsgId             int64         `json:"NewMsgId"`
	OriContent           string        `json:"OriContent"`
}

type TextMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
}

// MediaMessage
type MediaMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
	MediaId      string
}

// EmotionMessage: gif/emoji message struct
type EmotionMessage struct {
	ClientMsgId  int
	EmojiFlag    int
	FromUserName string
	LocalID      int
	MediaId      string
	ToUserName   string
	Type         int
}

type SendMessageResponse struct {
	BaseResponse BaseResponse
	MsgID        string `json:"MsgId"`
	LocalID      int    `json:"LocalID"`
}

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
type AppInfo struct {
	AppID string `json:"AppID"`
	Type  int    `json:"Type"`
}

/* I think this should be put into Member ---- Contact struct
type ModifiedContact struct {
	UserName          string
	NickName          string
	Sex               string
	HeadImgUpdateFlag int
	ContactType       int
	Alias             string
	ChatRoomOwner     string
	HeadImgUrl        string
	ContactFlag       int
	MemberCount       int
	MemberList        []Member
	HideInputBarFlag  int
	Signature         string
	VerifyFlag        int
	RemarkName        string
	Statues           int
	AttrStatus        int
	Province          string
	City              string
	SnsFlag           string
	KeyWor            string
}
*/

type UserName struct {
	Buff string `json:"Buff"`
}

type NickName struct {
	Buff string `json:"Buff"`
}

type BindEmail struct {
	Buff string `json:"Buff"`
}

type BindMobile struct {
	Buff string `json:"Buff"`
}

type Profile struct {
	BitFlag           int        `json:"BitFlag"`
	UserName          UserName   `json:"UserName"`
	NickName          NickName   `json:"NickName"`
	BindEmail         BindEmail  `json:"BindEmail"`
	BindMobile        BindMobile `json:"BindMobile"`
	Status            int        `json:"Status"`
	Sex               int        `json:"Sex"`
	PersonalCard      int        `json:"PersonalCard"`
	Alias             string     `json:"Alias"`
	HeadImgUpdateFlag int        `json:"HeadImgUpdateFlag"`
	HeadImgUrl        string     `json:"HeadImgUrl"`
	Signature         string     `json:"Signature"`
}

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

// GroupContactResponse: get batch contact response struct
type GetGroupMemberListResponse struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	Count        int           `json:"Count"`
	ContactList  []Member      `json:"ContactList"`
}

// VerifyUser: verify user request body struct
type VerifyUser struct {
	Value            string `json:"Value"`
	VerifyUserTicket string `json:"VerifyUserTicket"`
}

// ReceivedMessage: for received message
type ReceivedMessage struct {
	IsGroup      bool   `json:"IsGroup"`
	MsgId        string `json:"MsgId"`
	Content      string `json:"Content"`
	FromUserName string `json:"FromUserName"`
	ToUserName   string `json:"ToUserName"`
	Who          string `json:"Who"`
	MsgType      int    `json:"MsgType"`
}
type GetContactListResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	MemberCount  int          `json:"MemberCount"`
	MemberList   []Member     `json:"MemberList"`
	Seq          float64      `json:"Seq"`
}

type Response struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
}

type MsgResp struct {
	Response
}

type BaseResponse struct {
	Ret    int
	ErrMsg string
}

type SyncKey struct {
	Count int      `json:"Count"`
	List  []KeyVal `json:"List"`
}

func (sk *SyncKey) String() string {
	keys := make([]string, 0)
	for _, v := range sk.List {
		keys = append(keys, strconv.Itoa(v.Key)+"_"+strconv.Itoa(v.Val))
	}
	return strings.Join(keys, "|")
}

type KeyVal struct {
	Key int `json:"Key"`
	Val int `json:"Val"`
}

type GroupRequest struct {
	UserName        string
	EncryChatRoomId string
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

type MPArticle struct {
	Titile string `json:"Titile"`
	Digest string `json:"Digest"`
	Cover  string `json:"Cover"`
	Url    string `json:"Url"`
}

type MPSubscribeMsgList struct {
	UserName       string      `json:"UserName"`
	MPArticleCount int         `json:"MPArticleCount"`
	MPArticleList  []MPArticle `json:"MPArticleList"`
	Time           int         `json:"Time"`
	NickName       string      `json:"NickName"`
}

type initBaseRequest struct {
	BaseRequest *BaseRequest
}

type initBaseResp struct {
	Response
	User    Contact
	Skey    string
	SyncKey map[string]interface{}
}

var FetchTicket = regexp.MustCompile(`ticket=(?P<ticket>[[:word:]_\$@#\?\-=]+)&uuid=(?P<uuid>[[:word:]_\$@#\?\-=]+)&lang=(?P<lang>[[:word:]_\$@#\?\-=]+)&scan=(?P<lang>[[:word:]_\$@#\?\-=]+)`)

var FetchRedirectLink = regexp.MustCompile(`window.redirect_uri=\"(?P<redirect>[[:word:]=_\.\?@#&\-+%/\:]+)\"`)

func (wc *wechat) FastLogin() error {
	_, err := os.Stat(wc.cacheName)
	if os.IsNotExist(err) {
		log.Println("Cache not exist, log from scratch")
		return errors.New("Cache not exist")
	}

	data, err := ioutil.ReadFile(wc.cacheName)
	if err != nil {
		log.Println("Cannot read cache file: ", err.Error())
		return errors.New("Cannot read cache file")
	}

	err = json.Unmarshal(data, wc)
	if err != nil {
		log.Println("Cannot decode cache info: ", err.Error())
		return errors.New("Cannot decode cache")
	}

	u, ue := url.Parse(wc.BaseURL)
	if ue != nil {
		return errors.New("Cannot parse base url")
	}
	log.Println(wc)
	wc.client.Jar.SetCookies(u, wc.Cookies)
	return nil
}
func (wc *wechat) GetUUID() error {

	//wxeb7ec651dd0aefa9
	//wx782c26e4c19acffb //Wechat web version
	params := url.Values{}
	params.Set("appid", "wx782c26e4c19acffb")
	params.Set("fun", "new")
	params.Set("lang", "zh_CN")
	params.Set("_", strconv.FormatInt(time.Now().Unix(), 10))

	/*
		//This Appid is defined by Tencent
		req.Header.Add("appid", "wx782c26e4c19acffb")
		req.Header.Add("fun", "new")
		req.Header.Add("lang", "zh_CN")
		//Time Stamp
		req.Header.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	*/

	//This is how to do http post
	resp, err := wc.client.PostForm("https://login.weixin.qq.com/jslogin", params)
	if err != nil {
		log.Println("Erro happened when do http Request: ", err.Error())
		return errors.New("Post error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error happened when get respose body")
		return errors.New("Get body error")
	}
	defer resp.Body.Close()

	log.Println(string(body))
	log.Println(resp.Header)
	log.Println(resp.Status)
	log.Println(resp.StatusCode)
	log.Println(resp.Proto)
	log.Println(resp.Location())
	log.Println(resp.Cookies())

	//Fetch the UUID from response body
	re := regexp.MustCompile(`\"(?P<uuid>[[:word:]%=\-_#]+)\"`)
	match := re.FindStringSubmatch(string(body))
	log.Println(match[1])
	wc.uuid = match[1]
	return nil
}

//We can also use the webxwpushloginurl API to login after we get the UIN.

//If you want to display QR on terminal, use this link
//"https://login.weixin.qq.com/l/"+uuid,

//If you want to save QR code locally, use this link, USE POST method and with this params {"t": "webwx", "_": strconv.FormatInt(time.Now().Unix(), 10)}, Also remember to set the http Header.
//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//req.Header.Set("Cache-Control", "no-cache")
//"https://login.weixin.qq.com/qrcode/" + uuid

/*
<br>

| API | 绑定登陆（webwxpushloginurl） |
| --- | --------- |
| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxpushloginurl |
| method | GET |
| params | **uin**: xxx |

返回数据(String):
```
{'msg': 'all ok', 'uuid': 'xxx', 'ret': '0'}

通过这种方式可以省掉扫二维码这步操作，更加方便
```
<br>

| API | 生成二维码 |
| --- | --------- |
| url | https://login.weixin.qq.com/l/ `uuid` |
| method | GET |
<br>
*/

func (wc *wechat) GetQRCode() error {
	/*
		qrURL := `https://login.weixin.qq.com/qrcode/` + uuid
		params := url.Values{}
		params.Set("t", "webwx")
		params.Set("_", strconv.FormatInt(time.Now().Unix(), 10))

		req, err := http.NewRequest("POST", qrURL, strings.NewReader(params.Encode()))
		if err != nil {
			return ``, err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cache-Control", "no-cache")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return ``, err
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ``, err
		}

		path := `qrcode.png`
		if err = createFile(path, data, false); err != nil {
			return ``, err
		}
	*/

	//If you want to display QR on terminal, use this link
	//"https://login.weixin.qq.com/l/"+uuid,
	qrterminal.Generate("https://login.weixin.qq.com/l/"+wc.uuid, qrterminal.L, os.Stdout)
	return nil
}

//We use Get method with the following parameters to get the QR code san status.
//"tip: scanned 0, unscanned 1
//"uuid: our uuid"
//Time stamp
/*
| API | 二维码扫描登录 |
| --- | --------- |
| url | https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login |
| method | GET |
| params | **tip**: 1 `未扫描` 0 `已扫描` <br> **uuid**: xxx <br> **_**: `时间戳` |

*/
func (wc *wechat) WaitForQRCodeScan() {
	tick := time.Tick(time.Second * 5)
	params := url.Values{}
	params.Add("tip", "1")
	params.Add("uuid", wc.uuid)
	params.Add("_", strconv.FormatInt(time.Now().Unix(), 10))

	//It seems that both these two method can work
	uri := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?" + params.Encode()
	//uri := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?tip=" + "1" + "&uuid=" + wc.uuid + "&_=" + strconv.FormatInt(time.Now().Unix(), 10)
	log.Println(uri)
	for {
		<-tick
		resp, err := wc.client.Get(uri)
		if err != nil {
			log.Println("Get QRCode status scan failed: ", err.Error())
			continue
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		strb := string(body)
		codeRe := regexp.MustCompile(`window.code=(?P<state>[[:word:]]+)`)
		// In the response of this API:
		// window.code ==:
		//	408 Login Expired.
		//	201 Scan successed.
		//	200 Login Confirmed. When return code is 200, we also can get the redirect link.
		//	We need to get the "ticket" , "uuid", "lang" and "scan" from the redirect link for furture use.
		match := codeRe.FindStringSubmatch(strb)
		log.Println("Get QR Code Scan Status: ", match[1])

		if len(match) == 0 {
			log.Println("Cannot get the QR scan state code")
			continue
		}

		code := match[1]

		if code == "200" {
			log.Println(strb)
			redirect := FetchRedirectLink.FindStringSubmatch(strb)
			log.Println(redirect)
			if len(redirect) == 0 {
				continue
			}

			wc.Redirect = redirect[1]

			urls, err := url.Parse(wc.Redirect)
			if err != nil {
				log.Println("Invalide redirect link: ", err.Error())
				continue
			}
			wc.BaseURL = urls.Scheme + "://" + urls.Host + "/"

			fields := FetchTicket.FindStringSubmatch(strings.Split(wc.Redirect, "?")[1])
			if len(fields) != 5 {
				log.Println("Error happened when get fields")
				continue
			}

			wc.Ticket = fields[1]
			wc.UUID = fields[2]
			wc.Lang = fields[3]
			wc.Scan = fields[4]
			log.Println(wc)
			wc.login <- true
			break
		} else {
			if code == "201" {
				log.Println("Please confirm login on your telephone")
			} else {
				log.Println("Please Scan the QR Code.")
			}
			continue
		}
	}
}

//We can get the ticket,uuid,lang, scan from the REDIRECT link.
//After that we add "fun=new" to the parameter list to GET the
//redirect link. The return value is and XML which contains the
//Base information which is necessary for future use.
/*
| API | webwxnewloginpage |
| --- | --------- |
| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage |
| method | GET |
| params | **ticket**: xxx <br> **uuid**: xxx <br> **lang**: zh_CN `语言` <br> **scan**: xxx <br> **fun**: new |

<error>
	<ret>0</ret>
	<message>OK</message>
	<skey>xxx</skey>
	<wxsid>xxx</wxsid>
	<wxuin>xxx</wxuin>
	<pass_ticket>xxx</pass_ticket>
	<isgrayscale>1</isgrayscale>
</error>
*/
func (wc *wechat) GetBaseRequest() {
	resp, err := wc.client.Get(wc.Redirect + "&fun=new")
	if err != nil {
		log.Println("Cannot get redirect link: ", err.Error())
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error happened when read body: ", err.Error())
		return
	}

	wc.SaveToFile("BaseRequestBody.json", data)
	log.Println(string(data))
	reader := bytes.NewReader(data)
	if err = xml.NewDecoder(reader).Decode(wc.BaseRequest); err != nil {
		log.Println("Error happend when decode xml: ", err.Error())
		return
	}

	if wc.BaseRequest.Ret != 0 { // 0 is success
		log.Println("login failed message: ", wc.BaseRequest.Message)
		return
	}

	br, err := json.Marshal(wc.BaseRequest)
	if err != nil {
		log.Println("Unable to encode BaseRequest: ", err.Error())
		return
	}
	wc.SaveToFile("BaseRequest.json", br)

	log.Printf("%q", wc.BaseRequest)
	log.Println("++++++++++++++++++++++++++++++++++++++++++")
	log.Println(wc.BaseRequest.Message)
	log.Println(wc.BaseRequest.DeviceID)
	log.Println(wc.BaseRequest.PassTicket)
	log.Println(wc.BaseRequest.Skey)
	log.Println(wc.BaseRequest.Wxsid)
	log.Println(wc.BaseRequest.Wxuin)
	log.Println("++++++++++++++++++++++++++++++++++++++++++")
	log.Println("Login successfully")

	return
}

// Wechat init
//| API | webwxinit |
//| --- | --------- |
//| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxinit?pass_ticket=xxx&skey=xxx&r=xxx |
//| method | POST |
//| data | JSON |
//| header | ContentType: application/json; charset=UTF-8 |
//| params | { BaseRequest: { Uin: xxx, Sid: xxx, Skey: xxx, DeviceID: xxx} } |
//注意这里目前测得的结果是"https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit?pass_ticket=xxx&skey=xxx&r=xxx "
//同时注意BaseURL 并不是根域名，而是API根:https://wx2.qq.com/cgi-bin/mmwebwx-bin/
func (wc *wechat) WeChatInit() {
	/*
		params := url.Values{}
		params.Add("pass_ticket", wc.BaseRequest.PassTicket)
		params.Add("skey", wc.BaseRequest.Skey)
		params.Add("r", strconv.FormatInt(time.Now().Unix(), 10))

		uri := wc.BaseURL + "/webwxinit?" + params.Encode()
	*/
	uri := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/" + "webwxinit?pass_ticket=" + wc.BaseRequest.PassTicket + "&skey=" + wc.BaseRequest.Skey + "&r=" + strconv.FormatInt(time.Now().Unix(), 10)

	log.Println(uri)
	log.Println(wc.BaseRequest)
	data, err := json.Marshal(initBaseRequest{
		BaseRequest: wc.BaseRequest,
	})
	if err != nil {
		log.Println("Error happened when do json marshal: ", err.Error())
		return
	}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", wc.userAgent)

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Println("Error happened when do weixin init request: ", err.Error())
		return
	}
	defer resp.Body.Close()

	log.Println("Cookies: ", resp.Cookies())

	data, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Println("Error happend when read json: ", err.Error())
		return
	}

	//Pay more attention to here
	log.Println(string(data))
	reader := bytes.NewReader(data)
	log.Println("========================================================")
	wc.SaveToFile("InitRespBody.json", data)

	var result InitResponse
	if err = json.NewDecoder(reader).Decode(&result); err != nil {
		log.Println("Error happened when decode init information: ", err.Error())
		return
	}

	log.Println(result)
	log.Println(resp.Cookies())

	save, _ := json.Marshal(result)
	wc.SaveToFile("InitResp.json", save)
	wc.SyncKey = &result.SyncKey
	wc.UserName = result.User.UserName
	wc.NickName = result.User.NickName
	return
}

/*

| API | webwxgetcontact |
| --- | --------- |
| url | https://wx.qq.com/cgi-bin/mmwebwx-bin//webwxgetcontact?pass_ticket=xxx&skey=xxx&r=xxx&seq=xxx |
| method | POST |
| data | JSON |
| header | ContentType: application/json; charset=UTF-8 |
// 注意这里要带上BaseRequest头部 不然获得的列表时空的。
// 注意这里Get到的只是通讯录中的联系人，如果是群的话，这个API get到的是群ID并不能GET到群成员。
// 群成员是通过getbatchcontact来实现的。


### 账号类型

| 类型 | 说明 |
| :--: | --- |
| 个人账号 | 以`@`开头，例如：`@xxx` |
| 群聊 | 以`@@`开头，例如：`@@xxx` |
| 公众号/服务号 | 以`@`开头，但其`VerifyFlag` & 8 != 0
	`VerifyFlag`:
		一般公众号/服务号：8
		微信自家的服务号：24
		微信官方账号`微信团队`：56
		特殊账号 | 像文件传输助手之类的账号，有特殊的ID，目前已知的有：`filehelper`, `newsapp`, `fmessage`, `weibo`, `qqmail`, `fmessage`, `tmessage`, `qmessage`, `qqsync`, `floatbottle`, `lbsapp`, `shakeapp`, `medianote`, `qqfriend`, `readerapp`, `blogapp`, `facebookapp`, `masssendapp`, `meishiapp`, `feedsapp`, `voip`, `blogappweixin`, `weixin`, `brandsessionholder`, `weixinreminder`, `officialaccounts`, `notification_messages`, `wxitil`, `userexperience_alarm`, `notification_messages`
*/
func (wc *wechat) GetContactList() {
	//uri := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/" + "webwxgetcontact?pass_ticket=" + wc.BaseRequest.PassTicket + "&skey=" + wc.BaseRequest.Skey + "&r=" + strconv.FormatInt(time.Now().Unix(), 10) + "&seq=0"
	uri := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/" + "webwxgetcontact?pass_ticket=" + wc.BaseRequest.PassTicket + "&skey=" + wc.BaseRequest.Skey + "&r=" + strconv.FormatInt(time.Now().Unix(), 10)
	data, err := json.Marshal(initBaseRequest{
		BaseRequest: wc.BaseRequest,
	})
	if err != nil {
		log.Println("Error happened when do json marshal: ", err.Error())
		return
	}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", wc.userAgent)

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Println("Error happened when do weixin init request: ", err.Error())
		return
	}
	defer resp.Body.Close()

	log.Println("Cookies: ", resp.Cookies())

	data, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Println("Error happend when read json: ", err.Error())
		return
	}

	wc.SaveToFile("GetContactRespBody.json", data)
	//Pay more attention to here
	//log.Println(string(data))
	reader := bytes.NewReader(data)
	var crsp GetContactListResponse
	if err = json.NewDecoder(reader).Decode(&crsp); err != nil {
		log.Println("Error happened when decode init information: ", err.Error())
		return
	}

	wc.ContactDB = make(map[string]Member, len(crsp.MemberList))
	for _, c := range crsp.MemberList {
		if strings.HasPrefix(c.UserName, "@@") {
			wc.GroupList = append(wc.GroupList, c)
		}
		log.Println("Get new contact: ", c.NickName, " -> ", c.UserName)
		wc.ContactDB[c.NickName] = c
	}

	save, _ := json.Marshal(crsp)
	//log.Println(crsp)
	wc.SaveToFile("GetContactResp.json", save)
	log.Println(resp.Cookies())

	//I try to build a cache for fast login without QR code
	cache, _ := json.Marshal(wc)
	wc.SaveToFile("Cache.json", cache)
}

/*
| API | webwxbatchgetcontact |
| --- | --------- |
| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?type=ex&r=xxx&pass_ticket=xxx |
| method | POST |
| data | JSON |
| header | ContentType: application/json; charset=UTF-8 |
| params | { BaseRequest: {
    		Uin: xxx,
		Sid: xxx,
		Skey: xxx,
		DeviceID: xxx
	    },
	    Count: `群数量`,
	    List: [{ UserName: `群ID`, EncryChatRoomId: "" }, ...],
	}
注意这里的返回值和 getcontact的返回值是不同的. 留意两个结构体的内容。
目前测试的结果API 应该为： https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?type=ex&r=xxx&pass_ticket=xxx |
*/
func (wc *wechat) GetGroupMemberList() {
	params := url.Values{}
	params.Add("pass_ticket", wc.BaseRequest.PassTicket)
	params.Add("type", "ex")
	params.Add("r", strconv.FormatInt(time.Now().Unix(), 10))

	uri := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?" + params.Encode()

	//@liwei: Here the group list can be get from the result of webwxgetcontact
	data, err := json.Marshal(CommonReqBody{
		BaseRequest: wc.BaseRequest,
		Count:       len(wc.GroupList),
		List:        wc.GroupList,
	},
	)
	if err != nil {
		log.Println("Error happened when do json marshal: ", err.Error())
		return
	}

	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", wc.userAgent)

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Println("Error happened when do weixin init request: ", err.Error())
		return
	}
	defer resp.Body.Close()

	log.Println("+=================================+")
	log.Println(resp.Status)
	log.Println(resp.Header)
	log.Println(resp.StatusCode)
	log.Println(len(wc.GroupList))
	log.Println(wc.GroupList[0].UserName)

	log.Println("Cookies: ", resp.Cookies())

	data, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Println("Error happend when read json: ", err.Error())
		return
	}

	//Pay more attention to here
	log.Println(string(data))
	reader := bytes.NewReader(data)
	var crsp GetGroupMemberListResponse
	if err = json.NewDecoder(reader).Decode(&crsp); err != nil {
		log.Println("Error happened when decode init information: ", err.Error())
		return
	}

	log.Println(crsp)
	save, _ := json.Marshal(crsp)
	wc.SaveToFile("GetAllGroupMemberResp.json", save)
	log.Println(resp.Cookies())

	cache, _ := json.Marshal(wc)
	wc.SaveToFile("Cache.json", cache)
}

/*
| API | synccheck |
| --- | --------- |
| protocol | https |
| host | webpush.weixin.qq.com webpush.wx2.qq.com webpush.wx8.qq.com webpush.wx.qq.com webpush.web2.wechat.com webpush.web.wechat.com |
| path | /cgi-bin/mmwebwx-bin/synccheck |
| method | GET |
| data | URL Encode |
| params | **r**: `时间戳`  **sid**: xxx **uin**: xxx  **skey**: xxx  **deviceid**: xxx **synckey**: xxx **_**: `时间戳` |

返回数据(String):
```
window.synccheck={retcode:"xxx",selector:"xxx"}

retcode:
	0 正常
	1100 失败/登出微信
selector:
	0 正常
	2 新的消息
	7 进入/离开聊天界面

//心跳函数
*/
func (wc *wechat) SyncCheck() {
	tick := time.Tick(time.Second * 5)
	for {
		<-tick
		params := url.Values{}
		params.Add("r", strconv.FormatInt(time.Now().Unix()*1000, 10))
		params.Add("sid", wc.BaseRequest.Wxsid)
		params.Add("uin", wc.BaseRequest.Wxuin)
		params.Add("skey", wc.BaseRequest.Skey)
		params.Add("deviceid", wc.DeviceID)
		params.Add("synckey", wc.SyncKey.String())
		params.Add("_", strconv.FormatInt(time.Now().Unix()*1000, 10))

		uri := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/synccheck?" + params.Encode()
		resp, err := wc.client.Get(uri)
		if err != nil {
			log.Println("Failed to Sync for server: ", uri, " with error: ", err.Error())
			return
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to Read body for server: ", uri, " with error: ", err.Error())
			return
		}

		log.Println("++++++++++++++++++++++++++++++++++++++++++++++++")
		log.Println(string(data))
		log.Println(resp.Cookies())
		wc.SaveToFile("SyncCheckResponse.txt", data)
		log.Println("++++++++++++++++++++++++++++++++++++++++++++++++")
	}
}

//@liwei: 这里并不是随机的，目前感觉只要更所有的一般性请求使用相同的域名就可以了。
func (wc *wechat) GetSyncServer() bool {
	servers := [...]string{
		`webpush.wx.qq.com`,
		`wx2.qq.com`,
		`webpush.wx2.qq.com`,
		`wx8.qq.com`,
		`webpush.wx8.qq.com`,
		`qq.com`,
		`webpush.wx.qq.com`,
		`web2.wechat.com`,
		`webpush.web2.wechat.com`,
		`wechat.com`,
		`webpush.web.wechat.com`,
		`webpush.weixin.qq.com`,
		`webpush.wechat.com`,
		`webpush1.wechat.com`,
		`webpush2.wechat.com`,
		`webpush2.wx.qq.com`}

	for _, server := range servers {
		<-time.Tick(time.Second * 5)
		log.Printf("Attempt connect: %s ... ... ", server)
		wc.SyncServer = server
		wc.SyncCheck()
		log.Printf("%s connect failed", server)
	}

	return false
}

/*
| API | webwxsync |
| --- | --------- |
| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsync?sid=xxx&skey=xxx&pass_ticket=xxx |
| method | POST |
| data | JSON |
| header | ContentType: application/json; charset=UTF-8 |
| params | { BaseRequest: { Uin: xxx, Sid: xxx, Skey: xxx, DeviceID: xxx }, SyncKey: xxx, rr: `时间戳取反`}

返回数据(JSON):
```
{
	'BaseResponse': {'ErrMsg': '', 'Ret': 0},
	'SyncKey': {
		'Count': 7,
		'List': [
			{'Val': 636214192, 'Key': 1},
			...
		]
	},
	'ContinueFlag': 0,
	'AddMsgCount': 1,
	'AddMsgList': [
		{
			'FromUserName': '',
			'PlayLength': 0,
			'RecommendInfo': {...},
			'Content': "",
			'StatusNotifyUserName': '',
			'StatusNotifyCode': 5,
			'Status': 3,
			'VoiceLength': 0,
			'ToUserName': '',
			'ForwardFlag': 0,
			'AppMsgType': 0,
			'AppInfo': {'Type': 0, 'AppID': ''},
			'Url': '',
			'ImgStatus': 1,
			'MsgType': 51,
			'ImgHeight': 0,
			'MediaId': '',
			'FileName': '',
			'FileSize': '',
			...
		},
		...
	],
	'ModChatRoomMemberCount': 0,
	'ModContactList': [],
	'DelContactList': [],
	'ModChatRoomMemberList': [],
	'DelContactCount': 0,
	...
}
*/
func (wc *wechat) MessageSync() {
	tick := time.Tick(time.Second * 5)
	for _ = range tick {
		params := url.Values{}
		params.Add("skey", wc.BaseRequest.Skey)
		params.Add("sid", wc.BaseRequest.Wxsid)
		params.Add("lang", wc.Lang)
		params.Add("pass_ticket", wc.BaseRequest.PassTicket)

		uri := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsync?" + params.Encode()

		data, err := json.Marshal(CommonReqBody{
			BaseRequest: wc.BaseRequest,
			SyncKey:     wc.SyncKey,
			rr:          ^int(time.Now().Unix()) + 1,
		})
		if err != nil {
			log.Println("Error happend when marshal: ", err.Error())
			continue
		}

		req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
		if err != nil {
			log.Println("Error happend when create http request: ", err.Error())
			continue
		}
		req.Header.Add("Content-Type", "application/json; charset=UTF-8")
		req.Header.Add("User-Agent", wc.userAgent)

		resp, err := wc.client.Do(req)
		if err != nil {
			log.Println("Error happend when do request: ", err.Error())
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
		wc.SaveToFile("MessageSyncResponseBody.json", body)

		reader := bytes.NewReader(body)
		var ms MessageSyncResponse
		if err = json.NewDecoder(reader).Decode(&ms); err != nil {
			log.Println("Error happened when decode Response information: ", err.Error())
			continue
		}

		log.Println(ms)
		save, _ := json.Marshal(ms)
		wc.SaveToFile("MessageSyncResponse.json", save)

		wc.Cookies = resp.Cookies()
		cookie, err := json.Marshal(wc.Cookies)
		if err != nil {
			log.Println("Unable to encoding cookies: ", err.Error())
			continue
		}
		wc.SaveToFile("Cookies.json", cookie)

		cache, _ := json.Marshal(wc)
		wc.SaveToFile("Cache.json", cache)
	}
}

/*
消息一般格式：
{
	"FromUserName": "",
	"ToUserName": "",
	"Content": "",
	"StatusNotifyUserName": "",
	"ImgWidth": 0,
	"PlayLength": 0,
	"RecommendInfo": {...},
	"StatusNotifyCode": 4,
	"NewMsgId": "",
	"Status": 3,
	"VoiceLength": 0,
	"ForwardFlag": 0,
	"AppMsgType": 0,
	"Ticket": "",
	"AppInfo": {...},
	"Url": "",
	"ImgStatus": 1,
	"MsgType": 1,
	"ImgHeight": 0,
	"MediaId": "",
	"MsgId": "",
	"FileName": "",
	"HasProductId": 0,
	"FileSize": "",
	"CreateTime": 1454602196,
	"SubMsgType": 0
}

*/

/*

| MsgType | 说明 |
| ------- | --- |
| 1  | 文本消息 |
| 3  | 图片消息 |
| 34 | 语音消息 |
| 37 | 好友确认消息 |
| 40 | POSSIBLEFRIEND_MSG |
| 42 | 共享名片 |
| 43 | 视频消息 |
| 47 | 动画表情 |
| 48 | 位置消息 |
| 49 | 分享链接 |
| 50 | VOIPMSG |
| 51 | 微信初始化消息 |
| 52 | VOIPNOTIFY |
| 53 | VOIPINVITE |
| 62 | 小视频 |
| 9999 | SYSNOTICE |
| 10000 | 系统消息 |
| 10002 | 撤回消息 |

*/

func (wc *wechat) SaveToFile(name string, content []byte) {
	var file *os.File

	file, err := os.Create(name)
	if err != nil {
		log.Println("Cannot open/create file: ", err.Error())
		return
	}

	file.Write(content)
}

/*

| API | webwxstatusnotify |
| --- | --------- |
| url | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify?lang=zh_CN&pass_ticket=xxx |
| method | POST |
| data | JSON |
| header | ContentType: application/json; charset=UTF-8 |
| params | {
    BaseRequest: {
	Uin: xxx,
	Sid: xxx,
	Skey: xxx,
	DeviceID: xxx
    },
    Code: 3,
    FromUserName: `自己ID`,
    ToUserName: `自己ID`,
    ClientMsgId: `时间戳` <br> }
*/

func (wc *wechat) StatusNotify() {
	params := url.Values{}
	params.Add("lang", wc.Lang)
	params.Add("pass_ticket", wc.BaseRequest.PassTicket)

	uri := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify?" + params.Encode()

	data, err := json.Marshal(CommonReqBody{
		BaseRequest:  wc.BaseRequest,
		Code:         3,
		FromUserName: wc.UserName,
		ToUserName:   wc.UserName,
		ClientMsgId:  int(time.Now().Unix()) + 1,
	})
	if err != nil {
		log.Println("Error happend when marshal: ", err.Error())
		return
	}

	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	if err != nil {
		log.Println("Error happend when create http request: ", err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", wc.userAgent)

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Println("Error happend when do request: ", err.Error())
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	wc.SaveToFile("StatusNotifyBody.json", body)

	reader := bytes.NewReader(body)
	var sn StatusNotifyResponse
	if err = json.NewDecoder(reader).Decode(&sn); err != nil {
		log.Println("Error happened when decode Response information: ", err.Error())
		return
	}

	save, _ := json.Marshal(sn)
	wc.SaveToFile("StatusNotify.json", save)
	log.Println(sn)

}

/*

| API | webwxsendmsg |
| --- | ------------ |
| url | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg?pass_ticket=xxx |
| method | POST |
| data | JSON |
| header | ContentType: application/json; charset=UTF-8 |
| params |
	{
    		BaseRequest: {
			Uin: xxx,
			Sid: xxx,
			Skey: xxx,
			DeviceID: xxx
		},
		Msg: {
			Type: 1 `文字消息`,
			Content: `要发送的消息`,
			FromUserName: `自己ID`,
			ToUserName: `好友ID`,
			LocalID: `与clientMsgId相同`,
			ClientMsgId: `时间戳左移4位随后补上4位随机数`
		}
	}
*/

func (wc *wechat) SendTextMessage() {
	params := url.Values{}
	params.Add("pass_ticket", wc.BaseRequest.PassTicket)

	uri := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg?" + params.Encode()

	data, err := json.Marshal(CommonReqBody{
		BaseRequest: wc.BaseRequest,
		Msg: TextMessage{
			Type:         1,
			Content:      "Hello world",
			FromUserName: wc.UserName,
			ToUserName:   wc.ContactDB["cp4"].UserName,
			LocalID:      int(time.Now().Unix() * 1e4),
			ClientMsgId:  int(time.Now().Unix() * 1e4),
		},
	})

	log.Println(wc.ContactDB)
	log.Println("Sending message from: ", wc.UserName, " to ", wc.ContactDB["cp4"].UserName)
	if err != nil {
		log.Println("Error happend when marshal: ", err.Error())
		return
	}

	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	if err != nil {
		log.Println("Error happend when create http request: ", err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", wc.userAgent)

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Println("Error happend when do request: ", err.Error())
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	wc.SaveToFile("SendTextMessageResponseBody.json", body)
}

/*
| API | webwxsendmsgemotion |
| --- | ------------ |
| url | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsendemoticon?fun=sys&f=json&pass_ticket=xxx |
| method | POST |
| data | JSON |
| header | ContentType: application/json; charset=UTF-8 |
| params | {
    		BaseRequest: {
		    Uin: xxx,
		    Sid: xxx,
		    Skey: xxx,
		    DeviceID: xxx
		},
		Msg: {
		    Type: 47 `emoji消息`,
		    EmojiFlag: 2,
		    MediaId: `表情上传后的媒体ID`,
		    FromUserName: `自己ID`,
		    ToUserName: `好友ID`,
		    LocalID: `与clientMsgId相同`,
		    ClientMsgId: `时间戳左移4位随后补上4位随机数`
		}
	  }
*/
func (wc *wechat) SendEmotionMessage() {

}

/*
| API | webwxrevokemsg |
| --- | ------------ |
| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxrevokemsg |
| method | POST |
| data | JSON |
| header | ContentType: application/json; charset=UTF-8 |
| params | {
    	BaseRequest: {
	    Uin: xxx,
	    Sid: xxx,
	    Skey: xxx,
	    DeviceID: xxx
	},
	SvrMsgId: msg_id,
	ToUserName: user_id,
	ClientMsgId: local_msg_id
    }
*/
func (wc *wechat) RevokeMessage() {

}

/*
### 图片接口

| API | webwxgeticon |
| --- | ------------ |
| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgeticon |
| method | GET |
| params | **seq**: `数字，可为空` <br> **username**: `ID` <br> **skey**: xxx |
<br>

| API | webwxgetheadimg |
| --- | --------------- |
| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetheadimg |
| method | GET |
| params | **seq**: `数字，可为空` <br> **username**: `群ID` <br> **skey**: xxx |
<br>

| API | webwxgetmsgimg |
| --- | --------------- |
| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetmsgimg |
| method | GET |
| params | **MsgID**: `消息ID` <br> **type**: slave `略缩图` or `为空时加载原图` <br> **skey**: xxx |
<br>


*/

/*
### 多媒体接口

| API | webwxgetvideo |
| --- | --------------- |
| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetvideo |
| method | GET |
| params | **msgid**: `消息ID` <br> **skey**: xxx |
<br>

| API | webwxgetvoice |
| --- | --------------- |
| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetvoice |
| method | GET |
| params | **msgid**: `消息ID` <br> **skey**: xxx |
<br>

"https://file.wx.qq.com/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json"

'/webwxverifyuser?lang=zh_CN&r=%s&pass_ticket=%s', time() * 1000, server()->passTicket);
        $data = [
            'BaseRequest'        => server()->baseRequest,
            'Opcode'             => $code,
            'VerifyUserListSize' => 1,
            'VerifyUserList'     => [$ticket ?: $this->verifyTicket()],
            'VerifyContent'      => '',
            'SceneListCount'     => 1,
            'SceneList'          => [33],
            'skey'               => server()->skey,
        ];

*/

func main() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}

	wc := &wechat{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Jar: jar,
		},
		login:       make(chan bool),
		cacheName:   "Cache.json",
		BaseRequest: &BaseRequest{},
		DeviceID:    "0x1234567890",
		userAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36",
		GroupList:   make([]Member, 0, 10),
		BaseURL:     "https://wx2.qq.com/cgi-bin/mmwebwx-bin/",
	}
	if err := wc.FastLogin(); err != nil {
		wc.GetUUID()
		wc.GetQRCode()
		go wc.WaitForQRCodeScan()
		<-wc.login
	}
	wc.GetBaseRequest()
	wc.WeChatInit()
	go wc.SyncCheck()
	wc.StatusNotify()
	wc.GetContactList()
	wc.GetGroupMemberList()
	// wc.GetSyncServer()
	wc.SendTextMessage()
	wc.MessageSync()
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

/*

| MsgType | 说明 |
| ------- | --- |
| 1  | 文本消息 |
| 3  | 图片消息 |
| 34 | 语音消息 |
| 37 | 好友确认消息 |
| 40 | POSSIBLEFRIEND_MSG |
| 42 | 共享名片 |
| 43 | 视频消息 |
| 47 | 动画表情 |
| 48 | 位置消息 |
| 49 | 分享链接 |
| 50 | VOIPMSG |
| 51 | 微信初始化消息 |
| 52 | VOIPNOTIFY |
| 53 | VOIPINVITE |
| 62 | 小视频 |
| 9999 | SYSNOTICE |
| 10000 | 系统消息 |
| 10002 | 撤回消息 |


**微信初始化消息**
```html
MsgType: 51
FromUserName: 自己ID
ToUserName: 自己ID
StatusNotifyUserName: 最近联系的联系人ID
Content:
	<msg>
	    <op id='4'>
	        <username>
	        	// 最近联系的联系人
	            filehelper,xxx@chatroom,wxid_xxx,xxx,...
	        </username>
	        <unreadchatlist>
	            <chat>
	                <username>
	                	// 朋友圈
	                    MomentsUnreadMsgStatus
	                </username>
	                <lastreadtime>
	                    1454502365
	                </lastreadtime>
	            </chat>
	        </unreadchatlist>
	        <unreadfunctionlist>
	        	// 未读的功能账号消息，群发助手，漂流瓶等
	        </unreadfunctionlist>
	    </op>
	</msg>
```

**文本消息**
```
MsgType: 1
FromUserName: 发送方ID
ToUserName: 接收方ID
Content: 消息内容
```

**图片消息**
```html
MsgType: 3
FromUserName: 发送方ID
ToUserName: 接收方ID
MsgId: 用于获取图片
Content:
	<msg>
		<img length="6503" hdlength="0" />
		<commenturl></commenturl>
	</msg>
```

**小视频消息**
```html
MsgType: 62
FromUserName: 发送方ID
ToUserName: 接收方ID
MsgId: 用于获取小视频
Content:
	<msg>
		<img length="6503" hdlength="0" />
		<commenturl></commenturl>
	</msg>
```

**地理位置消息**
```
MsgType: 1
FromUserName: 发送方ID
ToUserName: 接收方ID
Content: http://weixin.qq.com/cgi-bin/redirectforward?args=xxx
// 属于文本消息，只不过内容是一个跳转到地图的链接
```

**名片消息**
```js
MsgType: 42
FromUserName: 发送方ID
ToUserName: 接收方ID
Content:
	<?xml version="1.0"?>
	<msg bigheadimgurl="" smallheadimgurl="" username="" nickname=""  shortpy="" alias="" imagestatus="3" scene="17" province="" city="" sign="" sex="1" certflag="0" certinfo="" brandIconUrl="" brandHomeUrl="" brandSubscriptConfigUrl="" brandFlags="0" regionCode="" />

RecommendInfo:
	{
		"UserName": "xxx", // ID
		"Province": "xxx",
		"City": "xxx",
		"Scene": 17,
		"QQNum": 0,
		"Content": "",
		"Alias": "xxx", // 微信号
		"OpCode": 0,
		"Signature": "",
		"Ticket": "",
		"Sex": 0, // 1:男, 2:女
		"NickName": "xxx", // 昵称
		"AttrStatus": 4293221,
		"VerifyFlag": 0
	}
```

**语音消息**
```html
MsgType: 34
FromUserName: 发送方ID
ToUserName: 接收方ID
MsgId: 用于获取语音
Content:
	<msg>
		<voicemsg endflag="1" cancelflag="0" forwardflag="0" voiceformat="4" voicelength="1580" length="2026" bufid="216825389722501519" clientmsgid="49efec63a9774a65a932a4e5fcd4e923filehelper174_1454602489" fromusername="" />
	</msg>
```

**动画表情**
```html
MsgType: 47
FromUserName: 发送方ID
ToUserName: 接收方ID
Content:
	<msg>
		<emoji fromusername = "" tousername = "" type="2" idbuffer="media:0_0" md5="e68363487d8f0519c4e1047de403b2e7" len = "86235" productid="com.tencent.xin.emoticon.bilibili" androidmd5="e68363487d8f0519c4e1047de403b2e7" androidlen="86235" s60v3md5 = "e68363487d8f0519c4e1047de403b2e7" s60v3len="86235" s60v5md5 = "e68363487d8f0519c4e1047de403b2e7" s60v5len="86235" cdnurl = "http://emoji.qpic.cn/wx_emoji/eFygWtxcoMF8M0oCCsksMA0gplXAFQNpiaqsmOicbXl1OC4Tyx18SGsQ/" designerid = "" thumburl = "http://mmbiz.qpic.cn/mmemoticon/dx4Y70y9XctRJf6tKsy7FwWosxd4DAtItSfhKS0Czr56A70p8U5O8g/0" encrypturl = "http://emoji.qpic.cn/wx_emoji/UyYVK8GMlq5VnJ56a4GkKHAiaC266Y0me0KtW6JN2FAZcXiaFKccRevA/" aeskey= "a911cc2ec96ddb781b5ca85d24143642" ></emoji>
		<gameext type="0" content="0" ></gameext>
	</msg>
```

**普通链接或应用分享消息**
```html
MsgType: 49
AppMsgType: 5
FromUserName: 发送方ID
ToUserName: 接收方ID
Url: 链接地址
FileName: 链接标题
Content:
	<msg>
		<appmsg appid=""  sdkver="0">
			<title></title>
			<des></des>
			<type>5</type>
			<content></content>
			<url></url>
			<thumburl></thumburl>
			...
		</appmsg>
		<appinfo>
			<version></version>
			<appname></appname>
		</appinfo>
	</msg>
```

**音乐链接消息**
```html
MsgType: 49
AppMsgType: 3
FromUserName: 发送方ID
ToUserName: 接收方ID
Url: 链接地址
FileName: 音乐名

AppInfo: // 分享链接的应用
	{
		Type: 0,
		AppID: wx485a97c844086dc9
	}

Content:
	<msg>
		<appmsg appid="wx485a97c844086dc9"  sdkver="0">
			<title></title>
			<des></des>
			<action></action>
			<type>3</type>
			<showtype>0</showtype>
			<mediatagname></mediatagname>
			<messageext></messageext>
			<messageaction></messageaction>
			<content></content>
			<contentattr>0</contentattr>
			<url></url>
			<lowurl></lowurl>
			<dataurl>
				http://ws.stream.qqmusic.qq.com/C100003i9hMt1bgui0.m4a?vkey=6867EF99F3684&amp;guid=ffffffffc104ea2964a111cf3ff3edaf&amp;fromtag=46
			</dataurl>
			<lowdataurl>
				http://ws.stream.qqmusic.qq.com/C100003i9hMt1bgui0.m4a?vkey=6867EF99F3684&amp;guid=ffffffffc104ea2964a111cf3ff3edaf&amp;fromtag=46
			</lowdataurl>
			<appattach>
				<totallen>0</totallen>
				<attachid></attachid>
				<emoticonmd5></emoticonmd5>
				<fileext></fileext>
			</appattach>
			<extinfo></extinfo>
			<sourceusername></sourceusername>
			<sourcedisplayname></sourcedisplayname>
			<commenturl></commenturl>
			<thumburl>
				http://imgcache.qq.com/music/photo/album/63/180_albumpic_143163_0.jpg
			</thumburl>
			<md5></md5>
		</appmsg>
		<fromusername></fromusername>
		<scene>0</scene>
		<appinfo>
			<version>29</version>
			<appname>摇一摇搜歌</appname>
		</appinfo>
		<commenturl></commenturl>
	</msg>
```

**群消息**
```
MsgType: 1
FromUserName: @@xxx
ToUserName: @xxx
Content:
	@xxx:<br/>xxx
```

**红包消息**
```
MsgType: 49
AppMsgType: 2001
FromUserName: 发送方ID
ToUserName: 接收方ID
Content: 未知
```
注：根据网页版的代码可以看到未来可能支持查看红包消息，但目前走的是系统消息，见下。

**系统消息**
```
MsgType: 10000
FromUserName: 发送方ID
ToUserName: 自己ID
Content:
	"你已添加了 xxx ，现在可以开始聊天了。"
	"如果陌生人主动添加你为朋友，请谨慎核实对方身份。"
	"收到红包，请在手机上查看"

*/
