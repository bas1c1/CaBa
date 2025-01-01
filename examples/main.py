import socket
conn = socket.socket()
conn.connect( ("127.0.0.1", 8080) )
conn.send(b"set{\"\\\"kek4\\\"\",\"lol\"}")
#conn.send(b"get{\"\\\"kek4\\\"\"}")
#conn.send(b"del{\"\\\"kek4\\\"\"}")
data = b""
tmp = conn.recv(1024)
while tmp:
    data += tmp
    tmp = conn.recv(1024)
print( data.decode("utf-8") )
conn.close()
