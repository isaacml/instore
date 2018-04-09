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
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), ents$, "</option><option", valores())
              ClearGadgetItems(Entidades)
              ClearGadgetItems(Almacenes)
              ClearGadgetItems(Paises)
              ClearGadgetItems(Regiones)
              ClearGadgetItems(Provincias)
              ClearGadgetItems(Tiendas)
              ForEach Valores()
                AddGadgetItem(Entidades, 0, Valores())
                SetGadgetItemData(Entidades, 0, Val(MapKey(Valores())))
              Next
            EndIf
          Else
            info_login = TextGadget(#PB_Any, 220, 220, 180, 25, "Fallo de login", #PB_Text_Center)
            SetGadgetColor(info_login, #PB_Gadget_FrontColor,RGB(200, 1, 0))
          EndIf
        Case Entidades
          Select EventType()
            Case #PB_EventType_Change
              valor = GetGadgetItemData(Entidades, GetGadgetState(Entidades))
              alms$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "entidad=" + valor + "&action=almacen")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), alms$, "</option><option", valores())
              ClearGadgetItems(Almacenes)
              ClearGadgetItems(Paises)
              ClearGadgetItems(Regiones)
              ClearGadgetItems(Provincias)
              ClearGadgetItems(Tiendas)
              ForEach Valores()
                AddGadgetItem(Almacenes, 0, Valores())
                SetGadgetItemData(Almacenes, 0, Val(MapKey(Valores())))
              Next
          EndSelect
        Case Almacenes
          Select EventType()
            Case #PB_EventType_Change
              valor = GetGadgetItemData(Almacenes, GetGadgetState(Almacenes))
              pais$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "almacen=" + valor + "&action=pais")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), pais$, "</option><option", valores())
              ClearGadgetItems(Paises)
              ClearGadgetItems(Regiones)
              ClearGadgetItems(Provincias)
              ClearGadgetItems(Tiendas)
              ForEach Valores()
                AddGadgetItem(Paises, 0, Valores())
                SetGadgetItemData(Paises, 0, Val(MapKey(Valores())))
              Next
          EndSelect
        Case Paises
          Select EventType()
            Case #PB_EventType_Change
              valor = GetGadgetItemData(Paises, GetGadgetState(Paises))
              reg$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "pais=" + valor + "&action=region")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), reg$, "</option><option", valores())
              ClearGadgetItems(Regiones)
              ClearGadgetItems(Provincias)
              ClearGadgetItems(Tiendas)
              ForEach Valores()
                AddGadgetItem(Regiones, 0, Valores())
                SetGadgetItemData(Regiones, 0, Val(MapKey(Valores())))
              Next
          EndSelect
        Case Regiones
          Select EventType()
            Case #PB_EventType_Change
              valor = GetGadgetItemData(Regiones, GetGadgetState(Regiones))
              prov$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "region=" + valor + "&action=provincia")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), prov$, "</option><option", valores())
              ClearGadgetItems(Provincias)
              ClearGadgetItems(Tiendas)
              ForEach Valores()
                AddGadgetItem(Provincias, 0, Valores())
                SetGadgetItemData(Provincias, 0, Val(MapKey(Valores())))
              Next
         EndSelect
        Case Provincias
          Select EventType()
            Case #PB_EventType_Change
              valor = GetGadgetItemData(Provincias, GetGadgetState(Provincias))
              prov$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "provincia=" + valor + "&action=tienda")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), prov$, "</option><option", valores())
              ClearGadgetItems(Tiendas)
              ForEach Valores()
                AddGadgetItem(Tiendas, 0, Valores())
                SetGadgetItemData(Tiendas, 0, Val(MapKey(Valores())))
              Next
       EndSelect
      EndSelect
    Case #PB_Event_CloseWindow
        eventClose = #True
  EndSelect
Until eventClose = #True
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 41
; FirstLine = 1
; EnableXP