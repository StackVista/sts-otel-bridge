package stackstate

type StackStateClient struct {
	ApiKey   string
	HostName string
	Port     int
	Prefix   string
}

func NewStackStateClient(apiKey string, hostName string, port int, prefix string) *StackStateClient {
	return &StackStateClient{
		ApiKey:   apiKey,
		HostName: hostName,
		Port:     port,
		Prefix:   prefix,
	}
}

// func (c *StackStateClient) SendTopology(topology *TopologyMessage) error {
// 	client := resty.New()

// 	resp, err := client.R().SetBody(topology).SetQueryParam("api_key", c.ApiKey).Post(fmt.Sprintf("https://%s:%s", c.HostName, c.Prefix))
// 	if err != nil {
// 		return err
// 	}

// 	if resp.IsError() {
// 		return fmt.Errorf("error sending topology: %s", resp.Status())
// 	}

// 	return nil
// }
