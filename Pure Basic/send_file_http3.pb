URL$ = "192.168.4.22" ; the main domain "posttestserver.com" is a good test domain
PATH$ = "/down_probe.cgi"         ; /post.php works fine at "posttestserver.com" or what ever script that accepts the enctype="multipart/form-data"    

FullFileName$ = "C:\Users\Isaac\Desktop\pajaro.txt"

#DQ$ = #DQUOTE$
InitNetwork()

Procedure.s HttpPostMultipart(Server$,Path$,PostVariables$,FileID$,File$,Cookies$="") ;for now only image files
  Protected Request$,Result$,BytesRead,*RecieveBuffer = AllocateMemory(40000)
  Protected FormData$, *Buffer, FileLength, ContentLength, i
  Protected Name$,Value$,String$,FileExt$=LCase(GetExtensionPart(File$))
  Protected EndString$,Text$,TextLength
  Protected ServerID = OpenNetworkConnection(Server$,9999)
  If FileExt$ = "jpg": FileExt$ = "jpeg": EndIf ;freeimagehosting.net fix
  
  If ServerID
    ;{ Make data for multipart
    If PostVariables$
      For i=1 To CountString(PostVariables$,"&")+1
        String$ = StringField(PostVariables$,i,"&")
        Name$   = StringField(String$,1,"=")
        Value$  = StringField(String$,2,"=")
        FormData$ + "--AaB03x"+#CRLF$
        FormData$ + "content-disposition: form-data; name="+#DQ$+Name$+#DQ$+#CRLF$+#CRLF$
        FormData$ + Value$+#CRLF$
        Debug Name$+"="+Value$
      Next
    EndIf
    FormData$ + "--AaB03x"+#CRLF$
    FormData$ + "content-disposition: form-data; name="+#DQ$+FileID$+#DQ$+"; filename="+#DQ$+GetFilePart(File$)+#DQ$+#CRLF$
    FormData$ + "Content-Type: image/"+FileExt$+#CRLF$;++#CRLF$
    FormData$ + "Content-Transfer-Encoding: binary"+#CRLF$+#CRLF$
    ;}
    If OpenFile(0,File$)
      FileLength = Lof(0)
      ContentLength = Len(FormData$)+FileLength+12
      ;{ Make post request
      Request$ = "POST "+Path$+" HTTP/1.1"+#CRLF$
      Request$ + "Host: "+Server$+#CRLF$
      If Cookies$
        Request$ + "Cookie: "+Cookies$+#CRLF$
      EndIf
      Request$ + "User-Agent: "+UAgent$+#CRLF$
      Request$ + "Content-Length: "+Str(ContentLength)+#CRLF$
      Request$ + "Content-Type: multipart/form-data, boundary=AaB03x"+#CRLF$
      Request$ + #CRLF$
      ;}
      Text$ = Request$+FormData$
      TextLength = Len(Text$)
      EndString$ = #CRLF$+"--AaB03x--" ;12
      ;{ Create send buffer
      *Buffer = AllocateMemory(TextLength+FileLength+12)
      CopyMemory(@Text$,*Buffer,TextLength)
      ReadData(0,*Buffer+TextLength,FileLength)
      CopyMemory(@EndString$,*Buffer+TextLength+FileLength,12)
      CloseFile(0)
      ;}
    Else
      Debug "Error opening file!"
    EndIf
    ;{ Send data and recieve answer
    If SendNetworkData(ServerID,*Buffer,MemorySize(*Buffer))
      FreeMemory(*Buffer)
      Repeat
        Delay(2)
      Until NetworkClientEvent(ServerID) = #PB_NetworkEvent_Data
      Repeat
        BytesRead = ReceiveNetworkData(ServerID,*RecieveBuffer,40000)
        Result$ + PeekS(*RecieveBuffer,BytesRead)
        ;Debug Result$
        Delay(600)
      Until NetworkClientEvent(ServerID) <> #PB_NetworkEvent_Data
      FreeMemory(*RecieveBuffer)
    Else
      Debug "Error sending data!"
    EndIf
    ;}
    CloseNetworkConnection(ServerID)
  Else
    Debug "Connection failed!"
  EndIf
  
  ProcedureReturn Result$
EndProcedure

Debug HttpPostMultipart(URL$, PATH$, "saludo=hola", "file", FullFileName$)
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 1
; Folding = -
; EnableXP