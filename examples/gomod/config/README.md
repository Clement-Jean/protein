# Config

This example shows how you can enable/disable features at compile time.

## Generating SourceCodeInfo

### enable

```
$ go run -tags=protein_generate_source_code_info main.go
```

### disable

```
$ go run main.go
```

## Keeping Comments

### enable

```
$ go run -tags=protein_keep_comments main.go
```

### disable

```
$ go run main.go
```