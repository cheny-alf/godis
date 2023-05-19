package database

import (
	"github.com/hdt3213/godis/datastruct/rax"
	"github.com/hdt3213/godis/interface/database"
	"github.com/hdt3213/godis/interface/redis"
	"github.com/hdt3213/godis/redis/protocol"
	"strconv"
	"strings"
	"time"
)

// streamID key-value
type streamID struct {
	Key   int64 //Unix time in milliseconds
	Value int64 //Sequence number
}

func (s *streamID) toString() string {
	key := strconv.FormatInt(s.Key, 10)
	value := strconv.FormatInt(s.Value, 10)
	return key + "-" + value
}

type Stream struct {
	MsgList rax.Rax  //Stream list
	Length  int      //Stream length
	LastID  streamID //Zero if there are yet no items.
	Cgroups rax.Rax  //Consumers group list
}

// StreamCG Consumer group
type StreamCG struct {
	LastID    streamID //Last delivered (not acknowledged) ID for this group.
	Pel       rax.Rax  //Pending entries list.
	Consumers rax.Rax  //Consumer list
}

type StreamConsumer struct {
	Name       string    //Consumer name
	ActiveTime time.Time //Last time this consumer was active.
	Pel        rax.Rax   //Pending entries list.
}
type StreamNACK struct {
	consumer      StreamConsumer //The consumer this message was delivered toin the last delivery.
	DeliveryTime  time.Time      //Last time this message was delivered.
	deliveryCount int            //Number of times this message was delivered
}

func NewStream() *Stream {
	return &Stream{
		MsgList: rax.NewRax(),
		Length:  0,
		LastID:  streamID{Key: 0, Value: 0},
		Cgroups: rax.NewRax(),
	}
}

func execAdd(db *DB, args [][]byte) redis.Reply {
	idx := 0
	_, isExist := db.GetEntity(string(args[idx]))
	if string(args[idx+1]) == "nomkstream" {
		if !isExist {
			return protocol.MakeNullBulkReply()
		}
	} else {
		if !isExist {
			stream := NewStream()
			db.PutEntity(string(args[idx]), &database.DataEntity{
				Data: stream,
			})
		}
	}
	entryKey := string(args[idx])
	dataEntity, isExist := db.GetEntity(entryKey)
	stream := dataEntity.Data.(*Stream)
	idx++
	var maxlen int64
	var isHard bool
	var err error
	if string(args[idx]) == "maxlen" {
		idx++
		if string(args[idx]) == "=" {
			isHard = true
		} else if string(args[idx]) == "~" {
			isHard = false
		} else {
			return protocol.MakeErrReply("ERR value is not an integer or out of range")
		}
		idx++
		maxlen, err = strconv.ParseInt(string(args[idx]), 0, 64)
		if err != nil {
			return protocol.MakeErrReply("ERR value is not an integer or out of range")
		}
		idx++
	}
	//拷贝问题
	var IDDetail streamID
	arr := strings.Split(string(args[idx]), "-")
	if len(arr) > 2 {
		return protocol.MakeErrReply("ERR Invalid stream ID specified as stream command argument")
	} else if len(arr) == 1 {
		if arr[0] == "*" {
			IDDetail.Key = time.Now().Unix()
			if stream.LastID.Key == IDDetail.Key {
				IDDetail.Value = stream.LastID.Value + 1
			} else {
				IDDetail.Value = 0
			}
		} else {
			tmpKey, err := strconv.ParseInt(arr[0], 10, 64)
			if err != nil {
				return protocol.MakeErrReply("ERR value is not an integer or out of range")
			}
			IDDetail.Key = tmpKey
			IDDetail.Value = 0
		}
	} else if len(arr) == 2 {
		tmpKey, err := strconv.ParseInt(arr[0], 10, 64)
		if err != nil {
			return protocol.MakeErrReply("ERR value is not an integer or out of range")
		}
		tmpVal, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			return protocol.MakeErrReply("ERR value is not an integer or out of range")
		}
		if stream.LastID.Key != 0 && stream.LastID.Key > tmpKey {
			return protocol.MakeErrReply("ERR The ID specified in XADD is equal or smaller than the target stream top item")
		} else if stream.LastID.Key != 0 && stream.LastID.Key == tmpKey && stream.LastID.Value >= tmpVal {
			return protocol.MakeErrReply("ERR The ID specified in XADD is equal or smaller than the target stream top item")
		}
		IDDetail.Key = tmpKey
		IDDetail.Value = tmpVal
	}
	idx++
	if (len(args)-idx)%2 != 0 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'xadd' command")
	}
	for i := idx; i <= len(args)-2; i += 2 {
		stream.Length++
		stream.MsgList.Insert(string(args[i]), string(args[i+1]))
	}
	for isHard && int64(stream.Length) > maxlen {
		//尾删法
		stream.MsgList.RemoveFirstNode()
		stream.Length--
	}
	stream.LastID = IDDetail
	db.PutIfExists(entryKey, &database.DataEntity{Data: stream})
	return protocol.MakeBulkReply([]byte(IDDetail.toString()))
}

func init() {
	registerCommand("XAdd", execAdd, writeFirstKey, rollbackFirstKey, -4, flagWrite)
}
