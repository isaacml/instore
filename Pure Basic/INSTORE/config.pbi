server$  = "192.168.4.22"                        ; Server Externo
port.l   = 8080                                  ; Port
domain_file$  = "configshop.reg"

Procedure explodeStringArray(Array a$(1), s$, delimeter$)
  Protected count, i
  count = CountString(s$,delimeter$) + 1
  Dim a$(count)
  For i = 1 To count
    a$(i - 1) = StringField(s$,i,delimeter$)
  Next
  ProcedureReturn count ;return count of substrings
EndProcedure
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 12
; Folding = -
; EnableXP