// Copyright 2015-2019 Brett Vickers.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package etree provides XML services through an Element Tree
// abstraction.
package etree

import (
	"encoding/xml"
	"errors"
	"io"
	"iter"
)

const (
	// NoIndent is used with the IndentSettings record to remove all
	// indenting.
	NoIndent = -1
)

// ErrXML is returned when XML parsing fails due to incorrect formatting.
var ErrXML = errors.New("etree: invalid XML format")

// cdataPrefix is used to detect CDATA text when ReadSettings.PreserveCData is
// true.
var cdataPrefix = []byte("<![CDATA[")

// ReadSettings determine the default behavior of the Document's ReadFrom*
// functions.
type ReadSettings struct {
	// CharsetReader, if non-nil, defines a function to generate
	// charset-conversion readers, converting from the provided non-UTF-8
	// charset into UTF-8. If nil, the ReadFrom* functions will use a
	// "pass-through" CharsetReader that performs no conversion on the reader's
	// data regardless of the value of the "charset" encoding string. Default:
	// nil.
	CharsetReader func(charset string, input io.Reader) (io.Reader, error)

	// Permissive allows input containing common mistakes such as missing tags
	// or attribute values. Default: false.
	Permissive bool

	// Preserve CDATA character data blocks when decoding XML (instead of
	// converting it to normal character text). This entails additional
	// processing and memory usage during ReadFrom* operations. Default:
	// false.
	PreserveCData bool

	// When an element has two or more attributes with the same name,
	// preserve them instead of keeping only one. Default: false.
	PreserveDuplicateAttrs bool

	// ValidateInput forces all ReadFrom* functions to validate that the
	// provided input is composed of "well-formed"(*) XML before processing it.
	// If invalid XML is detected, the ReadFrom* functions return an error.
	// Because this option requires the input to be processed twice, it incurs a
	// significant performance penalty. Default: false.
	//
	// (*) Note that this definition of "well-formed" is in the context of the
	// go standard library's encoding/xml package. Go's encoding/xml package
	// does not, in fact, guarantee well-formed XML as specified by the W3C XML
	// recommendation. See: https://github.com/golang/go/issues/68299
	ValidateInput bool

	// Entity to be passed to standard xml.Decoder. Default: nil.
	Entity map[string]string

	// When Permissive is true, AutoClose indicates a set of elements to
	// consider closed immediately after they are opened, regardless of
	// whether an end element is present. Commonly set to xml.HTMLAutoClose.
	// Default: nil.
	AutoClose []string
}

// defaultCharsetReader is used by the xml decoder when the ReadSettings
// CharsetReader value is nil. It behaves as a "pass-through", ignoring
// the requested charset parameter and skipping conversion altogether.
func defaultCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	_ = "STUB: not implemented"

	// dup creates a duplicate of the ReadSettings object.
	return *new(io.Reader), nil
}

func (s *ReadSettings) dup() ReadSettings { _ = "STUB: not implemented"; return *new(ReadSettings) }

// WriteSettings determine the behavior of the Document's WriteTo* functions.
type WriteSettings struct {
	// CanonicalEndTags forces the production of XML end tags, even for
	// elements that have no child elements. Default: false.
	CanonicalEndTags bool

	// CanonicalText forces the production of XML character references for
	// text data characters &, <, and >. If false, XML character references
	// are also produced for " and '. Default: false.
	CanonicalText bool

	// CanonicalAttrVal forces the production of XML character references for
	// attribute value characters &, < and ". If false, XML character
	// references are also produced for > and '. Ignored when AttrSingleQuote
	// is true. Default: false.
	CanonicalAttrVal bool

	// AttrSingleQuote causes attributes to use single quotes (attr='example')
	// instead of double quotes (attr = "example") when set to true. Default:
	// false.
	AttrSingleQuote bool

	// UseCRLF causes the document's Indent* functions to use a carriage return
	// followed by a linefeed ("\r\n") when outputting a newline. If false,
	// only a linefeed is used ("\n"). Default: false.
	//
	// Deprecated: UseCRLF is deprecated. Use IndentSettings.UseCRLF instead.
	UseCRLF bool
}

// dup creates a duplicate of the WriteSettings object.
func (s *WriteSettings) dup() WriteSettings {
	_ = "STUB: not implemented"

	// IndentSettings determine the behavior of the Document's Indent* functions.
	return *new(WriteSettings)
}

