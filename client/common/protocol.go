package common

import (
	"fmt"
	"net"
	"strings"
	"bufio"
	"strconv"
	"encoding/binary"
)

const EXPECTED_NUM_PARTS = 2
const DOCUMENT_INDEX = 0
const NUMBER_INDEX = 1
const FINISHED_MESSAGE = "Finished"
const WINNERS_REQUEST_MESSAGE = "Winners"

type AckAnswer struct {
	Amount  int
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

	amount_bets, err := strconv.Atoi(msg)
	if err != nil {
		log.Errorf("action: read_ack | result: fail | invalid amount_bets: %v", err)
		return AckAnswer{}, err
	}

	ack := AckAnswer{
		Amount: amount_bets,
	}
	return ack, nil
}

func send_batch(batch []Bet, conn net.Conn) error {
	// Formateo de las apuestas
	var builder strings.Builder
	for _, bet := range batch {
		line := fmt.Sprintf("%s|%s|%d|%s|%d|%d\n",
			bet.Name,
			bet.Surname,
			bet.Document,
			bet.Birthday,
			bet.Number,
			bet.Agency,
		)
		builder.WriteString(line)
	}

	formatted_batch := []byte(builder.String())

	// Envío del largo del mensaje (batch) en uint16 en Big Endian
	batch_lenght := uint16(len(formatted_batch))
	err := binary.Write(conn, binary.BigEndian, batch_lenght)
	if err != nil {
		log.Errorf("action: send_batch_lenght | result: fail | error: %v", err)
		return err
	}

	// Envío del batch
	sent_bytes := 0
	for sent_bytes < len(formatted_batch) {
		num_bytes, err := conn.Write(formatted_batch[sent_bytes:])
		if err != nil {
			log.Errorf("action: send_batch | result: fail | error: %v", err)
			return err
		}
		sent_bytes += num_bytes
	}
	return nil
}

func send_message(conn net.Conn, message string) error {

	// Envio el largo del mensaje de fin en uint16
	err := binary.Write(conn, binary.BigEndian, uint16(len(message)))
	if err != nil {
		log.Errorf("action: send_message_leght | result: fail | error: %v | message: %v",
			err,
			message,
		)
		return err
	}

	// Envío del mensaje en si
	serialized_message := []byte(message)
	sent_bytes := 0

	//For para asegurarnos que se envían todos los bytes del mensaje
	for sent_bytes < len(serialized_message) {
		num_bytes, err := conn.Write(serialized_message[sent_bytes:])
		if err != nil {
			log.Errorf("action: send_message | result: fail | error: %v | message: %v",
			err,
			message,
			)
			return err
		}
		sent_bytes += num_bytes
	}
	return nil
}