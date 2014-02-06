// Copyright 2014 SteelSeries ApS.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package implements a basic LISP interpretor for embedding in a go program for scripting.
// This file contains the list manipulation primitive functions.

package golisp

import (
    "errors"
)

func RegisterListManipulationPrimitives() {
    MakePrimitiveFunction("list", -1, MakeListImpl)
    MakePrimitiveFunction("length", 1, ListLengthImpl)
    MakePrimitiveFunction("cons", 2, ConsImpl)
    MakePrimitiveFunction("reverse", 1, ReverseImpl)
    MakePrimitiveFunction("flatten", 1, FlattenImpl)
    MakePrimitiveFunction("flatten*", 1, RecursiveFlattenImpl)
    MakePrimitiveFunction("append", 2, AppendImpl)
    MakePrimitiveFunction("append!", 2, AppendBangImpl)
    MakePrimitiveFunction("copy", 1, CopyImpl)
    MakePrimitiveFunction("partition", 2, PartitionImpl)
    MakePrimitiveFunction("sublist", 3, SublistImpl)
    MakePrimitiveFunction("take", 2, TakeImpl)
    MakePrimitiveFunction("drop", 2, DropImpl)
}

func MakeListImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    var items []*Data = make([]*Data, 0, Length(args))
    var item *Data
    for cell := args; NotNilP(cell); cell = Cdr(cell) {
        item, err = Eval(Car(cell), env)
        if err != nil {
            return
        }
        items = append(items, item)
    }
    result = ArrayToList(items)
    return
}

func ListLengthImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    d, err := Eval(Car(args), env)
    if err != nil {
        return
    }
    return NumberWithValue(uint32(Length(d))), nil
}

func ConsImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    var car *Data
    car, err = Eval(Car(args), env)
    if err != nil {
        return
    }

    var cdr *Data
    cdr, err = Eval(Cadr(args), env)
    if err != nil {
        return
    }

    result = Cons(car, cdr)
    return
}

func ReverseImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    var val *Data
    val, err = Eval(Car(args), env)
    if err != nil {
        return
    }
    result = Reverse(val)
    return
}

func FlattenImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    var val *Data
    val, err = Eval(Car(args), env)
    if err != nil {
        return
    }
    result, err = Flatten(val)
    return
}

func RecursiveFlattenImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    var val *Data
    val, err = Eval(Car(args), env)
    if err != nil {
        return
    }
    result, err = RecursiveFlatten(val)
    return
}

func AppendBangImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    firstList, err := Eval(Car(args), env)
    if err != nil {
        return
    }

    secondList, err := Eval(Cadr(args), env)
    if err != nil {
        return
    }

    result = AppendBangList(firstList, secondList)

    if SymbolP(Car(args)) {
        env.BindTo(Car(args), result)
    }

    return
}

func AppendImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    firstList, err := Eval(Car(args), env)
    if err != nil {
        return
    }

    secondList, err := Eval(Cadr(args), env)
    if err != nil {
        return
    }

    result = AppendList(Copy(firstList), secondList)
    return
}

func CopyImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    d, err := Eval(Car(args), env)
    if err != nil {
        return
    }

    return Copy(d), nil
}

func PartitionImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    n, err := Eval(Car(args), env)
    if err != nil {
        return
    } 
    if !NumberP(n) {
        err = errors.New("partition requires a number as it's first argument.")
    }
    size := int(NumericValue(n))

    l, err := Eval(Cadr(args), env)
    if err != nil {
        return
    }
    if !ListP(l) {
        err = errors.New("partition requires a list as it's second argument.")
    }

    var pieces []*Data = make([]*Data, 0, 5)
    var chunk []*Data = make([]*Data, 0, 5)
    for c := l; NotNilP(c); c = Cdr(c) {
        if len(chunk) < size {
            chunk = append(chunk, Car(c))
        } else {
            pieces = append(pieces, ArrayToList(chunk))
            chunk = make([]*Data, 0, 5)
            chunk = append(chunk, Car(c))
        }
    }
    if len(chunk) > 0 {
        pieces = append(pieces, ArrayToList(chunk))
    }

    return ArrayToList(pieces), nil
}

func SublistImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    return
}

func TakeImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    n, err := Eval(Car(args), env)
    if err != nil {
        return
    }
    if !NumberP(n) {
        err = errors.New("take requires a number as it's first argument.")
    }
    size := int(NumericValue(n))

    l, err := Eval(Cadr(args), env)
    if err != nil {
        return
    }
    if !ListP(l) {
        err = errors.New("take requires a list as it's second argument.")
    }

    var items []*Data = make([]*Data, 0, Length(args))
    for i, cell := 0, l; i < size && NotNilP(cell); i, cell = i+1, Cdr(cell) {
        items = append(items, Car(cell))
    }
    result = ArrayToList(items)
    return
}

func DropImpl(args *Data, env *SymbolTableFrame) (result *Data, err error) {
    n, err := Eval(Car(args), env)
    if err != nil {
        return
    }
    if !NumberP(n) {
        err = errors.New("drop requires a number as it's first argument.")
    }
    size := int(NumericValue(n))

    l, err := Eval(Cadr(args), env)
    if err != nil {
        return
    }
    if !ListP(l) {
        err = errors.New("drop requires a list as it's second argument.")
    }

    var cell *Data
    var i int
    for i, cell = 0, l; i < size && NotNilP(cell); i, cell = i+1, Cdr(cell) {
    }
    result = cell
    return
}