type IndentSettings struct {
	// Spaces indicates the number of spaces to insert for each level of
	// indentation. Set to etree.NoIndent to remove all indentation. Ignored
	// when UseTabs is true. Default: 4.
	Spaces int

	// UseTabs causes tabs to be used instead of spaces when indenting.
	// Default: false.
	UseTabs bool

	// UseCRLF causes newlines to be written as a carriage return followed by
	// a linefeed ("\r\n"). If false, only a linefeed character is output
	// for a newline ("\n"). Default: false.
	UseCRLF bool

	// PreserveLeafWhitespace causes indent functions to preserve whitespace
	// within XML elements containing only non-CDATA character data. Default:
	// false.
	PreserveLeafWhitespace bool

	// SuppressTrailingWhitespace suppresses the generation of a trailing
	// whitespace characters (such as newlines) at the end of the indented
	// document. Default: false.
	SuppressTrailingWhitespace bool
}

// NewIndentSettings creates a default IndentSettings record.
func NewIndentSettings() *IndentSettings { _ = "STUB: not implemented"; return nil }

type indentFunc func(depth int) string

func getIndentFunc(s *IndentSettings) indentFunc {
	_ = "STUB: not implemented"
	return *new(indentFunc)
}

// Writer is the interface that wraps the Write* functions called by each token
// type's WriteTo function.
type Writer interface {
	io.StringWriter
	io.ByteWriter
	io.Writer
}

// A Token is an interface type used to represent XML elements, character
// data, CDATA sections, XML comments, XML directives, and XML processing
// instructions.
type Token interface {
	Parent() *Element
	Index() int
	WriteTo(w Writer, s *WriteSettings)
	dup(parent *Element) Token
	setParent(parent *Element)
	setIndex(index int)
}

// A Document is a container holding a complete XML tree.
//
// A document has a single embedded element, which contains zero or more child
// tokens, one of which is usually the root element. The embedded element may
// include other children such as processing instruction tokens or character
// data tokens. The document's embedded element is never directly serialized;
// only its children are.
//
// A document also contains read and write settings, which influence the way
// the document is deserialized, serialized, and indented.
type Document struct {
	Element
	ReadSettings  ReadSettings
	WriteSettings WriteSettings
}

// An Element represents an XML element, its attributes, and its child tokens.
type Element struct {
	Space, Tag string   // namespace prefix and tag
	Attr       []Attr   // key-value attribute pairs
	Child      []Token  // child tokens (elements, comments, etc.)
	parent     *Element // parent element
	index      int      // token index in parent's children
}

// An Attr represents a key-value attribute within an XML element.
type Attr struct {
	Space, Key string   // The attribute's namespace prefix and key
	Value      string   // The attribute value string
	element    *Element // element containing the attribute
}

// charDataFlags are used with CharData tokens to store additional settings.
type charDataFlags uint8

const (
	// The CharData contains only whitespace.
	whitespaceFlag charDataFlags = 1 << iota

	// The CharData contains a CDATA section.
	cdataFlag
)

// CharData may be used to represent simple text data or a CDATA section
// within an XML document. The Data property should never be modified
// directly; use the SetData function instead.
type CharData struct {
	Data   string // the simple text or CDATA section content
	parent *Element
	index  int
	flags  charDataFlags
}

// A Comment represents an XML comment.
type Comment struct {
	Data   string // the comment's text
	parent *Element
	index  int
}

// A Directive represents an XML directive.
type Directive struct {
	Data   string // the directive string
	parent *Element
	index  int
}

// A ProcInst represents an XML processing instruction.
type ProcInst struct {
	Target string // the processing instruction target
	Inst   string // the processing instruction value
	parent *Element
	index  int
}

// NewDocument creates an XML document without a root element.
func NewDocument() *Document { _ = "STUB: not implemented"; return nil }

// NewDocumentWithRoot creates an XML document and sets the element 'e' as its
// root element. If the element 'e' is already part of another document, it is
// first removed from its existing document.
func NewDocumentWithRoot(e *Element) *Document { _ = "STUB: not implemented"; return nil }

// Copy returns a recursive, deep copy of the document.
func (d *Document) Copy() *Document { _ = "STUB: not implemented"; return nil }

// Root returns the root element of the document. It returns nil if there is
// no root element.
func (d *Document) Root() *Element { _ = "STUB: not implemented"; return nil }

// SetRoot replaces the document's root element with the element 'e'. If the
// document already has a root element when this function is called, then the
// existing root element is unbound from the document. If the element 'e' is
// part of another document, then it is unbound from the other document.
func (d *Document) SetRoot(e *Element) { _ = "STUB: not implemented"; return }

