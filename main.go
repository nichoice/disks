package main

//"encoding/json"
//"fmt"

//"github.com/rubiojr/go-udisks"

// // need install dbus and  udisks2
// func main() {
// 	client, err := udisks.NewClient()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// List all block devices available to UDisks2
// 	devs, err := client.BlockDevices()
// 	if err != nil {
// 		panic(err)
// 	}
// 	pretty(devs)
// }

// func pretty(dev interface{}) {
// 	prettyString, _ := json.MarshalIndent(dev, "", "  ")
// 	fmt.Println(string(prettyString))
// }

import (
	"encoding/json"
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/rubiojr/go-udisks"
)

type Client struct {
	conn *dbus.Conn
}

func main() {
	c := &Client{}
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		fmt.Println("dbus connection error", err)
		panic(err)
	}
	c.conn = conn
	defer conn.Close()

	/**
	// var filter map[string]interface{}
	// obj := conn.Object("com.redhat.lvmdbus1", "/com/redhat/lvmdbus1")
	// getErr := obj.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&filter)
	// if getErr != nil {
	// 	fmt.Println("get object error", getErr)
	// }

	// fmt.Println(filter)
	**/

	client, err := udisks.NewClient()
	if err != nil {
		panic(err)
	}

	devs, err := client.BlockDevices()
	if err != nil {
		panic(err)
	}
	pretty(devs)

	node, lverr := introspect.Call(conn.Object("com.redhat.lvmdbus1", "/com/redhat/lvmdbus1/Lv"))
	if lverr != nil {
		fmt.Println("sssss", lverr)
	}

	for _, ch := range node.Children {
		path := "/com/redhat/lvmdbus1/Lv/" + ch.Name
		obj := conn.Object("com.redhat.lvmdbus1", dbus.ObjectPath(path))

		pv, bberr := c.buildPv(obj)
		if err != nil {
			fmt.Println("sfsfsdfsf: ", bberr)
		}
		fmt.Println("-----", pv)
	}
}

func pretty(dev interface{}) {
	prettyString, _ := json.MarshalIndent(dev, "", "  ")
	fmt.Println(string(prettyString))
}

func (c *Client) buildPv(objPv dbus.BusObject) (*Pv, error) {
	pv := &Pv{
		VolumeType: "",
		Name:       "",
		Uuid:       "",
		SizeBytes:  0,
	}

	stringProperty("com.redhat.lvmdbus1.LvCommon.VolumeType", objPv, &pv.VolumeType)
	stringProperty("com.redhat.lvmdbus1.LvCommon.Name", objPv, &pv.Name)
	stringProperty("com.redhat.lvmdbus1.LvCommon.Uuid", objPv, &pv.Uuid)
	uint64Property("com.redhat.lvmdbus1.LvCommon.SizeBytes", objPv, &pv.SizeBytes)

	return pv, nil
}

type Pv struct {
	VolumeType string
	Name       string
	Uuid       string
	SizeBytes  uint64
}

func stringProperty(path string, obj dbus.BusObject, p *string) error {
	v, err := obj.GetProperty(path)
	if err != nil {
		return err
	}

	var ok bool
	*p, ok = v.Value().(string)
	if !ok {
		return err
	}

	return nil
}

// func boolProperty(path string, obj dbus.BusObject, p *bool) error {
// 	v, err := obj.GetProperty(path)
// 	if err != nil {
// 		return err
// 	}

// 	var ok bool
// 	*p, ok = v.Value().(bool)
// 	if !ok {
// 		return err
// 	}

// 	return nil
// }

func uint64Property(path string, obj dbus.BusObject, p *uint64) error {
	v, err := obj.GetProperty(path)
	if err != nil {
		return err
	}

	var ok bool
	*p, ok = v.Value().(uint64)
	if !ok {
		return err
	}

	return nil
}
