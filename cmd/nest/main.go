package main

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"dxkite.cn/log"
	"dxkite.cn/meownest/src/model"
	"dxkite.cn/meownest/src/repository"
	"dxkite.cn/meownest/src/server"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/glebarez/sqlite"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func init() {
	initLogger()
	initBinding()
}

func initLogger() {
	log.SetOutput(log.NewColorWriter(os.Stdout))
	log.SetLogCaller(true)
	log.SetAsync(false)
	log.SetLevel(log.LMaxLevel)
}

func initBinding() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			log.Error("[panic error]", r)
			log.Error(string(buf[:n]))
			name := fmt.Sprintf("crash-%s.log", time.Now().Format("20060102150405"))
			panicErr := string(buf[:n])
			_ = os.WriteFile(name, []byte(panicErr), os.ModePerm)
		}
	}()

	db, err := gorm.Open(sqlite.Open("data.db"))
	if err != nil {
		panic(err)
	}

	db = db.Debug()
	db.AutoMigrate(model.ServerName{}, model.Certificate{})

	certificateRepository := repository.NewCertificate(db)
	certificateService := service.NewCertificate(certificateRepository)
	certificateServer := server.NewCertificate(certificateService)

	nameServerRepo := repository.NewServerName(db)
	serverNameService := service.NewServerName(nameServerRepo, certificateRepository)
	serverNameServer := server.NewServerName(serverNameService)

	httpServer := gin.Default()
	apiV1 := httpServer.Group("/api/v1")

	serverNameApi := apiV1.Group("/server_name")
	{
		serverNameApi.POST("", serverNameServer.Create)
		serverNameApi.GET("/:id", serverNameServer.Get)
	}

	certificateApi := apiV1.Group("/certificate")
	{
		certificateApi.POST("", certificateServer.Create)
	}

	httpServer.Run(":2333")
}
