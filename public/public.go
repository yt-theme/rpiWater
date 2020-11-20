package public

var Chan_start     = make(chan int) // start
var Chan_stop      = make(chan int) // stop
var Chan_R1_isActi     = make(chan int) // relay1 is active
var Chan_R1_isInActi   = make(chan int) // relay1 is inactive
var Chan_R2_isActi     = make(chan int) // relay2 is active
var Chan_R2_isInActi   = make(chan int) // relay2 is inactive
// var Chan_S1_isActi = make(chan int) // sensor1 is active

var Chan_S1_stopCheck = make(chan int) // stop sensor1 check
var S1_isChecking = false

var Chan_sendMsg = make(chan string, 20) // send msg gpioControl to socket