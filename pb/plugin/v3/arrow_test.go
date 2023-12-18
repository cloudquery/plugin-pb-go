package plugin

import (
	"testing"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
)

func TestSchemaRoundTrip(t *testing.T) {
	md := arrow.NewMetadata([]string{"foo", "bar"}, []string{"baz", "quux"})
	sc := arrow.NewSchema([]arrow.Field{
		{Name: "a", Type: arrow.PrimitiveTypes.Int64},
		{Name: "b", Type: arrow.PrimitiveTypes.Float64},
	}, &md)
	b, err := SchemaToBytes(sc)
	if err != nil {
		t.Fatal(err)
	}
	sc2, err := NewSchemaFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if !sc.Equal(sc2) {
		t.Errorf("expected %v, got %v", sc, sc2)
	}

	schemasBytes, err := SchemasToBytes([]*arrow.Schema{sc})
	if err != nil {
		t.Fatal(err)
	}

	schemas, err := NewSchemasFromBytes(schemasBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !sc.Equal(schemas[0]) {
		t.Errorf("expected %v, got %v", sc, schemas[0])
	}
}

func TestRecordRoundTrip(t *testing.T) {
	md := arrow.NewMetadata([]string{"foo", "bar"}, []string{"baz", "quux"})
	sc := arrow.NewSchema([]arrow.Field{
		{Name: "a", Type: arrow.PrimitiveTypes.Int64},
		{Name: "b", Type: arrow.PrimitiveTypes.Float64},
	}, &md)
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
	bldr.Field(0).(*array.Int64Builder).AppendValues([]int64{1, 2, 3}, nil)
	bldr.Field(1).(*array.Float64Builder).AppendValues([]float64{1.1, 2.2, 3.3}, nil)
	rec := bldr.NewRecord()
	defer rec.Release()

	b, err := RecordToBytes(rec)
	if err != nil {
		t.Fatal(err)
	}
	rec2, err := NewRecordFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if !array.RecordEqual(rec, rec2) {
		t.Errorf("expected %v, got %v", rec, rec2)
	}
}
