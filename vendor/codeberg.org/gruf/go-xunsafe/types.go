package xunsafe

import (
	"reflect"
)

// TypeIter provides a simple wrapper for
// a means of following reflected types.
type TypeIter struct {
	TypeInfo
	Parent *TypeIter
}

// TypeInfo wraps reflect type information
// along with flags specifying further details
// necessary due to type nesting.
type TypeInfo struct {
	Type reflect.Type
	Flag Reflect_flag
}

// ToTypeIter creates a new TypeIter{} from reflect type and flags.
func ToTypeIter(rtype reflect.Type, flags Reflect_flag) TypeIter {
	return TypeIter{TypeInfo: TypeInfo{rtype, flags}}
}

// TypeIterFrom creates new TypeIter from interface value type.
// Note this will always assume the initial value passed to you
// will be coming from an interface.
func TypeIterFrom(a any) TypeIter {
	rtype := reflect.TypeOf(a)
	flags := ReflectIfaceElemFlags(0, rtype)
	return ToTypeIter(rtype, flags)
}

// Indirect returns whether Reflect_flagIndir is set on receiving TypeInfo{}.Flag.
func (t TypeInfo) Indirect() bool { return t.Flag.Indirect() }

// IfaceIndir calls Abi_Type_IfaceIndir() on receiving TypeInfo{}.Type.
func (t TypeInfo) IfaceIndir() bool { return Abi_Type_IfaceIndir(t.Type) }

// PointerElem returns the reflect.Type{}.Elem() with
// pointer elem flags as child of receiving TypeIter{}.
func (i TypeIter) PointerElem() TypeIter {
	e := i.Type.Elem()
	f := ReflectPointerElemFlags(i.Flag, e)
	return i.child(e, f)
}

// SliceElem returns the reflect.Type{}.Elem() with
// slice elem flags as child of receiving TypeIter{}.
func (i TypeIter) SliceElem() TypeIter {
	e := i.Type.Elem()
	f := ReflectSliceElemFlags(i.Flag, e)
	return i.child(e, f)
}

// ArrayElem returns the reflect.Type{}.Elem() with
// array elem flags as child of receiving TypeIter{}.
func (i TypeIter) ArrayElem() TypeIter {
	e := i.Type.Elem()
	f := ReflectArrayElemFlags(i.Flag, i.Type, e)
	return i.child(e, f)
}

// MapKey returns the reflect.Type{}.Key() with
// map key flags as child of receiving TypeIter{}.
func (i TypeIter) MapKey() TypeIter {
	k := i.Type.Key()
	f := ReflectMapKeyFlags(i.Flag, k)
	return i.child(k, f)
}

// MapElem returns the reflect.Type{}.Elem() with
// map elem flags as child of receiving TypeIter{}.
func (i TypeIter) MapElem() TypeIter {
	e := i.Type.Elem()
	f := ReflectMapElemFlags(i.Flag, e)
	return i.child(e, f)
}

// StructField returns the reflect.Type{}.Field() with
// struct field flags as child of receiving TypeIter{}.
func (i TypeIter) StructField(idx int) (TypeIter, reflect.StructField) {
	s := i.Type.Field(idx)
	f := ReflectStructFieldFlags(i.Flag, i.Type, s)
	return i.child(s.Type, f), s
}

// StructField returns the reflect.Type{}.FieldByName()
// with struct field flags as child of receiving TypeIter{}.
func (i TypeIter) StructFieldByName(name string) (TypeIter, reflect.StructField, bool) {
	s, ok := i.Type.FieldByName(name)
	if !ok {
		return TypeIter{}, reflect.StructField{}, false
	}
	f := ReflectStructFieldFlags(i.Flag, i.Type, s)
	return i.child(s.Type, f), s, true
}

// child returns a new TypeIter{} for given type and flags, with parent pointing to receiver.
func (i TypeIter) child(rtype reflect.Type, flags Reflect_flag) TypeIter {
	child := ToTypeIter(rtype, flags)
	child.Parent = &i
	return child
}
