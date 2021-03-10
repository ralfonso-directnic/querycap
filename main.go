package main

import (
        "fmt"
        "github.com/google/gopacket"
        "github.com/google/gopacket/layers"
        "github.com/google/gopacket/pcap"
        "flag"
        "github.com/ralfonso-directnic/sqlparser"
        "github.com/ralfonso-directnic/sqlparser/query"
        "strings"
        "log"
        "os"
        "io"
        "regexp"
)

var device string
var port int
var promport string
var logfile string


func main() {

        flag.StringVar(&device,"device","any","Device to listen to (any is default)")
        flag.IntVar(&port,"port",3306,"Mysql Port")
        flag.StringVar(&promport,"promport","9224","Prometheus Port")
        flag.StringVar(&logfile,"logfile","/var/log/querycap.log","Log file location")
        flag.Parse()

        setupLog()

        startProm()

        if handle, err := pcap.OpenLive("any", 1600, true, pcap.BlockForever); err != nil {
                panic(err)
        } else if err := handle.SetBPFFilter("tcp"); err != nil { // optional
                panic(err)
        } else {
                packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
                for packet := range packetSource.Packets() {
                        handlePacket(packet) // Do something with a packet here.
                }
        }
}

func setupLog(){

        f, err := os.OpenFile(logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
        if err != nil {
                log.Fatalf("error opening file: %v", err)
        }

        mw := io.MultiWriter(os.Stdout, f)
        log.SetOutput(mw)



}

func handlePacket(packet gopacket.Packet) {
        if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
              //  ip, _ := ipLayer.(*layers.IPv4)
                pport := uint16(port)


                lp := layers.TCPPort(pport)

                if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
                        // Get actual TCP data from this layer
                        tcp, _ := tcpLayer.(*layers.TCP)
			
			            if(tcp.DstPort==lp) {

                        handlePayload(string(tcp.LayerPayload()))

                        }
                    }
        }
}

func handlePayload(pl string){

        qry_valid:=[]string{"SELECT","UPDATE","INSERT","DELETE","select","update","insert","delete"}

        qry := stripCtlAndExtFromBytes(pl)


        for _,q:= range qry_valid {

                if strings.Contains(qry,q) {

                        parts := strings.SplitN(qry,q,2)


                        if(len(parts)>1){

                                finalqry := fmt.Sprintf("%s %s",q,parts[1])

                                handleQuery(finalqry)


                        }

                }


        }




}

func handleQuery(qry string){

        //parser chokes on order by

        qry_orig := qry

        parts:=strings.SplitN(qry,"ORDER",2)

        qry = parts[0]

        parts=strings.SplitN(qry,"LIMIT",2)

        qry = parts[0]

        parts=strings.SplitN(qry,"HAVING",2)

        qry = parts[0]

        qry=strings.Replace(qry,"(","",-1)
        qry=strings.Replace(qry,")","",-1)

        //fix for non quoted value
        qry = strings.Replace(qry,"= ?",`= ''`,-1)
        qry = strings.Replace(qry,"=?",`= ''`,-1)

        query, err := sqlparser.Parse(qry)
        if err != nil {
                //fallback
               query = fallbackQueryHandle(qry_orig)


        }

        log.Printf("%s %s\n",query.Database,query.TableName)
        go qryTotal.Inc()
        go qryDbTable.WithLabelValues(query.Database,query.TableName).Inc()

}

func fallbackQueryHandle(qry string) (query.Query){

        var q query.Query

        check_regex:=[]string{`FROM (.[^\s]*)`,`UPDATE (.[^\s]*)`,`INSERT INTO (.[^\s]*)`}

        for _,reraw := range check_regex {

                re := regexp.MustCompile(reraw)
                parts := re.FindStringSubmatch(qry)

                if (len(parts) > 1) {

                        item := parts[1]

                        item = strings.ReplaceAll(item, "`", "")
                        item = strings.ReplaceAll(item, "'", "")
                        item = strings.ReplaceAll(item, `"`, "")

                        tbl_parts := strings.Split(item,".")

                        if(len(tbl_parts)==2){



                                q.Database = tbl_parts[0]
                                q.TableName = tbl_parts[1]


                        }else{


                               q.TableName=item
                        }

                        //split up!
                        break

                }

        }

        return q




}

func stripCtlFromBytes(str string) string {
        b := make([]byte, len(str))
        var bl int
        for i := 0; i < len(str); i++ {
                c := str[i]
                if c >= 32 && c != 127 {
                        b[bl] = c
                        bl++
                }
        }
        return string(b[:bl])
}

func stripCtlAndExtFromBytes(str string) string {
        b := make([]byte, len(str))
        var bl int
        for i := 0; i < len(str); i++ {
                c := str[i]
                if c >= 32 && c < 127 {
                        b[bl] = c
                        bl++
                }
        }
        return string(b[:bl])
}
