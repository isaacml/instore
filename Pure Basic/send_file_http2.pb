EOL$ = Chr(13)+Chr(10) 
QT$=Chr(34)
URL$ = "192.168.1.173" ; the main domain "posttestserver.com" is a good test domain
PATH$ = "/info.cgi"         ; /post.php works fine at "posttestserver.com" or what ever script that accepts the enctype="multipart/form-data"    

FullFileName$ = "C:\Users\0oIsa\Desktop\perrillo.txt"
ActionName$ = "filesend"            ; this is important!! this action must be the same as  <form ... name="filesend">
Border$ = "232323RANDOMLETTERSNUMBERS23232" ; Border to the file data (Check RFC for more info)

; added an extra POST variable because sometimes you need this when you are posting like a security password etc.
; so the post on this will be "varnonumber" and the data will be "dataforvarnonumber 1" I called it varnonumber because "var2" would not work!!
FileHeader$ = "Content-Disposition: form-Data; name="+QT$+"varnonumber"+ QT$+EOL$+EOL$+"dataforvarnonumber 1"+EOL$+"--"+Border$+EOL$; Silly var names no numbers

FileHeader$ + "Content-Disposition: form-Data; name="+QT$+ActionName$ + QT$ +"; filename="+QT$+ FullFileName$+ QT$ +EOL$ 
FileHeader$ + "Content-Type: text/plain" ; <= Here change the content type regarding your file! (text,image etc...) we go on text
; ^^^ note: Havn't been tested with binary files.

If InitNetwork()
  conid.l = OpenNetworkConnection(URL$,9999)
  If conid
      Debug "Connected"
      *Buffer = AllocateMemory(100000) ; some memory for our file buffer
      POST$ = "POST "+ PATH$ +" HTTP/1.1"+#CRLF$+"Host: "+URL$+#CRLF$+"Accept: */*"+#CRLF$+"User-Agent: Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.2.1) Gecko/20021204"
   
      OpenFile(1,FullFileName$)
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
      Repeat
        Server$ = PeekS(*Buffer)
        res.l = ReceiveNetworkData(conid, *Buffer, 1000)
      Until Server$ = PeekS(*Buffer)
      Debug Server$
  Else 
    Debug "NO CONNECTION"
  EndIf
EndIf
; IDE Options = PureBasic 5.61 (Windows - x64)
; CursorPosition = 22
; FirstLine = 6
; EnableXP