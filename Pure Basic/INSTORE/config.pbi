server$  = "192.168.4.22"                        ; Server Externo
port.l   = 8080                                  ; Port
domain_file$  = "configshop.reg"
  
Procedure.s obtainIdName(Array a$(1), s$, delimeter$, Map Valores.s())
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
  
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 70
; FirstLine = 20
; Folding = -
; EnableXP