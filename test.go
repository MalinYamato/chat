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
	_publishers["X"] = Targets{"B": true,"C":true}
	//_publishers["B"] = Targets{"A": true, "F": true, "X": true}
	//_publishers["C"] = Targets{"D": true, "X": true, "F": true}
	allTargs, _ :=  _publishers.collectAllTargets("X")
	for k, v := range allTargs {
		log.Printf("Target %s %s \n", k, v)

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
	allTargs, _ =  _publishers.collectAllTargets("X")
	for at:= range allTargs{
		pp, col := _publishers.Status("X", at)
		log.Println("Col ", col.Status)
		for t, _ := range pp {
			for tt, _ := range pp[t] {
				log.Printf(">> %s %s\n", t, tt)
			}
		}

	}

}
