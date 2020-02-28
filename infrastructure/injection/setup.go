package injection

import (
	"fmt"
	"git.anphabe.net/event/anphabe-event-hub/app/sourceKeeper"
	"git.anphabe.net/event/anphabe-event-hub/config"
	"git.anphabe.net/event/anphabe-event-hub/domain/service"
	"git.anphabe.net/event/anphabe-event-hub/infrastructure/assets"
	"git.anphabe.net/event/anphabe-event-hub/infrastructure/controller"
	"git.anphabe.net/event/anphabe-event-hub/infrastructure/repository/bowDb"
	"git.anphabe.net/event/anphabe-event-hub/infrastructure/repository/memDb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var (
	conf         *config.ConfigurationInfo
	dbConnection service.DbConnectionInterface
	repoRegistry service.RepositoryRegistryInterface
	keeper       *sourceKeeper.Keeper
	logger       *zap.Logger
	outboundAddr string
)

func init() {
	// init random seed
	rand.Seed(time.Now().UTC().UnixNano())
}

func InitConfig() *config.ConfigurationInfo {
	if nil == conf {
		conf = config.NewConfig().GetConfig()
		GetOutboundAddress(conf)
	}

	return conf
}

func GetOutboundAddress(cfg *config.ConfigurationInfo) string {

	if "" == outboundAddr {
		if nil == cfg {
			cfg = InitConfig()
		}

		conn, err := net.Dial("udp", "8.8.8.8:80")

		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()

		localAddr := conn.LocalAddr().(*net.UDPAddr)
		outboundAddr = "https://" + localAddr.IP.String() + cfg.Server

		_ = qrcode.WriteFile(outboundAddr + "/vueScanAgent", qrcode.Highest, 256, "public/qr.png")
	}

	return outboundAddr
}

func InitLogger(cfg *config.ConfigurationInfo) *zap.Logger {

	if nil == logger {
		logConfig := cfg.Logging

		// First, define our level-handling logic.
		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})

		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})

		jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

		// High-priority output should also go to standard error, and low-priority
		// output should also go to standard out.
		consoleDebuggingOutput := zapcore.Lock(os.Stdout)
		consoleErrorsOutput := zapcore.Lock(os.Stderr)

		// lumberjack.Logger is already safe for concurrent use, so we don't need to lock it.
		fileOutput := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logConfig.Filename,
			MaxSize:    logConfig.MaxSize, // megabytes
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge, // days
		})

		// Join the outputs, encoders, and level-handling functions into
		// zapcore.Cores, then tee the four cores together.
		core := zapcore.NewTee(
			zapcore.NewCore(jsonEncoder, fileOutput, zap.InfoLevel),
			zapcore.NewCore(consoleEncoder, consoleErrorsOutput, highPriority),
			zapcore.NewCore(consoleEncoder, consoleDebuggingOutput, lowPriority),
		)

		// From a zapcore.Core, it's easy to construct a Logger.
		logger = zap.New(core)
	}

	return logger
}

func InitSourceKeeper() *sourceKeeper.Keeper {
	if nil == keeper {
		cfg := InitConfig()
		keeper = sourceKeeper.NewSourceKeeper(cfg.DbSources, InitRepositoryRegistry(cfg), InitLogger(cfg))
	}

	return keeper
}

func InitDBConnection(cfg *config.ConfigurationInfo) service.DbConnectionInterface {
	if nil == dbConnection {
		storageCfg := cfg.Storage
		switch storageCfg.Adapter {
		case "mem":
			dbConnection = memDb.NewMemDbConnection(storageCfg.Folder)
		default:
			dbConnection = bowDb.NewBowDbConnection(storageCfg.Folder)
		}
	}

	return dbConnection
}

func InitRepositoryRegistry(cfg *config.ConfigurationInfo) service.RepositoryRegistryInterface {
	if nil == repoRegistry {
		if nil == cfg {
			cfg = InitConfig()
		}

		repoRegistry = service.NewRepositoryRegistry(InitDBConnection(cfg))
	}

	return repoRegistry
}

func SignalsHandle() <-chan struct{} {
	quit := make(chan struct{})

	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

		defer func() {
			close(signals)
			close(quit)
		}()

		<-signals

		logger.Info("Interrupt signal Received")
	}()

	return quit
}

func InitGin() *gin.Engine {

	router := gin.Default()

	//t, err := loadTemplate()
	//if err != nil {
	//	panic(err)
	//}
	//router.SetHTMLTemplate(t)

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"*"}

	router.Use(cors.New(corsCfg))
	//router.Use(controller.ScanAgentRedirectMiddleWare())

	router.LoadHTMLGlob("templates/*")
	router.Static("/public/", "public")
	router.Static("/vueScanAgent", "vueScanAgent")

	router.NoRoute(func (c *gin.Context) {
		var scanAgentPath = regexp.MustCompile(`^/vueScanAgent/.+$`)

		if scanAgentPath.MatchString(c.Request.URL.Path) {
			c.File("vueScanAgent/index.html")
		}
	})


	//router.StaticFile("/admin/scanner", "scanner/index.html")

	router.GET("/hello", func(c *gin.Context) {controller.ShowHello(c, GetOutboundAddress(nil))})

	keeper := InitSourceKeeper()
	admin := router.Group("/admin")
	{
		admin.GET("/qr-check/:dbName", func(c *gin.Context) {controller.ScanCheckHTML(c, keeper)})
		admin.GET("/qr-check/:dbName/:itemKey", func(c *gin.Context) {controller.ScanCheckHTML(c, keeper)})
		admin.GET("/db/:dbName", func(c *gin.Context) {controller.ShowRepositoryHTML(c)})
		admin.GET("/item/:dbName/:itemKey", func(c *gin.Context) {controller.ShowItemDetailHTML(c, keeper)})
	}

	api := router.Group("/api")
	{
		api.GET("/qr-check/:dbName", func(c *gin.Context) {controller.ScanCheckJSON(c, keeper)})
		api.GET("/qr-check/:dbName/:itemKey", func(c *gin.Context) {controller.ScanCheckJSON(c, keeper)})
		api.GET("/db/:dbName", func(c *gin.Context) {controller.ShowRepositoryJSON(c, keeper)})
		api.GET("/db/:dbName/import", func(c *gin.Context) {controller.StartImport(c, keeper)})
		api.GET("/item/:dbName/:itemKey", func(c *gin.Context) {controller.ShowItemDetailJSON(c, keeper)})
	}

	return router
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range assets.Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}