package dashTrigger

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
        "encoding/json"
	"net/http"
	//"bytes"
	"github.com/google/gopacket"
	//"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	//"net"
	"fmt"
	"os"
	"io/ioutil"
)

type Button struct {
	Name    string
	Url     string
	Mac     string
}

// MyTriggerFactory My Trigger factory
type MyTriggerFactory struct{
	metadata *trigger.Metadata
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &MyTriggerFactory{metadata:md}
}

//New Creates a new trigger instance for a given id
func (t *MyTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &MyTrigger{metadata: t.metadata, config:config}
}

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
}

// Init implements trigger.Trigger.Init
func (t *MyTrigger) Init(runner action.Runner) {

	//inputs: dash mac (not needed for this example)
		//NIC interface name

	
	//log.Printf("Starting up on interface[%v]...", configuration.Nic)

	h, err := pcap.OpenLive("eth0", 65536, true, pcap.BlockForever)

	if err != nil || h == nil {
		log.Fatalf("Error opening interface: %s\nPerhaps you need to run as root?\n", err)
	}
	defer h.Close()
	
	var filter = "(arp or (udp and src port 68 and dst port 67)) and src host 0.0.0.0";
	err = h.SetBPFFilter(filter)
	if err != nil {
		log.Fatalf("Unable to set filter! %s\n", err)
	}
	log.Println("Listening for Dash buttons...")

	packetSource := gopacket.NewPacketSource(h, h.LinkType())


	t.runner = runner


	for packet := range packetSource.Packets() {

		fmt.Println(packet)
		fmt.Println("DASH PRESSED")

	}
}

// Metadata implements trigger.Trigger.Metadata
func (t *MyTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *MyTrigger) Start() error {
	// start the trigger
	return t.server.Start()
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *MyTrigger) Stop() error {
	// stop the trigger
	return t.server.Stop()
	return nil
}
