from .utils import Bet
import struct

MESSAGE_LENGHT_BYTES = 2
AGENCY = "1"
FIRST_NAME_POSITION = 0
LAST_NAME_POSITION = 1
DOCUMENT_POSITION = 2
BIRTHDATE_POSITION = 3
NUMBER_POSITION = 4

class Protocol:
    @staticmethod
    def receive_bet(socket):
        msg_lenght_recv = Protocol._receive_message(socket, MESSAGE_LENGHT_BYTES)
        msg_lenght = struct.unpack(">H", msg_lenght_recv)[0]

        bet_recv = Protocol._receive_message(socket, msg_lenght)
        bet = bet_recv.decode()
        bet_parts = bet.split("|")
        return Bet(
            agency=AGENCY,
            first_name=bet_parts[FIRST_NAME_POSITION],
            last_name=bet_parts[LAST_NAME_POSITION],
            document=bet_parts[DOCUMENT_POSITION],
            birthdate=bet_parts[BIRTHDATE_POSITION],
            number=bet_parts[NUMBER_POSITION],
        )

    @staticmethod
    def _receive_message(socket, bytes):
        message = b""
        while len(message) < bytes:
            recv_bytes = socket.recv(bytes - len(message))
            if not recv_bytes:
                raise ConnectionError("Client socket closed")
            message += recv_bytes
        return message