package common

import (
	"fmt"
	"net"
	"encoding/binary"
)

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