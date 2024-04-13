package model

type Message struct {
	Mqtt      string `json:"mqtt"`
	Invent    string `json:"invent"`
	UnitGUID  string `json:"unit_guid"`
	MsgID     string `json:"msg_id"`
	Text      string `json:"text"`
	Context   string `json:"context"`
	Class     string `json:"class"`
	Level     int    `json:"level"`
	Area      string `json:"area"`
	Addr      string `json:"addr"`
	Block     string `json:"block"`
	Type      string `json:"type"`
	Bit       int    `json:"bit"`
	InvertBit int    `json:"invert_bit"`
}
