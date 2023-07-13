package utils

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
)

func ReaderToMessage(src io.Reader, dst proto.Message) error {
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	return BytesToMessage(data, dst)
}

func BytesToMessage(src []byte, dst proto.Message) error {
	unmarshalOptions := proto.UnmarshalOptions{}

	//if !IsProduction() {
	//	// 指示是否允许反序列化部分消息。如果设置为 true，则当遇到未知的字段或子消息时，反序列化操作将继续执行而不是返回错误。
	//	unmarshalOptions.AllowPartial = true
	//	// 指示是否应该丢弃未知字段。如果设置为 true，则在处理未知字段时不会报告任何错误或警告。
	//	unmarshalOptions.DiscardUnknown = true
	//}

	return unmarshalOptions.Unmarshal(src, dst)
}

func MessageToBytes(msg proto.Message) ([]byte, error) {
	options := proto.MarshalOptions{}

	//if !IsProduction() {
	//	// 指示是否允许反序列化部分消息。如果设置为 true，则当遇到未知的字段或子消息时，反序列化操作将继续执行而不是返回错误。
	//	options.AllowPartial = true
	//}

	return options.Marshal(msg)
}

func AnyToMessage(src *anypb.Any, dst proto.Message) error {
	unmarshalOptions := proto.UnmarshalOptions{}

	//if !IsProduction() {
	//	// 指示是否允许反序列化部分消息。如果设置为 true，则当遇到未知的字段或子消息时，反序列化操作将继续执行而不是返回错误。
	//	unmarshalOptions.AllowPartial = true
	//	// 指示是否应该丢弃未知字段。如果设置为 true，则在处理未知字段时不会报告任何错误或警告。
	//	unmarshalOptions.DiscardUnknown = true
	//}

	return anypb.UnmarshalTo(src, dst, unmarshalOptions)
}

func MessageToAny(dst *anypb.Any, src proto.Message) error {
	marshalOptions := proto.MarshalOptions{}

	//if !IsProduction() {
	//	// 指示是否允许反序列化部分消息。如果设置为 true，则当遇到未知的字段或子消息时，反序列化操作将继续执行而不是返回错误。
	//	marshalOptions.AllowPartial = true
	//}

	return anypb.MarshalFrom(dst, src, marshalOptions)
}

func ToAny(src proto.Message) (*anypb.Any, error) {
	_any := &anypb.Any{
		TypeUrl: "",
		Value:   nil,
	}

	err := anypb.MarshalFrom(_any, src, proto.MarshalOptions{})
	if err != nil {
		return nil, err
	}

	return _any, err
}

// json to message

func JSONReaderToMessage(src io.Reader, dst proto.Message) error {
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	return JSONBytesToMessage(data, dst)
}

func JSONBytesToMessage(src []byte, dst proto.Message) error {
	unmarshalOptions := protojson.UnmarshalOptions{}

	//if !IsProduction() {
	//	// 客户端发过来的 request body 少字段
	//	// 指示是否允许反序列化部分消息。如果设置为 true，则当遇到未知的字段或子消息时，反序列化操作将继续执行而不是返回错误。
	//	unmarshalOptions.AllowPartial = true
	//	// 客户端发过来的 request body 多字段
	//	// 指示是否应该丢弃未知字段。如果设置为 true，则在处理未知字段时不会报告任何错误或警告。
	//	unmarshalOptions.DiscardUnknown = true
	//}

	return unmarshalOptions.Unmarshal(src, dst)
}

func MessageToJSONBytes(dst proto.Message) ([]byte, error) {
	options := protojson.MarshalOptions{}
	return options.Marshal(dst)
}
