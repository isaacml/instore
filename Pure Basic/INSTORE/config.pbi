﻿server$  = "192.168.4.22"                        ; Server Externo
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
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 32
; Folding = -
; EnableXP