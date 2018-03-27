
;XML$ = "<div class='panel-heading'>Entidad</div><div class='panel-body'><Select name='entidad'><option value='' selected>Selecciona una entidad</option><option value='2'>Dinosol</option></Select></div>"


; CatchXML(0, @XML$,Len(XML$))
;  
; If IsXML(0)
;   *algo = MainXMLNode(0)
;   Debug GetXMLNodeName(ChildXMLNode(*algo))
;   FreeXML(0)
; EndIf 

XML$ = "<div class='panel-body'><Select name='entidad'><option value='' selected='selected'>Selecciona una entidad</option><option value='2'>Dinosol</option></Select></div>"
Test$ = "<test><head/><body><greeting><hello>world</hello><trace>1</trace></greeting></body></test>"
If ParseXML(0, XML$)
   ;Debug GetXMLNodeText(XMLNodeFromPath(MainXMLNode(0), "Select/option")) ; --> "world" --> OK
   ;Debug GetXMLNodeText(XMLNodeFromPath(MainXMLNode(0), "body/greeting/trace")) ; --> "1" --> OK
   If CreateXML(1, #PB_UTF8)
      CopyXMLNode(XMLNodeFromPath(MainXMLNode(0), "Select"), RootXMLNode(1))
      Debug ComposeXML(1, #PB_XML_NoDeclaration)
      FreeXML(1)
   EndIf
   FreeXML(0)
EndIf

Procedure explodeStringArray(Array a$(1), s$, delimeter$)
  Protected count, i
  count = CountString(s$,delimeter$) + 1
  
  Debug Str(count) + " substrings found"
  Dim a$(count)
  For i = 1 To count
    a$(i - 1) = StringField(s$,i,delimeter$)
  Next
  ProcedureReturn count ;return count of substrings
EndProcedure

Dim output.s(0) ;this will be resized later

explodeStringArray(output(), "<div class='panel-heading'>Entidad</div><div class='panel-body'><Select name='entidad'><option value='' selected>Selecciona una entidad</option><option value='2'>Dinosol</option><option value='3'>Mercadona</option></Select></div>", "<option")
For i = 1 To ArraySize(output())
  value.s = output(i)
  Debug value
  cut1 = FindString(value, ">") + 1
  first_key.s = Right(value, cut1)
  cut2 = FindString(first_key, "</") + 1 
  Debug cut2
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
; CursorPosition = 42
; FirstLine = 8
; Folding = -
; EnableXP