;#File$ = "c:\tmp\1.jpg"
#File$ = "C:/Users/Isaac/Desktop/Cuñas/Supersol_Andalucia.mp3"
#Boundary$ = "---------------------------47972514120"

Define.l Bytes_read
Define.i SendFile, ReceiveFile, Length, Ptr
Define.i Open_handle, Connect_handle, Request_handle, Send_handle
Define.s Host, Page, Header, Post, Buffer, Html

SendFile = ReadFile(#PB_Any, #File$)
If SendFile
  Length = Lof(SendFile)                            
  *MemoryID = AllocateMemory(Length + 1024)
  
  If *MemoryID
    
    Post = "--" + #Boundary$ + #CRLF$
    Post + "Content-Disposition: form-data; name=" + #DQUOTE$ + "attached" + #DQUOTE$ + "; filename=" + #DQUOTE$ + "publicidad.mp3" + #DQUOTE$ + #CRLF$
    Post + "Content-Type: audio/mpeg" + #CRLF$
    Post + #CRLF$
    
    PokeS(*MemoryID, Post, -1, #PB_Ascii)
    
    Ptr = Len(Post)
    
    If ReadData(SendFile, *MemoryID + Ptr, Length) = Length
      Ptr + Length
      
      Post = #CRLF$ + "--" + #Boundary$ + "--" + #CRLF$
      
      PokeS(*MemoryID + Ptr, Post, -1, #PB_Ascii)
      Ptr + Len(Post)
      
      
      ;Wininet stuff
      Host = "192.168.4.22"
      Page = "/info.cgi"
      Open_handle = InternetOpen_("Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)", 1, "", "", 0)
      InternetSetOption_(Open_handle, 2, 1000, 4)
      Connect_handle = InternetConnect_(Open_handle, Host, 9999, "", "", 3, 0, 0)
      Request_handle = HttpOpenRequest_(Connect_handle, "POST", Page, "", "", 0, $00080000|$00000100|$04000000, 0)
      
      Header = "Content-Type: multipart/form-data; boundary=" + #Boundary$ + #CRLF$ + #CRLF$
      HttpAddRequestHeaders_(Request_handle, Header, Len(Header), $80000000|$20000000)
      
      Send_handle = HttpSendRequest_(Request_handle, "", 0, *MemoryID, Ptr)
      
      Buffer = Space(1024)
      Repeat
        InternetReadFile_(Request_handle, @Buffer, 1024, @Bytes_read)
        Html + Left(PeekS(@Buffer, -1, #PB_Ascii), Bytes_read)
        Buffer = Space(1024)
      Until Bytes_read = 0
      
      InternetCloseHandle_(Send_handle)
      InternetCloseHandle_(Request_handle)
      InternetCloseHandle_(Connect_handle)
      InternetCloseHandle_(Open_handle)
      
      ; write received html to file
      ReceiveFile = CreateFile(#PB_Any, "C:\Output.html")
      If ReceiveFile
        WriteString(ReceiveFile, Html)
        CloseFile(ReceiveFile)
        RunProgram("C:\Output.html")
      EndIf
    EndIf
    FreeMemory(*MemoryID)
  EndIf
  CloseFile(SendFile)
EndIf
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 45
; FirstLine = 18
; EnableXP