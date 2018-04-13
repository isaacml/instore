IncludeFile "panel_inicio.pbf"
IncludeFile "menu.pbf"
IncludeFile "config.pbi"
IncludeFile "config_shop.pbf"
IncludeFile "mensajes.pbf"
IncludeFile  "../LIBS/libs.pb"
Define output.s

Global user.s

Openpanel_login()
InitNetwork()
InitSound()
UseSQLiteDatabase()

ImportC ""
  time(*tloc = #Null)
EndImport

ConnectionID = OpenNetworkConnection(server$, port.l) 
DatabaseFile$ = "C:\Users\Isaac\Documents\Prueba Compilado PB\shop.db"
DirectoryMsg$ = "C:\Users\Isaac\Documents\PB\INSTORE\Messages"

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
                AddGadgetItem(Entidades, 0, valores())
                SetGadgetItemData(Entidades, 0, Val(MapKey(valores())))
              Next
            EndIf
          Else
            info_login = TextGadget(#PB_Any, 220, 220, 180, 25, "Fallo de login", #PB_Text_Center)
            SetGadgetColor(info_login, #PB_Gadget_FrontColor,RGB(200, 1, 0))
          EndIf
        Case Entidades
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(Enviar, 1)
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
                AddGadgetItem(Almacenes, 0, valores())
                SetGadgetItemData(Almacenes, 0, Val(MapKey(valores())))
              Next
          EndSelect
        Case Almacenes
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(Enviar, 1)
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
                AddGadgetItem(Paises, 0, valores())
                SetGadgetItemData(Paises, 0, Val(MapKey(valores())))
              Next
          EndSelect
        Case Paises
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(Enviar, 1)
              valor = GetGadgetItemData(Paises, GetGadgetState(Paises))
              reg$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "pais=" + valor + "&action=region")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), reg$, "</option><option", valores())
              ClearGadgetItems(Regiones)
              ClearGadgetItems(Provincias)
              ClearGadgetItems(Tiendas)
              ForEach Valores()
                AddGadgetItem(Regiones, 0, valores())
                SetGadgetItemData(Regiones, 0, Val(MapKey(valores())))
              Next
          EndSelect
        Case Regiones
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(Enviar, 1)
              valor = GetGadgetItemData(Regiones, GetGadgetState(Regiones))
              prov$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "region=" + valor + "&action=provincia")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), prov$, "</option><option", valores())
              ClearGadgetItems(Provincias)
              ClearGadgetItems(Tiendas)
              ForEach Valores()
                AddGadgetItem(Provincias, 0, valores())
                SetGadgetItemData(Provincias, 0, Val(MapKey(valores())))
              Next
         EndSelect
        Case Provincias
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(Enviar, 1)
              valor = GetGadgetItemData(Provincias, GetGadgetState(Provincias))
              shop$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "provincia=" + valor + "&action=tienda")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), shop$, "</option><option", valores())
              ClearGadgetItems(Tiendas)
              ForEach Valores()
                AddGadgetItem(Tiendas, 0, valores())
                SetGadgetItemData(Tiendas, 0, Val(MapKey(valores())))
              Next
          EndSelect
        Case Tiendas
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(Enviar, 0) ;Se habilita el boton de envio de configuracion
              valor = GetGadgetItemData(Tiendas, GetGadgetState(Tiendas))
              POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "tienda=" + valor + "&action=cod_tienda")
          EndSelect
        Case Enviar ;Se envian los datos de configuracion de la tienda
          Select EventType()
            Case #PB_EventType_LeftClick
              res$ = POST_PB(ConnectionID, server$, "/acciones.cgi", "action=save_domain")
              ok$ = StringField(res$, 1, ";") 
              If ok$ = "OK"
                dom$ = StringField(res$, 2, ";")
                ;Creamos el fichero de configuracion: guardamos el dominio de la tienda
                If CreateFile(0, domain_file$)
                  WriteString(0, "shopdomain = " + dom$ + Chr(10))
                  CloseFile(0)
                  ;Se hace un guardado del dominio en base de datos
                  If OpenDatabase(0, DatabaseFile$, "", "")
                    err = DatabaseUpdate(0, "INSERT INTO tienda (dominio, last_connect) VALUES ('"+ dom$ +"',"+ time() +")")
                    If Not err = 0 
                      Openmenu()
                      CloseWindow(config_shop)
                    EndIf
                    CloseDatabase(0)
                  Else
                    Debug "Can't open database!"
                  EndIf
                Else
                  MessageRequester("Information","May not create the file!")
                EndIf
              EndIf
          EndSelect
       Case msg_normal
          Select EventType()
            Case #PB_EventType_LeftClick
              CloseWindow(menu)
              Openpanel_mensajes()
              NewList msgfiles.s()
              obtainMsgFiles(DirectoryMsg$, msgfiles())
              ForEach msgfiles()
                AddGadgetItem(show_msg, 0, msgfiles())
              Next
         EndSelect   
       Case logout
          Select EventType()
            Case #PB_EventType_LeftClick
              CloseWindow(EventWindow())
              Openpanel_login()
          EndSelect
      Case play_msg
        Select EventType()
          Case #PB_EventType_LeftClick
            If DirectoryMsg$
              FullMsg$ = DirectoryMsg$ + "\" + GetGadgetText(show_msg)
              Debug FullMsg$
              If LoadSound(0, FullMsg$)
                PlaySound(0)
              EndIf
            EndIf
        EndSelect
      EndSelect
    Case #PB_Event_CloseWindow
        eventClose = #True
  EndSelect
Until eventClose = #True
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 208
; FirstLine = 162
; EnableXP