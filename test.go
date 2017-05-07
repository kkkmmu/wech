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
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type wechat struct {
	client      *http.Client
	uuid        string
	Redirect    string
	UUID        string
	Ticket      string
	Lang        string
	Scan        string
	Login       chan bool
	BaseURL     string
	BaseRequest *BaseRequest
	DeviceID    string
	userAgent   string
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
	Wxuin      int64  `xml:"wxuin" json:"Uin"`
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
	Uin              int64  `json:"Uin"`
	UserName         string `json:"UserName"`
	NickName         string `json:"NickName"`
	HeadImgUrl       string `json:"HeadImgUrl"`
	ContactFlag      int    `json:"ContactFlag"`
	MemberCount      int    `json:"MemberCount"`
	MemberList       []User `json:"MemberList"`
	RemarkName       string `json:"RemarkName"`
	HideInputBarFlag int    `json:"HideInputBarFlag"`
	Sex              int    `json:"Sex"`
	Signature        string `json:"Signature"`
	VerifyFlag       int    `json:"VerifyFlag"`
	OwnerUin         int    `json:"OwnerUin"`
	PYInitial        string `json:"PYInitial"`
	PYQuanPin        string `json:"PYQuanPin"`
	RemarkPYInitial  string `json:"RemarkPYInitial"`
	RemarkPYQuanPin  string `json:"RemarkPYQuanPin"`
	StarFriend       int    `json:"StarFriend"`
	AppAccountFlag   int    `json:"AppAccountFlag"`
	Statues          int    `json:"Statues"`
	AttrStatus       int    `json:"AttrStatus"`
	Province         string `json:"Province"`
	City             string `json:"City"`
	Alias            string `json:"Alias"`
	SnsFlag          int    `json:"SnsFlag"`
	UniFriend        int    `json:"UniFriend"`
	DisplayName      string `json:"DisplayName"`
	ChatRoomId       int    `json:"ChatRoomId"`
	KeyWord          string `json:"KeyWord"`
	EncryChatRoomId  string `json:"EncryChatRoomId"`
	IsOwner          int    `json:"IsOwner"`
}

type contactResponse struct {
	Response
	MemberCount int
	MemberList  []Member
	Seq         float64
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

type KeyVal struct {
	Key int `json:"Key"`
	Val int `json:"Val"`
}

type InitResp struct {
	Response
	Count               int      `json:"Count"`
	User                User     `json:"User"`
	ContactList         []Member `json:"ContactList"`
	SyncKey             SyncKey  `json:"SyncKey"`
	ChatSet             string   `json:"ChatSet"`
	SKey                string   `json:"SKey"`
	ClientVersion       int      `json:"ClientVersion"`
	SystemTime          int      `json:"SystemTime"`
	GrayScale           int      `json:"GrayScale"`
	InviteStartCount    int      `json:"InviteStartCount"`
	MPSubscribeMsgCount int      `json:"MPSubscribeMsgCount"`
	MPSubscribeMsgList  []string `json:"MPSubscribeMsgList"`
	ClickReportInterval int      `json:"ClickReportInterval"`
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
			wc.Login <- true
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
	<-wc.Login
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

	log.Println("+=================================+")
	log.Println(resp.Status)
	log.Println(resp.StatusCode)

	log.Println("Cookies: ", resp.Cookies())

	data, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Println("Error happend when read json: ", err.Error())
		return
	}

	//Pay more attention to here
	log.Println(string(data))
	reader := bytes.NewReader(data)

	var result InitResp
	if err = json.NewDecoder(reader).Decode(&result); err != nil {
		log.Println("Error happened when decode init information: ", err.Error())
		return
	}

	log.Println(result)
	log.Println(resp.Cookies())

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
*/
func (wc *wechat) GetContactList() {
	//uri := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/" + "webwxgetcontact?pass_ticket=" + wc.BaseRequest.PassTicket + "&skey=" + wc.BaseRequest.Skey + "&r=" + strconv.FormatInt(time.Now().Unix(), 10) + "&seq=10"
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

	log.Println("+=================================+")
	log.Println(resp.Status)
	log.Println(resp.StatusCode)

	log.Println("Cookies: ", resp.Cookies())

	data, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Println("Error happend when read json: ", err.Error())
		return
	}

	//Pay more attention to here
	log.Println(string(data))
	reader := bytes.NewReader(data)
	var crsp contactResponse
	if err = json.NewDecoder(reader).Decode(&crsp); err != nil {
		log.Println("Error happened when decode init information: ", err.Error())
		return
	}

	log.Println(crsp)
	log.Println(resp.Cookies())

}

func main() {
	wc := &wechat{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		Login:       make(chan bool),
		BaseRequest: &BaseRequest{},
		DeviceID:    "0x1234567890",
		userAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36",
	}
	wc.GetUUID()
	wc.GetQRCode()
	go wc.WaitForQRCodeScan()
	wc.GetBaseRequest()
	wc.WeChatInit()
	wc.GetContactList()
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
