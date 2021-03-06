// Code generated by protoc-gen-go.
// source: mvt.proto
// DO NOT EDIT!

/*
Package mapDrawer is a generated protocol buffer package.

It is generated from these files:
	mvt.proto

It has these top-level messages:
	Tile
*/
package mapDrawer

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

// GeomType is described in section 4.3.4 of the specification
type Tile_GeomType int32

const (
	Tile_UNKNOWN    Tile_GeomType = 0
	Tile_POINT      Tile_GeomType = 1
	Tile_LINESTRING Tile_GeomType = 2
	Tile_POLYGON    Tile_GeomType = 3
)

var Tile_GeomType_name = map[int32]string{
	0: "UNKNOWN",
	1: "POINT",
	2: "LINESTRING",
	3: "POLYGON",
}
var Tile_GeomType_value = map[string]int32{
	"UNKNOWN":    0,
	"POINT":      1,
	"LINESTRING": 2,
	"POLYGON":    3,
}

func (x Tile_GeomType) Enum() *Tile_GeomType {
	p := new(Tile_GeomType)
	*p = x
	return p
}
func (x Tile_GeomType) String() string {
	return proto.EnumName(Tile_GeomType_name, int32(x))
}
func (x *Tile_GeomType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Tile_GeomType_value, data, "Tile_GeomType")
	if err != nil {
		return err
	}
	*x = Tile_GeomType(value)
	return nil
}

type Tile struct {
	Layers           []*Tile_Layer             `protobuf:"bytes,3,rep,name=layers" json:"layers,omitempty"`
	XXX_extensions   map[int32]proto.Extension `json:"-"`
	XXX_unrecognized []byte                    `json:"-"`
}

func (m *Tile) Reset()         { *m = Tile{} }
func (m *Tile) String() string { return proto.CompactTextString(m) }
func (*Tile) ProtoMessage()    {}

var extRange_Tile = []proto.ExtensionRange{
	{16, 8191},
}

func (*Tile) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_Tile
}
func (m *Tile) ExtensionMap() map[int32]proto.Extension {
	if m.XXX_extensions == nil {
		m.XXX_extensions = make(map[int32]proto.Extension)
	}
	return m.XXX_extensions
}

func (m *Tile) GetLayers() []*Tile_Layer {
	if m != nil {
		return m.Layers
	}
	return nil
}

// Variant type encoding
// The use of values is described in section 4.1 of the specification
type Tile_Value struct {
	// Exactly one of these values must be present in a valid message
	StringValue      *string                   `protobuf:"bytes,1,opt,name=string_value" json:"string_value,omitempty"`
	FloatValue       *float32                  `protobuf:"fixed32,2,opt,name=float_value" json:"float_value,omitempty"`
	DoubleValue      *float64                  `protobuf:"fixed64,3,opt,name=double_value" json:"double_value,omitempty"`
	IntValue         *int64                    `protobuf:"varint,4,opt,name=int_value" json:"int_value,omitempty"`
	UintValue        *uint64                   `protobuf:"varint,5,opt,name=uint_value" json:"uint_value,omitempty"`
	SintValue        *int64                    `protobuf:"zigzag64,6,opt,name=sint_value" json:"sint_value,omitempty"`
	BoolValue        *bool                     `protobuf:"varint,7,opt,name=bool_value" json:"bool_value,omitempty"`
	XXX_extensions   map[int32]proto.Extension `json:"-"`
	XXX_unrecognized []byte                    `json:"-"`
}

func (m *Tile_Value) Reset()         { *m = Tile_Value{} }
func (m *Tile_Value) String() string { return proto.CompactTextString(m) }
func (*Tile_Value) ProtoMessage()    {}

var extRange_Tile_Value = []proto.ExtensionRange{
	{8, 536870911},
}

func (*Tile_Value) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_Tile_Value
}
func (m *Tile_Value) ExtensionMap() map[int32]proto.Extension {
	if m.XXX_extensions == nil {
		m.XXX_extensions = make(map[int32]proto.Extension)
	}
	return m.XXX_extensions
}

func (m *Tile_Value) GetStringValue() string {
	if m != nil && m.StringValue != nil {
		return *m.StringValue
	}
	return ""
}

func (m *Tile_Value) GetFloatValue() float32 {
	if m != nil && m.FloatValue != nil {
		return *m.FloatValue
	}
	return 0
}

func (m *Tile_Value) GetDoubleValue() float64 {
	if m != nil && m.DoubleValue != nil {
		return *m.DoubleValue
	}
	return 0
}

func (m *Tile_Value) GetIntValue() int64 {
	if m != nil && m.IntValue != nil {
		return *m.IntValue
	}
	return 0
}

