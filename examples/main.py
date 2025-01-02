import socket
conn = socket.socket()
conn.connect( ("127.0.0.1", 8080) )
conn.send(b"set{\"\\\"kek6\\\"\",\"lol\"}")
#conn.send(b"get{\"\\\"kek6\\\"\"}")
#conn.send(b"del{\"\\\"kek4\\\"\"}")
#conn.send(b"update{}")
#conn.send(b"load{}")
#conn.send(b"save{}")
#conn.send(b"loadfrom{\"backup 2025-01-02\"}")
data = b""
tmp = conn.recv(1024)
while tmp:
    data += tmp
    tmp = conn.recv(1024)
print( data.decode("utf-8") )
conn.close()
