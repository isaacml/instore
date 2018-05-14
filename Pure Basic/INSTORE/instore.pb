IncludeFile "config_instore.pbi"
IncludeFile "mensajes.pbf"
IncludeFile "dominios.pbf"
IncludeFile "tienda.pbf"

Procedure Instore(id_con, server$, user.s, domain_file$)
  DatabaseFile$ = "C:\Users\0oIsa\Documents\PRUEBAS_INSTORE\shop.db"
  DirectoryMsg$ = "C:\Users\0oIsa\Documents\PRUEBAS_INSTORE\Messages"
  
  panel_main = OpenWindow(#PB_Any, 0, 0, 800, 600, "Panel Cliente", #PB_Window_SystemMenu | #PB_Window_ScreenCentered | #PB_Window_WindowCentered)
  
  If CreateMenu(0, WindowID(panel_main))
    MenuTitle("Menu")
    MenuItem(#prog, "Programar Música")
    MenuItem(#msg, "Mensajes")
    MenuItem(#dom, "Dominios")
    MenuItem(#shop, "Tienda")
  EndIf
  
  Repeat
    Event = WaitWindowEvent()
    Select Event
      Case #PB_Event_Gadget
        Select EventGadget()
          Case play_msg ;Boton que reproduce un mensaje
            MP3_Free(1)
            Select EventType()
              Case #PB_EventType_LeftClick
                FullMsg$ = DirectoryMsg$ + "\" + GetGadgetText(show_msg)
                If MP3_Load(1, FullMsg$)
                  MP3_Play(1)
                EndIf
            EndSelect
          Case enviar_dom ;Ventana de configuración adiccional (dominios)
            res$ = POST_PB_STORE(id_con, server$, "/acciones.cgi", "action=save_domain")
            ok$ = StringField(res$, 1, ";") 
            If ok$ = "OK"
              count = 0 
              NewList dat.s()
              NewList actualizar.s()
              dom$ = StringField(res$, 2, ";")
              load_Domains(domain_file$, dat())
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
                  load_Domains(domain_file$, actualizar())
                  ForEach actualizar()
                    AddGadgetItem(lista_dominios, -1, actualizar())
                  Next
                EndIf
              Else
                SetGadgetText(err_dom, "Ese dominio ya existe")  
              EndIf  
            EndIf  
          Case Entidades
            Select EventType()
              Case #PB_EventType_Change
                DisableGadget(enviar_dom, 1)
                valor = GetGadgetItemData(Entidades, GetGadgetState(Entidades))
                alms$ = POST_PB_STORE(id_con, server$, "/transf_orgs.cgi", "entidad=" + valor + "&action=almacen")
                Dim output.s(0)
                NewMap valores.s()
                obtain_Id_Name(output(), alms$, "</option><option", valores())
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
                DisableGadget(enviar_dom, 1)
                valor = GetGadgetItemData(Almacenes, GetGadgetState(Almacenes))
                pais$ = POST_PB_STORE(id_con, server$, "/transf_orgs.cgi", "almacen=" + valor + "&action=pais")
                Dim output.s(0)
                NewMap valores.s()
                obtain_Id_Name(output(), pais$, "</option><option", valores())
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
                DisableGadget(enviar_dom, 1)
                valor = GetGadgetItemData(Paises, GetGadgetState(Paises))
                reg$ = POST_PB_STORE(id_con, server$, "/transf_orgs.cgi", "pais=" + valor + "&action=region")
                Dim output.s(0)
                NewMap valores.s()
                obtain_Id_Name(output(), reg$, "</option><option", valores())
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
                DisableGadget(enviar_dom, 1)
                valor = GetGadgetItemData(Regiones, GetGadgetState(Regiones))
                prov$ = POST_PB_STORE(id_con, server$, "/transf_orgs.cgi", "region=" + valor + "&action=provincia")
                Dim output.s(0)
                NewMap valores.s()
                obtain_Id_Name(output(), prov$, "</option><option", valores())
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
                DisableGadget(enviar_dom, 1)
                valor = GetGadgetItemData(Provincias, GetGadgetState(Provincias))
                shop$ = POST_PB_STORE(id_con, server$, "/transf_orgs.cgi", "provincia=" + valor + "&action=tienda")
                Dim output.s(0)
                NewMap valores.s()
                obtain_Id_Name(output(), shop$, "</option><option", valores())
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
                DisableGadget(enviar_dom, 0) ;Se habilita el boton de envio de configuracion
                SetGadgetText(err_dom, "")
                valor = GetGadgetItemData(Tiendas, GetGadgetState(Tiendas))
                POST_PB_STORE(id_con, server$, "/transf_orgs.cgi", "tienda=" + valor + "&action=cod_tienda")
            EndSelect
        EndSelect
      Case #PB_Event_Menu
        Select EventMenu()
          Case #prog, #prog_dom: Debug "voy pa programar"
          Case #msg, #msg_dom
            CloseWindow(EventWindow())
            Openpanel_mensajes()
            NewList msgfiles.s()
            obtainMsgFiles(DirectoryMsg$, msgfiles())
            ForEach msgfiles()
              AddGadgetItem(show_msg, 0, msgfiles())
            Next
          Case #dom, #dom_dom;(Boton de Menu: Dominios)
            CloseWindow(EventWindow())
            Opendominios()
            Dim output.s(0)
            NewMap valores.s()
            NewList dat2.s()
            ents$ = POST_PB_STORE(id_con, server$, "/transf_orgs.cgi", "user=" + user + "&action=entidad")
            obtain_Id_Name(output(), ents$, "</option><option", valores())
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
            load_Domains(domain_file$, dat2())
            ForEach dat2()
              AddGadgetItem(lista_dominios, -1, dat2())
            Next
          Case #shop, #shop_dom: Debug "Voy pa tienda"
        EndSelect
    EndSelect
  Until Event = #PB_Event_CloseWindow
EndProcedure
; IDE Options = PureBasic 5.61 (Windows - x64)
; CursorPosition = 34
; FirstLine = 24
; Folding = -
; EnableXP