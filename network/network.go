package network

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	DEFAULT_PORT = 8080
	MAX_PORT     = 65535
)

// Network interface to deal with http server and IP
type Network interface {
	StartHttpServer(port uint16, useFile bool) error
	AllUnavailablePorts() PortList
	AllUnavailablePortsFromList(pl *PortList) PortList
	PortIsAvailable(port uint16) (status bool, err error)
	InternalIP() (net.IP, error)
	ExternalIP() (net.IP, error)
	Forwarding(target string, port uint16) error
}

// NetworkHelper implements Network interface
type networkHandler struct {
	httpClient HttpServer
	netClient  LocalNet
}

/*
 * Return a new instance of Network
 */
func NewNetwork(httpClient HttpServer, netClient LocalNet) Network {
	if httpClient == nil {
		httpClient = &httpServer{}
	}
	if netClient == nil {
		netClient = &localNet{}
	}

	return &networkHandler{
		httpClient: httpClient,
		netClient:  netClient,
	}
}

/*
 * Start a local http server with the port number provided
 */
func (n *networkHandler) StartHttpServer(port uint16, useFile bool) error {
	if useFile {
		n.httpClient.Handle("/", http.FileServer(http.Dir("./")))
	} else {
		n.httpClient.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "You're now on port %v [%s]", port, r.URL.Path[0:])
		})
	}

	//go func() {
	if err := n.httpClient.ListenAndServe(portStringify(port), nil); err != nil {
		return err
	}
	//}()
	return nil
}

/*
 * Return all the unavailable port number from the machine
 */
func (n *networkHandler) AllUnavailablePorts() PortList {
	var unavailablePorts PortList
	for i := 0; i <= MAX_PORT; i++ {
		if status, _ := n.PortIsAvailable(uint16(i)); !status {
			unavailablePorts = append(unavailablePorts, uint16(i))
		}
	}
	return unavailablePorts
}

/*
 * Return all the unavailable port number from the given list of port numbers
 */
func (n *networkHandler) AllUnavailablePortsFromList(pl *PortList) PortList {
	var unavailablePorts PortList
	for _, port := range *pl {
		if status, _ := n.PortIsAvailable(port); !status {
			unavailablePorts = append(unavailablePorts, port)
		}
	}
	return unavailablePorts
}

/*
 * Check the given port number is availble to be used for the machine
 */
func (n *networkHandler) PortIsAvailable(port uint16) (status bool, err error) {
	host := ":" + strconv.Itoa(int(port))
	server, err := n.netClient.Listen("tcp", host)
	if err != nil {
		return false, err
	}
	server.Close()
	return true, nil
}

/*
 * Retrieve machine local IP address
 */
func (n *networkHandler) InternalIP() (net.IP, error) {
	// Dial to connect to local server
	conn, err := n.netClient.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return nil, err
	}

	// Close the connection when the response is read
	defer conn.Close()

	// Get the local IP address
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

/*
 * Retrieve machine external IP address
 */
func (n *networkHandler) ExternalIP() (net.IP, error) {
	// Get request on myexternalip.com to retrieve external IP address for the machine
	resp, err := n.httpClient.Get("http://myexternalip.com/raw")

	if err != nil {
		return nil, err
	}

	// Close the connection when response body is read
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	ipStr := string(body)

	// Convert body byte array to string and return it
	return net.ParseIP(strings.TrimSuffix(ipStr, "\n")), nil
}

/*
 * Port forwarding
 */
func (n *networkHandler) Forwarding(target string, port uint16) error {
	// Declare listener to the origin port
	listener, err := n.netClient.Listen("tcp", portStringify(port))
	if err != nil {
		return err
	}

	errChan := make(chan bool)
	for {
		select {
		case <-errChan:
			return fmt.Errorf("TCP connection failed")
		default:
			// Start connection to listener
			conn, err := listener.Accept()
			if err != nil {
				return err
			}
			go n.forward(conn, target, errChan)
		}
	}
}

/*
 * Forward connection to then given port number
 */
func (n *networkHandler) forward(conn LocalNetConn, target string, c chan bool) {
	// Declare client to the forwarding port
	client, err := n.netClient.Dial("tcp", target)
	if err != nil {
		c <- true
		return
	}

	// Copy IO
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

/*
 * Make port number to string format ":xxx"
 */
func portStringify(port uint16) string {
	return fmt.Sprintf(":%d", port)
}
