IncludeFile "panel_inicio.pbf"
IncludeFile "config.pbi"
IncludeFile  "../LIBS/libs.pb"

Openpanel_login()
InitNetwork()
ConnectionID = OpenNetworkConnection(server$, port.l) 

Repeat
  event = WaitWindowEvent()
  eventClose.i = #False
  Select event
    Case #PB_Event_Gadget
      evGadget = EventGadget()
      Select evGadget
        Case send
          username$ = GetGadgetText(username)
          password$ = GetGadgetText(password)
          parameters$ = "user=" + username$ + "&pass=" + password$
          Debug POST_PB(ConnectionID, server$, cgi$, parameters$)
      EndSelect
    Case #PB_Event_CloseWindow
        eventClose = #True
  EndSelect
  
Until eventClose = #True
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 19
; EnableXP