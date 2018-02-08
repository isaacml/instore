Procedure.s POST_PB(ConnectionID, host$, path$, parameters$)
  lenstr$ = Str(Len(parameters$)) ;Longitud de los parámetros
  request$  = "POST " + path$ + " HTTP/1.1" + Chr(13) + Chr(10)
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
  request$  + "Referer: http://www.google.de/" + Chr(13) + Chr(10) 
  request$  + "Cache-Control: max-age=0" + Chr(13) + Chr(10) 
  request$  + "Content-Type: application/x-www-form-urlencoded" + Chr(13) + Chr(10)
  request$  + "Content-Length: " + lenstr$ + Chr(13) + Chr(10) 
  request$  + Chr(13) + Chr(10) 
  request$  + parameters$
  
  SendNetworkString(ConnectionID, request$)
  
  While NetworkClientEvent(ConnectionID) <> 2 
  Delay(1) 
  Wend 
  *Buffer = AllocateMemory(50000) 
  ReceiveNetworkData(ConnectionID, *Buffer, 50000)
  res$ = PeekS(*Buffer, -1, #PB_UTF8)
  
  res$ = StringField(res$, 2, Chr(13) + Chr(10) + Chr(13) + Chr(10))
  
  ProcedureReturn res$ ;Devuelve la cadena de respuesta
  
  FreeMemory(*Buffer)
  CloseNetworkConnection(ConnectionID) 
EndProcedure

;Librería que copia una página HTML completa
Procedure.s DOWN_PAGE(host$, port.l, origen$, destino$)
  If ReceiveHTTPFile("http://"+host$+":"+port.l+"/"+origen$, destino$)
    res$ = "Success"
  Else
    res$ = "Failed"
  EndIf
  ProcedureReturn res$ ;Devuelve respuesta
EndProcedure

;Librería que toma un fichero MP3 del servidor y realiza una copia
Procedure.s DOWN_MP3(host$, port.l, origen$, action$, destino$)
  If ReceiveHTTPFile("http://"+host$+":"+port.l+"/"+origen$+"?accion="+action$, destino$)
    res$ = "Success"
  Else
    res$ = "Failed"
  EndIf
  ProcedureReturn res$ ;Devuelve respuesta
EndProcedure
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 55
; FirstLine = 4
; Folding = -
; EnableXP