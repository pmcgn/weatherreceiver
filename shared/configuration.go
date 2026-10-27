package shared

type Configuration struct {
	BrokerUrl string
	Prefix    string
	HttpPort  string
	Retained  bool
	Qos       byte
}
