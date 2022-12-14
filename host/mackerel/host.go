package mackerel

import (
	"errors"
	"os"

	"github.com/mackerelio-labs/mackerel-lambda-extension-agent/host"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/sirupsen/logrus"
)

type Host struct {
	client *mackerel.Client
	ID     string
}

var _ host.Host = &Host{}

var Logger *logrus.Entry

const hostIDFilePath = "/tmp/mackerel-lambda-extension-agent.id"

type CreateHostParam struct {
	MackerelApiKey string
	RoleFullnames  []string
	FunctionName   string
	EnvironmentID  string
}
type CreateOrGetHostParam = CreateHostParam

func CreateOrGetHost(param *CreateOrGetHostParam) (*Host, error) {
	host, err := GetHost(param.MackerelApiKey)
	if err != nil {
		return nil, err
	}
	if host != nil {
		Logger.Info("host already exists. hostID =", host.ID)
		return host, nil
	}

	host, err = CreateHost(param)
	if err != nil {
		return nil, err
	}
	Logger.Info("created a new host. hostID =", host.ID)

	return host, nil
}

func CreateHost(param *CreateHostParam) (*Host, error) {
	if param.MackerelApiKey == "" {
		return nil, errors.New("MackerelApiKey is not set")
	}
	client := mackerel.NewClient(param.MackerelApiKey)

	hostID, err := client.CreateHost(&mackerel.CreateHostParam{
		Name:          param.EnvironmentID,
		DisplayName:   param.FunctionName,
		Memo:          "",
		Meta:          mackerel.HostMeta{},
		Interfaces:    []mackerel.Interface{},
		RoleFullnames: param.RoleFullnames,
		Checks:        []mackerel.CheckConfig{},
	})
	if err != nil {
		return nil, err
	}
	host := &Host{
		client: client,
		ID:     hostID,
	}
	return host, nil
}

func GetHost(apiKey string) (*Host, error) {
	if apiKey == "" {
		return nil, errors.New("MackerelApiKey is not set")
	}
	client := mackerel.NewClient(apiKey)

	_, err := os.Stat(hostIDFilePath)
	if err != nil {
		return nil, nil
	}

	storedHostIDBytes, err := os.ReadFile(hostIDFilePath)
	if err != nil {
		return nil, nil
	}

	host := &Host{
		client: client,
		ID:     string(storedHostIDBytes),
	}
	return host, nil
}

func (h *Host) Retire() error {
	Logger.Info("retiring the host")
	return h.client.RetireHost(h.ID)
}

func (h *Host) CreateGraphDefs() error {
	Logger.Info("creating graph defs")
	return h.client.CreateGraphDefs(GraphDefs)
}

func (h *Host) PostMetrics(metrics []*mackerel.MetricValue) error {
	Logger.Info("posting metrics")
	return h.client.PostHostMetricValuesByHostID(h.ID, metrics)
}
