/*
 * Copyright (C) 2025 Ian M. Fink.
 * All rights reserved.
 *
 * File:	main.go
 *
 * Purpose:	Experimental Golang ticker using a go function
 *
 * This program is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option)
 * any later version.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
 * or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for
 * more details.
 * 
 * You should have received a copy of the GNU General Public License along
 * with this program. If not, please see: https://www.gnu.org/licenses
 *
 */

package main

/*
 * Imports
 */

import (
	"fmt"
	"time"
	"os"
)

/*
 * Global Consts
 */

const secondsToRun = 6

/******************************************************************************/

/**
 * Name:	myGoFunc
 *
 * @brief	a go function that receives a golang ticker channel
 *
 * @param	myTicker		a pointer to a time.Ticker instance
 * @param	stopItNow		a channel that communicates boolean values
 *
 */


func tickMeOff(theTicksPtr *int) {
	if *theTicksPtr == 1 {
		fmt.Println("tickMeOff(): Been ticked off 1 time!")
	} else {
		fmt.Printf("tickMeOff(): Been ticked off %d times!\n", *theTicksPtr)
	}

	*theTicksPtr++
} /* tickMeOff */

/******************************************************************************/

func myGoFunc(myTicker *time.Ticker, stopItNow chan bool) {
	var	(
		err			error
		theTicks	int
	)

	theTicks = 1
	fmt.Println("myGoFunc():  Ready to be ticked off!!!")

	for {
		select {
			case <-myTicker.C:
				// call the func to be executed/run periodically
				tickMeOff(&theTicks)

			case <-stopItNow:
				fmt.Println("myGoFunc(): Exiting.")

				// there may be a race condition between main exiting
				// and clearing stdout, so flush stdout
				// *** while this should NOT produce an error, it does
				// sometimes, at least on MacOS v12.7.4. sometimes it
				// does not work at all.
				err = os.Stdout.Sync()
				if err != nil {
					fmt.Println("Error flushing stdout:", err)
				}
				return
		}
	}
} /* myGoFunc */

/******************************************************************************/

func main() {
	var (
		theTicker	*time.Ticker
		stopIt		chan bool
	)

	fmt.Println("main():  Get ready to be ticked off!!!")

	// allocate the channel
	stopIt = make(chan bool)

	// set a "tick" of every 1 second
	theTicker = time.NewTicker(1 * time.Second)
	defer theTicker.Stop()

	go myGoFunc(theTicker, stopIt)

	// Let things run for "secondsToRun" seconds,
	// otherwise myGoFunc() will no longer run when
	// main() exits.  Sleep is used to simulate
	// doing other processings/running other functions.
	time.Sleep(secondsToRun * time.Second)

	fmt.Println("main():  Woken up from sleep.")

	// send a value to the channel to cease operations of the go function
	stopIt <- false

	fmt.Println("main():  Exiting.")
} /* main */

/******************************************************************************/

/*
 * End of file:	main.go
 */

