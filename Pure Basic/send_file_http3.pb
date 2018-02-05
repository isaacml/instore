host$ = "192.168.4.22"
path$ = "/info.cgi"
fichero.s = ""
mp3.s = "C:/Users/Isaac/Desktop/Cuñas/Supersol_Andalucia.mp3"

; Define.i FileLength, File, Base64Length
; Define *FileBuffer
; 
; InitNetwork()
; ConnectionID = OpenNetworkConnection(host$, 9999) 
; 
; If ReadFile(0, mp3)
;   While Eof(0) = 0           ; loop as long the 'end of file' isn't reached
;     fichero + ReadString(0)      ; display line by line in the debug window
;   Wend
;       request$  = "POST " + path$ + " HTTP/1.1" + Chr(13) + Chr(10)
;       request$  + "Host: " + host$ + Chr(13) + Chr(10) 
;       request$  + "User-Agent: Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.2.1) Gecko/20021204" + Chr(13) + Chr(10)
;       request$  + "Accept: text/xml,application/xml,application/xhtml+xml," 
;       request$  + "text/html;q=0.9,text/plain;q=0.8,audio/mpeg,video/x-mng,image/png," 
;       request$  + "image/jpeg,image/gif;q=0.2,text/css,*/*;q=0.1" + Chr(13) + Chr(10) 
;       request$  + "Accept-Language: en-us, en;q=0.50" + Chr(13) + Chr(10) 
;       request$  + "Accept-Encoding: gzip, deflate, compress;q=0.9" + Chr(13) + Chr(10) 
;       request$  + "Accept-Charset: ISO-8859-1, utf-8;q=0.66, *;q=0.66" + Chr(13) + Chr(10) 
;       request$  + "Keep-Alive: 300" + Chr(13) + Chr(10) 
;       request$  + "Connection: keep-alive" + Chr(13) + Chr(10) 
;       request$  + "Referer: http://www.google.de/" + Chr(13) + Chr(10) 
;       request$  + "Cache-Control: max-age=0" + Chr(13) + Chr(10) 
;       request$  + "Content-Type: application/x-www-form-urlencoded" + Chr(13) + Chr(10)
;       ;request$  + "Content-Disposition: form-Data; name="+Chr(34)+"filesend"+Chr(34)+"; filename="+Chr(34)+mp3+ Chr(34) + Chr(13) + Chr(10)
;       request$  + "Content-Length: " + FileLength + Chr(13) + Chr(10) 
;       request$  + Chr(13) + Chr(10)
;       request$  + "perico=pedro"
;       SendNetworkString(ConnectionID, request$)
;   CloseFile(0)
; Else
;   MessageRequester("Information","No se puede abrir el fichero!")
; EndIf

If InitNetwork()
  conid.l = OpenNetworkConnection(host$,9999)
  If conid
      Debug "Connected"
      *Buffer = AllocateMemory(1000000000) ; some memory for our file buffer
      POST$ = "POST "+ path$ +" HTTP/1.1"+#CRLF$+"Host: "+host$+#CRLF$+"Accept: */*"+#CRLF$+"User-Agent: Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.2.1) Gecko/20021204"+#CRLF$
   
      OpenFile(1,mp3)
      Repeat
        Text$ = ReadString(1)
        FILE$ + EOL$+Text$
      Until Eof(1)
      ; This is the border header for uploading
      FILE$ = "--"+Border$ + EOL$ + FileHeader$ +EOL$ + FILE$ +EOL$+ "--" + Border$ + "--" 
      ; Back to post, while sending header with the correct content length (border+file+border)
      POST$ + EOL$ + "Content-Type: multipart/form-Data, boundary="+Border$ + EOL$ + "Content-Length: " + Str(Len(FILE$)) 
      POST$ + EOL$ + EOL$ + FILE$ 
      CloseFile(1)
      Debug POST$
      Debug "+++++++++++++++++++"
      PokeS(*Buffer,"",0)
      PokeS(*Buffer,POST$,Len(POST$))
      SendNetworkData(conid,*Buffer,Len(POST$))
  Else 
    Debug "NO CONNECTION"
  EndIf

EndIf
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 59
; FirstLine = 16
; EnableXP