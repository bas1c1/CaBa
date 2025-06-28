##"set{\"abc2\",\"qwerty123\"}")
##"get{\"abc\"}"
##"multiget{\"abc\"}"
##"del{\"abc\"}"
##"updatecache{}"
##"loadcache{}"
##"save{}"
##"asave{}"
##"loadfrom{\"backup 2025-02-11-0\"}"
##"create_db{\"name/lol\"}"
##"choose_db{\"name\"}"
##"zip{}"
##"unzip{}"
##"list{}"
##"choose_db{\"name\"}list{}choose_db{\"name/lol\"}list{}"

import struct
import socket

conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
conn.connect( ("127.0.0.1", 8080) )

def _recv_exact(n: int) -> bytes:
    data = b""
    while len(data) < n:
        chunk = conn.recv(n - len(data))
        if not chunk:
            break
        data += chunk
    return data

def send_message(msg):
    d = msg.encode()
    conn.send(struct.pack(">I", len(d)) + d)

def recieve_message():
    length_bytes = _recv_exact(4)
    if len(length_bytes) != 4:
        raise Exception("Connection closed")
    
    length = struct.unpack(">I", length_bytes)[0]
    
    return _recv_exact(length)

send_message("choose_db{\"name\"}")
recieve_message()

send_message("get{\"" + str(55) + "\"}")
print(recieve_message())