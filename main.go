//go:generate bash -c "go get github.com/swaggo/swag/cmd/swag && swag init"
//go:generate bash -c "cd web && rm -rf ./web/dist && npm install --legacy-peer-deps && npm run build && cd .. && go get github.com/rakyll/statik && statik -src ./web/dist -f"

package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/maxiepax/go-via/api"
	"github.com/maxiepax/go-via/config"
	ca "github.com/maxiepax/go-via/crypto"
	"github.com/maxiepax/go-via/db"
	"github.com/maxiepax/go-via/dhcpd"
	"github.com/maxiepax/go-via/models"
	"github.com/maxiepax/go-via/secrets"
	"github.com/maxiepax/go-via/websockets"
	"github.com/rakyll/statik/fs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"

	_ "github.com/maxiepax/go-via/docs"
	_ "github.com/maxiepax/go-via/statik"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// @title go-via
// @version 0.1
// @description VMware Imaging Appliances written in GO with full HTTP-REST API

// @BasePath /v1

func main() {

	logServer := websockets.NewLogServer()
	logrus.AddHook(logServer.Hook)
	logrus.WithFields(logrus.Fields{
		"version": version,
		"commit":  commit,
		"date":    date,
	}).Infof("Startup")

	// load config file
	conf := config.Load()

	//connect to database
	if conf.Debug {
		db.Connect(true)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		db.Connect(false)
		gin.SetMode(gin.ReleaseMode)
	}

	//migrate all models
	err := db.DB.AutoMigrate(&models.Pool{}, &models.Host{}, &models.Option{}, &models.DeviceClass{}, &models.Group{}, &models.Image{}, &models.User{})
	if err != nil {
		logrus.Fatal(err)
	}

	//create admin user if it doesn't exist
	var adm models.User
	hp := api.HashAndSalt([]byte("VMware1!"))
	if res := db.DB.Where(models.User{UserForm: models.UserForm{Username: "admin"}}).Attrs(models.User{UserForm: models.UserForm{Password: hp}}).FirstOrCreate(&adm); res.Error != nil {
		logrus.Warning(res.Error)
	}

	// load secrets key
	key := secrets.Init()

	// DHCPd
	if !conf.DisableDhcp {
		for _, v := range conf.Network.Interfaces {
			go dhcpd.Init(v)
		}
	}

	// TFTPd
	go TFTPd(conf)

	//REST API
	r := gin.New()
	r.Use(cors.Default())

	statikFS, err := fs.New()
	if err != nil {
		logrus.Fatal(err)
	}

	// ks.cfg is served at top to not place it behind BasicAuth
	r.GET("ks.cfg", api.Ks(key))

	// middleware to check if user is logged in
	r.Use(func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth {
			logrus.WithFields(logrus.Fields{
				"login": "unauthorized request",
			}).Info("auth")
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//get the user that is trying to authenticate
		var user models.User
		if res := db.DB.Select("username", "password").Where("username = ?", username).First(&user); res.Error != nil {
			logrus.WithFields(logrus.Fields{
				"username": username,
				"status":   "supplied username does not exist",
			}).Info("auth")
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//check if passwords match
		if api.ComparePasswords(user.Password, []byte(password), username) {
			logrus.WithFields(logrus.Fields{
				"username": username,
				"status":   "successfully authenticated",
			}).Debug("auth")
		} else {
			logrus.WithFields(logrus.Fields{
				"username": username,
				"status":   "invalid password supplied",
			}).Info("auth")
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	})

	r.NoRoute(func(c *gin.Context) {
		c.Request.URL.Path = "/web/" // force us to always return index.html and not the requested page to be compatible with HTML5 routing
		http.FileServer(statikFS).ServeHTTP(c.Writer, c.Request)
	})

	ui := r.Group("/")
	{
		ui.GET("/web/*all", gin.WrapH(http.FileServer(statikFS)))

		ui.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	v1 := r.Group("/v1")
	{
		//v1.GET("log", logServer.Handle)

		pools := v1.Group("/pools")
		{
			pools.GET("", api.ListPools)
			pools.GET(":id", api.GetPool)
			pools.POST("/search", api.SearchPool)
			pools.POST("", api.CreatePool)
			pools.PATCH(":id", api.UpdatePool)
			pools.DELETE(":id", api.DeletePool)
		}

		relay := v1.Group("/relay")
		{
			relay.GET(":relay", api.GetPoolByRelay)
		}

		hosts := v1.Group("/hosts")
		{
			hosts.GET("", api.ListHosts)
			hosts.GET(":id", api.GetHost)
			hosts.POST("/search", api.SearchHost)
			hosts.POST("", api.CreateHost)
			hosts.PATCH(":id", api.UpdateHost)
			hosts.DELETE(":id", api.DeleteHost)
		}

		options := v1.Group("/options")
		{
			options.GET("", api.ListOptions)
			options.GET(":id", api.GetOption)
			options.POST("/search", api.SearchOption)
			options.POST("", api.CreateOption)
			options.PATCH(":id", api.UpdateOption)
			options.DELETE(":id", api.DeleteOption)
		}

		deviceClass := v1.Group("/device_classes")
		{
			deviceClass.GET("", api.ListDeviceClasses)
			deviceClass.GET(":id", api.GetDeviceClass)
			deviceClass.POST("/search", api.SearchDeviceClass)
			deviceClass.POST("", api.CreateDeviceClass)
			deviceClass.PATCH(":id", api.UpdateDeviceClass)
			deviceClass.DELETE(":id", api.DeleteDeviceClass)
		}

		groups := v1.Group("/groups")
		{
			groups.GET("", api.ListGroups)
			groups.GET(":id", api.GetGroup)
			groups.POST("", api.CreateGroup(key))
			groups.PATCH(":id", api.UpdateGroup(key))
			groups.DELETE(":id", api.DeleteGroup)
		}

		images := v1.Group("/images")
		{
			images.GET("", api.ListImages)
			images.GET(":id", api.GetImage)
			images.POST("", api.CreateImage)
			images.PATCH(":id", api.UpdateImage)
			images.DELETE(":id", api.DeleteImage)
		}

		users := v1.Group("/users")
		{
			users.GET("", api.ListUsers)
			users.GET(":id", api.GetUser)
			users.POST("", api.CreateUser)
			users.PATCH(":id", api.UpdateUser)
			users.DELETE(":id", api.DeleteUser)
		}

		postconfig := v1.Group("/postconfig")
		{
			postconfig.GET("", api.PostConfig(key))
			postconfig.GET(":id", api.PostConfigID(key))
		}

		v1.GET("log", logServer.Handle)

		v1.GET("version", api.Version(version, commit, date))
	}

	/*	r.GET("postconfig", api.PostConfig) */

	// check if ./cert/server.crt exists, if not we will create the folder, and initiate a new CA and a self-signed certificate
	crt, err := os.Stat("./cert/server.crt")
	if os.IsNotExist(err) {
		// create folder for certificates
		logrus.WithFields(logrus.Fields{
			"certificate": "server.crt does not exist, initiating new CA and creating self-signed ceritificate server.crt",
		}).Info("cert")
		os.MkdirAll("cert", os.ModePerm)
		ca.CreateCA()
		ca.CreateCert("./cert", "server", "server")
	} else {
		logrus.WithFields(logrus.Fields{
			crt.Name(): "server.crt found",
		}).Info("cert")
	}
	//enable HTTPS
	listen := ":" + strconv.Itoa(conf.Port)
	logrus.WithFields(logrus.Fields{
		"port": listen,
	}).Info("Webserver")
	err = r.RunTLS(listen, "./cert/server.crt", "./cert/server.key")

	logrus.WithFields(logrus.Fields{
		"error": err,
	}).Error("Webserver")

}
