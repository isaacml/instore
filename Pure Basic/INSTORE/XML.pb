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

;explodeStringArray(output(), "<div class='panel-heading'>Entidad</div><div class='panel-body'><Select name='entidad'><option value='' selected>Selecciona una entidad</option><option value='2'>Dinosol</option><option value='5'>Moraleja</option><option value='3'>Mercadona</option></Select></div>", "</option><option")
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
Dim output.s(0) ;this will be resized later
NewMap Valores.s()
obtainIdName(output(), "<div class='panel-heading'>Entidad</div><div class='panel-body'><Select name='entidad'><option value='' selected>Selecciona una entidad</option><option value='2'>Dinosol</option><option value='5'>Moraleja</option><option value='3'>Mercadona</option></Select></div>", "</option><option", Valores())
ForEach Valores()
  Debug MapKey(Valores())+ " - "+Valores()
Next

; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 47
; Folding = -
; EnableXP