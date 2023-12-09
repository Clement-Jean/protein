# Config

This is compile time configuration exposed to users.

For now, this is handled with build tags, however, in the future we might take advantage of compiler flags if [this proposal](https://github.com/golang/go/issues/63372) gets accepted and implemented.

## Flags

- `protein_keep_comments`: As its name suggests, it doesn't skip comments. This might be something you want to enable if you want to generate documentation out of comments, for example.
- `protein_generate_source_code_info`: This keeps the spaces/newline information and generate a [SourceCodeInfo](https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/descriptor.proto) object in the final AST.