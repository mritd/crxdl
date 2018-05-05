package downloader

import (
	"text/template"

	"net/http"

	"fmt"

	"io/ioutil"

	"github.com/mritd/crxdl/utils"
)

type Downloader struct {
	OsType      string
	Arch        string
	Prod        string
	ProdChannel string
	ProdVersion string
	ExtensionID string
}

const (
	downloadURLTpl = "https://clients2.google.com/service/update2/crx?response=redirect&os={{ .OsType }}&arch={{ .Arch }}&prod={{ .Prod }}&prodchannel={{ .ProdChannel }}&prodversion={{ .ProdVersion }}&x=id%3D{{ .ExtensionID }}%26uc"
	refererTpl     = "https://chrome.google.com/webstore/detail/{{ .ExtensionID }}?hl=en"
)

func (dl Downloader) Download() []byte {

	var buf []byte

	dlTpl, err := template.New("").Parse(downloadURLTpl)
	utils.CheckAndExit(err)
	refTpl, err := template.New("").Parse(refererTpl)
	utils.CheckAndExit(err)

	downloadURL := string(utils.Render(dlTpl, dl))
	referer := string(utils.Render(refTpl, dl))

	client := &http.Client{}
	req, err := http.NewRequest("GET", downloadURL, nil)
	req.Header.Set("Referer", referer)
	fmt.Printf("Dwonloading %s...\n", dl.ExtensionID)
	resp, err := client.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		fmt.Printf("Dwonload %s failed!\n", dl.ExtensionID)
		fmt.Println(err)
	} else if resp.StatusCode != http.StatusOK {
		fmt.Printf("Extension %s not found!\n", dl.ExtensionID)
		return buf
	}

	fmt.Printf("Dwonload %s success!\n", dl.ExtensionID)
	buf, err = ioutil.ReadAll(resp.Body)
	utils.CheckAndExit(err)

	return buf

}
