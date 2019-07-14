package main

import (
    "fmt"
    "os"
    "net/http"
    "strings"

    "rakuten.co.jp/chatroom/controllers"
    "rakuten.co.jp/chatroom/services"
    "rakuten.co.jp/chatroom/clients/mysql"
    "rakuten.co.jp/chatroom/clients/elastic"
    "rakuten.co.jp/chatroom/clients/redis"

	"github.com/gorilla/mux"
    "github.com/spf13/viper"
    log "github.com/sirupsen/logrus"
)

func main() {
    viper.SetConfigName("env")

    env := os.Getenv("ENV")
    switch strings.ToUpper(env) {
        case "PROD":
            viper.AddConfigPath("./configs/prod/")
        case "DEV":
        fallthrough
        default:
            viper.AddConfigPath("./configs/dev/")
    }

    err := viper.ReadInConfig()
    if err != nil {
        log.Error("Fail to load config")
        return
    }

    logLevel := viper.GetString("server.loglevel")
    switch strings.ToLower(logLevel) {
    case "trace":
        log.SetLevel(log.TraceLevel)
    case "warn":
        log.SetLevel(log.WarnLevel)
    case "debug":
        log.SetLevel(log.DebugLevel)
    case "info":
        fallthrough
    default:
        log.SetLevel(log.InfoLevel)
    }

    r := mux.NewRouter().StrictSlash(true)

    for _, route := range controllers.Routes {
        path := fmt.Sprintf("/v%d/%s", route.Version, route.Endpoint)
        r.Methods(route.Method).Path(path).HandlerFunc(route.Handler)
    }

    // Init MySQL client.
    db, err := mysql.InitClient()
    if err != nil {
        return
    }

    // Init Elasticsearch client.
    esClient, err := elastic.InitClient()
    if err != nil {
        return
    }

    // Init Redis client.
    redisClient := redis.InitClient()

    // Init services.
    services.Init(db, esClient, redisClient)

    port := viper.GetInt("server.port")
    serverURL := fmt.Sprintf("0.0.0.0:%d", port)

    log.WithField("url", serverURL).Info("Starting server")

    err = http.ListenAndServe(serverURL, r)
    if err != nil {
        log.WithError(err).Error("Start server failed")
    }
}
