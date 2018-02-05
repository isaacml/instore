URL$ = "192.168.4.22" ; the main domain "posttestserver.com" is a good test domain
PATH$ = "/info.cgi"   ; /post.php works fine at "posttestserver.com" or what ever script that accepts the enctype="multipart/form-data"    

fileToUpload.s = "C:/Users/Isaac/Desktop/Cuñas/Supersol_Andalucia.mp3" ; Full path+filename, choose your own path
fileLen = FileSize(fileToUpload)
boundary.s = "aA7hhfk4"
actionType.s = "POST " + PATH$ + " HTTP/1.0" + #CRLF$
host.s = "Host: " + URL$ + #CRLF$
agent.s = "User-Agent: UploaderFiles 1.0" + #CRLF$ + #CRLF$
majorContentType.s = "Content-type: multipart/form-data, boundary=" + boundary + #CRLF$
majorContentLength.s = "Content-Length: <length>" + #CRLF$
submit.s = "Content-Disposition: form-data; name=" + #DOUBLEQUOTE$ + "submit" + #DOUBLEQUOTE$ + #CRLF$ + "Submit"
boundary = "--" + boundary + #CRLF$
fileUploadField.s = boundary + "Content-Disposition: form-data; name=" + #DOUBLEQUOTE$ + "file" + #DOUBLEQUOTE$ +"; filename=" + #DOUBLEQUOTE$ + fileToUpload + #DOUBLEQUOTE$ + #CRLF$
minorContentType.s = "Content-Type: audio/mpeg3" + #CRLF$

If InitNetwork()
  conid.l = OpenNetworkConnection(URL$,9999)
  If conid
    Debug "Connected"
    *buff = AllocateMemory(fileLen)
    ReadFile(0,fileToUpload);
    ReadData(0, *buff, fileLen)
    finalBoundary.s = ReplaceString(boundary,#CRLF$,"--",#PB_String_NoCase,3)
    requestHeader.s = actionType + host + agent
    finisher.s = #CRLF$ + boundary + submit + #CRLF$+ finalBoundary
    postDataLength = Len(fileUploadField) + Len(minorContentType) + Len(#CRLF$)+ fileLen + Len(finisher)
    majorContentLength = ReplaceString(majorContentLength,"<length>",Str(5490))
    inputData.s = requestHeader + majorContentType +  majorContentLength
    inputData = inputData + fileUploadField + minorContentType + #CRLF$
    
    OpenFile(1,"C:/Users/Isaac/Desktop/file.dat")
    sendBuffer = AllocateMemory(Len(inputData) + fileLen + 512)
    PokeS(sendBuffer,inputData)
    CopyMemory(*buff,sendBuffer+Len(inputData),fileLen)
    WriteData(1,sendBuffer,Len(inputData) + fileLen)
    CloseFile(1)
    Debug(Str(FileSize("C:/Users/Isaac/Desktop/file.dat")))
    PokeS(sendBuffer + Len(inputData) + fileLen,finisher)
    res.l = SendNetworkData(conid,sendBuffer,postDataLength)
  Else 
    Debug "NO CONNECTION"
  EndIf
EndIf
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 39
; EnableXP