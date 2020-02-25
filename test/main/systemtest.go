// +build test

package main

import (
	"fmt"
	core "github.com/nuts-foundation/nuts-go-core"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const vendorId = "urn:oid:1.3.6.1.4.1.54851.4:00000001"

var serverHttpAddress string
var serverDataDirectory string
var serverEventsDirectory string
var serverCommand *exec.Cmd
var workingDir string
var clientEnvs map[string]string

func main() {
	var err error
	workingDir, err = os.Getwd()
	if err != nil {
		logrus.Fatalf("Unable to get working directory: %v", err)
		return
	}
	serverDataDirectory, err = ioutil.TempDir("", "nuts-systemtest")
	if err != nil {
		logrus.Fatalf("Unable to create temp test directory: %v", err)
		return
	}
	defer func() {
		logrus.Info("Cleaning up data directory")
		os.RemoveAll(serverDataDirectory)
	}()
	serverEventsDirectory = filepath.Join(serverDataDirectory, "events")

	serverHttpPort, err := getFreePort()
	if err != nil {
		panic(err)
	}
	serverHttpAddress = fmt.Sprintf("localhost:%d", serverHttpPort)

	go func() {
		err := startServer()
		logrus.Infof("Server stopped (%v)", err)
	}()

	waitForServerRunning()
	defer shutdownServer()

	clientEnvs = map[string]string{
		"NUTS_MODE":             core.GlobalCLIMode,
		"NUTS_ADDRESS":          serverHttpAddress,
		"NUTS_REGISTRY_ADDRESS": serverHttpAddress,
	}

	const vendorName = "BecauseWeCare B.V."
	const organisationId = "1232456"
	const organisationName = "Kunstgebit Thuiszorg"
	if err := execClientCommand("register-vendor", vendorId, vendorName); err != nil {
		logrus.Fatalf("Error while registering vendor: %v", err)
	}
	if err := execClientCommand("vendor-claim", vendorId, organisationId, organisationName); err != nil {
		logrus.Fatalf("Error while claiming organisation: %v", err)
	}
	if err := execClientCommand("register-endpoint", organisationId, "AwesomeEndpoint", "http://foobar.nl"); err != nil {
		logrus.Fatalf("Error while registering endpoint: %v", err)
	}
}

func execClientCommand(args ...string) error {
	cliCmd := exec.Command(filepath.Join(workingDir, "nuts-go"))
	cliCmd.Env = envMapToSlice(clientEnvs)
	cliCmd.Stdout = os.Stdout
	cliCmd.Stderr = os.Stderr
	cliCmd.Args = append([]string{"nuts", "registry"}, args...)
	return cliCmd.Run()
}

func shutdownServer() {
	err := serverCommand.Process.Signal(os.Kill)
	if err != nil {
		logrus.Fatalf("Unable to shutdown server: %v", err)
		return
	}
	serverCommand.Process.Wait()
}

func waitForServerRunning() {
	for ; ; {
		if isServerRunning() {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	logrus.Info("Server running!")
}

func isServerRunning() bool {
	response, err := http.Get(fmt.Sprintf("http://%s/status/diagnostics", serverHttpAddress))
	if err != nil {
		logrus.Debugf("Is server running check failed: %v", err)
		return false
	}
	return response.StatusCode == http.StatusOK
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func startServer() error {
	serverNatsPort, err := getFreePort()
	if err != nil {
		return err
	}
	var serverEnvs = map[string]string{
		"NUTS_IDENTITY":                       vendorId,
		"NUTS_ADDRESS":                        serverHttpAddress,
		"NUTS_EVENTS_NATSPORT":                fmt.Sprintf("%d", serverNatsPort),
		"NUTS_AUTH_ACTINGPARTYCN":             "Foo",
		"NUTS_AUTH_PUBLICURL":                 "Bar",
		"NUTS_AUTH_IRMACONFIGPATH":            filepath.Join(workingDir, "testdata", "irma"),
		"NUTS_AUTH_SKIPAUTOUPDATEIRMASCHEMAS": "true",
		"NUTS_REGISTRY_DATADIR":               serverDataDirectory,
		"NUTS_CRYPTO_FSPATH":                  filepath.Join(serverDataDirectory, "keys"),
	}
	serverCommand = exec.Command(filepath.Join(workingDir, "nuts-go"))
	serverCommand.Env = envMapToSlice(serverEnvs)
	serverCommand.Stdout = os.Stdout
	serverCommand.Stderr = os.Stderr
	logrus.Info("Starting server")
	return serverCommand.Run()
}

func envMapToSlice(serverEnvs map[string]string) []string {
	var env []string
	for key, value := range serverEnvs {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}
