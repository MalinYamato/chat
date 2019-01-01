//
// Copyright 2018 Malin Yamato --  All rights reserved.
// https://github.com/MalinYamato
//
// MIT License
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Rakuen. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.



package main

import "log"

//      expected return value
//
//     [publisherA] --> [targetA],[targetB],[targetC], n
//     [targetA]  --> [target = PublisherA],[target], n
//     [targetB]  --> [target = publiherA],[target],[target], n
//     [targetC]  --> [target], n
//

func testA() {
	_publishers = make(PublishersTargets)
	_publishers["X"] = Targets{"B": true, "X": true}

	//_publishers["B"] = Targets{"A": true, "F": true, "X": true}
	//_publishers["C"] = Targets{"D": true, "X": true, "F": true}
	allTargs, _ := _publishers.collectAllTargets("X")
	for k, _ := range allTargs {
		log.Printf("Target %s  \n", k)

	}
	log.Println("TEST A #############")
	for at := range _publishers["X"] {
		pp, col := _publishers.Status("X", at)
		log.Println("Col ", col.Status)
		for t, _ := range pp {
			for tt, _ := range pp[t] {
				log.Printf(">> %s %s\n", t, tt)
			}
		}
	}

	log.Println("TEST B ##############")
	allTargs, _ = _publishers.collectAllTargets("X")
	for at := range allTargs {
		pp, col := _publishers.Status("X", at)
		log.Println("Col ", col.Status)
		for t, _ := range pp {
			for tt, _ := range pp[t] {
				log.Printf(">> %s %s\n", t, tt)
			}
		}

	}

}
