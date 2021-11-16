package master_conf

import (
	"encoding/json"
	terrors "github.com/pkg/errors"
	"io/ioutil"
)

type Config struct {
	ApiPort               string   `json:"api_port"`
	EtcdEndpoints         []string `json:"etcd_endpoints"`
	EtcdDailTimeout       int      `json:"etcd_dail_timeout"`
	MongodbUri            string   `json:"mongodb_uri"`
	MongodbConnectTimeout int      `json:"mongodb_connect_timeout"`
	WebRoot               string   `json:"web_root"`
}

var (
	GConfig *Config
)

// InitialConfig 加载配置
func InitialConfig(filename string) (err error) {
	// 读配置文件
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return terrors.Wrap(err, "read conf file bytes failed")
	}

	// json 反序列化
	conf := &Config{}
	if err = json.Unmarshal(bytes, conf); err != nil {
		return terrors.Wrap(err, "unmarshal conf file bytes failed")
	}

	// 单例赋值
	GConfig = conf

	return nil
}
