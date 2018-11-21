package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jchorl/watchdog"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"

	"github.com/jchorl/dnsserv/common"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the DNS server",
	Long:  `Starts the DNS server. The server listens for DNS requests (on port 53 by default) and also starts an HTTPS server (on port 643 by default) to get updates.`,
	Run:   server,
}

var (
	dnsPort   int
	httpsPort int
)

func init() {
	serverCmd.Flags().StringVar(&caPath, "ca-path", "./certs/root.pem", "Path to root ca cert")
	serverCmd.Flags().StringVar(&certPath, "cert-path", "./certs/cert.pem", "Path to cert for mTLS")
	serverCmd.Flags().StringVar(&keyPath, "key-path", "./certs/cert.key", "Path to key for mTLS")
	serverCmd.Flags().IntVar(&dnsPort, "dns-port", 53, "Port to serve the dns server on")
	serverCmd.Flags().IntVar(&httpsPort, "https-port", 443, "Port to serve the https server on")
	rootCmd.AddCommand(clientCmd)
}

var domainsToAddresses = map[string]string{}
var mapLock = sync.RWMutex{}

type dnsHandler struct{}

func (*dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name

		var address string
		var ok bool
		mapLock.RLock()
		address, ok = domainsToAddresses[domain]
		mapLock.RUnlock()

		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
		}
	default:
		log.Printf("Unknown question type: %d\n", r.Question[0].Qtype)
	}
	w.WriteMsg(&msg)
}

type updateHandler struct{}

func (*updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ipAddrStr := r.RemoteAddr
	if host, _, err := net.SplitHostPort(ipAddrStr); err == nil {
		ipAddrStr = host
	}

	// validate that the ip can be parsed correctly
	ipAddr := net.ParseIP(ipAddrStr)
	if ipAddr == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to get client ip address. Found remote addr: %s", ipAddr)
		log.Printf("Unable to get client ip address. Found remote addr: %s", ipAddr)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var updateRequest common.UpdateRequest
	err := decoder.Decode(&updateRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to unmarshal body: %s", err)
		log.Printf("Unable to unmarshal body: %s", err)
		return
	}

	// need to append a period to the domain for the lookup
	domainEntry := updateRequest.Domain + "."
	mapLock.Lock()
	domainsToAddresses[domainEntry] = ipAddr.String()
	mapLock.Unlock()

	w.WriteHeader(http.StatusOK)
}

func server(cmd *cobra.Command, args []string) {
	wg := sync.WaitGroup{}

	dnsServer := &dns.Server{Addr: ":" + strconv.Itoa(dnsPort), Net: "udp"}
	dnsServer.Handler = &dnsHandler{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Printf("Serving DNS server on %d\n", dnsPort)
			if err := dnsServer.ListenAndServe(); err != nil {
				log.Printf("DNS server failed %s\n", err.Error())
			}
		}
	}()

	tlsConfig := common.LoadTLSConfigOrPanic(caPath, certPath, keyPath)
	updateServer := http.Server{
		Addr:      ":" + strconv.Itoa(httpsPort),
		Handler:   &updateHandler{},
		TLSConfig: tlsConfig,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Printf("Serving HTTPS server on %d\n", httpsPort)
			if err := updateServer.ListenAndServeTLS(certPath, keyPath); err != nil {
				log.Printf("HTTP server failed %s\n", err.Error())
			}
		}
	}()

	wdClient := watchdog.Client{"https://watchdog.joshchorlton.com"}
	ticker := time.NewTicker(time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				wdClient.Ping("dnsserv", watchdog.Watch_DAILY)
			}
		}
	}()

	wg.Wait()
}