// If there is already a root element, replace it.

// No existing root element, so add it.

// ReadFrom reads XML from the reader 'r' into this document. The function
// returns the number of bytes read and any error encountered.
func (d *Document) ReadFrom(r io.Reader) (n int64, err error) {
	_ = "STUB: not implemented"
	return 0, nil
}

// ReadFromFile reads XML from a local file at path 'filepath' into this
// document.
func (d *Document) ReadFromFile(filepath string) error { _ = "STUB: not implemented"; return nil }

// ReadFromBytes reads XML from the byte slice 'b' into the this document.
func (d *Document) ReadFromBytes(b []byte) error { _ = "STUB: not implemented"; return nil }

// ReadFromString reads XML from the string 's' into this document.
func (d *Document) ReadFromString(s string) error { _ = "STUB: not implemented"; return nil }

// validateXML determines if the data read from the reader 'r' contains
// well-formed XML according to the rules set by the go xml package.
func validateXML(r io.Reader, settings ReadSettings) error { _ = "STUB: not implemented"; return nil }

// If there are any trailing tokens after unmarshalling with Decode(),
// then the XML input didn't terminate properly.

// newDecoder creates an XML decoder for the reader 'r' configured using
// the provided read settings.
func newDecoder(r io.Reader, settings ReadSettings) *xml.Decoder {
	_ = "STUB: not implemented"
	return nil
}

// WriteTo serializes the document out to the writer 'w'. The function returns
// the number of bytes written and any error encountered.
func (d *Document) WriteTo(w io.Writer) (n int64, err error) {
	_ = "STUB: not implemented"
	return 0, nil
}

// WriteToFile serializes the document out to the file at path 'filepath'.
func (d *Document) WriteToFile(filepath string) error { _ = "STUB: not implemented"; return nil }

// WriteToBytes serializes this document into a slice of bytes.
func (d *Document) WriteToBytes() (b []byte, err error) { _ = "STUB: not implemented"; return nil, nil }

// WriteToString serializes this document into a string.
func (d *Document) WriteToString() (s string, err error) { _ = "STUB: not implemented"; return "", nil }

// Indent modifies the document's element tree by inserting character data
// tokens containing newlines and spaces for indentation. The amount of
// indentation per depth level is given by the 'spaces' parameter. Other than
// the number of spaces, default IndentSettings are used.
func (d *Document) Indent(spaces int) { _ = "STUB: not implemented"; return }

// IndentTabs modifies the document's element tree by inserting CharData
// tokens containing newlines and tabs for indentation. One tab is used per
// indentation level. Other than the use of tabs, default IndentSettings
// are used.
func (d *Document) IndentTabs() { _ = "STUB: not implemented"; return }

// IndentWithSettings modifies the document's element tree by inserting
// character data tokens containing newlines and indentation. The behavior
// of the indentation algorithm is configured by the indent settings.
func (d *Document) IndentWithSettings(s *IndentSettings) {
	_ = "STUB: not implemented"
	// WriteSettings.UseCRLF is deprecated. Until removed from the package, it
	// overrides IndentSettings.UseCRLF when true.
	return
}

// Unindent modifies the document's element tree by removing character data
// tokens containing only whitespace. Other than the removal of indentation,
// default IndentSettings are used.
func (d *Document) Unindent() { _ = "STUB: not implemented"; return }

// NewElement creates an unparented element with the specified tag (i.e.,
// name). The tag may include a namespace prefix followed by a colon.
func NewElement(tag string) *Element { _ = "STUB: not implemented"; return nil }

// newElement is a helper function that creates an element and binds it to
// a parent element if possible.
func newElement(space, tag string, parent *Element) *Element { _ = "STUB: not implemented"; return nil }

// Copy creates a recursive, deep copy of the element and all its attributes
// and children. The returned element has no parent but can be parented to a
// another element using AddChild, or added to a document with SetRoot or
// NewDocumentWithRoot.
func (e *Element) Copy() *Element { _ = "STUB: not implemented"; return nil }

// FullTag returns the element e's complete tag, including namespace prefix if
// present.
func (e *Element) FullTag() string { _ = "STUB: not implemented"; return "" }

// NamespaceURI returns the XML namespace URI associated with the element. If
// the element is part of the XML default namespace, NamespaceURI returns the
// empty string.
func (e *Element) NamespaceURI() string { _ = "STUB: not implemented"; return "" }

