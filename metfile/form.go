package metfile

//go:generate stringer -type=TagCode,TagValueType -output=strings.go

import (
	"encoding/binary"
	"net"
)

var Endian = binary.LittleEndian
var Extensions = []string{"met"}

type MetForm struct {
	Servers []*MetServer `json:"servers"`
}

type MetServer struct {
	IP   net.IP          `json:"ip"`
	Port int             `json:"port"`
	Tags []*MetServerTag `json:"tags"`
}
type TagValueType uint8

const (
	TagValueTypeString TagValueType = 0x2
	TagValueTypeUint32 TagValueType = 0x3
)

type TagCode uint8

const (
	TagCodeServerName  TagCode = 0x01 // Name of the server
	TagCodeDescription TagCode = 0x0B // Short description about the server
	TagCodePing        TagCode = 0x0C // Time (in ms) it takes to communicate with the server
	TagCodeFail        TagCode = 0x0D // How many times connecting to the server failed
	TagCodePreference  TagCode = 0x0E // Priority given to this server among the others (Normal=0, High=1, Low=2)

	TagCodeDNS       TagCode = 0x85 // DNS of the server
	TagCodeMaxUsers  TagCode = 0x87 // Maximum number of users the server allows to simoultaneously connect to it
	TagCodeSoftFiles TagCode = 0x88 // Soft files number
	TagCodeHardFiles TagCode = 0x89 // Hard files number
	TagCodeLastPing  TagCode = 0x90 // Last time the server was pinged
	TagCodeVersion   TagCode = 0x91 // Version and name of the software the server is running to support the ed2k network
	// UDP flags (0x92)	Unsigned 32 bits number	Informs of the actions the server accepts through UDP connections. This flags are:
	// 0x01: Get sources
	// 0x02: Get files
	// 0x08: New tags
	// 0x10: Unicode
	// 0x20: Get extended sources info
	TagCodeUDPFlags           TagCode = 0x92
	TagCodeAuxiliaryPortsList TagCode = 0x93 // Some servers have additional ports open for those users who cannot connect to the standard one (usually because they have a firewall which tries to stop P2P connections). This servers tell in this field which additional ports they have open. Each additional port is separated from the others by a coma (,).
	TagCodeLowIDClients       TagCode = 0x94 // Number of users connected with a LowID

)

func (code TagCode) ValueType(predefinedValueType TagValueType) TagValueType {
	switch code {
	default:
		// TODO Default type to string is this ok ?
		return TagValueTypeString
	case TagCodeVersion:
		// TODO Should I use predefined type for version ?
		return predefinedValueType
	case TagCodePing, TagCodeFail, TagCodePreference, TagCodeMaxUsers, TagCodeSoftFiles, TagCodeHardFiles, TagCodeLastPing, TagCodeUDPFlags, TagCodeLowIDClients:
		return TagValueTypeUint32
	}
}

type MetServerTag struct {
	Name        string       `json:"name,omitempty"`
	Code        TagCode      `json:"code,omitempty"`
	ValueType   TagValueType `json:"type"`
	StringValue string       `json:"string_value,omitempty"`
	Uint32Value uint32       `json:"uint32_value,omitempty"`
}
