XML$ = "<div class='panel-body'><Select name='entidad'><option value='' selected='selected'>Selecciona una entidad</option><option value='2'>Dinosol</option></Select></div>"
Test$ = "<test><head/><body><greeting><hello>world</hello><trace>1</trace></greeting></body></test>"

Procedure explodeStringArray(Array a$(1), s$, delimeter$)
  Protected count, i
  count = CountString(s$,delimeter$) + 1
  Dim a$(count)
  For i = 1 To count
    a$(i - 1) = StringField(s$,i,delimeter$)
  Next
  ProcedureReturn count ;return count of substrings
EndProcedure

Dim output.s(0) ;this will be resized later

explodeStringArray(output(), "<div class='panel-heading'>Entidad</div><div class='panel-body'><Select name='entidad'><option value='' selected>Selecciona una entidad</option><option value='2'>Dinosol</option><option value='5'>Moraleja</option><option value='3'>Mercadona</option></Select></div>", "</option><option")
For i = 1 To ArraySize(output())-1
  value.s = output(i)
  ocurrencias = CountString(value, "</")
  firstappear = FindString(value, ">")+1
  onlyname.s = Mid(value, firstappear)
  If ocurrencias > 0
    lastappear = FindString(onlyname, "</") -1
    name.s = Left(onlyname, lastappear) 
  Else
    name.s = onlyname
  EndIf 
  Debug name
;   If ocurrencias > 0
;     id = FindString(value, "'") +1
;     sup.s = Mid(value, id, id-8)
;     cut2 = FindString(value, "</") -1
;     first.s = Left(value, cut2)
;     name = FindString(first, "'>")-1
;     del.s = Right(first, name)
;   Else
;     id = FindString(value, "'") +1
;     sup.s + Mid(value, id, id-8)
;     name = FindString(value, "'>")-3
;     nombre.s = Right(value, name)
;     del.s = nombre
;   EndIf
;   Debug del
;   Debug sup
  ;cut1 = FindString(value, ">") + 1
  ;first_key.s = Mid(value, cut1)
  ;Debug first_key
  ;cut2 = FindString(first_key, "</") -1
  ;final_key.s = Left(first_key, cut2)
  ;Debug final_key
Next
; IDE Options = PureBasic 5.61 (Windows - x64)
; CursorPosition = 20
; FirstLine = 7
; Folding = -
; EnableXP