// findLocalNamespaceURI finds the namespace URI corresponding to the
// requested prefix.
func (e *Element) findLocalNamespaceURI(prefix string) string { _ = "STUB: not implemented"; return "" }

// findDefaultNamespaceURI finds the default namespace URI of the element.
func (e *Element) findDefaultNamespaceURI() string { _ = "STUB: not implemented"; return "" }

// namespacePrefix returns the namespace prefix associated with the element.
func (e *Element) namespacePrefix() string {
	_ = "STUB: not implemented"

	// name returns the tag associated with the element.
	return ""
}

func (e *Element) name() string {
	_ = "STUB: not implemented"

	// ReindexChildren recalculates the index values of the element's child
	// tokens. This is necessary only if you have manually manipulated the
	// element's `Child` array.
	return ""
}

func (e *Element) ReindexChildren() { _ = "STUB: not implemented"; return }

// Text returns all character data immediately following the element's opening
// tag.
func (e *Element) Text() string { _ = "STUB: not implemented"; return "" }

// ignore

// SetText replaces all character data immediately following an element's
// opening tag with the requested string.
func (e *Element) SetText(text string) { _ = "STUB: not implemented"; return }

// SetCData replaces all character data immediately following an element's
// opening tag with a CDATA section.
func (e *Element) SetCData(text string) { _ = "STUB: not implemented"; return }

// Tail returns all character data immediately following the element's end
// tag.
func (e *Element) Tail() string { _ = "STUB: not implemented"; return "" }

// SetTail replaces all character data immediately following the element's end
// tag with the requested string.
func (e *Element) SetTail(text string) { _ = "STUB: not implemented"; return }

// replaceText is a helper function that replaces a series of chardata tokens
// starting at index i with the requested text.
func (e *Element) replaceText(i int, text string, flags charDataFlags) {
	_ = "STUB: not implemented"
	return
}

// insert a new chardata token at index i

// remove the chardata token at index i

// replace the first and only character token at index i

// remove all chardata tokens starting from index i

// replace the first chardata token at index i and remove all
// subsequent chardata tokens

// findTermCharDataIndex finds the index of the first child token that isn't
// a CharData token. It starts from the requested start index.
func (e *Element) findTermCharDataIndex(start int) int { _ = "STUB: not implemented"; return 0 }

// CreateElement creates a new element with the specified tag (i.e., name) and
// adds it as the last child of element 'e'. The tag may include a prefix
// followed by a colon.
func (e *Element) CreateElement(tag string) *Element { _ = "STUB: not implemented"; return nil }

// CreateChild performs the same task as CreateElement but calls a
// continuation function after the child element is created, allowing
// additional actions to be performed on the child element before returning.
//
// This method of element creation is particularly useful when building nested
// XML documents from code. For example:
//
//	org := doc.CreateChild("organization", func(e *Element) {
//		e.CreateComment("Mary")
//		e.CreateChild("person", func(e *Element) {
//			e.CreateAttr("name", "Mary")
//			e.CreateAttr("age", "30")
//			e.CreateAttr("hair", "brown")
//		})
//	})
func (e *Element) CreateChild(tag string, cont func(e *Element)) *Element {
	_ = "STUB: not implemented"
	return nil
}

// AddChild adds the token 't' as the last child of the element. If token 't'
// was already the child of another element, it is first removed from its
// parent element.
func (e *Element) AddChild(t Token) { _ = "STUB: not implemented"; return }

// InsertChild inserts the token 't' into this element's list of children just
// before the element's existing child token 'ex'. If the existing element
// 'ex' does not appear in this element's list of child tokens, then 't' is
// added to the end of this element's list of child tokens. If token 't' is
// already the child of another element, it is first removed from the other
// element's list of child tokens.
//
// Deprecated: InsertChild is deprecated. Use InsertChildAt instead.
func (e *Element) InsertChild(ex Token, t Token) { _ = "STUB: not implemented"; return }

// InsertChildAt inserts the token 't' into this element's list of child
// tokens just before the requested 'index'. If the index is greater than or
// equal to the length of the list of child tokens, then the token 't' is
// added to the end of the list of child tokens.
func (e *Element) InsertChildAt(index int, t Token) { _ = "STUB: not implemented"; return }

// RemoveChild attempts to remove the token 't' from this element's list of
// child tokens. If the token 't' was a child of this element, then it is
// removed and returned. Otherwise, nil is returned.
func (e *Element) RemoveChild(t Token) Token { _ = "STUB: not implemented"; return *new(Token) }

