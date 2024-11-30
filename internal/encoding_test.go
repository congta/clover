package internal

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

type BaseModel struct {
	ID string `clover:"_id" json:"cloverId"`
}

type TestStruct struct {
	BaseModel
	IntField    int                    `clover:"int,omitempty"`
	UintField   uint                   `clover:"uint,omitempty"`
	StringField string                 `clover:",omitempty"`
	FloatField  float32                `clover:",omitempty"`
	BoolField   bool                   `clover:",omitempty"`
	TimeField   time.Time              `clover:",omitempty"`
	IntPtr      *int                   `clover:",omitempty"`
	SliceField  []int                  `clover:",omitempty"`
	MapField    map[string]interface{} `clover:",omitempty"`
	Data        []byte                 `clover:",omitempty"`
}

type TestStruct2 struct {
	BaseModel
	IntField    int  `clover:"int_field,omitempty"`
	UintField   uint `clover:"uint,omitempty" json:"uint_field"`
	StringField string
	FloatField  float32   `clover:",omitempty" json:"float_field"`
	BoolField   bool      `clover:"bool_field,omitempty" json:"bool_field"`
	TimeField   time.Time `clover:",omitempty"`
	Data        []byte    `clover:",omitempty"`
}

type BsonStruct struct {
	CamelCase int    `bson:"camel_case" json:"camelCase"`
	Normal    string `bson:"normal" json:"normal"`
}

type TestStruct3 struct {
	CamelCase   int         `bson:"camel_case" json:"camelCase"`
	Normal      string      `bson:"normal" json:"normal"`
	NestedCamel *BsonStruct `bson:"nested_camel" json:"nestedCamel"`
	Ignore      string      `bson:"-" json:"-"`
}

func TestNormalize(t *testing.T) {
	date := time.Date(2020, 01, 1, 0, 0, 0, 0, time.UTC)

	var x = 100

	s := &TestStruct{
		BaseModel: BaseModel{
			ID: "UID",
		},
		TimeField:   date,
		IntField:    10,
		FloatField:  0.1,
		StringField: "aString",
		BoolField:   true,
		IntPtr:      &x,
		SliceField:  []int{1, 2, 3, 4},
		Data:        []byte("hello, clover!"),
		MapField: map[string]interface{}{
			"hello": "clover",
		},
	}

	ns, err := Normalize(s)
	require.NoError(t, err)

	require.IsType(t, ns, map[string]interface{}{})

	m := ns.(map[string]interface{})
	require.IsType(t, m["Data"], []byte{})

	require.Nil(t, m["uint"]) // testing omitempty
	require.Equal(t, m["IntPtr"], int64(100))
	require.Equal(t, m["_id"], "UID")

	s1 := &TestStruct{}
	err = Convert(m, s1)
	require.NoError(t, err)

	require.Equal(t, s, s1)
	require.Equal(t, s1.ID, "UID")

	err = Convert(m, 10)
	require.Error(t, err)
}

func TestNormalize2(t *testing.T) {
	norm, err := Normalize(nil)
	require.NoError(t, err)
	require.Nil(t, norm)

	_, err = Normalize(make(chan struct{}))
	require.Error(t, err)

	_, err = Normalize(make(map[int]interface{}))
	require.Error(t, err)
}

func TestNormalize3(t *testing.T) {
	date := time.Date(2020, 01, 1, 0, 0, 0, 0, time.UTC)

	s := &TestStruct{
		TimeField:   date,
		UintField:   0,
		IntField:    0,
		FloatField:  0,
		StringField: "",
		BoolField:   false,
		IntPtr:      nil,
		SliceField:  []int{},
		Data:        nil,
		MapField:    map[string]interface{}{},
	}

	ns, err := Normalize(s)
	require.NoError(t, err)

	require.IsType(t, ns, map[string]interface{}{})
	m := ns.(map[string]interface{})

	require.Nil(t, m["int"])
	require.Nil(t, m["uint"])
	require.Nil(t, m["FloatField"])
	require.Nil(t, m["BoolField"])
	require.Nil(t, m["SliceField"])
	require.Nil(t, m["Data"])
	require.Nil(t, m["MapField"])
	require.Nil(t, m["IntPtr"])
}

func TestNormalize4(t *testing.T) {
	var x = 100

	s := &TestStruct3{
		CamelCase: x,
		Normal:    "normal",
		NestedCamel: &BsonStruct{
			CamelCase: x + 1,
			Normal:    "normal-b",
		},
		Ignore: "hello",
	}

	ns, err := Normalize(s)
	require.NoError(t, err)

	require.IsType(t, ns, map[string]interface{}{})

	m := ns.(map[string]interface{})
	require.Equal(t, m["camel_case"], int64(x))
	require.Equal(t, m["normal"], "normal")
	require.Nil(t, m["Ignore"])
	require.Nil(t, m["-"])
	require.IsType(t, m["nested_camel"], map[string]interface{}{})

	nested, _ := m["nested_camel"].(map[string]interface{})
	require.Equal(t, nested["camel_case"], int64(x+1))
	require.Equal(t, nested["normal"], "normal-b")

	s1 := &TestStruct3{}
	err = Convert(m, s1)
	require.NoError(t, err)

	require.Equal(t, s1.Normal, "normal")
	require.Equal(t, s1.CamelCase, x)
	require.NotNil(t, s1.NestedCamel)
	require.Equal(t, s1.NestedCamel.Normal, "normal-b")
	require.Equal(t, s1.NestedCamel.CamelCase, x+1)
}

func TestEncodeDecode(t *testing.T) {
	s := &TestStruct{}

	require.NoError(t, gofakeit.Struct(s))

	norm, err := Normalize(s)
	require.NoError(t, err)

	data, err := Encode(norm.(map[string]interface{}))
	require.NoError(t, err)

	var m map[string]interface{}

	require.NoError(t, Decode(data, &m))

	require.Equal(t, m, norm)
}

func TestJsonTag(t *testing.T) {
	date := time.Date(2020, 01, 1, 0, 0, 0, 0, time.UTC)

	s := &TestStruct2{
		TimeField:   date,
		UintField:   100,
		IntField:    200,
		FloatField:  300,
		StringField: "json",
		BoolField:   true,
	}

	normalized, err := Normalize(s)
	require.NoError(t, err)

	require.IsType(t, normalized, map[string]interface{}{})
	fields, _ := normalized.(map[string]interface{})
	require.NotNil(t, fields)

	var ns TestStruct2
	err = Convert(fields, &ns)
	require.NoError(t, err)

	require.Equal(t, s, &ns)
}
