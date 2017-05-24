package wechat

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
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"sync/atomic"

	filetype "gopkg.in/h2non/filetype.v1"
	//"gopkg.in/h2non/filetype.v1/types"
	//"net/http/httputil"
	"db"
	"handler"
	"message"
	"net/url"
	"os"
	"os/signal"
	"processor"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"util"
)

type WeChat struct {
	client            *http.Client         `json:"-"`
	cacheName         string               `json:"-"`
	UserAgent         string               `json:"UserAgent"`
	uuid              string               `json:"-"`
	Redirect          string               `json:"Redirect"`
	UUID              string               `json:"UUID"`
	Scan              string               `json:"-"`
	Login             chan bool            `json:"-"`
	Name              string               `json:"Name"`
	Ticket            string               `json:"Ticket"`
	Lang              string               `json:"Lang"`
	BaseURL           string               `json:"BaseURL"`
	BaseRequest       *BaseRequest         `json:"BaseRequest"`
	APIPath           string               `json:"APIPath"`
	DeviceID          string               `json:"DeviceID"`
	GroupList         []Member             `json:"GroupList"`
	SyncServer        string               `json:"SyncServer"`
	SyncKey           *SyncKey             `json:"SyncKey"`
	UserName          string               `json:"UserName"`
	NickName          string               `json:"NickName"`
	ContactDBNickName map[string]Member    `json:"ContactDBNickName"`
	ContactDBUserName map[string]Member    `json:"ContactDBUserName"`
	Cookies           []*http.Cookie       `json:"Cookies"`
	DB                *db.FileDB           `json:"DB"`
	Cookie            map[string]string    `json:"Cookie"`
	MediaCount        uint32               `json:"MediaCount"`
	Shutdown          chan bool            `json:"-"`
	Signal            chan os.Signal       `json:"-"`
	Consumer          *processor.Processor `json:"-"`
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

type SendMessageResponse struct {
	BaseResponse BaseResponse
	MsgID        string `json:"MsgId"`
	LocalID      int64  `json:"LocalID"`
}

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
	BaseResponse           BaseResponse      `json:"BaseResponse"`
	AddMsgCount            int               `json:"AddMsgCount"` //New message count
	AddMsgList             []message.Message `json:"AddMsgList"`
	ModContactCount        int               `json:"ModContactCount"` //Changed Contact count
	ModContactList         []Member          `json:"ModContactList"`
	DelContactCount        int               `json:"DelContactCount"` //Delete Contact count
	DelContactList         []Member          `json:"DelContactList"`
	ModChatRoomMemberCount int               `json:"ModChatRoomMemberCount"`
	ModChatRoomMemberList  []Member          `json:"ModChatRoomMemberList"`
	Profile                Profile           `json:"Profile"`
	ContinueFlag           int               `json:"ContinueFlag"`
	SyncKey                SyncKey           `json:"SyncKey"`
	Skey                   string            `json:"Skey"`
	SyncCheckKey           SyncKey           `json:"SyncCheckKey"`
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

func (wc *WeChat) FastLogin() error {
	data, err := wc.DB.Get(wc.cacheName)
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

	err = wc.SyncCheck()
	if err != nil {
		return errors.New("Cannot do fast login: " + err.Error())
	}
	return nil
}

func (wc *WeChat) SetReqCookies(req *http.Request) {
	for _, c := range wc.Cookies {
		req.AddCookie(c)
	}
}

//| url | https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxpushloginurl |
func (wc *WeChat) GetUUID() error {

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

func (wc *WeChat) GetQRCode() error {
	qrterminal.Generate("https://login.weixin.qq.com/l/"+wc.uuid, qrterminal.L, os.Stdout)
	return nil
}

func (wc *WeChat) WaitForQRCodeScan() {
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
			wc.BaseURL = urls.Scheme + "://" + urls.Host

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
			wc.Login <- true
			break
		} else if code == "201" {
			log.Println("Please confirm login on your telephone")
		} else if code == "408" {
			wc.GetQRCode()
		}
	}
}