// RemoveChildAt removes the child token appearing in slot 'index' of this
// element's list of child tokens. The removed child token is then returned.
// If the index is out of bounds, no child is removed and nil is returned.
func (e *Element) RemoveChildAt(index int) Token { _ = "STUB: not implemented"; return *new(Token) }

// autoClose analyzes the stack's top element and the current token to decide
// whether the top element should be closed.
func (e *Element) autoClose(stack *stack[*Element], t xml.Token, tags []string) {
	_ = "STUB: not implemented"
	return
}

// ReadFrom reads XML from the reader 'ri' and stores the result as a new
// child of this element.
func (e *Element) readFrom(ri io.Reader, settings ReadSettings) (n int64, err error) {
	_ = "STUB: not implemented"
	return 0, nil
}

// SelectAttr finds an element attribute matching the requested 'key' and, if
// found, returns a pointer to the matching attribute. The function returns
// nil if no matching attribute is found. The key may include a namespace
// prefix followed by a colon.
func (e *Element) SelectAttr(key string) *Attr { _ = "STUB: not implemented"; return nil }

// SelectAttrValue finds an element attribute matching the requested 'key' and
// returns its value if found. If no matching attribute is found, the function
// returns the 'dflt' value instead. The key may include a namespace prefix
// followed by a colon.
func (e *Element) SelectAttrValue(key, dflt string) string { _ = "STUB: not implemented"; return "" }

// ChildElements returns all elements that are children of this element.
func (e *Element) ChildElements() []*Element { _ = "STUB: not implemented"; return nil }

// ChildElementsSeq returns an iterator over all child elements of this
// element.
func (e *Element) ChildElementsSeq() iter.Seq[*Element] { _ = "STUB: not implemented"; return nil }

// SelectElement returns the first child element with the given 'tag' (i.e.,
// name). The function returns nil if no child element matching the tag is
// found. The tag may include a namespace prefix followed by a colon.
func (e *Element) SelectElement(tag string) *Element { _ = "STUB: not implemented"; return nil }

// SelectElements returns a slice of all child elements with the given 'tag'
// (i.e., name). The tag may include a namespace prefix followed by a colon.
func (e *Element) SelectElements(tag string) []*Element { _ = "STUB: not implemented"; return nil }

// SelectElementsSeq returns an iterator over all child elements with the
// given 'tag' (i.e., name). The tag may include a namespace prefix followed
// by a colon.
func (e *Element) SelectElementsSeq(tag string) iter.Seq[*Element] {
	_ = "STUB: not implemented"
	return nil
}

// FindElement returns the first element matched by the XPath-like 'path'
// string. The function returns nil if no child element is found using the
// path. It panics if an invalid path string is supplied.
func (e *Element) FindElement(path string) *Element { _ = "STUB: not implemented"; return nil }

// FindElementPath returns the first element matched by the 'path' object. The
// function returns nil if no element is found using the path.
func (e *Element) FindElementPath(path Path) *Element { _ = "STUB: not implemented"; return nil }

// FindElements returns a slice of elements matched by the XPath-like 'path'
// string. The function returns nil if no child element is found using the
// path. It panics if an invalid path string is supplied.
func (e *Element) FindElements(path string) []*Element { _ = "STUB: not implemented"; return nil }

// FindElementsSeq returns an iterator over elements matched by the XPath-like
// 'path' string. This function uses Go's iterator support for
// memory-efficient traversal. It panics if an invalid path string is
// supplied.
func (e *Element) FindElementsSeq(path string) iter.Seq[*Element] {
	_ = "STUB: not implemented"
	return nil
}

// FindElementsPath returns a slice of elements matched by the 'path' object.
func (e *Element) FindElementsPath(path Path) []*Element { _ = "STUB: not implemented"; return nil }

// FindElementsPathSeq returns an iterator over elements matched by the 'path'
// object.
func (e *Element) FindElementsPathSeq(path Path) iter.Seq[*Element] {
	_ = "STUB: not implemented"
	return nil

	// NotNil returns the receiver element if it isn't nil; otherwise, it returns
	// an unparented element with an empty string tag. This function simplifies
	// the task of writing code to ignore not-found results from element queries.
	// For example, instead of writing this:
	//
	//	if e := doc.SelectElement("enabled"); e != nil {
	//		e.SetText("true")
	//	}
	//
	// You could write this:
	//
	//	doc.SelectElement("enabled").NotNil().SetText("true")
}

