package bootstrap

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"fmt"
	"io/ioutil"
	"log"
	"time"
)

const (
	// Path to the AWS CA file
	caFilePath = ""

	// Timeout operations after N seconds
	connectTimeout  = 5
	queryTimeout    = 30
	username        = ""
	password        = ""
	clusterEndpoint = ""

	// Which instances to read from
	readPreference = "secondaryPreferred"

	// connectionStringTemplate = "mongodb://%s:%s@%s/test?tls=false&replicaSet=rs0&readpreference=%s"
	connectionStringTemplate = ""
)

func InitMongoDatabase() mongo.Client {

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb://root:example@mongo:27017").
		SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)

	/*
			   _, b, _, _ := runtime.Caller(0)
			   basepath := filepath.Dir(b)


			caFilePath := "rds-combined-ca-bundle.pem"
			tlsConfig, err := getCustomTLSConfig(caFilePath)
			if err != nil {
				log.Fatalf("Failed getting TLS configuration: %v", err)
			}


		dbHost := App.Config.GetString(`mongo.host`)
		dbPort := App.Config.GetString(`mongo.port`)
		dbUser := App.Config.GetString(`mongo.user`)
		dbPass := App.Config.GetString(`mongo.pass`)
		dbName := App.Config.GetString(`mongo.name`)
	*/
	//	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	//mongodbURI := connectionStringTemplate
	//mongodbURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint, readPreference)
	// mongodbURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint, readPreference)

	//if dbUser == "" || dbPass == "" {
	//	mongodbURI = fmt.Sprintf("mongodb://%s:%s/%s", dbHost, dbPort, dbName)
	//}

	//client, err := mongo.NewClient(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Client success")
	log.Print(client)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Ping success")

	return *client
}

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("Failed parsing pem file")
	}

	return tlsConfig, nil

}
