package monkey_patch

import (
	"context"
	"encoding/json"
	"github.com/weilianlow/monkey-patch/util"
	"io/ioutil"
	"os"
	"path"
)

type DataList struct {
	Config  []string        `json:"config"`
	Data    []Data          `json:"data"`
	DataMap map[string]Data `json:"data_map, omitempty"`
}

type Data struct {
	Name    string      `json:"name"`
	Value   interface{} `json:"value"`
	Enabled bool        `json:"enabled"`
	Guard   GuardFunc   `json:"guard, omitempty"`
}

type GuardFunc struct {
	Unpatch func()
	Restore func()
}

func New(ctx context.Context) *DataList {
	var (
		d          = &DataList{}
		logger = util.GetLoggerFromContext(ctx).WithField("func", "New")
		jsonFile   *os.File
		byteValue  []byte
		err        error
	)
	d.DataMap = make(map[string]Data, 0)
	// open json file
	pwd, _ := os.Getwd()
	pathList := []string{path.Join(path.Dir(pwd+"/"), "data.json"),
		path.Join(path.Dir(pwd+"/"), path.Dir("/patch/"), "data.json"),
		path.Join(path.Dir(pwd+"/"), path.Dir("/etc/"), "data.json")}
	for _, filePath := range pathList {
		if jsonFile, err = os.Open(filePath); jsonFile != nil && err == nil {
			break
		}
	}
	if err != nil {
		logger.WithError(err).Error("unable to locate data.json")
		return d
	}
	// open file
	if byteValue, err = ioutil.ReadAll(jsonFile); err != nil {
		logger.WithError(err).Error("unable to read data.json")
		return d
	}
	// unmarshal json
	if err = json.Unmarshal(byteValue, d); err != nil {
		logger.WithError(err).Error("unable to unmarshal data.json")
		return d
	}
	// init data map
	for _, data := range d.Data {
		d.DataMap[data.Name] = data
	}
	return d
}
