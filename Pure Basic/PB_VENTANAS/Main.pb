IncludeFile "Ventana1.pbf"
IncludeFile "Ventana2.pbf"

Declare W1_btnClose_Click()
Declare W1_btnOpen_Click()

OpenAdmin()
DisableGadget(Cerrar, #True)

Repeat
  event = WaitWindowEvent()
  eventClose.i = #False
  Select event
    Case #PB_Event_Gadget
      evGadget = EventGadget()
      Select evGadget
        Case Cerrar
          W1_btnClose_Click()
        Case Abrir
          W1_btnOpen_Click()
      EndSelect
    Case #PB_Event_CloseWindow
      If GetActiveWindow() = Tienda
        DisableGadget(Abrir, #False)
        DisableGadget(Cerrar, #True)
        CloseWindow(GetActiveWindow())
      ElseIf GetActiveWindow() = Admin
        eventClose = #True
      EndIf
  EndSelect
  
Until eventClose = #True


Procedure W1_btnOpen_Click()
  DisableGadget(Cerrar, #False)
  DisableGadget(Abrir, #True)
  OpenTienda()
  
EndProcedure

Procedure W1_btnClose_Click()
  DisableGadget(Abrir, #False)
  DisableGadget(Cerrar, #True)
  CloseWindow(Tienda)  
EndProcedure
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 27
; Folding = -
; EnableXP