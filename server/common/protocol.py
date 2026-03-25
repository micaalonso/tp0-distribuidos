from .utils import Bet
import struct
import logging


MESSAGE_LENGHT_BYTES = 2
FIRST_NAME_POSITION = 0
LAST_NAME_POSITION = 1
DOCUMENT_POSITION = 2
BIRTHDATE_POSITION = 3
NUMBER_POSITION = 4
AGENCY_POSITION = 5
AMOUNT_BET_PARTS = 6

class Protocol:
    @staticmethod
    def receive_client_request(socket):
        msg_lenght_recv = Protocol._receive_message(socket, MESSAGE_LENGHT_BYTES)
        msg_lenght = struct.unpack(">H", msg_lenght_recv)[0]

        msg = Protocol._receive_message(socket, msg_lenght)
        decoded_msg = msg.decode()
        return decoded_msg

    @staticmethod
    def deserialize_bets(bets):
        formatted_bets = []
        bets_lines = bets.strip().split("\n") # separo las bets
        for bet_line in bets_lines:
            parts = bet_line.split("|")

            if len(parts) != AMOUNT_BET_PARTS:
                logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(formatted_bets)} | error: Invalid bet data")
                raise ValueError("Invalid bet data")
        
            bet = Bet(
                first_name=parts[FIRST_NAME_POSITION],
                last_name=parts[LAST_NAME_POSITION],
                document=int(parts[DOCUMENT_POSITION]),
                birthdate=parts[BIRTHDATE_POSITION],
                number=int(parts[NUMBER_POSITION]),
                agency=parts[AGENCY_POSITION],
            )

            formatted_bets.append(bet)

        return formatted_bets
    
    @staticmethod
    def _receive_message(socket, bytes):
        message = b""
        while len(message) < bytes:
            recv_bytes = socket.recv(bytes - len(message))
            if not recv_bytes:
                raise ConnectionError("Client socket closed")
            message += recv_bytes
        return message
    
    @staticmethod
    def send_ack(socket, amount_bets):
        socket.sendall("{}\n".format(amount_bets).encode('utf-8'))