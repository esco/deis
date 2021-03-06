package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/utils"
)

func runDeisLoggerTest(
	t *testing.T, testSessionUID string, etcdPort string, servicePort string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	err := dockercliutils.BuildImage(t, "../", "deis/logger:"+testSessionUID)
	if err != nil {
		t.Fatal(err)
	}
	dockercliutils.RunDeisDataTest(t, "--name", "deis-logger-data",
		"-v", "/var/log/deis", "deis/base", "true")
	ipaddr := utils.GetHostIPAddress()
	done <- true
	go func() {
		<-done
		err = dockercliutils.RunContainer(cli,
			"--name", "deis-logger-"+testSessionUID,
			"--rm",
			"-p", servicePort+":514/udp",
			"-e", "PUBLISH="+servicePort,
			"-e", "HOST="+ipaddr,
			"-e", "ETCD_PORT="+etcdPort,
			"--volumes-from", "deis-logger-data",
			"deis/logger:"+testSessionUID)
	}()
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "deis-logger running")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogger(t *testing.T) {
	var testSessionUID = utils.NewUuid()
	etcdPort := utils.GetRandomPort()
	servicePort := utils.GetRandomPort()
	fmt.Println("UUID for the session logger Test :" + testSessionUID)
	dockercliutils.RunEtcdTest(t, testSessionUID, etcdPort)
	fmt.Println("starting logger component test:")
	runDeisLoggerTest(t, testSessionUID, etcdPort, servicePort)
	dockercliutils.DeisServiceTest(
		t, "deis-logger-"+testSessionUID, servicePort, "udp")
	dockercliutils.ClearTestSession(t, testSessionUID)
}
