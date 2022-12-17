package services

import (
	"fmt"
	"io/fs"
	"net"
	"os/exec"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/oxodao/photobooth/config"
	"github.com/oxodao/photobooth/models"
	"github.com/oxodao/photobooth/orm"
	"github.com/oxodao/photobooth/utils"
)

var GET *Provider

type Provider struct {
	Sockets    Sockets
	MqttClient *mqtt.Client

	Admin      Admin
	Photobooth Photobooth

	WebappFS   *fs.FS
	AdminappFS *fs.FS
}

func (p *Provider) GetFrontendSettings() *models.FrontendSettings {

	settings := models.FrontendSettings{
		AppState:     p.Photobooth.CurrentState,
		Photobooth:   config.GET.Photobooth,
		DebugDisplay: p.Photobooth.DisplayDebug,
		CurrentMode:  p.Photobooth.CurrentMode,

		IPAddress:  map[string][]string{},
		KnownModes: config.MODES,
	}

	events, err := orm.GET.Events.GetEvents()
	if err != nil {
		fmt.Println("Failed to get events: ", err)
		return nil
	}

	settings.KnownEvents = events

	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		shouldSkip := false
		for _, ignored := range []string{"lo", "br-", "docker", "vmnet", "veth"} { // Ignoring docker / vmware networks for in-dev purposes
			if strings.HasPrefix(inter.Name, ignored) {
				shouldSkip = true
				break
			}
		}

		if shouldSkip {
			continue
		}

		settings.IPAddress[inter.Name] = []string{}

		addrs, _ := inter.Addrs()
		for _, ip := range addrs {
			settings.IPAddress[inter.Name] = append(settings.IPAddress[inter.Name], ip.String())
		}
	}

	return &settings
}

func SetSystemDate(newTime time.Time) error {
	_, lookErr := exec.LookPath("sudo")
	if lookErr != nil {
		fmt.Printf("Sudo binary not found, cannot set system date: %s\n", lookErr.Error())
		return lookErr
	} else {
		dateString := newTime.Format("2 Jan 2006 15:04:05")
		fmt.Printf("Setting system date to: %s\n", dateString)
		args := []string{"date", "--set", dateString}
		return exec.Command("sudo", args...).Run()
	}
}

func (prv *Provider) loadState() error {
	state, err := orm.GET.AppState.GetState()
	if err != nil {
		return err
	}

	if state.CurrentEvent != nil {
		evt, err := orm.GET.Events.GetEvent(*state.CurrentEvent)
		if err != nil {
			return err
		}

		state.CurrentEventObj = evt
	}

	prv.Photobooth.CurrentState = state

	return nil
}

func Load(webapp, adminapp *fs.FS) error {
	for _, folder := range []string{"images"} {
		if err := utils.MakeOrCreateFolder(folder); err != nil {
			return err
		}
	}

	err := orm.Load()
	if err != nil {
		return err
	}

	opts := mqtt.NewClientOptions().AddBroker(config.GET.Mosquitto.Address).SetClientID("photobooth").SetPingTimeout(10 * time.Second).SetKeepAlive(10 * time.Second)
	opts.SetAutoReconnect(true).SetMaxReconnectInterval(10 * time.Second)
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		fmt.Printf("[MQTT] Connection lost: %s\n" + err.Error())
	})
	opts.SetReconnectingHandler(func(c mqtt.Client, options *mqtt.ClientOptions) {
		fmt.Println("[MQTT] Reconnecting...")
	})

	prv := &Provider{
		Sockets:    []*Socket{},
		WebappFS:   webapp,
		AdminappFS: adminapp,
	}

	prv.Admin = Admin{prv: prv}
	prv.Photobooth = Photobooth{
		prv:         prv,
		CurrentMode: config.GET.DefaultMode,
	}

	err = prv.loadState()
	if err != nil {
		return err
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	initButtonHandler(&prv.Photobooth)

	actions := map[string]mqtt.MessageHandler{
		"photobooth/button_press":   BPH.OnButtonPress,
		"photobooth/sync":           prv.Photobooth.OnSyncRequested,
		"photobooth/admin/set_mode": prv.Admin.OnSetMode,
	}

	for topic, action := range actions {
		if token := client.Subscribe(topic, 2, action); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	prv.MqttClient = &client
	GET = prv

	return nil
}

func (p *Provider) Shutdown() error {
	return orm.GET.DB.Close()
}
