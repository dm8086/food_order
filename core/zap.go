package core

import (
	"fmt"
	"order_food/core/internal"
	"order_food/global"
	"order_food/utils"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zap 获取 zap.Logger
// Author [SliverHorn](https://github.com/SliverHorn)
func Zap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.GVA_CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.GVA_CONFIG.Zap.Director)
		_ = os.Mkdir(global.GVA_CONFIG.Zap.Director, os.ModePerm)
	}

	cores := internal.Zap.GetZapCores()
	// cores := make([]zapcore.Core, 0)

	mInfoCore, err := NewMongoDBWriteSyncer("jihai-apilog-Info")
	if err == nil {
		infoCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(mInfoCore),
			zapcore.InfoLevel,
		)
		cores = append(cores, infoCore)
	}

	mErrorCore, err := NewMongoDBWriteSyncer("jihai-apilog-Error")
	if err == nil {
		errorCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(mErrorCore),
			zapcore.ErrorLevel,
		)
		cores = append(cores, errorCore)
	}

	mDebugCore, err := NewMongoDBWriteSyncer("jihai-apilog-Debug")
	if err == nil {
		debugCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(mDebugCore),
			zapcore.DebugLevel,
		)
		cores = append(cores, debugCore)
	}

	logger = zap.New(zapcore.NewTee(cores...))

	if global.GVA_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
