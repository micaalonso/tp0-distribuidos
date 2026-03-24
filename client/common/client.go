package common

import (
	"net"
	"time"
	"os"
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
        }

        batch, err := next_batch(reader, batch_size, c.config.ID)
        if err == io.EOF {
            if err := send_finished(c.conn); err != nil {
                log.Errorf("action: send_finished | result: fail | error: %v |client_id: %v", err, c.config.ID)
            } else {
                log.Infof("action: send_finished | result: success | client_id: %v", c.config.ID)
            }
            return
        }
        if err != nil {
            log.Errorf("action: next_batch | result: fail | error: %v", err)
            return
        }

        if err := send_batch(batch, c.conn); err != nil {
            log.Errorf("action: send_batch | result: fail | error: %v", err)
            return
        }
        log.Infof("action: send_batch | result: success | client_id: %v | amount: %v", c.config.ID, len(batch))

		ack, err := read_ack(c.conn)
		if err != nil {
			return
		}
		log.Infof("action: read_ack | result: success | client_id: %v | amount: %v", c.config.ID, ack.Amount)
    }

	select {
        case <-sigterm_channel:
            log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
            return
        default:
        }

	if c.conn != nil {
        c.conn.Close()
    }
	log.Infof("action: finished_client | result: success | client_id: %v", c.config.ID)
}

