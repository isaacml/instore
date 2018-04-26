IncludeFile "panel_inicio.pbf"
IncludeFile "menu.pbf"
IncludeFile "config.pbi"
IncludeFile "config_shop.pbf"
IncludeFile "mensajes.pbf"
IncludeFile "dominios.pbf"
IncludeFile  "../LIBS/libs.pb"
Define output.s
Global user.s

InitNetwork()
Openpanel_login()
UseSQLiteDatabase()

ImportC ""
  time(*tloc = #Null)
EndImport

ConnectionID = OpenNetworkConnection(server$, port.l) 
DatabaseFile$ = "C:\Users\Isaac\Documents\Prueba Compilado PB\shop.db"
DirectoryMsg$ = "C:\Users\Isaac\Documents\Prueba Compilado PB\Messages"

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
              CloseWindow(EventWindow())
            Else
              Openconfig_shop()
              SetWindowData(GetActiveWindow(), 12)
              CloseWindow(EventWindow())
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
            SetGadgetText(info_login, "Fallo de login")
          EndIf
        Case Entidades
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(enviar_config, 1)
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
              SetGadgetText(err_dom, "")
              ForEach Valores()
                AddGadgetItem(Almacenes, 0, valores())
                SetGadgetItemData(Almacenes, 0, Val(MapKey(valores())))
              Next
          EndSelect
        Case Almacenes
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(enviar_config, 1)
              valor = GetGadgetItemData(Almacenes, GetGadgetState(Almacenes))
              pais$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "almacen=" + valor + "&action=pais")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), pais$, "</option><option", valores())
              ClearGadgetItems(Paises)
              ClearGadgetItems(Regiones)
              ClearGadgetItems(Provincias)
              ClearGadgetItems(Tiendas)
              SetGadgetText(err_dom, "")
              ForEach Valores()
                AddGadgetItem(Paises, 0, valores())
                SetGadgetItemData(Paises, 0, Val(MapKey(valores())))
              Next
          EndSelect
        Case Paises
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(enviar_config, 1)
              valor = GetGadgetItemData(Paises, GetGadgetState(Paises))
              reg$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "pais=" + valor + "&action=region")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), reg$, "</option><option", valores())
              ClearGadgetItems(Regiones)
              ClearGadgetItems(Provincias)
              ClearGadgetItems(Tiendas)
              SetGadgetText(err_dom, "")
              ForEach Valores()
                AddGadgetItem(Regiones, 0, valores())
                SetGadgetItemData(Regiones, 0, Val(MapKey(valores())))
              Next
          EndSelect
        Case Regiones
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(enviar_config, 1)
              valor = GetGadgetItemData(Regiones, GetGadgetState(Regiones))
              prov$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "region=" + valor + "&action=provincia")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), prov$, "</option><option", valores())
              ClearGadgetItems(Provincias)
              ClearGadgetItems(Tiendas)
              SetGadgetText(err_dom, "")
              ForEach Valores()
                AddGadgetItem(Provincias, 0, valores())
                SetGadgetItemData(Provincias, 0, Val(MapKey(valores())))
              Next
         EndSelect
        Case Provincias
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(enviar_config, 1)
              valor = GetGadgetItemData(Provincias, GetGadgetState(Provincias))
              shop$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "provincia=" + valor + "&action=tienda")
              Dim output.s(0)
              NewMap valores.s()
              obtainIdName(output(), shop$, "</option><option", valores())
              ClearGadgetItems(Tiendas)
              SetGadgetText(err_dom, "")
              ForEach Valores()
                AddGadgetItem(Tiendas, 0, valores())
                SetGadgetItemData(Tiendas, 0, Val(MapKey(valores())))
              Next
          EndSelect
        Case Tiendas
          Select EventType()
            Case #PB_EventType_Change
              DisableGadget(enviar_config, 0) ;Se habilita el boton de envio de configuracion
              SetGadgetText(err_dom, "")
              valor = GetGadgetItemData(Tiendas, GetGadgetState(Tiendas))
              POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "tienda=" + valor + "&action=cod_tienda")
          EndSelect
        Case enviar_config  ;Se envian los datos de configuracion de la tienda
          Select EventType()
            Case #PB_EventType_LeftClick
              res$ = POST_PB(ConnectionID, server$, "/acciones.cgi", "action=save_domain")
              what_win = GetWindowData(EventWindow()) ;Nos indica en que ventana nos encontramos (config_shop o config_dom)
              If what_win = 12 ;Ventana de configuración principal (config_shop)
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
                        CloseWindow(EventWindow())
                      EndIf
                      CloseDatabase(0)
                    Else
                      Debug "Can't open database!"
                    EndIf
                  Else
                    MessageRequester("Information","May not create the file!")
                  EndIf
                EndIf
              ElseIf what_win = 16 ;Ventana de configuración adiccional (dominios)
                ok$ = StringField(res$, 1, ";") 
                If ok$ = "OK"
                 count = 0 
                 NewList dat.s()
                 NewList actualizar.s()
                 dom$ = StringField(res$, 2, ";")
                 loadDomains(domain_file$, dat())
                 ;Recorremos el listado de dominios
                 ResetList(dat())
                 While NextElement(dat())
                   If dat() = dom$
                     count + 1
                   EndIf
                 Wend
                 ;Contador igual a 0: añadimos dominios extra
                 If count = 0
                   If OpenFile(0, domain_file$)
                     FileSeek(0, Lof(0))
                     WriteStringN(0, "extradomain = " + dom$)
                     CloseFile(0)
                     ClearGadgetItems(lista_dominios)
                     ;Mostramos los dominios actualizados
                     loadDomains(domain_file$, actualizar())
                     ForEach actualizar()
                       AddGadgetItem(lista_dominios, -1, actualizar())
                     Next
                   EndIf
                 Else
                   SetGadgetText(err_dom, "Ese dominio ya existe")  
                 EndIf  
                EndIf  
              EndIf
          EndSelect
       Case shop_status
       Case msg_normal
          Select EventType()
            Case #PB_EventType_LeftClick
              CloseWindow(EventWindow())
              Openpanel_mensajes()
              NewList msgfiles.s()
              obtainMsgFiles(DirectoryMsg$, msgfiles())
              ForEach msgfiles()
                AddGadgetItem(show_msg, 0, msgfiles())
              Next
          EndSelect
       Case doms
          Select EventType()
            Case #PB_EventType_LeftClick
              Dim output.s(0)
              NewMap valores.s()
              NewList dat2.s()
              CloseWindow(EventWindow())
              Opendominios()
              SetWindowData(GetActiveWindow(), 16)
              ents$ = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "user=" + user + "&action=entidad")
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
              ;Mostramos los dominios actuales
              loadDomains(domain_file$, dat2())
              ForEach dat2()
                AddGadgetItem(lista_dominios, -1, dat2())
              Next
         EndSelect 
       Case play_msg
         MP3_Free(0)
         Select EventType()
           Case #PB_EventType_LeftClick
             FullMsg$ = DirectoryMsg$ + "\" + GetGadgetText(show_msg)
             If MP3_Load(0, FullMsg$)
               MP3_Play(0)
             EndIf
         EndSelect
     EndSelect
    Case #PB_Event_Menu
     Select EventMenu()
       Case #back_msg, #back_dom
         CloseWindow(EventWindow())
         Openmenu()
       Case #logout_menu, #logout_msg, #logout_dom
         CloseWindow(EventWindow())
         Openpanel_login()
    EndSelect
    Case #PB_Event_CloseWindow
        eventClose = #True
    EndSelect
Until eventClose = #True
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 220
; FirstLine = 187
; EnableXP