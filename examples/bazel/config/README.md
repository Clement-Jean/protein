# Config

This example shows how you can enable/disable features at compile time.

## Generating SourceCodeInfo

### enable

```starlark
go_binary(
	//...
	gotags = [
		"protein_generate_source_code_info",
	],
)
```

You can try running with the following command:

```
$ bazel run //:lexer_with_sci_example
```

### disable

Simply do not add the gotags.

You can try running with the following command:

```
$ bazel run //:lexer_without_sci_example
```

## Keeping Comments

### enable

```starlark
go_binary(
	//...
	gotags = [
		"protein_keep_comments",
	],
)
```

You can try running with the following command:

```
$ bazel run //:lexer_with_comments_example
```

### disable

Simply do not add the gotags.

You can try running with the following command:

```
$ bazel run //:lexer_without_comments_example
```