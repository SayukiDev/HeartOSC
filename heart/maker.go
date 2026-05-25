package heart

import (
	_ "embed"
	"runtime"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
	"resty.dev/v3"
)

var makerList map[uint16]string

//go:embed company_identifiers.yaml
var makerListBytes []byte

const listUrl = "https://bitbucket.org/bluetooth-SIG/public/raw/a87138721ab82f2b69436603c0534532029be72a/assigned_numbers/company_identifiers/company_identifiers.yaml"

func pullMakerList() error {
	resp, err := resty.New().SetTimeout(time.Second * 20).R().Get(listUrl)
	if err == nil {
		makerListBytes = resp.Bytes()
	}
	var temp = new(struct {
		CompanyIdentifiers []struct {
			Value uint16 `yaml:"value"`
			Name  string `yaml:"name"`
		} `yaml:"company_identifiers"`
	})
	err = yaml.Unmarshal(makerListBytes, &temp)
	if err != nil {
		return err
	}
	makerList = make(map[uint16]string, len(temp.CompanyIdentifiers))
	for _, v := range temp.CompanyIdentifiers {
		makerList[v.Value] = v.Name
	}
	makerListBytes = nil
	runtime.GC()
	return nil
}

func getMakerName(id uint16) string {
	if makerList == nil {
		return strconv.Itoa(int(id))
	}
	if v, ok := makerList[id]; ok {
		return v
	}

	return strconv.Itoa(int(id))
}
