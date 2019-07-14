package elastic

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/olivere/elastic.v6"
)

var esClient *elastic.Client

// InitClient initialize Elasticsearch client.
// Note: This not thread-safe and should be called at server startup.
func InitClient() (*elastic.Client, error) {
	if esClient == nil {
		host := viper.GetString("elasticsearch.host")
		port := viper.GetInt("elasticsearch.port")
		url := fmt.Sprintf("http://%s:%d", host, port)

		log.WithFields(log.Fields{
			"url": url,
		}).Info("Initialize Elasticsearch client")

		// Turn-off sniffing.
		client, err := elastic.NewClient(elastic.SetURL(url),
			elastic.SetSniff(false))
		if err != nil {
			log.WithError(err).Error("Fail to initialize Elasticsearch client")
			return nil, err
		}

		esClient = client
	}

	return esClient, nil
}
