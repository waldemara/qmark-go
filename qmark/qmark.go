/* Simple CPU benchmark test

Go servers read messages from their queues, update trace information, and send
them back to the originators.

Go clients create a single message and pass it through all servers
sequentially.  Clients exit when the message passes through all
servers.  The test ends when all clients complete.

 Message format

	CMD:TRACE

	CMD		- command
				exit	- exit
				queue	- update trace and continue test

	TRACE	- list of visited clients and servers, the last item is the
			  originator of the message

 For example:

	server 1:	queue:client(1)
	client 1:	queue:client(1)-server(1)
	server 2:	queue:client(1)-server(1)-client(1)
	client 1:	queue:client(1)-server(1)-client(1)-server(2)
	...
*/

package qmark

import "fmt"
import s "strings"
import "strconv"
import "time"
import "runtime"

var Qmark = qmark
var RunQmark = run_qmark

const debug = false
const SERVERS = 151
const CLIENTS = 1109
const RUNS = 7

type Data struct {
	num_clients int
	num_servers int
	num_runs    int
	clientqs    []chan string
	serverqs    []chan string
	client_exit chan int
	server_exit chan int
}

var qm = &Data{}

func client(cid int) {

	count := qm.num_servers
	dstid := cid % qm.num_servers
	msgout := fmt.Sprintf("queue:client(%d)", cid)
	qm.serverqs[dstid] <- msgout
	for msg := range qm.clientqs[cid] {
		if debug {
			fmt.Printf("client(%d):  %s\n", cid, msg)
		}
		if count--; count < 1 {
			break
		}
		dstid = (dstid + 1) % qm.num_servers
		msgout = fmt.Sprintf("%s-client(%d)", msg, cid)
		qm.serverqs[dstid] <- msgout
		runtime.Gosched()
	}
	if debug {
		fmt.Printf("client(%d):  exit\n", cid)
	}
	qm.client_exit <- cid
}

func extract_srcid(trace string) int {

	src := trace[s.LastIndex(trace, "(")+1 : len(trace)-1]
	srcid, _ := strconv.Atoi(src)
	return srcid
}

func server(sid int) {

	for msg := range qm.serverqs[sid] {
		if debug {
			fmt.Printf("server(%d):  %s\n", sid, msg)
		}
		toks := s.Split(msg, ":")
		if toks[0] == "exit" {
			break
		}
		dstid := extract_srcid(toks[1])
		msgout := fmt.Sprintf("%s:%s-server(%d)", toks[0], toks[1], sid)
		qm.clientqs[dstid] <- msgout
		runtime.Gosched()
	}
	if debug {
		fmt.Printf("server(%d):  exit\n", sid)
	}
	qm.server_exit <- sid
}

func run(num_clients, num_servers int) time.Duration {

	start_time := time.Now()
	for ii := 0; ii < num_clients; ii++ {
		qm.clientqs[ii] = make(chan string)
	}
	for ii := 0; ii < num_servers; ii++ {
		qm.serverqs[ii] = make(chan string, num_clients)
	}
	// start the test
	for ii := 0; ii < num_servers; ii++ {
		go server(ii)
	}
	for ii := 0; ii < num_clients; ii++ {
		go client(ii)
	}
	// wait for clients to complete
	for ii := 0; ii < num_clients; ii++ {
		cid := <-qm.client_exit
		if debug {
			fmt.Printf("exit client(%d)\n", cid)
		}
	}
	for ii := 0; ii < num_servers; ii++ {
		qm.serverqs[ii] <- "exit:adm"
	}
	// wait for servers to complete
	for ii := 0; ii < num_servers; ii++ {
		<-qm.server_exit
	}
	result := time.Since(start_time)
	return result
}

func run_qmark(num_clients, num_servers, num_runs int) []time.Duration {

	if num_runs < 1 {
		num_runs = 1
	}
	runs := make([]time.Duration, num_runs) // List of bench mark results
	for ix := range runs {
		qm = &Data{
			num_clients: num_clients,
			num_servers: num_servers,
			num_runs:    num_runs,
			clientqs:    make([]chan string, num_clients),
			serverqs:    make([]chan string, num_servers),
			client_exit: make(chan int),
			server_exit: make(chan int),
		}
		result := run(num_clients, num_servers)
		runs[ix] = result
	}
	return runs
}

func qmark() int {

	var sum time.Duration
	results := run_qmark(CLIENTS, SERVERS, RUNS)
	for _, res := range results {
		sum += res
	}
	avg := (int(sum) / len(results)) / 1000 // [µs]
	return int(1000.0 / (float64(avg) / 1000000.0))
}
