package common

import (
	"fmt"
	"net"
	"strings"
	"bufio"
	"encoding/binary"
)

const EXPECTED_NUM_PARTS = 2
const DOCUMENT_INDEX = 0
const NUMBER_INDEX = 1

type AckAnswer struct {
	Document  string
	Number 	  string
}

func send_bet(bet Bet, conn net.Conn) error {
	formatted_bet := fmt.Sprintf("%s|%s|%d|%s|%d",
		bet.Name,
		bet.Surname,
		bet.Document,
		bet.Birthday,
		bet.Number,
	)

	// Envio el largo de la apuesta en uint16
	err := binary.Write(conn, binary.BigEndian, uint16(len(formatted_bet)))
	if err != nil {
		log.Errorf("action: send_leght_bet | result: fail | client_document: %v | error: %v",
			bet.Document,
			err,
		)
		return err
	}

	// Envío de la apuesta en si
	serialized_bet := []byte(formatted_bet)
	sent_bytes := 0

	//For para asegurarnos que se envían todos los bytes de la apuesta
	for sent_bytes < len(serialized_bet) {
		num_bytes, err := conn.Write(serialized_bet[sent_bytes:])
		if err != nil {
			log.Errorf("action: send_bet | result: fail | client_document: %v | error: %v",
			bet.Document,
			err,
			)
			return err
		}
		sent_bytes += num_bytes
	}

	log.Infof("action: send_bet | result: success | client_document: %v",
		bet.Document,
	)
	return nil

}

func read_ack(conn net.Conn) (AckAnswer, error) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Errorf("action: read_ack | result: fail | error: %v", err)
		return AckAnswer{}, err
	}

	msg = strings.TrimSpace(msg)

	parts := strings.Split(msg, "|")
	if len(parts) != EXPECTED_NUM_PARTS {
		log.Errorf("action: read_ack | result: fail | invalid ack format: %v", msg)
		return AckAnswer{}, fmt.Errorf("Invalid ACK format")
	}

	ack := AckAnswer{
		Document: parts[DOCUMENT_INDEX],
		Number:   parts[NUMBER_INDEX],
	}
	return ack, nil
}