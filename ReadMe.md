## Tello

###### Being an adventure into interfacing with an DJI Tello drone and learning some Golang at the same time

#### Overview

This is intended to be rudimentary flight-control software for a DJI Tello drone.  The goal for
v1 is to

* have streaming video to the console

* be able to construct and run rudimentary flight plans of the format

      # pseudocode
      t=0 TAKE OFF
      t=5 translate 10m forward
      t=10 rotate 180 degrees
      t=15 translate 10m forward
      t=20 rotate 180 degrees
      t=25 LAND
    
    and have these be carried out by the drone.
    
The control software should also validate that flight plans are sane, and that while
horizontal displacement is ok at termination, vertical displacement isn't.

###### ToDo

1. integrate https://github.com/gizak/termui         