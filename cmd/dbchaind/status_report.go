package main

import (
    "fmt"
    "time"
    "crypto/tls"
    "net"
    "strings"
)

func statusReport() {
    for true {
        sendReport()
        time.Sleep(30 * time.Minute)
    }
}
     
func sendReport() {
    conf := &tls.Config{
         //InsecureSkipVerify: true,
    }

    conn, err := tls.Dial("tcp", "www.dbchain.cloud:443", conf)
    if err != nil {
        return
    }
    defer conn.Close()

    var ips = getIp()
    var msg = fmt.Sprintf("hello /dbchain_status_report/%s/", ips)
    _, err = conn.Write([]byte(msg))
    if err != nil {
        return
    }

    buf := make([]byte, 100)
    _, err = conn.Read(buf)
    if err != nil {
        return
    }
}


func getIp() string {
    ifaces, err:= net.Interfaces()
    if err != nil {
        return "failed to get intefaces"
    }

    var ips []string
    for _, i := range ifaces {
        addrs, err := i.Addrs()
        if err != nil {
            return "failed to get addresses"
        }
        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }
            var strIp = ip.String()
            if (len(strIp) > 6) && strIp[:4] != "127." {
                ips = append(ips, ip.String())
            }
        }
    }
    return strings.Join(ips, "-")
}
