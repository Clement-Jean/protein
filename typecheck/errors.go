package typecheck

import (
	"fmt"
	"strings"
)

type Warning interface {
	Warning() string
}

type UnknownSyntaxError struct {
	File      string
	Line, Col int
	Value     string
}

func (e *UnknownSyntaxError) Error() string {
	return fmt.Sprintf("%s:%d:%d: error: %q is not a known protobuf syntax", e.File, e.Line, e.Col, e.Value)
}

type ImportCycleError struct {
	Files []string
	// TODO LINES, COLS and change error message
}

func (e *ImportCycleError) Error() string {
	var msg strings.Builder

	for _, file := range e.Files {
		msg.WriteString(file)
		msg.WriteString(" -> ")
	}
	msg.WriteString(e.Files[0])

	return fmt.Sprintf("cycle found: %s", msg.String())
}

type ImportFileNotFoundError struct {
	File string
	// TODO LINE, COL and change error message
}

func (e *ImportFileNotFoundError) Error() string {
	return fmt.Sprintf("file %s was not found or had errors", e.File)
}

type PackageMultipleDefError struct {
	File string
	// TODO LINES, COLS and change error message
}

func (e *PackageMultipleDefError) Error() string {
	return fmt.Sprintf("multiple package definitions in %s", e.File)
}

type NotTypeError struct {
	Name      string
	File      string
	Line, Col int
}

func (e *NotTypeError) Error() string {
	return fmt.Sprintf("%s:%d:%d: error: %s is not a type", e.File, e.Line, e.Col, e.Name)
}

type NotMessageTypeError struct {
	Name      string
	File      string
	Line, Col int
}

func (e *NotMessageTypeError) Error() string {
	return fmt.Sprintf("%s:%d:%d: error: %s is not a message type", e.File, e.Line, e.Col, e.Name)
}

type TypeResolvedNotDefinedError struct {
	Name, ResolvedName string
	File               string
	Line, Col          int
}

func (e *TypeResolvedNotDefinedError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d: error: %s is resolved to %s, which is not defined. The innermost scope is searched first in name resolution. Consider using a leading '.' (i.e., %s) to start from the outermost scope",
		e.File, e.Line, e.Col,
		e.Name, e.ResolvedName, "."+e.Name,
	)
}

type TypeNotDefinedError struct {
	Name      string
	File      string
	Line, Col int
}

func (e *TypeNotDefinedError) Error() string {
	return fmt.Sprintf("%s:%d:%d: error: %s is not defined", e.File, e.Line, e.Col, e.Name)
}

type TypeNotImportedError struct {
	Name             string
	DefFile, RefFile string
	Line, Col        int
}

func (e *TypeNotImportedError) Error() string {
	return fmt.Sprintf("%s:%d:%d: error: %s seems to be defined in %s, which is not imported by %s. To use it here, please add the necessary import.",
		e.RefFile, e.Line, e.Col,
		e.Name,
		e.DefFile,
		e.RefFile,
	)
}

type TypeRedefinedError struct {
	Name        string
	Files       []string
	Lines, Cols []int
}

func (e *TypeRedefinedError) Error() string {
	if len(e.Files) < 2 {
		return ""
	}

	var sb strings.Builder

	msg := fmt.Sprintf("%s:%d:%d: error: %s is redefined\n", e.Files[0], e.Lines[0], e.Cols[0], e.Name)
	sb.WriteString(msg)

	// FIX: we can do better for letting tools pickup the filepath
	for i := 1; i < len(e.Files); i++ {
		msg = fmt.Sprintf("\tredefined here: %s:%d:%d\n", e.Files[i], e.Lines[i], e.Cols[i])
		sb.WriteString(msg)
	}
	return strings.TrimRight(sb.String(), "\n")
}

type TypeUnusedWarning struct {
	Name      string
	File      string
	Line, Col int
}

func (w *TypeUnusedWarning) Warning() string {
	return fmt.Sprintf(
		"%s:%d:%d: warning: %s is defined but not used",
		w.File, w.Line, w.Col, w.Name,
	)
}

func (w *TypeUnusedWarning) Error() string {
	return w.Warning()
}

type ImportAlreadyImportedWarning struct {
	ImportingFile, ImportedFile string
	Line, Col                   int
}

func (w *ImportAlreadyImportedWarning) Warning() string {
	return fmt.Sprintf(
		"%s:%d:%d: warning: %s is already imported",
		w.ImportingFile, w.Line, w.Col, w.ImportedFile,
	)
}

func (w *ImportAlreadyImportedWarning) Error() string {
	return w.Warning()
}

type WeakImportNoEffectWarning struct {
	File      string
	Line, Col int
}

func (w *WeakImportNoEffectWarning) Warning() string {
	return fmt.Sprintf(
		"%s:%d:%d: warning: weak imports have no effect in Protein",
		w.File, w.Line, w.Col,
	)
}

func (w *WeakImportNoEffectWarning) Error() string {
	return w.Warning()
}

type FieldNameReusedError struct {
	ParentName, Name string
	File             string
	Line, Col        int
}

func (e *FieldNameReusedError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d: error: field name %q is already defined in %s",
		e.File, e.Line, e.Col,
		e.Name,
		e.ParentName,
	)
}

type FieldTagReusedError struct {
	ParentName string
	Tag        int64
	File       string
	Line, Col  int
}

func (e *FieldTagReusedError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d: error: field tag %d is already defined in %s",
		e.File, e.Line, e.Col,
		e.Tag,
		e.ParentName,
	)
}

type EnumCannotBeEmptyError struct {
	File      string
	Name      string
	Line, Col int
}

func (e *EnumCannotBeEmptyError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d: error: enum %q must contain at least one value.",
		e.File, e.Line, e.Col,
		e.Name,
	)
}

type OneofCannotHaveLessThan2FieldsError struct {
	File      string
	Name      string
	Line, Col int
}

func (e *OneofCannotHaveLessThan2FieldsError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d: error: oneof %q should contain at least two values.",
		e.File, e.Line, e.Col,
		e.Name,
	)
}

type MaxFieldTagError struct {
	File      string
	Line, Col int
}

func (e *MaxFieldTagError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d: error: field tags cannot be greater than 536,870,911.",
		e.File, e.Line, e.Col,
	)
}

type MaxEnumValueTagError struct {
	File      string
	Line, Col int
}

func (e *MaxEnumValueTagError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d: error: enum value tags cannot be greater than 2,147,483,647.",
		e.File, e.Line, e.Col,
	)
}

type EnumFirstValueTagZeroError struct {
	File      string
	Line, Col int
}

func (e *EnumFirstValueTagZeroError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d: error: the first enum value must be zero for enums.",
		e.File, e.Line, e.Col,
	)
}

type OptionUnknownError struct {
	File      string
	Line, Col int
	Name      string
}

func (e *OptionUnknownError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d: error: option %s is unknown.",
		e.File, e.Line, e.Col, e.Name,
	)
}

type ExtendMapNotAllowedError struct {
	File      string
	Line, Col int
}

func (e *ExtendMapNotAllowedError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d: error: map fields are not allowed to be extensions.",
		e.File, e.Line, e.Col,
	)
}
