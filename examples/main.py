import socket
conn = socket.socket()
conn.connect( ("127.0.0.1", 8080) )
#conn.send(b"set{\"kek7\",\"lol213\"}")
#conn.send(b"get{\"kek7\"}")
#conn.send(b"multiget{\"kek7\",\"ke6\"}")
#conn.send(b"del{\"\\\"kek6\\\"\"}")
#conn.send(b"updatecache{}")
#conn.send(b"loadcache{}")
#conn.send(b"save{}")
#conn.send(b"asave{}")
#conn.send(b"loadfrom{\"backup 2025-01-26\"}")
#conn.send(b"create_db{\"name\"}")
#conn.send(b"choose_db{\"name\"}")
#conn.send(b"zip{}")
#conn.send(b"unzip{}")
data = b""
tmp = conn.recv(1024)
while tmp:
    data += tmp
    tmp = conn.recv(1024)
print( data.decode() )
conn.close()
