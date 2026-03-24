package common

import (
	"net"
	"time"
	"os"
	"strconv"
	"bufio"
	"io"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	MaxBatchSize  int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	bet    Bet
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(sigterm_channel chan os.Signal) {
	if len(sigterm_channel) > 0 {
		<-sigterm_channel

		log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
		return
	}

	// Create the connection the server in every loop iteration. Send an
	c.createClientSocket()
	defer c.conn.Close()

	file, err := os.Open("/data/agency-" + c.config.ID + ".csv")
	if err != nil {
		log.Errorf("action: open_csv | result: fail | error: %v", err)
		return
	}
	log.Infof("action: open_csv | result: success")
	defer file.Close()

	reader := bufio.NewReader(file)
	batch_size := c.config.MaxBatchSize

	for {
		select {
		case <-sigterm_channel:
			log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
			return
		default:
			batch, err := next_batch(reader, batch_size, c.config.ID)
			if err == io.EOF {
				// Error de EOF es esperado, me indica que temriné de leer el csv
				break
			}
			if err != nil {
				log.Errorf("action: next_batch | result: fail | error: %v", err)
				return
			}

			send_batch(batch, c.conn)

			log.Infof("action: send_batch | result: success | client_id: %v", c.config.ID)
		}
	}

	// if send_bet(c.bet, c.conn) != nil {
	// 	return
	// }

	// ack, err := read_ack(c.conn)
	// if err != nil {
	// 	return
	// }

	// if check_ack(ack, c.bet.Document, c.bet.Number) != nil {
	// 	return
	// }

	// c.conn.Close()
}

func check_ack(ack AckAnswer, expected_document int, expected_number int) error {
	recv_document, err := strconv.Atoi(ack.Document)
	if err != nil {
		log.Errorf("action: check_ack | result: fail | invalid document: %v", err)
		return err
	}

	recv_number, err := strconv.Atoi(ack.Number)
	if err != nil {
		log.Errorf("action: check_ack | result: fail | invalid number: %v", err)
		return err
	}

	if recv_document == expected_document && recv_number == expected_number {
		log.Infof("action: apuesta_enviada | result: success | dni: %v | number: %v",
			expected_document, expected_number)
	} else {
		log.Errorf("action: apuesta_enviada | result: fail | expected: %v|%v | got: %v|%v",
			expected_document, expected_number, recv_document, recv_number)
	}
	return nil
}
