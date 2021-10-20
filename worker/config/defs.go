package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	EtcdEndpoints         []string `json:"etcd_endpoints"`
	EtcdDailTimeout       int      `json:"etcd_dail_timeout"`
	MongodbUri            string   `json:"mongodb_uri"`
	MongodbConnectTimeout int      `json:"mongodb_connect_timeout"`
	JobLogBatchSize       int      `json:"job_log_batch_size"`
	JobLogCommitTimeout   int      `json:"job_log_commit_timeout"`
}

var (
	GConfig *Config
)

// InitialConfig 加载配置
func InitialConfig(filename string) (err error) {
	// 读配置文件
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// json 反序列化
	conf := &Config{}
	if err = json.Unmarshal(bytes, conf); err != nil {
		return err
	}

	// 单例赋值
	GConfig = conf

	return nil
}
