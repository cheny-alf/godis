package database

import (
	"github.com/hdt3213/godis/lib/utils"
	"github.com/hdt3213/godis/redis/protocol"
	"testing"
)

func TestXAdd(t *testing.T) {
	actual := testDB.Exec(nil, utils.ToCmdLine("xadd", "stream1", "12-0", "field1", "value1"))
	expected := protocol.MakeBulkReply([]byte("12-0"))
	if !utils.BytesEquals(actual.ToBytes(), expected.ToBytes()) {
		t.Error("expected: " + string(expected.ToBytes()) + ", actual: " + string(actual.ToBytes()))
	}
	actual = testDB.Exec(nil, utils.ToCmdLine("xadd", "stream1", "12-1", "field2", "value2"))
	expected = protocol.MakeBulkReply([]byte("12-1"))
	if !utils.BytesEquals(actual.ToBytes(), expected.ToBytes()) {
		t.Error("expected: " + string(expected.ToBytes()) + ", actual: " + string(actual.ToBytes()))
	}
	actual = testDB.Exec(nil, utils.ToCmdLine("xadd", "stream1", "maxlen", "=", "2", "12-2", "field3", "value3"))
	expected = protocol.MakeBulkReply([]byte("12-2"))
	if !utils.BytesEquals(actual.ToBytes(), expected.ToBytes()) {
		t.Error("expected: " + string(expected.ToBytes()) + ", actual: " + string(actual.ToBytes()))
	}

}
