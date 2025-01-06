package proxy

import (
	"io"
	"log"
	"net"
	"net/http"
	"runtime"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func Run() {
	// 启动HTTP代理服务器
	http.HandleFunc("/", handleRequest)
	log.Println("启动HTTP代理服务器，监听端口: 8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("启动HTTP代理服务器失败: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		handleHTTPS(w, r)
	} else {
		handleHTTP(w, r)
	}
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleHTTPS(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	go copyData(destConn, clientConn)
	copyData(clientConn, destConn)
}

func copyData(dst io.Writer, src io.Reader) {
	io.Copy(dst, src)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func modifyPacket(packet gopacket.Packet) {
	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		ip.SrcIP = net.IPv4(1, 2, 3, 4) // 修改源IP地址
	}
}

func capturePackets() {
	if runtime.GOOS == "windows" {
		log.Println("Windows系统不支持pcap捕获，跳过此功能")
		return
	}

	handle, err := pcap.OpenLive("eth0", 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		modifyPacket(packet)
	}
}

func init() {
	go capturePackets()
}
