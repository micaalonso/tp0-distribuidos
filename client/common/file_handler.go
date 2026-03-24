package common

import (
	"bufio"
	"io"
	"strings"
	"fmt"
	"strconv"
)

const CSV_PARTS = 5
const NAME_INDEX_IN_CSV = 0
const SURNAME_INDEX_IN_CSV = 1
const DOCUMENT_INDEX_IN_CSV = 2
const BIRTHDATE_INDEX_IN_CSV = 3
const NUMBER_INDEX_IN_CSV = 4

func next_batch(reader *bufio.Reader, batch_size int, agency_id string) ([]Bet, error) {
	batch := make([]Bet, 0, batch_size)

	for i := 0; i < batch_size; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) == 0 {
					break
				}
			} else {
				log.Errorf("action: read_line | result: fail | error: %v", err)
				return nil, err
			}
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Armado de la apuesta
		bet, err := make_bet(line, agency_id)
		if err != nil {
			log.Errorf("action: make_bet | result: fail | error: %v", err)
			return nil, err
		}

		batch = append(batch, bet)
	}

	if len(batch) == 0 {
		return nil, io.EOF
	}
	return batch, nil
}

func make_bet(line string, agency_id string) (Bet, error) {
	parts := strings.Split(line, ",")

	if len(parts) != CSV_PARTS {
		return Bet{}, fmt.Errorf("Invalid CSV format")
	}

	document, err := strconv.Atoi(parts[DOCUMENT_INDEX_IN_CSV])
	if err != nil {
		log.Errorf("action: make_bet | result: fail | invalid document: %v", err)
		return Bet{}, err
	}

	number, err := strconv.Atoi(parts[NUMBER_INDEX_IN_CSV])
	if err != nil {
		log.Errorf("action: make_bet | result: fail | invalid number: %v", err)
		return Bet{}, err
	}

	agency, err := strconv.Atoi(agency_id)
	if err != nil {
		log.Errorf("action: make_bet | result: fail | invalid agency_id: %v", err)
		return Bet{}, err
	}

	return Bet{
		Name:     parts[NAME_INDEX_IN_CSV],
		Surname:  parts[SURNAME_INDEX_IN_CSV],
		Document: document,
		Birthday: parts[BIRTHDATE_INDEX_IN_CSV],
		Number:   number,
		Agency:   agency,
	}, nil
}