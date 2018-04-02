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

explodeStringArray(output(), "<div class='panel-heading'>Entidad</div><div class='panel-body'><Select name='entidad'><option value='' selected>Selecciona una entidad</option><option value='2'>Dinosol</option><option value='3'>Mercadona</option></Select></div>", "</option><option")
For i = 1 To ArraySize(output())-1
  value.s = output(i)
  final_key$ = value
  ocurrencias = CountString(value, "</")
  If ocurrencias > 0 
    cut2 = FindString(value, "</") -1
    del.s = Left(value, cut2)
  EndIf
  Debug final_key
  ;cut1 = FindString(value, ">") + 1
  ;first_key.s = Mid(value, cut1)
  ;Debug first_key
  ;cut2 = FindString(first_key, "</") -1
  ;final_key.s = Left(first_key, cut2)
  ;Debug final_key
Next

var$ = output(3)

;Debug RTrim(variable, "</option>")
;CatchXML(0, @XML$,Len(XML$))
;*MainNode = MainXMLNode(0)   
;*SubNode = ChildXMLNode(*MainNode)
;children = XMLChildCount(*MainNode)

;Debug "Text: " + GetXMLNodeText(*SubNode)
;Debug "Name: " + GetXMLNodeName(*SubNode)

; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 22
; Folding = -
; EnableXP