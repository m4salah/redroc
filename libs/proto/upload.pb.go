// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: libs/proto/upload.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UploadImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ObjName string `protobuf:"bytes,1,opt,name=obj_name,json=objName,proto3" json:"obj_name,omitempty"`
	Image   []byte `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
}

func (x *UploadImageRequest) Reset() {
	*x = UploadImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_libs_proto_upload_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadImageRequest) ProtoMessage() {}

func (x *UploadImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_libs_proto_upload_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadImageRequest.ProtoReflect.Descriptor instead.
func (*UploadImageRequest) Descriptor() ([]byte, []int) {
	return file_libs_proto_upload_proto_rawDescGZIP(), []int{0}
}

func (x *UploadImageRequest) GetObjName() string {
	if x != nil {
		return x.ObjName
	}
	return ""
}

func (x *UploadImageRequest) GetImage() []byte {
	if x != nil {
		return x.Image
	}
	return nil
}

type UploadImageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UploadImageResponse) Reset() {
	*x = UploadImageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_libs_proto_upload_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadImageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadImageResponse) ProtoMessage() {}

func (x *UploadImageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_libs_proto_upload_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadImageResponse.ProtoReflect.Descriptor instead.
func (*UploadImageResponse) Descriptor() ([]byte, []int) {
	return file_libs_proto_upload_proto_rawDescGZIP(), []int{1}
}

type CreateMetadataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ObjName  string   `protobuf:"bytes,1,opt,name=obj_name,json=objName,proto3" json:"obj_name,omitempty"`
	User     string   `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	Hashtags []string `protobuf:"bytes,3,rep,name=hashtags,proto3" json:"hashtags,omitempty"`
}

func (x *CreateMetadataRequest) Reset() {
	*x = CreateMetadataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_libs_proto_upload_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateMetadataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMetadataRequest) ProtoMessage() {}

func (x *CreateMetadataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_libs_proto_upload_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMetadataRequest.ProtoReflect.Descriptor instead.
func (*CreateMetadataRequest) Descriptor() ([]byte, []int) {
	return file_libs_proto_upload_proto_rawDescGZIP(), []int{2}
}

func (x *CreateMetadataRequest) GetObjName() string {
	if x != nil {
		return x.ObjName
	}
	return ""
}

func (x *CreateMetadataRequest) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *CreateMetadataRequest) GetHashtags() []string {
	if x != nil {
		return x.Hashtags
	}
	return nil
}

type CreateMetadataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateMetadataResponse) Reset() {
	*x = CreateMetadataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_libs_proto_upload_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateMetadataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMetadataResponse) ProtoMessage() {}

