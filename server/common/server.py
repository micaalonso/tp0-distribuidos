import socket
import logging
import signal
from common.protocol import Protocol
from common.utils import Bet, store_bets


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.running = False

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
            # Recibo de la apuesta
            client_bet = Protocol.receive_bet(client_sock)
            logging.info(f'action: apuesta_recibida | result: success | dni: {client_bet.document} | number: {client_bet.number}')

            #Almacenamiento de la apuesta
            store_bets([client_bet])
            logging.info(f'action: apuesta_almacenada | result: success | dni: {client_bet.document} | number: {client_bet.number}')
            
            # Envío de la respuesta
            Protocol.send_ack(client_sock, client_bet.document)
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
        return 0
