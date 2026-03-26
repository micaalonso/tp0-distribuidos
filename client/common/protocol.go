package common

import (
	"fmt"
	"net"
	"strings"
	"bufio"
	"strconv"
	"io"
	"encoding/binary"
)

const FINISHED_MESSAGE = 0x02
const WINNERS_REQUEST_MESSAGE = 0x03
const BET_MESSAGE = 0x01

type AckAnswer struct {
	Amount  int
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

	// Envío del byte correspondiente a los mensajes del tipo BET
	send_type_byte(conn, BET_MESSAGE)
	
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

func send_message_lenght(conn net.Conn, msg_lenght int) error {
	err := binary.Write(conn, binary.BigEndian, uint16(msg_lenght))
	if err != nil {
		return err
	}
	return nil
}

func send_finished_transfer(conn net.Conn) error {
	if err := send_type_byte(conn, FINISHED_MESSAGE); err != nil {
		return err
	}
	return nil
}

func send_winners_request(conn net.Conn, agency_id string) error {
	if err := send_type_byte(conn, WINNERS_REQUEST_MESSAGE); err != nil {
		return err
	}
	
	if err := send_message_lenght(conn, len(agency_id)); err != nil {
		log.Errorf("action: send_message_lenght | result: fail | error: %v | message: %v",
			err,
			agency_id,
			)
	}

	if err := send_message(conn, agency_id); err != nil {
		log.Errorf("action: send_message | result: fail | error: %v | message: %v",
			err,
			agency_id,
			)
	}
	return nil
}

func send_type_byte(conn net.Conn, byte uint8) error {
	err := binary.Write(conn, binary.BigEndian, byte)
	if err != nil {
		log.Errorf("action: send_type_byte | result: fail | error: %v", err)
		return err
	}
	return nil
}

func read_winners(conn net.Conn) ([]string, error) {
    size := make([]byte, 2)
	_, err := io.ReadFull(conn, size)
	if err != nil {
		return nil, fmt.Errorf("failed to read winners size: %v", err)
	}

	winnersSize := int(binary.BigEndian.Uint16(size))
	winnersData := make([]byte, winnersSize)
	_, err = io.ReadFull(conn, winnersData)
	if err != nil {
		return nil, fmt.Errorf("failed to read winners data: %v", err)
	}

	if len(winnersData) == 0 {
		return []string{}, nil
	}

	return strings.Split(string(winnersData), "|"), nil
}
