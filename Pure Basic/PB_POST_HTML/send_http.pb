server$   = "x.x.x.x"              ; IP address 
host$     = "192.168.1.173"              ; Domain name 
target$   = "/info.cgi" ; Specific program 
referer$  = "http://www.google.de/"       ; Referer 
port.l    = 9999                            ; Port 

post$     = "test1=bla&test2=foo"
length = Len(post$)
lenstr$ = Str(length)

request$  = "POST " + target$ + " HTTP/1.1" + Chr(13) + Chr(10)
request$  + "Host: " + host$ + Chr(13) + Chr(10) 
request$  + "User-Agent: Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.2.1) Gecko/20021204" + Chr(13) + Chr(10)
request$  + "Accept: text/xml,application/xml,application/xhtml+xml," 
request$  + "text/html;q=0.9,text/plain;q=0.8,video/x-mng,image/png," 
request$  + "image/jpeg,image/gif;q=0.2,text/css,*/*;q=0.1" + Chr(13) + Chr(10) 
request$  + "Accept-Language: en-us, en;q=0.50" + Chr(13) + Chr(10) 
request$  + "Accept-Encoding: gzip, deflate, compress;q=0.9" + Chr(13) + Chr(10) 
request$  + "Accept-Charset: ISO-8859-1, utf-8;q=0.66, *;q=0.66" + Chr(13) + Chr(10) 
request$  + "Keep-Alive: 300" + Chr(13) + Chr(10) 
request$  + "Connection: keep-alive" + Chr(13) + Chr(10) 
request$  + "Referer: " + referer$ + Chr(13) + Chr(10) 
request$  + "Cache-Control: max-age=0" + Chr(13) + Chr(10) 
request$  + "Content-Type: application/x-www-form-urlencoded" + Chr(13) + Chr(10)
request$  + "Content-Length: " + lenstr$ + Chr(13) + Chr(10) 
request$  + Chr(13) + Chr(10) 
request$  + post$ 

InitNetwork()
ConnectionID = OpenNetworkConnection(host$, port.l) 

SendNetworkString(ConnectionID, request$) 

While NetworkClientEvent(ConnectionID) <> 2 
  Delay(1) 
Wend 
*Buffer = AllocateMemory(50000) 
ReceiveNetworkData(ConnectionID, *Buffer, 50000)
resultado$ = PeekS(*Buffer, -1, #PB_UTF8)

Debug StringField(resultado$, 1, "\r\n")

FreeMemory(*Buffer)

CloseNetworkConnection(ConnectionID) 

End
; IDE Options = PureBasic 5.61 (Windows - x64)
; CursorPosition = 1
; EnableXP