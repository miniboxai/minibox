package server

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var typeProtoMessage = reflect.TypeOf((*proto.Message)(nil)).Elem()

// JSONPb is a Marshaler which marshals/unmarshals into/from JSON
// with the "github.com/gogo/protobuf/jsonpb".
// It supports fully functionality of protobuf unlike JSONBuiltin.
type JSONPb jsonpb.Marshaler

// ContentType always returns "application/json".
func (*JSONPb) ContentType() string {
	return "application/json"
}

// Marshal marshals "v" into JSON
func (j *JSONPb) Marshal(v interface{}) ([]byte, error) {

	if val, ok := v.(proto.Message); ok {
		var buf bytes.Buffer

		marshalfn := (*jsonpb.Marshaler)(j).Marshal

		if err := marshalfn(&buf, val); err != nil {

			return nil, err
		}
		return buf.Bytes(), nil
	} else if val, ok := v.(map[string]proto.Message); ok {
		var buf bytes.Buffer

		// buf.WriteString("{")
		// count := len(val)
		// index := 0
		for _, v := range val {
			// key, _ := k.(string)
			// buf.WriteString(strconv.Quote(key))
			// buf.WriteString(":")
			if child, err := j.Marshal(v); err != nil {
				return nil, err
			} else {
				buf.Write(child)
			}
			// if index < count-1 {
			// 	buf.WriteString(",")
			// }
		}
		// buf.WriteString("}")
		return buf.Bytes(), nil
	}

	return nil, fmt.Errorf("unexpected type %T does not implement %s", v, typeProtoMessage)
}

// Unmarshal unmarshals JSON "data" into "v"
func (j *JSONPb) Unmarshal(data []byte, v interface{}) error {
	if pb, ok := v.(proto.Message); ok {
		return jsonpb.Unmarshal(bytes.NewReader(data), pb)
	}
	return fmt.Errorf("unexpected type %T does not implement %s", v, typeProtoMessage)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (j *JSONPb) NewDecoder(r io.Reader) gwruntime.Decoder {
	return gwruntime.DecoderFunc(func(v interface{}) error {
		if pb, ok := v.(proto.Message); ok {
			return jsonpb.Unmarshal(r, pb)
		}
		return fmt.Errorf("unexpected type %T does not implement %s", v, typeProtoMessage)
	})
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *JSONPb) NewEncoder(w io.Writer) gwruntime.Encoder {
	return gwruntime.EncoderFunc(func(v interface{}) error {
		if pb, ok := v.(proto.Message); ok {
			marshalFn := (*jsonpb.Marshaler)(j).Marshal
			return marshalFn(w, pb)
		}
		return fmt.Errorf("unexpected type %T does not implement %s", v, typeProtoMessage)
	})
}

// func (j *JSONPb) Delimiter() []byte {
// 	return []byte(",\n")
// }
