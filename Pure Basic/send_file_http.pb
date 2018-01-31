;################################################################
;# 22/2/2006 - Upload file to web trough POST command | Pantcho #
;################################################################
EOL$ = Chr(13)+Chr(10) 
URL$ = "192.168.4.22" ; the main domain
PATH$ = "/info.cgi" ; or what ever script that accepts the enctype="multipart/form-data"
FullFileName$ = "C:/Users/Isaac/Desktop/Cuñas/Supersol_Andalucia.mp3" ; Full path+filename 
ActionName$ = "file" ; this is important!! this action must be the same as  <form ... name="file">
FileHeader$ = "Content-Disposition: form-Data; name="+Chr(34)+ActionName$ + Chr(34) +"; filename="+Chr(34)+ FullFileName$+ Chr(34) +EOL$ 
FileHeader$ + "Content-Type: text/plain" ; <= Here change the content type regarding your file! (text,image etc...) we go on text
; ^^^ note: Havn't been tested with binary files.
Border$ = "23232323232" ; Border to the file data (Check RFC for more info)


If InitNetwork()
  conid.l = OpenNetworkConnection(URL$,9999)
  If conid
      Debug "Connected"
      *Buffer = AllocateMemory(100000) ; some memory for our file buffer
      POST$ = "POST "+ PATH$ +" HTTP/1.0"  ; the Post command we are going to send to the server
   
      OpenFile(1,FullFileName$)
      Repeat
        Text$ = ReadString(1)
        FILE$ + Chr(13) + Chr(10)+Text$
      Until Eof(1)
      ; This is the border header for uploading
      FILE$ = "------"+Border$ + EOL$ + FileHeader$ +EOL$ + FILE$ + "------" + Border$ + "--" 
      ; Back to post, while sending header with the correct content length (border+file+border)
      POST$ + EOL$ + "Content-Type: multipart/form-Data, boundary=----"+Border$ + EOL$ + "Content-Length: " + Str(Len(FILE$)) 
      POST$ + EOL$ + EOL$ + FILE$
      Debug POST$
      CloseFile(1)
      PokeS(*Buffer,"",0)
      PokeS(*Buffer,POST$,Len(POST$))
      SendNetworkData(conid,*Buffer,Len(POST$))
      Repeat
        Server$ = PeekS(*Buffer)
        Debug Server$
        res.l = ReceiveNetworkData(conid, *Buffer, 1000)
      Until Server$ = PeekS(*Buffer)
  Else 
    Debug "NO CONNECTION"
  EndIf
EndIf
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 31
; EnableXP