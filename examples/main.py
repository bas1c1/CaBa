import socket
conn = socket.socket()
conn.connect( ("127.0.0.1", 8080) )
#conn.send(b"set{\"abc2\",\"qwerty123\"}")
#conn.send(b"get{\"abc\"}")
#conn.send(b"multiget{\"abc\"}")
#conn.send(b"del{\"abc\"}")
#conn.send(b"updatecache{}")
#conn.send(b"loadcache{}")
#conn.send(b"save{}")
#conn.send(b"asave{}")
#conn.send(b"loadfrom{\"backup 2025-02-11-0\"}")
#conn.send(b"create_db{\"name/lol\"}")
#conn.send(b"choose_db{\"name\"}")
#conn.send(b"zip{}")
#conn.send(b"unzip{}")
#conn.send(b"list{}")
#conn.send(b"choose_db{\"name\"}list{}choose_db{\"name/lol\"}list{}")
data = b""
tmp = conn.recv(1024)
while tmp:
    data += tmp
    tmp = conn.recv(1024)
print( data.decode() )

conn.close()
