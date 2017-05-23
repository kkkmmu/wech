package spider

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	//"github.com/guotie/gogb2312"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Spider struct {
	name   string
	client *http.Client
}

func NewSpider(name string) *Spider {
	return &Spider{
		name: name,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (s *Spider) SpideContent(page, rule string) ([]string, error) {

	resp, err := s.client.Get(page)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	/*
		//Need convert the encode to UTF-8
		body, err, _, _ = gogb2312.ConvertGB2312(body)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	*/

	reader := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Cannot create Docuement by Response")
	}

	var res = make([]string, 0)
	doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		res = append(res, sl.Text())
	})
	return res, nil
}

func (s *Spider) SpideAttribute(page, rule, attr string) ([]string, error) {
	resp, err := s.client.Get(page)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	reader := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Cannot create Docuement by Response")
	}

	var res = make([]string, 0)
	doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		attr, ok := sl.Attr(attr)
		if ok {
			res = append(res, attr)
		}
	})
	return res, nil
}

func (s *Spider) SpideHTML(page, rule string) ([]string, error) {
	resp, err := s.client.Get(page)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	reader := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Cannot create Docuement by Response")
	}

	var res = make([]string, 0)
	doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		content, _ := sl.Html()
		res = append(res, content)
	})
	return res, nil
}