func (wc *WeChat) GetBaseRequest() {
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

	wc.DB.Save("BaseRequestBody.json", data)
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
	wc.DB.Save("BaseRequest.json", br)
	log.Println("Login successfully")

	return
}

func NewWeChatClient(name string) (*WeChat, error) {
	if name == "" {
		return nil, errors.New("Name cannot be empty")
	}

	database, err := db.NewFileDB("." + name)
	if err != nil {
		return nil, errors.New("Cannot create Cache DB: " + err.Error())
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.New("Cannot create cookie for Client: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jar: jar,
	}

	wc := &WeChat{
		Login:       make(chan bool),
		client:      client,
		cacheName:   "Cache.json",
		BaseRequest: &BaseRequest{},
		DeviceID:    util.GenerateDeviceID(),
		UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36",
		GroupList:   make([]Member, 0, 10),
		Shutdown:    make(chan bool),
		Signal:      make(chan os.Signal),
		APIPath:     "/cgi-bin/mmwebwx-bin/",
		DB:          database,
		Consumer:    processor.DefaultProcessor,
	}

	return wc, nil
}

func (wc *WeChat) WeChatInit() {
	/*
		params := url.Values{}
		params.Add("pass_ticket", wc.BaseRequest.PassTicket)
		params.Add("skey", wc.BaseRequest.Skey)
		params.Add("r", strconv.FormatInt(time.Now().Unix(), 10))

		uri := wc.BaseURL + "/webwxinit?" + params.Encode()
	*/
	uri := wc.BaseURL + wc.APIPath + "webwxinit?pass_ticket=" + wc.BaseRequest.PassTicket + "&skey=" + wc.BaseRequest.Skey + "&r=" + strconv.FormatInt(time.Now().Unix(), 10)

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
	req.Header.Add("User-Agent", wc.UserAgent)

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

	reader := bytes.NewReader(data)
	wc.DB.Save("InitRespBody.json", data)

	var result InitResponse
	if err = json.NewDecoder(reader).Decode(&result); err != nil {
		log.Println("Error happened when decode init information: ", err.Error())
		return
	}

	save, err := json.Marshal(result)
	if err != nil {
		log.Println("Error happened when marshal: ", err.Error())
		return
	}
	wc.DB.Save("InitResp.json", save)
	wc.SyncKey = &result.SyncKey
	wc.UserName = result.User.UserName
	wc.NickName = result.User.NickName
	return
}

func (wc *WeChat) GetContactList() {
	uri := wc.BaseURL + wc.APIPath + "webwxgetcontact?pass_ticket=" + wc.BaseRequest.PassTicket + "&skey=" + wc.BaseRequest.Skey + "&r=" + strconv.FormatInt(time.Now().Unix(), 10)
	data, err := json.Marshal(initBaseRequest{
		BaseRequest: wc.BaseRequest,
	})
	if err != nil {
		log.Println("Error happened when do json marshal: ", err.Error())
		return
	}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", wc.UserAgent)

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

	wc.DB.Save("GetContactRespBody.json", data)
	reader := bytes.NewReader(data)
	var crsp GetContactListResponse
	if err = json.NewDecoder(reader).Decode(&crsp); err != nil {
		log.Println("Error happened when decode init information: ", err.Error())
		return
	}

	wc.ContactDBNickName = make(map[string]Member, len(crsp.MemberList))
	wc.ContactDBUserName = make(map[string]Member, len(crsp.MemberList))
	for _, c := range crsp.MemberList {
		if strings.HasPrefix(c.UserName, "@@") {
			wc.GroupList = append(wc.GroupList, c)
		}
		log.Println("Get new contact: ", c.NickName, " -> ", c.UserName)
		wc.ContactDBNickName[c.NickName] = c
		wc.ContactDBUserName[c.UserName] = c
	}

	save, err := json.Marshal(crsp)
	if err != nil {
		log.Println("Error happened when marshal: ", err.Error())
		return
	}
	//log.Println(crsp)
	wc.DB.Save("GetContactResp.json", save)
	log.Println(resp.Cookies())

	//I try to build a cache for fast login without QR code
	cache, err := json.Marshal(wc)
	if err != nil {
		log.Println("Error happened when marshal: ", err.Error())
		return
	}
	wc.DB.Save("Cache.json", cache)
}

//请求群组
func (wc *WeChat) GetGroupMemberList() {
	params := url.Values{}
	params.Add("pass_ticket", wc.BaseRequest.PassTicket)
	params.Add("type", "ex")
	params.Add("r", strconv.FormatInt(time.Now().Unix(), 10))

	uri := wc.BaseURL + wc.APIPath + "webwxbatchgetcontact?" + params.Encode()

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
	req.Header.Add("User-Agent", wc.UserAgent)

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
	var crsp GetGroupMemberListResponse
	if err = json.NewDecoder(reader).Decode(&crsp); err != nil {
		log.Println("Error happened when decode init information: ", err.Error())
		return
	}

	save, err := json.Marshal(crsp)
	if err != nil {
		log.Println("Error happened when marshal: ", err.Error())
		return
	}
	wc.DB.Save("GetAllGroupMemberResp.json", save)

	cache, err := json.Marshal(wc)
	if err != nil {
		log.Println("Error happened when marshal: ", err.Error())
		return
	}
	wc.DB.Save("Cache.json", cache)
}

//@liwei: 这里应该在该函数返回2以后再去用消息同步接口来获取消息

var retCodeAndSelector = regexp.MustCompile(`window\.synccheck=\{retcode\:\"(?P<retcode>[0-9]+)\"\,selector\:\"(?P<selector>[0-9]+)\"\}`)

func (wc *WeChat) SyncCheck() error {
	params := url.Values{}
	params.Add("r", strconv.FormatInt(time.Now().Unix()*1000, 10))
	params.Add("sid", wc.BaseRequest.Wxsid)
	params.Add("uin", wc.BaseRequest.Wxuin)
	params.Add("skey", wc.BaseRequest.Skey)
	params.Add("deviceid", wc.DeviceID)
	params.Add("synckey", wc.SyncKey.String())
	params.Add("_", strconv.FormatInt(time.Now().Unix()*1000, 10))

	uri := wc.BaseURL + wc.APIPath + "synccheck?" + params.Encode()
	resp, err := wc.client.Get(uri)
	if err != nil {
		log.Println("Failed to Sync for server: ", uri, " with error: ", err.Error())
		return errors.New("Failed to Sync for server")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Failed to read response body")
	}

	wc.DB.Save("SyncCheckResponse.txt", data)

	matches := retCodeAndSelector.FindSubmatch(data)
	if len(matches) == 0 {
		return errors.New("SynCheck Result parse failed, please check!")
	}

	retCode, err := strconv.ParseInt(string(matches[1]), 10, 64)
	if err != nil {
		return errors.New("Cannot parse retcode")
	}
	selector, err := strconv.ParseInt(string(matches[2]), 10, 64)
	if err != nil {
		return errors.New("Cannot parse selector")
	}

	log.Println(retCode, selector)
	if retCode == 0 { //Normal
		if selector == 0 { //Normal State, Nothing Happened
			return nil
		} else if selector == 2 { // Received new message
			go wc.MessageSync()
			return nil
		} else if selector == 7 { // Enter/Leave ?
			go wc.MessageSync()
			log.Println("Currently I don't know what is this mean ? selector ==: ", selector)
			return nil
		} else if selector == 5 { // Enter/Leave ?
			go wc.MessageSync()
			log.Println("Currently I don't know what is this mean ? selector ==: ", selector)
			return nil
		} else {
			log.Println("Unknown selector returned during syncheck process: " + strconv.Itoa(int(selector)))
			return nil
		}
	} else if retCode == 1100 { //Error happend or logout
		return errors.New("You have already logout, Try to relogin!!!!!")
	} else { //Unknown retCode
		return errors.New("Received Unknown retCode during syncheck process: " + strconv.Itoa(int(retCode)))
	}
}

//@liwei: 这里并不是随机的，目前感觉只要更所有的一般性请求使用相同的域名就可以了。
func (wc *WeChat) GetSyncServer() bool {
	servers := [...]string{
		`webpush.wx.qq.com`,
		`wx2.qq.com`,
		`webpush.wx2.qq.com`,
		`wx8.qq.com`,
		`webpush.wx8.qq.com`,
		`qq.com`,
		`webpush.wx.qq.com`,
		`web2.WeChat.com`,
		`webpush.web2.WeChat.com`,
		`WeChat.com`,
		`webpush.web.WeChat.com`,
		`webpush.weixin.qq.com`,
		`webpush.WeChat.com`,
		`webpush1.WeChat.com`,
		`webpush2.WeChat.com`,
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

//@liwei: 注意，消息同步过程中synckey在最开始的时候使用init时获取的值。
//此后每次都要使用最近一次获取的synckey值来进行同步，否则每次获取到的都是从init到目前的消息
func (wc *WeChat) MessageSync() {
	params := url.Values{}
	params.Add("skey", wc.BaseRequest.Skey)
	params.Add("sid", wc.BaseRequest.Wxsid)
	params.Add("lang", wc.Lang)
	params.Add("pass_ticket", wc.BaseRequest.PassTicket)

	uri := wc.BaseURL + wc.APIPath + "webwxsync?" + params.Encode()

	data, err := json.Marshal(CommonReqBody{
		BaseRequest: wc.BaseRequest,
		SyncKey:     wc.SyncKey,
		rr:          ^int(time.Now().Unix()) + 1,
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
	req.Header.Add("User-Agent", wc.UserAgent)

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Println("Error happend when do request: ", err.Error())
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	wc.Cookies = resp.Cookies()
	//log.Println(string(body))
	wc.DB.Save("MessageSyncResponseBody.json", body)

	reader := bytes.NewReader(body)
	var ms MessageSyncResponse
	if err = json.NewDecoder(reader).Decode(&ms); err != nil {
		log.Println("Error happened when decode Response information: ", err.Error())
		return
	}

	//@liwei: We need to carefully handle the message
	//log.Println(ms)
	save, err := json.Marshal(ms)
	if err != nil {
		log.Println("Error happened when marshal: ", err.Error())
		return
	}
	wc.DB.Save("MessageSyncResponse.json", save)

	wc.Cookie = make(map[string]string, len(wc.Cookies))
	for _, c := range wc.Cookies {
		wc.Cookie[c.Name] = c.Value
	}
	cookie, err := json.Marshal(wc.Cookies)
	if err != nil {
		log.Println("Unable to encoding cookies: ", err.Error())
		return
	}

	//@liwei: We Must Update the SyncKey everytime. With this we can just get the recent message.
	wc.SyncKey = &ms.SyncKey
	wc.DB.Save("Cookies.json", cookie)

	cache, err := json.Marshal(wc)
	if err != nil {
		log.Println("Error happened when marshal: ", err.Error())
		return
	}
	wc.DB.Save("Cache.json", cache)

	wc.HandleMessageSyncResponse(&ms)
}

func (wc *WeChat) HandleMessageSyncResponse(resp *MessageSyncResponse) {
	if resp.AddMsgCount > 0 {
		go wc.HandleNewMessage(resp.AddMsgList)
	}

	if resp.ModContactCount > 0 {
		go wc.HandleModContact(resp.ModContactList)
	}

	if resp.DelContactCount > 0 {
		go wc.HandleDelContact(resp.DelContactList)
	}

	if resp.ModChatRoomMemberCount > 0 {
		go wc.HandleModChatRoomMember(resp.ModChatRoomMemberList)
	}
}

func (wc *WeChat) HandleNewMessage(msg []message.Message) {
	for _, m := range msg {
		var local message.LocalMessage
		from, err := wc.GetNickNameByUserName(m.From())
		if err != nil {
			log.Println("Cannot process new message: ", err)
			continue
		}

		content, err := wc.GetMessageContent(&m)
		if err != nil {
			log.Println("Cannot process new message: ", err)
			continue
		}

		local.SetFrom(from).SetTo(wc.NickName).SetType(m.MsgType).SetContent(content)

		resp, err := wc.Consumer.Process(&local)

		for _, r := range resp {
			if r != nil && r.IsValid() {
				wc.SendMessage(r)
			}
		}
	}
}

func (wc *WeChat) SendMessage(msg *message.LocalMessage) {
	to, err := wc.GetUserNameByNickName(msg.To())
	if err != nil {
		log.Println("Cannot find user: ", msg.To())
		return
	}

	switch msg.Type() {
	case message.TEXT:
		wc.SendTextMessage(wc.UserName, to, string(msg.Content()))
	case message.IMAGE:
		wc.SendImageMessage(wc.UserName, to, string(msg.Content()))
	case message.EMOJI:
		wc.SendEmojiMessage(wc.UserName, to, string(msg.Content()))
	default:
		log.Println("Send message of type: ", message.Type[msg.Type()], " is not emplemented!")
	}
}

func (wc *WeChat) GetNickNameByUserName(name string) (string, error) {
	if m, ok := wc.ContactDBUserName[name]; ok {
		return m.NickName, nil
	}

	return "", errors.New(name + " is not exist")
}

func (wc *WeChat) GetUserNameByNickName(name string) (string, error) {
	if m, ok := wc.ContactDBNickName[name]; ok {
		return m.UserName, nil
	}

	return "", errors.New(name + " is not exist")
}

func (wc *WeChat) GetMessageContent(msg *message.Message) ([]byte, error) {
	if msg.IsTextMessage() {
		return wc.GetTextMessageContent(msg), nil
	} else {
		return []byte("Not implemented"), errors.New("Not Implemented!")
	}
}

func (wc *WeChat) GetTextMessageContent(msg *message.Message) []byte {
	return []byte(msg.Content)
}

func (wc *WeChat) GetImageMessageContent(msg *message.Message) []byte {
	return []byte("Not implemented")
}

func (wc *WeChat) GetVideoMessageContent(msg *message.Message) []byte {
	return []byte("Not implemented")
}

func (wc *WeChat) GetEmojiMessageContent(msg *message.Message) []byte {
	return []byte("Not implemented")
}

func (wc *WeChat) GetVoiceMessageContent(msg *message.Message) []byte {
	return []byte("Not implemented")
}

func (wc *WeChat) GetShortVideoMessageContent(msg *message.Message) []byte {
	return []byte("Not implemented")
}

func (wc *WeChat) GetVerifyMessageContent(msg *message.Message) []byte {
	return []byte("Not implemented")
}

func (wc *WeChat) HandleModContact(members []Member) {
	log.Println("Modified Contact: ")
	for _, m := range members {
		log.Println(m)
	}
}

func (wc *WeChat) HandleDelContact(members []Member) {
	log.Println("Deleted Contact: ")
	for _, m := range members {
		log.Println(m)
	}
}

func (wc *WeChat) HandleModChatRoomMember(members []Member) {
	log.Println("Mod Chat Room Member: ")
	for _, m := range members {
		log.Println(m)
	}
}

//这个函数到底是用来干啥的？ ----> 开启微信状态通知
func (wc *WeChat) StatusNotify() {
	params := url.Values{}
	params.Add("lang", wc.Lang)
	params.Add("pass_ticket", wc.BaseRequest.PassTicket)

	uri := wc.BaseURL + wc.APIPath + "webwxstatusnotify?" + params.Encode()

	data, err := json.Marshal(CommonReqBody{
		BaseRequest:  wc.BaseRequest,
		Code:         3,
		FromUserName: wc.UserName,
		ToUserName:   wc.UserName,
		ClientMsgId:  int64(time.Now().Unix()) + 1,
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
	req.Header.Add("User-Agent", wc.UserAgent)

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Println("Error happend when do request: ", err.Error())
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	wc.DB.Save("StatusNotifyBody.json", body)

	reader := bytes.NewReader(body)
	var sn StatusNotifyResponse
	if err = json.NewDecoder(reader).Decode(&sn); err != nil {
		log.Println("Error happened when decode Response information: ", err.Error())
		return
	}

	save, err := json.Marshal(sn)
	if err != nil {
		log.Println("Error happened when marshal: ", err.Error())
		return
	}
	wc.DB.Save("StatusNotify.json", save)
	//log.Println(sn)
}

func (wc *WeChat) SendTextMessage(from, to, msg string) {
	params := url.Values{}
	params.Add("pass_ticket", wc.BaseRequest.PassTicket)

	uri := wc.BaseURL + wc.APIPath + "webwxsendmsg?" + params.Encode()

	data, err := json.Marshal(CommonReqBody{
		BaseRequest: wc.BaseRequest,
		Msg: message.TextMessage{
			Type:         1,
			Content:      msg,
			FromUserName: from,
			ToUserName:   to,
			LocalID:      int64(time.Now().Unix() * 1e4),
			ClientMsgId:  int64(time.Now().Unix() * 1e4),
		},
	})

	log.Println("Sending message from: ", from, " to ", to)
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
	req.Header.Add("User-Agent", wc.UserAgent)

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Println("Error happend when do request: ", err.Error())
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	wc.DB.Save("SendTextMessageResponseBody.json", body)
}

func (wc *WeChat) SendShareLinkMessage(from, to, msg string) {
	params := url.Values{}
	params.Add("pass_ticket", wc.BaseRequest.PassTicket)

	uri := wc.BaseURL + wc.APIPath + "webwxsendmsg?" + params.Encode()

	data, err := json.Marshal(CommonReqBody{
		BaseRequest: wc.BaseRequest,
		Msg: message.TextMessage{
			Type:         49,
			Content:      msg,
			FromUserName: from,
			ToUserName:   to,
			LocalID:      int64(time.Now().Unix() * 1e4),
			ClientMsgId:  int64(time.Now().Unix() * 1e4),
		},
	})

	log.Println("Sending message from: ", wc.UserName, " to ", to)
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
	req.Header.Add("User-Agent", wc.UserAgent)

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Println("Error happend when do request: ", err.Error())
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	wc.DB.Save("SendTextMessageResponseBody.json", body)
}

func (wc *WeChat) SendEmojiMessage(from, to, id string) {
	params := url.Values{}
	//params.Add("pass_ticket", wc.BaseRequest.PassTicket)
	params.Add("fun", "sys")
	params.Add("lang", wc.Lang)

	uri := wc.BaseURL + wc.APIPath + "webwxsendemoticon?" + params.Encode()

	data, err := json.Marshal(CommonReqBody{
		BaseRequest: wc.BaseRequest,
		Msg: message.EmotionMessage{
			Type:         47,
			EmojiFlag:    2,
			FromUserName: from,
			ToUserName:   wc.ContactDBNickName["cp4"].UserName,
			LocalID:      int64(time.Now().Unix() * 1e4),
			ClientMsgId:  int64(time.Now().Unix() * 1e4),
			MediaId:      id,
		},
		Scene: 0,
	})

	log.Println(wc.ContactDBNickName)
	log.Println("Sending message from: ", wc.UserName, " to ", wc.ContactDBNickName["cp4"].UserName)
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
	req.Header.Add("User-Agent", wc.UserAgent)
	wc.SetReqCookies(req)

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Println("Error happend when do request: ", err.Error())
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	wc.DB.Save("SendEmotionMessageResponseBody.json", body)
}

func (wc *WeChat) SendImageMessage(from, to, id string) {
	params := url.Values{}
	params.Add("pass_ticket", wc.BaseRequest.PassTicket)
	params.Add("fun", "async")
	params.Add("f", "json")
	params.Add("lang", wc.Lang)

	log.Println("++==", id, "==++")
	uri := wc.BaseURL + wc.APIPath + "webwxsendmsgimg?" + params.Encode()

	data, err := json.Marshal(CommonReqBody{
		BaseRequest: wc.BaseRequest,
		Msg: message.MediaMessage{
			Type:         3,
			Content:      "",
			FromUserName: from,
			ToUserName:   to,
			LocalID:      int64(time.Now().Unix() * 1e4),
			ClientMsgId:  int64(time.Now().Unix() * 1e4),
			MediaId:      id,
		},
		Scene: 0,
	})

	log.Println(wc.ContactDBUserName[wc.UserName].NickName)
	log.Println(wc.ContactDBNickName["cp4"].NickName)
	log.Println(string(data))
	log.Println("Sending message from: ", wc.UserName, " to ", wc.ContactDBNickName["cp4"].UserName)
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
	req.Header.Add("User-Agent", wc.UserAgent)
	wc.SetReqCookies(req)

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Println("Error happend when do request: ", err.Error())
		return
	}

	log.Println(resp.Status)
	log.Println(resp.StatusCode)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	wc.DB.Save("SendImageMessageResponseBody.json", body)
}

func (wc *WeChat) UploadMedia(name, to string) (string, error) {
	//func (wc *weChat) UploadMedia(buf []byte, kind types.Type, info os.FileInfo, to string) (string, error) {
	file, err := os.Stat(name)
	if err != nil {
		return "", err
	}

	buf, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	kind, _ := filetype.Get(buf)

	var mediatype string
	if filetype.IsImage(buf) {
		mediatype = `pic`
	} else if filetype.IsVideo(buf) {
		mediatype = `video`
	} else {
		mediatype = `doc`
	}

	atomic.AddUint32(&wc.MediaCount, 1)
	fields := map[string]string{
		`id`:                `WU_FILE_` + strconv.Itoa(int(wc.MediaCount)), //@liwei
		`name`:              file.Name(),
		`type`:              kind.MIME.Value,
		`lastModifiedDate`:  file.ModTime().UTC().String(),
		`size`:              strconv.FormatInt(file.Size(), 10),
		`mediatype`:         mediatype,
		`pass_ticket`:       wc.BaseRequest.PassTicket,
		`webwx_data_ticket`: wc.Cookie["webwx_data_ticket"],
	}

	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	fw, err := writer.CreateFormFile(`filename`, file.Name())
	if err != nil {
		return "", err
	}
	fw.Write(buf)

	for k, v := range fields {
		writer.WriteField(k, v)
	}

	data, err := json.Marshal(CommonReqBody{
		BaseRequest:   wc.BaseRequest,
		ClientMediaId: int(time.Now().Unix() * 1e4),
		TotalLen:      strconv.FormatInt(file.Size(), 10),
		StartPos:      0,
		DataLen:       strconv.FormatInt(file.Size(), 10),
		//UploadType:    2,
		MediaType:    4,
		ToUserName:   "cp4",
		FromUserName: wc.UserName,
		//FileMd5:       string(md5.New().Sum(buf)),
	})

	writer.WriteField(`uploadmediarequest`, string(data))
	writer.Close()

	//req, err := http.NewRequest("POST", "https://file.wx.qq.com/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json", buffer)
	req, err := http.NewRequest("POST", "https://file2.wx.qq.com/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json", buffer)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("User-Agent", wc.UserAgent)

	wc.SetReqCookies(req)

	resp, err := wc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))

	reader := bytes.NewReader(body)
	var result UploadMediaResponse
	if err = json.NewDecoder(reader).Decode(&result); err != nil {
		log.Println("Error happened when decode Response information: ", err.Error())
		return "", errors.New("Cannot decode")
	}

	save, err := json.Marshal(result)
	if err != nil {
		log.Println("Error happened when marshal: ", err.Error())
		return "", err
	}
	wc.DB.Save("UploadMediaResponse.json", save)
	return result.MediaID, nil
}

func (wc *WeChat) RevokeMessage() {

}

func (wc *WeChat) Run() {
	signal.Notify(wc.Signal, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		for s := range wc.Signal {
			log.Println("Received signal: ", s, " Shuttdown the instance.")
			wc.Shutdown <- true
		}
	}()

	tick := time.Tick(time.Second * 10)
	go func() {
		for _ = range tick {
			wc.SyncCheck()
		}
	}()
	<-wc.Shutdown
}

func (wc *WeChat) AddHandler(h *handler.Handler) {
	wc.Consumer.RegisterHandler(h)
}

func (wc *WeChat) DelHandler(h *handler.Handler) {
	wc.Consumer.DeRegisterHandler(h)
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


	uri := common.CgiUrl + "/webwxverifyuser?" + km.Encode()
	uri := common.CgiUrl + "/webwxcreatechatroom?" + km.Encode()
*/