func (e *Element) NotNil() *Element { _ = "STUB: not implemented"; return nil }

// GetPath returns the absolute path of the element. The absolute path is the
// full path from the document's root.
func (e *Element) GetPath() string { _ = "STUB: not implemented"; return "" }

// Reverse the path.

// GetRelativePath returns the path of this element relative to the 'source'
// element. If the two elements are not part of the same element tree, then
// the function returns the empty string.
func (e *Element) GetRelativePath(source *Element) string { _ = "STUB: not implemented"; return "" }

// Build a reverse path from the element toward the root. Stop if the
// source element is encountered.

// If we found the source element, reverse the path and compose the
// string.

// The source wasn't encountered, so climb from the source element toward
// the root of the tree until an element in the reversed path is
// encountered.

// truncate at found segment

// No element in the reversed path was encountered, so the two elements
// must not be part of the same tree.

// Reverse the (possibly truncated) path and prepend ".." segments to
// climb.

// IndentWithSettings modifies the element and its child tree by inserting
// character data tokens containing newlines and indentation. The behavior of
// the indentation algorithm is configured by the indent settings. Because
// this function indents the element as if it were at the root of a document,
// it is most useful when called just before writing the element as an XML
// fragment using WriteTo.
func (e *Element) IndentWithSettings(s *IndentSettings) { _ = "STUB: not implemented"; return }

// indent recursively inserts proper indentation between an XML element's
// child tokens.
func (e *Element) indent(depth int, indent indentFunc, s *IndentSettings) {
	_ = "STUB: not implemented"
	return
}

// Insert NL+indent before child if it's not character data.
// Exceptions: when it's the first non-character-data child, or when
// the child is at root depth.

// Recursively process child elements.

// Insert NL+indent before the last child.

// stripIndent removes any previously inserted indentation.
func (e *Element) stripIndent(s *IndentSettings) {
	_ = "STUB: not implemented"
	// Count the number of non-indent child tokens
	return
}

// Strip out indent CharData

// stripTrailingWhitespace removes any trailing whitespace CharData tokens
// from the element's children.
func (e *Element) stripTrailingWhitespace() { _ = "STUB: not implemented"; return }

// dup duplicates the element.
func (e *Element) dup(parent *Element) Token { _ = "STUB: not implemented"; return *new(Token) }

// NextSibling returns this element's next sibling element. It returns nil if
// there is no next sibling element.
func (e *Element) NextSibling() *Element { _ = "STUB: not implemented"; return nil }

// PrevSibling returns this element's preceding sibling element. It returns
// nil if there is no preceding sibling element.
func (e *Element) PrevSibling() *Element { _ = "STUB: not implemented"; return nil }

// Parent returns this element's parent element. It returns nil if this
// element has no parent.
func (e *Element) Parent() *Element {
	_ = "STUB: not implemented"

	// Index returns the index of this element within its parent element's
	// list of child tokens. If this element has no parent, then the function
	// returns -1.
	return nil
}

func (e *Element) Index() int {
	_ = "STUB: not implemented"

	// WriteTo serializes the element to the writer w.
	return 0
}

func (e *Element) WriteTo(w Writer, s *WriteSettings) { _ = "STUB: not implemented"; return }

// setParent replaces this element token's parent.
func (e *Element) setParent(parent *Element) {
	_ = "STUB: not implemented"

	// setIndex sets this element token's index within its parent's Child slice.
	return
}

func (e *Element) setIndex(index int) {
	_ = "STUB: not implemented"

	// addChild adds a child token to the element e.
	return
}

func (e *Element) addChild(t Token) { _ = "STUB: not implemented"; return }

// CreateAttr creates an attribute with the specified 'key' and 'value' and
// adds it to this element. If an attribute with same key already exists on
// this element, then its value is replaced. The key may include a namespace
// prefix followed by a colon.
func (e *Element) CreateAttr(key, value string) *Attr { _ = "STUB: not implemented"; return nil }

// addAttr is a helper function that adds an attribute to an element. Returns
// the index of the added attribute.
func (e *Element) addAttr(space, key, value string) int { _ = "STUB: not implemented"; return 0 }

// RemoveAttr removes the first attribute of this element whose key matches
// 'key'. It returns a copy of the removed attribute if a match is found. If
// no match is found, it returns nil. The key may include a namespace prefix
// followed by a colon.
func (e *Element) RemoveAttr(key string) *Attr { _ = "STUB: not implemented"; return nil }

