server$  = "192.168.0.102"                       ; Server Externo
port.l   = 8080                                  ; Port
settings_file$ = "SettingsShop.reg"

Procedure MP3_Load(Nb,file.s)
  i=mciSendString_("OPEN "+Chr(34)+file+Chr(34)+" Type MPEGVIDEO ALIAS MP3_"+Str(Nb),0,0,0)
  If i=0
    ProcedureReturn #True
  Else
    ProcedureReturn #False
  EndIf
EndProcedure

Procedure MP3_Play(Nb)
  i=mciSendString_("play MP3_"+Str(Nb),0,0,0)
  ProcedureReturn i
EndProcedure

Procedure MP3_Free(Nb)
  i=mciSendString_("close MP3_"+Str(Nb),0,0,0)
  ProcedureReturn i
EndProcedure

Procedure.s obtain_Id_Name(Array a$(1), s$, delimeter$, Map Valores.s())
  Protected count, i, f
  count = CountString(s$,delimeter$) + 1
  Dim a$(count)
  For i = 1 To count
    a$(i - 1) = StringField(s$,i,delimeter$)
  Next
  For f = 1 To ArraySize(a$())-1
    value.s = a$(f)
    ocurrencias = CountString(value, "</")
    firstappear = FindString(value, ">")+1
    onlyname.s = Mid(value, firstappear)
    If ocurrencias > 0
      firstidappear = FindString(value, "'")+1
      id.s = Mid(value, firstidappear)
      lastidappear = FindString(id, "'")-1
      endid.s = Left(id, lastidappear)
      lastappear = FindString(onlyname, "</") -1
      name.s = Left(onlyname, lastappear) 
    Else
      firstidappear = FindString(value, "'")+1
      id.s = Mid(value, firstidappear)
      lastidappear = FindString(id, "'")-1
      endid.s = Left(id, lastidappear)
      name.s = onlyname
    EndIf
    Valores(endid) = name
  Next
EndProcedure

Procedure.s obtainMsgFiles(directory$, List MsgFiles.s())
  If ExamineDirectory(0, directory$, "*.*")  
    While NextDirectoryEntry(0)
      If DirectoryEntryType(0) = #PB_DirectoryEntry_File
        is_mp3 = CountString(DirectoryEntryName(0), ".mp3")
        is_wma = CountString(DirectoryEntryName(0), ".wma")
        If is_mp3 = 1
          AddElement(MsgFiles())
          MsgFiles() = DirectoryEntryName(0)
        ElseIf is_wma = 1
          AddElement(MsgFiles())
          MsgFiles() = DirectoryEntryName(0)
        EndIf
      EndIf
    Wend
    FinishDirectory(0)
  EndIf
EndProcedure

;Genera el select de horas que necesita la Tienda
Procedure formar_horas(gadget)
  For a = 0 To 23
    If Len(Str(a)) = 1
      AddGadgetItem(gadget, -1, RSet(Str(a), 2, "0"))
    Else
      AddGadgetItem(gadget, -1, Str(a))
    EndIf
  Next
EndProcedure
;Genera el select de minutos que necesita la Tienda
Procedure formar_minutos(gadget)
  For a = 0 To 59
    If Len(Str(a)) = 1
      AddGadgetItem(gadget, -1, RSet(Str(a), 2, "0"))
    Else
      AddGadgetItem(gadget, -1, Str(a))
    EndIf
  Next
EndProcedure

Procedure.s load_Domains(file$, List doms.s())
  If ReadFile(0, file$)   ; if the file could be read, we continue...
    While Eof(0) = 0      ; loop as long the 'end of file' isn't reached
      dom$ = StringField(ReadString(0), 2, " = ")
      AddElement(doms())
      doms() = dom$
    Wend
    CloseFile(0)
  EndIf
EndProcedure

Procedure.s POST_PB_STORE(ConnectionID, host$, path$, parameters$)
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
; IDE Options = PureBasic 5.61 (Windows - x64)
; CursorPosition = 1
; Folding = --
; EnableXP