func (m *Tile_Value) GetUintValue() uint64 {
	if m != nil && m.UintValue != nil {
		return *m.UintValue
	}
	return 0
}

func (m *Tile_Value) GetSintValue() int64 {
	if m != nil && m.SintValue != nil {
		return *m.SintValue
	}
	return 0
}

func (m *Tile_Value) GetBoolValue() bool {
	if m != nil && m.BoolValue != nil {
		return *m.BoolValue
	}
	return false
}

// Features are described in section 4.2 of the specification
type Tile_Feature struct {
	Id *uint64 `protobuf:"varint,1,opt,name=id,def=0" json:"id,omitempty"`
	// Tags of this feature are encoded as repeated pairs of
	// integers.
	// A detailed description of tags is located in sections
	// 4.2 and 4.4 of the specification
	Tags []uint32 `protobuf:"varint,2,rep,packed,name=tags" json:"tags,omitempty"`
	// The type of geometry stored in this feature.
	Type *Tile_GeomType `protobuf:"varint,3,opt,name=type,enum=vector_tile.Tile_GeomType,def=0" json:"type,omitempty"`
	// Contains a stream of commands and parameters (vertices).
	// A detailed description on geometry encoding is located in
	// section 4.3 of the specification.
	Geometry         []uint32 `protobuf:"varint,4,rep,packed,name=geometry" json:"geometry,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Tile_Feature) Reset()         { *m = Tile_Feature{} }
func (m *Tile_Feature) String() string { return proto.CompactTextString(m) }
func (*Tile_Feature) ProtoMessage()    {}

const Default_Tile_Feature_Id uint64 = 0
const Default_Tile_Feature_Type Tile_GeomType = Tile_UNKNOWN

func (m *Tile_Feature) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return Default_Tile_Feature_Id
}

func (m *Tile_Feature) GetTags() []uint32 {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *Tile_Feature) GetType() Tile_GeomType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Default_Tile_Feature_Type
}

func (m *Tile_Feature) GetGeometry() []uint32 {
	if m != nil {
		return m.Geometry
	}
	return nil
}

// Layers are described in section 4.1 of the specification
type Tile_Layer struct {
	// Any compliant implementation must first read the version
	// number encoded in this message and choose the correct
	// implementation for this version number before proceeding to
	// decode other parts of this message.
	Version *uint32 `protobuf:"varint,15,req,name=version,def=1" json:"version,omitempty"`
	Name    *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	// The actual features in this tile.
	Features []*Tile_Feature `protobuf:"bytes,2,rep,name=features" json:"features,omitempty"`
	// Dictionary encoding for keys
	Keys []string `protobuf:"bytes,3,rep,name=keys" json:"keys,omitempty"`
	// Dictionary encoding for values
	Values []*Tile_Value `protobuf:"bytes,4,rep,name=values" json:"values,omitempty"`
	// Although this is an "optional" field it is required by the specification.
	// See https://github.com/mapbox/vector-tile-spec/issues/47
	Extent           *uint32                   `protobuf:"varint,5,opt,name=extent,def=4096" json:"extent,omitempty"`
	XXX_extensions   map[int32]proto.Extension `json:"-"`
	XXX_unrecognized []byte                    `json:"-"`
}

func (m *Tile_Layer) Reset()         { *m = Tile_Layer{} }
func (m *Tile_Layer) String() string { return proto.CompactTextString(m) }
func (*Tile_Layer) ProtoMessage()    {}

var extRange_Tile_Layer = []proto.ExtensionRange{
	{16, 536870911},
}

func (*Tile_Layer) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_Tile_Layer
}
func (m *Tile_Layer) ExtensionMap() map[int32]proto.Extension {
	if m.XXX_extensions == nil {
		m.XXX_extensions = make(map[int32]proto.Extension)
	}
	return m.XXX_extensions
}

const Default_Tile_Layer_Version uint32 = 1
const Default_Tile_Layer_Extent uint32 = 4096

func (m *Tile_Layer) GetVersion() uint32 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return Default_Tile_Layer_Version
}

func (m *Tile_Layer) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Tile_Layer) GetFeatures() []*Tile_Feature {
	if m != nil {
		return m.Features
	}
	return nil
}

func (m *Tile_Layer) GetKeys() []string {
	if m != nil {
		return m.Keys
	}
	return nil
}

func (m *Tile_Layer) GetValues() []*Tile_Value {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *Tile_Layer) GetExtent() uint32 {
	if m != nil && m.Extent != nil {
		return *m.Extent
	}
	return Default_Tile_Layer_Extent
}

func init() {
	proto.RegisterEnum("vector_tile.Tile_GeomType", Tile_GeomType_name, Tile_GeomType_value)
}