// SortAttrs sorts this element's attributes lexicographically by key.
func (e *Element) SortAttrs() { _ = "STUB: not implemented"; return }

// FullKey returns this attribute's complete key, including namespace prefix
// if present.
func (a *Attr) FullKey() string { _ = "STUB: not implemented"; return "" }

// Element returns a pointer to the element containing this attribute.
func (a *Attr) Element() *Element {
	_ = "STUB: not implemented"

	// NamespaceURI returns the XML namespace URI associated with this attribute.
	// The function returns the empty string if the attribute is unprefixed or
	// if the attribute is part of the XML default namespace.
	return nil
}

func (a *Attr) NamespaceURI() string { _ = "STUB: not implemented"; return "" }

// WriteTo serializes the attribute to the writer.
func (a *Attr) WriteTo(w Writer, s *WriteSettings) { _ = "STUB: not implemented"; return }

// NewText creates an unparented CharData token containing simple text data.
func NewText(text string) *CharData { _ = "STUB: not implemented"; return nil }

// NewCData creates an unparented XML character CDATA section with 'data' as
// its content.
func NewCData(data string) *CharData { _ = "STUB: not implemented"; return nil }

// NewCharData creates an unparented CharData token containing simple text
// data.
//
// Deprecated: NewCharData is deprecated. Instead, use NewText, which does the
// same thing.
func NewCharData(data string) *CharData { _ = "STUB: not implemented"; return nil }

// newCharData creates a character data token and binds it to a parent
// element. If parent is nil, the CharData token remains unbound.
func newCharData(data string, flags charDataFlags, parent *Element) *CharData {
	_ = "STUB: not implemented"
	return nil
}

// CreateText creates a CharData token containing simple text data and adds it
// to the end of this element's list of child tokens.
func (e *Element) CreateText(text string) *CharData { _ = "STUB: not implemented"; return nil }

// CreateCData creates a CharData token containing a CDATA section with 'data'
// as its content and adds it to the end of this element's list of child
// tokens.
func (e *Element) CreateCData(data string) *CharData { _ = "STUB: not implemented"; return nil }

// CreateCharData creates a CharData token containing simple text data and
// adds it to the end of this element's list of child tokens.
//
// Deprecated: CreateCharData is deprecated. Instead, use CreateText, which
// does the same thing.
func (e *Element) CreateCharData(data string) *CharData { _ = "STUB: not implemented"; return nil }

// SetData modifies the content of the CharData token. In the case of a
// CharData token containing simple text, the simple text is modified. In the
// case of a CharData token containing a CDATA section, the CDATA section's
// content is modified.
func (c *CharData) SetData(text string) { _ = "STUB: not implemented"; return }

// IsCData returns true if this CharData token is contains a CDATA section. It
// returns false if the CharData token contains simple text.
func (c *CharData) IsCData() bool { _ = "STUB: not implemented"; return false }

// IsWhitespace returns true if this CharData token contains only whitespace.
func (c *CharData) IsWhitespace() bool { _ = "STUB: not implemented"; return false }

// Parent returns this CharData token's parent element, or nil if it has no
// parent.
func (c *CharData) Parent() *Element {
	_ = "STUB: not implemented"

	// Index returns the index of this CharData token within its parent element's
	// list of child tokens. If this CharData token has no parent, then the
	// function returns -1.
	return nil
}

func (c *CharData) Index() int {
	_ = "STUB: not implemented"

	// WriteTo serializes character data to the writer.
	return 0
}

func (c *CharData) WriteTo(w Writer, s *WriteSettings) { _ = "STUB: not implemented"; return }

// dup duplicates the character data.
func (c *CharData) dup(parent *Element) Token { _ = "STUB: not implemented"; return *new(Token) }

// setParent replaces the character data token's parent.
func (c *CharData) setParent(parent *Element) {
	_ = "STUB: not implemented"

	// setIndex sets the CharData token's index within its parent element's Child
	// slice.
	return
}

func (c *CharData) setIndex(index int) {
	_ = "STUB: not implemented"

	// NewComment creates an unparented comment token.
	return
}

func NewComment(comment string) *Comment { _ = "STUB: not implemented"; return nil }

// NewComment creates a comment token and sets its parent element to 'parent'.
func newComment(comment string, parent *Element) *Comment { _ = "STUB: not implemented"; return nil }

