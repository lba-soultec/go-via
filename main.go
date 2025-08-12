//go:generate bash -c "swag init"
//go:generate bash -c "cd web && rm -rf ./web/dist && npm install --legacy-peer-deps && npm run build && cd .. && statik -src ./web/dist/web -f"

package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/maxiepax/go-via/api"
	"github.com/maxiepax/go-via/config"
	ca "github.com/maxiepax/go-via/crypto"
	"github.com/maxiepax/go-via/db"
	"github.com/maxiepax/go-via/models"
	"github.com/maxiepax/go-via/secrets"
	"github.com/maxiepax/go-via/websockets"

	"github.com/gin-contrib/static"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/maxiepax/go-via/docs"
	_ "github.com/maxiepax/go-via/statik"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
)

var (
	commit = "none"
	date   = "unknown"
)

// @title go-via
// @version 0.1
// @description VMware Imaging Appliances written in GO with full HTTP-REST API

// @BasePath /v1

func main() {

	logServer := websockets.NewLogServer()
	logrus.AddHook(logServer.Hook)
	ConfigureLogger()
	//setup logging
	logrus.WithFields(logrus.Fields{
		"commit": commit,
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

	// TFTPd
	go TFTPd(conf)

	//REST API
	r := gin.New()
	r.Use(cors.Default())

	// ks.cfg is served at top to not place it behind BasicAuth
	r.GET("ks.cfg", api.Ks(key))

	statikFS, err := fs.New()
	if err != nil {
		logrus.Fatal(err)
	}

	r.Use(static.Serve("/", NewMyServeFileSystem(statikFS)))

	r.NoRoute(func(c *gin.Context) {
		logrus.Debugf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	ui := r.Group("/")
	{

		ui.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	v1 := r.Group("/v1")
	{

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

		login := v1.Group("/login")
		{
			login.POST("", api.Login)
		}

		hostConfig := v1.Group("/hostconfig")
		{
			hostConfig.GET("", api.HostConfig)
		}
		ilohosts := v1.Group("/ilohosts")
		{
			ilohosts.POST(":id/setvlanID", api.SetVLANID)      // Set VLAN ID
			ilohosts.POST(":id/start", api.StartIloHost)       // Start the host
			ilohosts.POST(":id/shutdown", api.ShutdownIloHost) // Shutdown the host
			ilohosts.POST(":id/reboot", api.RebootIloHost)     // Reboot the host
			ilohosts.POST(":id/onetimeboot", api.OneTimeBoot)  // Set one time boot

			ilohosts.POST("/checkilo", api.CheckIP) // Check ILO IP
		}
		v1.GET("log", logServer.Handle)

		v1.GET("version", api.Version(commit, date))
	}

	/*	r.GET("postconfig", api.PostConfig) */

	// check if ./cert/server.crt exists, if not we will create the folder, and initiate a new CA and a self-signed certificate
	crt, err := os.Stat("./cert/server.crt")
	if os.IsNotExist(err) {
		// create folder for certificates
		logrus.WithFields(logrus.Fields{
			"certificate": "server.crt does not exist, initiating new CA and creating self-signed ceritificate server.crt",
		}).Info("cert")
		err := os.MkdirAll("cert", os.ModePerm)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("could not create cert directory")
		}
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

// ServeFileSystem implementation that wraps around http.FileSystem
type MyServeFileSystem struct {
	fs http.FileSystem
}

// NewMyServeFileSystem creates a new instance of MyServeFileSystem
func NewMyServeFileSystem(fs http.FileSystem) *MyServeFileSystem {
	return &MyServeFileSystem{fs: fs}
}

// Open implements the http.FileSystem interface
func (fs *MyServeFileSystem) Open(name string) (http.File, error) {
	return fs.fs.Open(name)
}

// Exists implements the Exists method to check if a file exists
func (fs *MyServeFileSystem) Exists(prefix string, path string) bool {
	// Join prefix and path to create the full file path
	fullPath := filepath.Join(prefix, path)

	// Check if the file exists in the wrapped file system
	_, err := fs.fs.Open(fullPath)
	return err == nil // If there's no error, the file exists
}
