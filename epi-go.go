// epi-go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func sum(array []int) int { // used to add up the initial people sick during epidemic seeding
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}

// Population is stored in a list.  The list size variable is "size". You can
// adjust this and all the other parameters are calculated automatically
// Variables:
//	size: the size of the population can be anynumber 1000 works good
//	seed: the initial number of sick people in the population to get the epidemic started
//	days: the number of days to run the simulation. Should start at least at 50
//	daystorecovery: average days for a normal person to recover and become immune
//	PeopleContact: typical number of people a sick person will come in contact with per day
// 	targetperson: a person choosen at random in the population and if not sick, we try to get them sick
//	ProbOfInfection: each non-sick person who encounters a sick person has a probability of getting sick
//	RO: the reproduction rate. How many people a sick person will get sick
//	History: An array where we store the day results in the epidemic lifecycle
// 	sick: A list that is variable size long. Stores a person's status 0: Not sick, 1: sick, 2: immune
//	dayssick: A list where the days each person is sick is stored. Same index as list: sick

func main() {
	var size, seed, days, daystorecovery, PeopleContact int
	var targetperson int
	var ProbOfInfection, R0 float32
	var history [150][3]int // the size of this array needs to at least as

	size = 1000
	seed = 2
	days = 75
	daystorecovery = 10
	R0 = 1.3
	PeopleContact = 10
	ProbOfInfection = 100 * R0 / (float32(PeopleContact) * float32(daystorecovery))

	// big as variable: days
	// [day][# not sick][# sick][# immune]
	sick := make([]int, size)
	dayssick := make([]int, size)

	fmt.Println("Prob:", ProbOfInfection)

	rand.Seed(time.Now().UnixNano()) // randomly get "seed" number of people sick
	for i := 0; i < seed; i++ {
		sick[rand.Intn(size)] = 1
	}

	history[0][0] = size - sum(sick) // set up history for day=0
	history[0][1] = sum(sick)
	// print out day=0 score
	// this is printed as a string suitable for importing to excel for graphing as a csv file
	fmt.Println("0,", history[0][0], ",", history[0][1], ",", history[0][2])

	for day := 1; day < days; day++ { // loop for the simulation days
		history[day][0] = history[day-1][0] // carry over yesterdays totals
		history[day][1] = history[day-1][1]
		history[day][2] = history[day-1][2]

		for i := 0; i < size; i++ { // loop through the population
			if sick[i] == 1 { // if sick
				dayssick[i]++                     // increment days sick
				if dayssick[i] > daystorecovery { // if reached recovery days
					sick[i] = 2       // then change status to recovered (2)
					history[day][2]++ // increment # of people recovered
					history[day][1]-- // decrement # of people sick
				}
				for x := 0; x < PeopleContact; x++ { // while with a sick person
					targetperson = rand.Intn(size) // find people whom they might infect
					// pick people in the sick list at random
					if i != targetperson && sick[targetperson] == 0 { // available to be sick, roll the dice
						tar := rand.Intn(100)           // get a number from 0 to 100
						if tar < int(ProbOfInfection) { // if < than ProbOfInfection calculated above
							sick[targetperson] = 1 // lucky them they are now sick
							history[day][0]--      // decrement # of people not sick
							history[day][1]++      // increment # of people sick
						}
					}
				}
			}
		}
		// print out the day's score
		fmt.Println(day, ",", history[day][0], ",", history[day][1], ",", history[day][2])
	}
}
