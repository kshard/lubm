//
// Copyright (C) 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/lubm
//

package adapter

import (
	"github.com/kshard/sigma/vm"
	"github.com/kshard/spock"
	"github.com/kshard/spock/store/ephemeral"
	"github.com/kshard/xsd"
)

type subQ struct {
	addr   []vm.Addr
	store  *ephemeral.Store
	stream spock.Stream
}

func NewStream(store *ephemeral.Store) func(addr []vm.Addr) vm.Stream {
	return func(addr []vm.Addr) vm.Stream {
		return &subQ{
			addr:  addr,
			store: store,
		}
	}
}

func (seq *subQ) Init(heap *vm.Heap) error {
	var (
		s *spock.Predicate[xsd.AnyURI]
		p *spock.Predicate[xsd.AnyURI]
		o *spock.Predicate[xsd.Value]
	)

	if !seq.addr[0].IsWritable() {
		v := heap.Get(seq.addr[0])
		s = &spock.Predicate[xsd.AnyURI]{
			Clause: spock.EQ,
			Value:  v.(xsd.AnyURI),
		}
	}

	if !seq.addr[1].IsWritable() {
		v := heap.Get(seq.addr[1])
		p = &spock.Predicate[xsd.AnyURI]{
			Clause: spock.EQ,
			Value:  v.(xsd.AnyURI),
		}
	}

	if !seq.addr[2].IsWritable() {
		v := heap.Get(seq.addr[2])
		o = &spock.Predicate[xsd.Value]{
			Clause: spock.EQ,
			Value:  v,
		}
	}

	q := spock.Query(s, p, o)
	stream, err := ephemeral.Match(seq.store, q)
	if err != nil {
		return err
	}
	seq.stream = stream
	return seq.Read(heap)
}

func (seq *subQ) Read(heap *vm.Heap) error {
	if !seq.stream.Next() {
		return vm.EndOfStream
	}

	spock := seq.stream.Head()
	if seq.addr[0].IsWritable() {
		heap.Put(seq.addr[0], spock.S)
	}
	if seq.addr[1].IsWritable() {
		heap.Put(seq.addr[1], spock.P)
	}
	if seq.addr[2].IsWritable() {
		heap.Put(seq.addr[2], spock.O)
	}

	return nil
}
