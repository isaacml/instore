server$  = "192.168.0.102"                        ; Server Externo
port.l   = 8080                                  ; Port
domain_file$  = "configshop.reg"
settings_file$ = "SettingsShop.reg"
  
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

Procedure.s loadDomains(file$, List doms.s())
  If ReadFile(0, file$)   ; if the file could be read, we continue...
    While Eof(0) = 0      ; loop as long the 'end of file' isn't reached
      dom$ = StringField(ReadString(0), 2, " = ")
      AddElement(doms())
      doms() = dom$
    Wend
    CloseFile(0)
  EndIf
EndProcedure
; IDE Options = PureBasic 5.61 (Windows - x64)
; Folding = -
; EnableXP