package grpc_proxy

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// Name is the name registered for the proto compressor.
const Name = "proxy"

type MessageCodec interface {
	// Marshal returns the wire format of v.
	Marshal() ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte) error
}

// ProxyCodec is a Codec implementation with protobuf. It is the default ProxyCodec for gRPC.
type ProxyCodec struct{}

func (ProxyCodec) Marshal(v interface{}) ([]byte, error) {
	pd, ok := v.(MessageCodec)
	if ok {
		return pd.Marshal()
	}
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	return proto.Marshal(vv)
}

func (ProxyCodec) Unmarshal(data []byte, v interface{}) error {
	pd, ok := v.(MessageCodec)
	if ok {
		return pd.Unmarshal(data)
	}
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return proto.Unmarshal(data, vv)
}

func (ProxyCodec) Name() string {
	return Name
}

func (ProxyCodec) String() string {
	return Name
}

type ProxyData struct {
	Data []byte
}

func (d *ProxyData) Unmarshal(data []byte) error {
	d.Data = make([]byte, len(data))
	copy(d.Data, data)
	return nil
}

func (d *ProxyData) Marshal() ([]byte, error) {
	if len(d.Data) > 0 {
		return d.Data, nil
	}
	return nil, errors.New("empty data")
}
