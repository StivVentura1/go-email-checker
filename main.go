package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func checkDomain(domain string) (bool, bool, string, bool, string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	// Controllo dei record MX
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Errore nel cercare record MX per %s: %v\n", domain, err)
	} else if len(mxRecords) > 0 {
		hasMX = true
	}

	// Controllo dei record SPF
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Errore nel cercare record TXT per %s: %v\n", domain, err)
	} else {
		for _, record := range txtRecords {
			if strings.HasPrefix(record, "v=spf1") {
				hasSPF = true
				spfRecord = record
				break
			}
		}
	}

	// Controllo dei record DMARC
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Errore nel cercare record DMARC per %s: %v\n", domain, err)
	} else {
		for _, record := range dmarcRecords {
			if strings.HasPrefix(record, "v=DMARC1") {
				hasDMARC = true
				dmarcRecord = record
				break
			}
		}
	}

	return hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain,hasMX,hasSPF,spfRecord,hasDMARC,dmarcRecord\n")

	// Legge i domini dalla riga di comando
	for scanner.Scan() {
		domain := scanner.Text()
		hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord := checkDomain(domain)

		// Stampa i risultati
		fmt.Printf("%s,%t,%t,%s,%t,%s\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
	}

	// Gestione degli errori dello scanner
	if err := scanner.Err(); err != nil {
		log.Fatalf("Errore nella lettura dell'input: %v\n", err)
	}
}
