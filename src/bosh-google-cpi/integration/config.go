package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"bosh-google-cpi/action"
	boshapi "bosh-google-cpi/api"
	boshdisp "bosh-google-cpi/api/dispatcher"
	"bosh-google-cpi/api/transport"
	boshcfg "bosh-google-cpi/config"
	"bosh-google-cpi/google/client"

	boshlogger "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/bosh-utils/uuid"
)

const (
	reusableVMName = "google-cpi-int-tests"
)

var (
	// Provided by user
	googleProject    = os.Getenv("GOOGLE_PROJECT")
	externalStaticIP = os.Getenv("EXTERNAL_STATIC_IP")
	keepResuableVM   = os.Getenv("KEEP_REUSABLE_VM")

	// Configurable defaults
	networkName          = envOrDefault("NETWORK_NAME", "cfintegration")
	customNetworkName    = envOrDefault("CUSTOM_NETWORK_NAME", "cfintegration-custom")
	customSubnetworkName = envOrDefault("CUSTOM_SUBNETWORK_NAME", "cfintegration-custom-us-central1")
	ip                   = envOrDefault("PRIVATE_IP", "192.168.100.102")
	stemcellURL          = envOrDefault("STEMCELL_URL", "https://storage.googleapis.com/evandbrown17/bosh-stemcell-3215-google-kvm-ubuntu-trusty-go_agent-raw.tar.gz")
	existingStemcell     = envOrDefault("EXISTING_STEMCELL", "stemcell-decdea81-a0a3-47b6-5d76-093d505a6de9")
	targetPool           = envOrDefault("TARGET_POOL", "cfintegration")
	backendService       = envOrDefault("BACKEND_SERVICE", "cfintegration")
	instanceGroup        = envOrDefault("BACKEND_SERVICE", "cfintegration-us-central1-a")
	zone                 = envOrDefault("ZONE", "us-central1-a")
	region               = envOrDefault("REGION", "us-central1")

	cfgContent = fmt.Sprintf(`{
	  "google": {
		"project": "%v",
		"default_zone": "%v"
	  },
	  "actions": {
		"agent": {
		  "mbus": "http://127.0.0.1",
		  "blobstore": {
			"type": "local"
		  }
		},
		"registry": {
		  "use_gce_metadata": true
		}
	  }
	}`, googleProject, zone)
)

func execCPI(request string) (boshdisp.Response, error) {
	var err error
	var config boshcfg.Config
	var in, out, errOut, errOutLog bytes.Buffer
	var boshResponse boshdisp.Response
	var googleClient client.GoogleClient

	if config, err = boshcfg.NewConfigFromString(cfgContent); err != nil {
		return boshResponse, err
	}

	multiWriter := io.MultiWriter(&errOut, &errOutLog)
	logger := boshlogger.NewWriterLogger(boshlogger.LevelDebug, multiWriter, multiWriter)
	multiLogger := boshapi.MultiLogger{Logger: logger, LogBuff: &errOutLog}
	uuidGen := uuid.NewGenerator()
	if googleClient, err = client.NewGoogleClient(config.Google, multiLogger); err != nil {
		return boshResponse, err
	}

	actionFactory := action.NewConcreteFactory(
		googleClient,
		uuidGen,
		config.Actions,
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

func envOrDefault(key, defaultVal string) (val string) {
	if val = os.Getenv(key); val == "" {
		val = defaultVal
	}
	return
}
