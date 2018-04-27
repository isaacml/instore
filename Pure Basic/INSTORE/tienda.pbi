IncludeFile "config.pbi"
IncludeFile "tienda.pbf"

Procedure.s shop_status()
  CloseWindow(EventWindow())
  Openpanel_tienda()
  NewList dats.s()
  ;Leemos el fichero de settings
  loadDomains(settings_file$, dats())
  ForEach dats()
    If CountString(dats(), "http:") >= 1
      ;Envíamos la IP de salida
      SetGadgetText(dir_ip, dats())
    EndIf
  Next
EndProcedure

Procedure.s send_ip()
  ip.s = GetGadgetText(dir_ip)
  Debug ip
EndProcedure
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 20
; Folding = -
; EnableXP