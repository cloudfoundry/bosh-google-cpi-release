package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"bosh-google-cpi/action"
	boshapi "bosh-google-cpi/api"
	boshdisp "bosh-google-cpi/api/dispatcher"
	"bosh-google-cpi/api/transport"
	boshcfg "bosh-google-cpi/config"

	boshlogger "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/bosh-utils/uuid"
)

var (
	// A stemcell that will be created in integration_suite_test.go
	existingStemcell string

	// Provided by user
	googleProject    = envRequired("GOOGLE_PROJECT")
	externalStaticIP = envRequired("EXTERNAL_STATIC_IP")
	serviceAccount   = envRequired("SERVICE_ACCOUNT")

	// Configurable defaults
	stemcellFile                  = envOrDefault("STEMCELL_FILE", "")
	stemcellVersion               = envOrDefault("STEMCELL_VERSION", "")
	networkName                   = envOrDefault("NETWORK_NAME", "cfintegration")
	customNetworkName             = envOrDefault("CUSTOM_NETWORK_NAME", "cfintegration-custom")
	customSubnetworkName          = envOrDefault("CUSTOM_SUBNETWORK_NAME", "cfintegration-custom-us-central1")
	ipAddrs                       = strings.Split(envOrDefault("PRIVATE_IP", "192.168.100.102,192.168.100.103,192.168.100.104"), ",")
	targetPool                    = envOrDefault("TARGET_POOL", "cfintegration")
	backendService                = envOrDefault("BACKEND_SERVICE", "cfintegration")
	regionBackendService          = envOrDefault("REGION_BACKEND_SERVICE", "cfintegration-r")
	collisionBackendService       = envOrDefault("COLLISION_BACKEND_SERVICE", "cfintegration-collision")
	collisionRegionBackendService = envOrDefault("COLLISION_REGION_BACKEND_SERVICE", "cfintegration-collision")
	instanceGroup                 = envOrDefault("BACKEND_SERVICE", "cfintegration")
	ilbInstanceGroup              = envOrDefault("ILB_INSTANCE_GROUP", "cfintegration-ilb")
	zone                          = envOrDefault("ZONE", "us-central1-a")
	region                        = envOrDefault("REGION", "us-central1")
	imageURL                      = envOrDefault("IMAGE_URL", "https://s3.amazonaws.com/bosh-core-stemcells/google/bosh-stemcell-170.23-google-kvm-ubuntu-xenial-go_agent.tgz")

	// Channel that will be used to retrieve IPs to use
	ips chan string

	// If true, CPI will not wait for delete to complete. Speeds up tests significantly.
	asyncDelete = envOrDefault("CPI_ASYNC_DELETE", "true")

	cfgContent = fmt.Sprintf(`{
	  "cloud": {
		"plugin": "google",
		"properties": {
		  "google": {
			"project": "%v"
		  },
		  "agent": {
			"mbus": "http://127.0.0.1",
			"blobstore": {
			  "provider": "local"
			}
		  },
		  "registry": {
			"use_gce_metadata": true
		  }
		}
	  }
	}`, googleProject)
)

func toggleAsyncDelete() {
	key := "CPI_ASYNC_DELETE"
	current := os.Getenv(key)
	if current == "" {
		os.Setenv(key, "true")
	} else {
		os.Setenv(key, "")
	}
}

func execCPI(request string) (boshdisp.Response, error) {
	var err error
	var cfg boshcfg.Config
	var in, out, errOut, errOutLog bytes.Buffer
	var boshResponse boshdisp.Response

	if cfg, err = boshcfg.NewConfigFromString(cfgContent); err != nil {
		return boshResponse, err
	}

	// We're going to convert the Google config to a map[string]interface{}
	googCfg, err := json.Marshal(cfg.Cloud.Properties.Google)
	if err != nil {
		return boshResponse, err
	}
	var ctx map[string]interface{}
	if err = json.Unmarshal(googCfg, &ctx); err != nil {
		return boshResponse, err
	}

	// Unmarshal the reqest string to a struct
	var req boshdisp.Request
	if err = json.Unmarshal([]byte(request), &req); err != nil {
		return boshResponse, err
	}
	req.Context = ctx

	// Marshal the modified request back to string
	requestByte, err := json.Marshal(req)
	if err != nil {
		return boshResponse, err
	}
	request = string(requestByte)

	multiWriter := io.MultiWriter(&errOut, &errOutLog)
	logger := boshlogger.NewWriterLogger(boshlogger.LevelDebug, multiWriter)
	multiLogger := boshapi.MultiLogger{Logger: logger, LogBuff: &errOutLog}
	uuidGen := uuid.NewGenerator()

	actionFactory := action.NewConcreteFactory(
		uuidGen,
		cfg,
		multiLogger,
	)

	caller := boshdisp.NewJSONCaller()
	dispatcher := boshdisp.NewJSON(actionFactory, caller, multiLogger)

	in.WriteString(request)
	cli := transport.NewCLI(&in, &out, dispatcher, multiLogger)

	var response []byte

	if err = cli.ServeOnce(); err != nil {
		return boshResponse, err
	}

	if response, err = ioutil.ReadAll(&out); err != nil {
		return boshResponse, err
	}

	if err = json.Unmarshal(response, &boshResponse); err != nil {
		return boshResponse, err
	}
	return boshResponse, nil
}

func envRequired(key string) (val string) {
	if val = os.Getenv(key); val == "" {
		panic(fmt.Sprintf("Could not find required environment variable '%s'", key))
	}
	return
}

func envOrDefault(key, defaultVal string) (val string) {
	if val = os.Getenv(key); val == "" {
		val = defaultVal
	}
	return
}