// CreateComment creates a comment token using the specified 'comment' string
// and adds it as the last child token of this element.
func (e *Element) CreateComment(comment string) *Comment { _ = "STUB: not implemented"; return nil }

// dup duplicates the comment.
func (c *Comment) dup(parent *Element) Token { _ = "STUB: not implemented"; return *new(Token) }

// Parent returns comment token's parent element, or nil if it has no parent.
func (c *Comment) Parent() *Element {
	_ = "STUB: not implemented"

	// Index returns the index of this Comment token within its parent element's
	// list of child tokens. If this Comment token has no parent, then the
	// function returns -1.
	return nil
}

func (c *Comment) Index() int {
	_ = "STUB: not implemented"

	// WriteTo serialies the comment to the writer.
	return 0
}

func (c *Comment) WriteTo(w Writer, s *WriteSettings) { _ = "STUB: not implemented"; return }

// setParent replaces the comment token's parent.
func (c *Comment) setParent(parent *Element) {
	_ = "STUB: not implemented"

	// setIndex sets the Comment token's index within its parent element's Child
	// slice.
	return
}

func (c *Comment) setIndex(index int) {
	_ = "STUB: not implemented"

	// NewDirective creates an unparented XML directive token.
	return
}

func NewDirective(data string) *Directive { _ = "STUB: not implemented"; return nil }

// newDirective creates an XML directive and binds it to a parent element. If
// parent is nil, the Directive remains unbound.
func newDirective(data string, parent *Element) *Directive { _ = "STUB: not implemented"; return nil }

// CreateDirective creates an XML directive token with the specified 'data'
// value and adds it as the last child token of this element.
func (e *Element) CreateDirective(data string) *Directive { _ = "STUB: not implemented"; return nil }

// dup duplicates the directive.
func (d *Directive) dup(parent *Element) Token { _ = "STUB: not implemented"; return *new(Token) }

// Parent returns directive token's parent element, or nil if it has no
// parent.
func (d *Directive) Parent() *Element {
	_ = "STUB: not implemented"

	// Index returns the index of this Directive token within its parent element's
	// list of child tokens. If this Directive token has no parent, then the
	// function returns -1.
	return nil
}

func (d *Directive) Index() int {
	_ = "STUB: not implemented"

	// WriteTo serializes the XML directive to the writer.
	return 0
}

func (d *Directive) WriteTo(w Writer, s *WriteSettings) { _ = "STUB: not implemented"; return }

// setParent replaces the directive token's parent.
func (d *Directive) setParent(parent *Element) {
	_ = "STUB: not implemented"

	// setIndex sets the Directive token's index within its parent element's Child
	// slice.
	return
}

func (d *Directive) setIndex(index int) {
	_ = "STUB: not implemented"

	// NewProcInst creates an unparented XML processing instruction.
	return
}

func NewProcInst(target, inst string) *ProcInst { _ = "STUB: not implemented"; return nil }

// newProcInst creates an XML processing instruction and binds it to a parent
// element. If parent is nil, the ProcInst remains unbound.
func newProcInst(target, inst string, parent *Element) *ProcInst {
	_ = "STUB: not implemented"
	return nil
}

// CreateProcInst creates an XML processing instruction token with the
// specified 'target' and instruction 'inst'. It is then added as the last
// child token of this element.
func (e *Element) CreateProcInst(target, inst string) *ProcInst {
	_ = "STUB: not implemented"
	return nil
}

// dup duplicates the procinst.
func (p *ProcInst) dup(parent *Element) Token { _ = "STUB: not implemented"; return *new(Token) }

// Parent returns processing instruction token's parent element, or nil if it
// has no parent.
func (p *ProcInst) Parent() *Element {
	_ = "STUB: not implemented"

	// Index returns the index of this ProcInst token within its parent element's
	// list of child tokens. If this ProcInst token has no parent, then the
	// function returns -1.
	return nil
}

func (p *ProcInst) Index() int {
	_ = "STUB: not implemented"

	// WriteTo serializes the processing instruction to the writer.
	return 0
}

func (p *ProcInst) WriteTo(w Writer, s *WriteSettings) { _ = "STUB: not implemented"; return }

// setParent replaces the processing instruction token's parent.
func (p *ProcInst) setParent(parent *Element) {
	_ = "STUB: not implemented"

	// setIndex sets the processing instruction token's index within its parent
	// element's Child slice.
	return
}

func (p *ProcInst) setIndex(index int) { _ = "STUB: not implemented"; return }
