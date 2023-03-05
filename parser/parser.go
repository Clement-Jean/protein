package parser

import pb "google.golang.org/protobuf/types/descriptorpb"

// Parser is protein's parser
type Parser interface {
	// Parse returns the representation of a file in Protobuf Descriptor
	Parse() pb.FileDescriptorProto
}
