package boot

import (
	"github.com/995933447/confloader"
	"github.com/gzjjyz/logger"
	"github.com/gzjjyz/srvlib/utils"
	"time"
)

type RpcSrvCfg struct {
	Ip    string `json:"ip"`
	Port  int    `json:"port"`
	Pport int    `json:"pport"`
}

func InitDynamicCfg(cfgFilePath string, cfg interface{}) error {
	refreshCfgInterval := time.Second * 10
	cfgLoader := confloader.NewLoader(cfgFilePath, refreshCfgInterval, cfg)
	if err := cfgLoader.Load(); err != nil {
		return err
	}

	errCh := make(chan error)
	utils.ProtectGo(func() {
		for {
			select {
			case err := <-errCh:
				logger.Errorf(err.Error())
			}
		}
	})

	utils.ProtectGo(func() {
		cfgLoader.WatchToLoad(errCh)
	})

	return nil
}
