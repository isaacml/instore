IncludeFile "panel_inicio.pbf"
IncludeFile "menu.pbf"
IncludeFile "config.pbi"
IncludeFile "config_shop.pbf"
IncludeFile  "../LIBS/libs.pb"
Define output.s

Global user.s
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
          st_login$ = POST_PB(ConnectionID, server$, "/auth.cgi", parameters$)
          If st_login$ = "OK"
            user = username$
            If ReadFile(0,domain_file$)
              Openmenu()
              CloseWindow(panel_login)
            Else
              Openconfig_shop()
              CloseWindow(panel_login)
              ents$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "user=" + user + "&action=entidad")
              Debug ents$
              Dim output.s(0) ;this will be resized later
              explodeStringArray(output(), ents$, "</option><option")
              For i = 1 To ArraySize(output())-1
                value.s = output(i)
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
                AddGadgetItem(Entidades, 0, name)
                SetGadgetItemData(Entidades, 0, Val(endid))
                Debug endid + " - " + name
              Next
            EndIf
          Else
            info_login = TextGadget(#PB_Any, 220, 220, 180, 25, "Fallo de login", #PB_Text_Center)
            SetGadgetColor(info_login, #PB_Gadget_FrontColor,RGB(200, 1, 0))
          EndIf
        Case Entidades
          valor = GetGadgetItemData(Entidades, 0)
          alms$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "entidad=" + valor + "&action=almacen")
          Debug alms$
          Dim output.s(0) ;this will be resized later
          explodeStringArray(output(), alms$, "</option><option")
          For i = 1 To ArraySize(output())-1
            value.s = output(i)
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
            AddGadgetItem(Entidades, 0, name)
            SetGadgetItemData(Entidades, 0, Val(endid))
            Debug endid + " - " + name
          Next
        Case Almacenes
          valor = GetGadgetItemData(Almacenes, 0)
      EndSelect
    Case #PB_Event_CloseWindow
        eventClose = #True
  EndSelect
Until eventClose = #True
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 95
; FirstLine = 46
; EnableXP