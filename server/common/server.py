import socket
import logging
import signal
from common.protocol import Protocol
from common.utils import Bet, store_bets

FINISHED_MSG = 2
WINNERS_REQUEST_MSG = 3
BET_MSG = 1

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.running = False
        self.clients_list = []
        self.amount_consulting_agencies = 0

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        self.running = True
        signal.signal(signal.SIGTERM, self.shutdown)

        while self.running:
            try:
                client_sock = self.__accept_new_connection()
                self.clients_list.append(client_sock)
                self.__handle_client_connection(client_sock)
            except:
                if not self.running:
                    break

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            while True:
                request = Protocol.receive_client_request_type(client_sock)

                if request == FINISHED_MSG:
                    self.amount_consulting_agencies += 1
                    break
                elif request == WINNERS_REQUEST_MSG:
                    message = Protocol.receive_client_request(client_sock)
                    logging.info(f'action: receive_client_request WINNERS | result: {message}')
                    if self.amount_consulting_agencies == 3:
                        logging.info(f'action: recv_WINNERS_req | result: success')
                    else:
                        logging.info(f'action: recv_WINNERS_req | result: FAIL')
                    break
                elif request == BET_MSG:
                    message = Protocol.receive_client_request(client_sock)
                    bets = Protocol.deserialize_bets(message)
                    Protocol.send_ack(client_sock, len(bets))
                    logging.info(f'action: send_ack | result: success | cantidad: {len(bets)}')
                    store_bets(bets)
                    logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
                else:
                    logging.error(f'action: receive_message | result: fail | error: Invalid request type: {request}')
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
    
    def shutdown(self, signum = None, frame = None):
        print("Received SIGTERM, shutting down...")
        self.running = False
        if self._server_socket:
            try:
                self._server_socket.close()
                logging.info('action: close_server_socket | result: success')
            except OSError as e:
                logging.error("action: close_server_socket | result: fail | error: {e}")
        
        for client_skt in self.clients_list:
            try:
                client_skt.close()
                logging.info('action: close_client_socket | result: success')
            except OSError as e:
                logging.error("action: close_client_socket | result: fail | error: {e}")
        logging.info('action: shutdown | result: success')
        return 0