func (x *CreateMetadataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_libs_proto_upload_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMetadataResponse.ProtoReflect.Descriptor instead.
func (*CreateMetadataResponse) Descriptor() ([]byte, []int) {
	return file_libs_proto_upload_proto_rawDescGZIP(), []int{3}
}

type ImageUploadedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ObjName  string   `protobuf:"bytes,1,opt,name=obj_name,json=objName,proto3" json:"obj_name,omitempty"`
	User     string   `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	Hashtags []string `protobuf:"bytes,3,rep,name=hashtags,proto3" json:"hashtags,omitempty"`
}

func (x *ImageUploadedRequest) Reset() {
	*x = ImageUploadedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_libs_proto_upload_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImageUploadedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageUploadedRequest) ProtoMessage() {}

func (x *ImageUploadedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_libs_proto_upload_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageUploadedRequest.ProtoReflect.Descriptor instead.
func (*ImageUploadedRequest) Descriptor() ([]byte, []int) {
	return file_libs_proto_upload_proto_rawDescGZIP(), []int{4}
}

func (x *ImageUploadedRequest) GetObjName() string {
	if x != nil {
		return x.ObjName
	}
	return ""
}

func (x *ImageUploadedRequest) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *ImageUploadedRequest) GetHashtags() []string {
	if x != nil {
		return x.Hashtags
	}
	return nil
}

type ImageUploadedResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ImageUploadedResponse) Reset() {
	*x = ImageUploadedResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_libs_proto_upload_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImageUploadedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageUploadedResponse) ProtoMessage() {}

func (x *ImageUploadedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_libs_proto_upload_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageUploadedResponse.ProtoReflect.Descriptor instead.
func (*ImageUploadedResponse) Descriptor() ([]byte, []int) {
	return file_libs_proto_upload_proto_rawDescGZIP(), []int{5}
}

var File_libs_proto_upload_proto protoreflect.FileDescriptor

var file_libs_proto_upload_proto_rawDesc = []byte{
	0x0a, 0x17, 0x6c, 0x69, 0x62, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x72, 0x70, 0x63, 0x22,
	0x45, 0x0a, 0x12, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x62, 0x6a, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x62, 0x6a, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x22, 0x15, 0x0a, 0x13, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x62, 0x0a,
	0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x62, 0x6a, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x62, 0x6a, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x61, 0x73, 0x68, 0x74, 0x61, 0x67,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x68, 0x61, 0x73, 0x68, 0x74, 0x61, 0x67,
	0x73, 0x22, 0x18, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x61, 0x0a, 0x14, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x62, 0x6a, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x62, 0x6a, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x61, 0x73, 0x68, 0x74, 0x61, 0x67, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x68, 0x61, 0x73, 0x68, 0x74, 0x61, 0x67, 0x73, 0x22, 0x17,
	0x0a, 0x15, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xe3, 0x01, 0x0a, 0x0b, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x12, 0x3d, 0x0a, 0x06, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x12, 0x18, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4b, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x0d, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x65, 0x64, 0x12, 0x1a, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x49, 0x6d, 0x61, 0x67,
	0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1b, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x24, 0x5a,
	0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x34, 0x73, 0x6c,
	0x61, 0x68, 0x2f, 0x72, 0x65, 0x64, 0x72, 0x6f, 0x63, 0x2f, 0x6c, 0x69, 0x62, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_libs_proto_upload_proto_rawDescOnce sync.Once
	file_libs_proto_upload_proto_rawDescData = file_libs_proto_upload_proto_rawDesc
)

func file_libs_proto_upload_proto_rawDescGZIP() []byte {
	file_libs_proto_upload_proto_rawDescOnce.Do(func() {
		file_libs_proto_upload_proto_rawDescData = protoimpl.X.CompressGZIP(file_libs_proto_upload_proto_rawDescData)
	})
	return file_libs_proto_upload_proto_rawDescData
}

var file_libs_proto_upload_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_libs_proto_upload_proto_goTypes = []interface{}{
	(*UploadImageRequest)(nil),     // 0: grpc.UploadImageRequest
	(*UploadImageResponse)(nil),    // 1: grpc.UploadImageResponse
	(*CreateMetadataRequest)(nil),  // 2: grpc.CreateMetadataRequest
	(*CreateMetadataResponse)(nil), // 3: grpc.CreateMetadataResponse
	(*ImageUploadedRequest)(nil),   // 4: grpc.ImageUploadedRequest
	(*ImageUploadedResponse)(nil),  // 5: grpc.ImageUploadedResponse
}
var file_libs_proto_upload_proto_depIdxs = []int32{
	0, // 0: grpc.UploadPhoto.Upload:input_type -> grpc.UploadImageRequest
	2, // 1: grpc.UploadPhoto.CreateMetadata:input_type -> grpc.CreateMetadataRequest
	4, // 2: grpc.UploadPhoto.ImageUploaded:input_type -> grpc.ImageUploadedRequest
	1, // 3: grpc.UploadPhoto.Upload:output_type -> grpc.UploadImageResponse
	3, // 4: grpc.UploadPhoto.CreateMetadata:output_type -> grpc.CreateMetadataResponse
	5, // 5: grpc.UploadPhoto.ImageUploaded:output_type -> grpc.ImageUploadedResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_libs_proto_upload_proto_init() }
func file_libs_proto_upload_proto_init() {
	if File_libs_proto_upload_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_libs_proto_upload_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadImageRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_libs_proto_upload_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadImageResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_libs_proto_upload_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateMetadataRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_libs_proto_upload_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateMetadataResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_libs_proto_upload_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImageUploadedRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_libs_proto_upload_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImageUploadedResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_libs_proto_upload_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_libs_proto_upload_proto_goTypes,
		DependencyIndexes: file_libs_proto_upload_proto_depIdxs,
		MessageInfos:      file_libs_proto_upload_proto_msgTypes,
	}.Build()
	File_libs_proto_upload_proto = out.File
	file_libs_proto_upload_proto_rawDesc = nil
	file_libs_proto_upload_proto_goTypes = nil
	file_libs_proto_upload_proto_depIdxs = nil
}
