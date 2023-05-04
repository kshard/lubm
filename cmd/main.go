//
// Copyright (C) 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/lubm
//

package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/kshard/lubm"
	"github.com/kshard/lubm/internal/adapter"
	"github.com/kshard/sigma"
	"github.com/kshard/sigma/asm"
	"github.com/kshard/sigma/lang"
	"github.com/kshard/spock"
	"github.com/kshard/spock/store/ephemeral"
)

// TODO: use cobra

func main() {
	n := 1
	store := ephemeral.New()

	//
	// Intake
	//
	size := 0
	ch := make(chan spock.Bag, 0)
	go func() {
		for x := range ch {
			ephemeral.Add(store, x)
			size = size + len(x)
		}
	}()

	t := time.Now()
	ds := lubm.NewDataSet(1683234740, n, ch)
	for i := 0; i < n; i++ {
		if err := ds.Generate(i); err != nil {
			panic(err)
		}
		fmt.Printf("==> university %d: %d in %v\n", i, size, time.Since(t))
	}
	fmt.Printf("==> loaded %d in %v\n", size, time.Since(t))

	close(ch)

	//
	// Benchmark
	//
	qs := []string{
		lubm.Query1(),
		lubm.Query2(),
		lubm.Query3(),
		lubm.Query4(),
		lubm.Query5(),
		lubm.Query6(),
		lubm.Query7(),
		lubm.Query8(),
		lubm.Query9(),
	}

	for i, qx := range qs {
		t = time.Now()
		if q, err := query(store, qx); err == nil {
			fmt.Printf("==> query #%d %8.d in %v\n", i+1, q, time.Since(t))
		} else {
			fmt.Printf("==> query #%d failed %s", i+1, err)
		}
	}
}

func query(store *ephemeral.Store, q string) (int, error) {
	buf := bytes.NewBuffer([]byte(q))

	parser := lang.NewParser(buf)
	rules, err := parser.Parse()
	if err != nil {
		return 0, err
	}

	machine, err := sigma.New("q", rules)
	if err != nil {
		return 0, err
	}

	ctx := asm.NewContext().Add("f", adapter.NewStream(store))
	reader := sigma.Stream(ctx, machine)

	seq := reader.ToSeq()
	return len(seq), nil
}
