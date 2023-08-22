package main

import (
  "fmt"
  "bufio"
  "log"
  "os"
  "strings"
  "net"
)

func main() {
  scanner := bufio.NewScanner(os.Stdin)
  fmt.Print("domain,hasMX,hasSPF,sprRecord,hasDMARC,dmarcRecord\n")

  for scanner.Scan() {
    checkDomain(scanner.Text())
  }

  if err := scanner.Err(); err != nil {
    log.Fatal("Error: COuld not read from input: %v\n", err)
  }
}

func checkDomain(domain string) {
  var hasMX, hasSPF, hasDMARC bool
  var spfRecord, dmarcRecord string

  mxRecords, err := net.LookupMX(domain)

  if err != nil {
    log.Fatal("Error: Could not lookup MX records for domain %s: %v\n", domain, err)
  }

  if len(mxRecords) > 0 {
    hasMX = true
  }

  textRecors, err := net.LookupTXT(domain)

  if err != nil {
    log.Fatal("Error: Could not lookup TXT records for domain %s: %v\n", domain, err)
  }

  for _, record := range textRecors {
    if strings.HasPrefix(record, "v=spf1") {
      hasSPF = true
      spfRecord = record
      break
    }
  }

  dmarcRecords, err := net.LookupTXT("_dmarc." + domain)

  if err != nil {
    log.Fatal("Error: Could not lookup dmarc records for domain %s: %v\n", domain, err)
  }

  for _, record := range dmarcRecords {
    if strings.HasPrefix(record, "v=DMARC1") {
      hasDMARC = true
      dmarcRecord = record
      break
    }
  }

  fmt.Printf("%s,%t,%t,%s,%t,%s\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
