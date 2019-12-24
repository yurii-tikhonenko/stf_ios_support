package main

import (
  "fmt"
  "net"
  "os"
  "strconv"
  "strings"
  log "github.com/sirupsen/logrus"
)

func construct_ports( config *Config, spec string ) ( map [int] *PortItem ) {
    ports := make( map [int] *PortItem )
    if strings.Contains( spec, "-" ) {
        parts := strings.Split( spec, "-" )
        from, _ := strconv.Atoi( parts[0] )
        to, _ := strconv.Atoi( parts[1] )
        for i := from; i <= to; i++ {
            portItem := PortItem{
                available: true,
            }
            ports[ i ] = &portItem
        }
    }
    return ports
}

func ifAddr( ifName string ) ( addrOut string ) {
    ifaces, err := net.Interfaces()
    if err != nil {
        fmt.Printf( err.Error() )
        os.Exit( 1 )
    }

    addrOut = ""
    for _, iface := range ifaces {
        addrs, err := iface.Addrs()
        if err != nil {
            fmt.Printf( err.Error() )
            os.Exit( 1 )
        }
        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
                case *net.IPNet:
                    ip = v.IP
                case *net.IPAddr:
                    ip = v.IP
                default:
                    fmt.Printf("Unknown type\n")
            }
            if iface.Name == ifName {
                addrOut = ip.String()
            }
        }
    }
    return addrOut
}

func get_net_info( config *Config ) ( string, string, bool ) {
    var vpnMissing bool = false

    // This information comes from Tunnelblick log
    // It may no longer be active
    tunName, curIP, err := vpn_info( config )

    if err != "" {
        log.WithFields( log.Fields{
            "type": "vpn_err",
            "err": err,
        } ).Info( err )
        return "", "", true
    }

    log.WithFields( log.Fields{
        "type": "info_vpn",
        "tunnel_name": tunName,
    } ).Info("Tunnel name")

    ipConfirm := getTunIP( tunName )
    if ipConfirm != curIP {
        // The tunnel is no longer active
        vpnMissing = true
    } else {
        log.WithFields( log.Fields{
            "type": "info_vpn",
            "tunnel_name": tunName,
            "ip": curIP,
        } ).Info("IP on VPN")
    }

    return tunName, curIP, vpnMissing
}

func ifaceCurIP( tunName string ) string {
    ipStr := ifAddr( tunName )
    if ipStr != "" {
        log.WithFields( log.Fields{
            "type": "net_interface_info",
            "interface_name": tunName,
            "ip": ipStr,
        } ).Debug("Interface Details")
    } else {
        log.WithFields( log.Fields{
            "type": "err_net_interface",
            "interface_name": tunName,
        } ).Fatal("Could not find interface")
    }

    return ipStr
}