package mongolet

import (
	"log"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
)

type config struct {
	Storage storage `yaml:"storage,omitempty"`
}

type storage struct {
	DbPath         string `yaml:"dbPath,omitempty"`
	DirectoryPerDB bool   `yaml:"directoryPerDB,omitempty"`
}

func TestMongoInstance(t *testing.T) {
	m1, err := NewMongoInstance("/Users/mgenov/temp/db1", 27017)
	check(err, t)

	m2, err := NewMongoInstance("/Users/mgenov/temp/db2", 27018)
	check(err, t)

	time.Sleep(time.Second)

	m1.Stop()
	m2.Stop()
}

func TestMarshallYaml(t *testing.T) {
	c := &config{}
	c.Storage = storage{"test", true}

	b, err := yaml.Marshal(c)
	check(err, t)

	log.Println(string(b))
}

func TestAnotherThing(t *testing.T) {
	var data = `
storage:
    dbPath: "/data/db"
    directoryPerDB: true
    journal:
        enabled: true
systemLog:
    destination: file
    path: "/data/db/mongodb.log"
    logAppend: true
    timeStampFormat: iso8601-utc
processManagement:
    fork: true
net:
    bindIp: 127.0.0.1
    port: 27017
    wireObjectCheck : false
    unixDomainSocket: 
        enabled : true
`
	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal([]byte(data), &m)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println(m["systemLog"])
}

func check(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
