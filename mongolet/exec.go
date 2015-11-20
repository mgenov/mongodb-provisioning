package mongolet

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// MongoInstance is representing an abstraction to different mongo instances
// such as mongod, mongos and mongoc
type MongoInstance struct {
	workDir  string // working directory
	mongoDir string // mongo directory

	port int

	cmd *exec.Cmd
	// quit channel used for stopping of the instance
	quit chan bool
}

// NewMongoInstance is creating a new mongo instance which will use
// the provided workDir and port.
func NewMongoInstance(workDir string, port int) (*MongoInstance, error) {
	mi := &MongoInstance{}
	mi.quit = make(chan bool)
	mi.port = port

	err := initWorkSpace(workDir)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("mongod", "--dbpath", "data", "--logpath", "db.log", "--fork", "--port", fmt.Sprintf("%d", port))
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	mi.cmd = cmd
	go func() {
		err = cmd.Run()

		if err != nil {
			mi.quit <- true
		}
	}()

	go mi.run()

	return mi, nil
}

func initWorkSpace(workDir string) error {
	dataDir := filepath.Join(workDir, "data")

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	return nil
}

func (m MongoInstance) Stop() {
	m.quit <- true
}

func (m MongoInstance) run() {
	for {
		select {
		case <-m.quit:
			err := m.cmd.Process.Kill()
			// (mgenov): failed to kill process? It cannot be done so much
			if err != nil {
				log.Println(err)
			}
			break
		}
	}
}
