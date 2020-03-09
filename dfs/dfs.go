package dfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

type Config struct {
	Master []string
	Volume []string
}

var (
	_ json.Marshaler   = (*Dfs)(nil)
	_ json.Unmarshaler = (*Dfs)(nil)
)

type Dfs struct {
	content  []byte
	dfsIndex struct {
		FID  string `json:"fid"`
		URL  string `json:"url"`
		Size int    `json:"size"`
	}
}

func (dfs *Dfs) MarshalJSON() ([]byte, error) {
	return json.Marshal(dfs.dfsIndex)
}

func (dfs *Dfs) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &dfs.dfsIndex)
}

func NewDfsWriter(content []byte) (*Dfs, error) {
	if len(content) == 0 {
		return nil, errors.New("Dfs Nil Content")
	}
	dfs := new(Dfs)
	dfs.content = content
	return dfs, nil
}

func NewDfsReader(dfsJson []byte) (*Dfs, error) {
	dfs := new(Dfs)
	err := json.Unmarshal(dfsJson, dfs)
	if err != nil {
		return nil, errors.Wrapf(err, "Json Unmarshal Dfs Index %s",
			string(dfsJson))
	}
	return dfs, nil
}

func NewDfsDeleter(dfsJson []byte) (*Dfs, error) {
	return NewDfsReader(dfsJson)
}

func (dfs *Dfs) Read() (content []byte, err error) {
	if dfs.dfsIndex.URL == "" || dfs.dfsIndex.FID == "" {
		return nil, errors.New("Dfs Reader URL or FID is Nil")
	}
	getterRequest, err := sling.New().
		Get(fmt.Sprintf("http://%s/%s",
			dfs.dfsIndex.URL, dfs.dfsIndex.FID)).
		Request()
	if err != nil {
		return nil, errors.Wrap(err,
			"Dfs Reader New sling request")
	}
	rsp, err := new(http.Client).Do(getterRequest)
	if err != nil {
		return nil, errors.Wrap(err,
			"Dfs Reader request")
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Dfs Reader response code %v",
			rsp.Status)
	}

	content, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrap(err,
			"Dfs Reader read response body")
	}
	dfs.content = content
	return content, nil
}

func (dfs *Dfs) Write() (dfsJson []byte, err error) {
	if len(dfs.content) == 0 {
		return nil, errors.New("Dfs Writer content is Nil")
	}
	if len(constConfig.Master) == 0 {
		return nil, errors.Errorf("Dfs Error Nil Config %+v",
			constConfig)
	}

	assignRequest, err := sling.New().
		Get(fmt.Sprintf("http://%s/dir/assign",
			constConfig.Master[0])).
		Request()

	if err != nil {
		return nil, errors.Wrap(err,
			"Dfs Writer New sling assign-request")
	}
	assignRsp, err := new(http.Client).Do(assignRequest)
	if err != nil {
		return nil, errors.Wrap(err,
			"Dfs Writer assign-request")
	}
	defer assignRsp.Body.Close()
	assignBody, err := ioutil.ReadAll(assignRsp.Body)
	if err != nil {
		return nil, errors.Wrap(err,
			"Dfs Writer read assign-reponse")
	}
	assignBodyStr := string(assignBody)

	assign := &struct {
		FID       string `json:"fid"`
		URL       string `json:"url"`
		Name      string `json:"name"`
		Size      int    `json:"count"`
		PublicUrl string `json:"publicUrl"`
	}{}

	err = json.Unmarshal(assignBody, assign)
	if err != nil {
		return nil, errors.Wrapf(err,
			"Dfs Writer Unmarshal assign %s",
			assignBodyStr)
	}

	dfs.dfsIndex.FID = assign.FID
	dfs.dfsIndex.URL = assign.URL
	dfs.dfsIndex.Size = len(dfs.content)

	body := bytes.NewBuffer(nil)
	bodyWriter := multipart.NewWriter(body)
	bodyWriter.WriteField("file", string(dfs.content))
	err = bodyWriter.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Parse Multipart")
	}

	writerRequest, err := sling.New().
		Post(fmt.Sprintf("http://%s/%s",
			assign.URL, assign.FID)).
		Set("Content-Type", bodyWriter.FormDataContentType()).
		Body(body).
		Request()
	if err != nil {
		return nil, errors.Wrap(err,
			"Dfs Writer New sling writer-request")
	}
	writerRsp, err := new(http.Client).Do(writerRequest)
	if err != nil {
		return nil, errors.Wrap(err,
			"Dfs Writer writer-request")
	}
	defer writerRsp.Body.Close()

	if writerRsp.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("Dfs Writer Writing %v",
			writerRsp.Status)
	}
	dfsJson, _ = json.Marshal(dfs)
	return dfsJson, nil
}

func (dfs *Dfs) Del() error {
	if dfs.dfsIndex.URL == "" || dfs.dfsIndex.FID == "" {
		return errors.New("Dfs Deleter URL or FID is Nil")
	}

	url := fmt.Sprintf("http://%s/%s",
		dfs.dfsIndex.URL, dfs.dfsIndex.FID)
	request, err := sling.New().
		Delete(url).
		Request()
	if err != nil {
		return errors.Wrap(err,
			"Dfs Deleter New sling request")
	}
	rsp, err := new(http.Client).Do(request)
	if err != nil {
		return errors.Wrapf(err, "Dfs Deleter request %s",
			url)
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return errors.Errorf("Dfs Deleter response %v",
			rsp.Status)
	}
	return nil
}
