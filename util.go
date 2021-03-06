// Copyright 2014 SteelSeries ApS.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package implements a basic LISP interpretor for embedding in a go program for scripting.
// This file implements some utilities.

package golisp

func ArrayToList(sexprs []*Data) *Data {
    head := EmptyCons()
    lastCell := head
    for _, element := range sexprs {
        newCell := Cons(element, nil)
        lastCell.Cdr = newCell
        lastCell = newCell
    }
    return head.Cdr
}

func ArrayToListWithTail(sexprs []*Data, tail *Data) *Data {
    head := EmptyCons()
    lastCell := head
    for _, element := range sexprs {
        newCell := Cons(element, nil)
        lastCell.Cdr = newCell
        lastCell = newCell
    }
    lastCell.Cdr = tail
    return head.Cdr
}

func ToArray(list *Data) []*Data {
    result := make([]*Data, 0)
    for c := list; NotNilP(c); c = Cdr(c) {
        result = append(result, Car(c))
    }
    return result
}
