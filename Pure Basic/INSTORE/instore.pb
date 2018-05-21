IncludeFile "config_instore.pbi"
IncludeFile "mensajes.pbf"
IncludeFile "dominios.pbf"
IncludeFile "tienda.pbf"

Procedure Instore(id_con, server$, user.s, domain_file$, settings_file$)
  DatabaseFile$ = "C:\Users\Isaac\Documents\Prueba Compilado PB\shop.db"
  DirectoryMsg$ = "C:\Users\Isaac\Documents\Prueba Compilado PB\Messages"
  
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
          Case ip_send
            ip.s = GetGadgetText(dir_ip)
            Debug ip
          Case play_msg ;Boton que reproduce un mensaje
            MP3_Free(1)
            Select EventType()
              Case #PB_EventType_LeftClick
                FullMsg$ = DirectoryMsg$ + "\" + GetGadgetText(show_msg)
                If MP3_Load(1, FullMsg$)
                  MP3_Play(1)
                EndIf
            EndSelect
          Case enviar_dom ;Boton de envio (dominios adiccionales)
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
          Case send_horario ;Boton de envio (horario de la tienda)
            Select EventType()
              Case #PB_EventType_LeftClick
                hora1.s = GetGadgetText(h1)
                min1.s = GetGadgetText(m1)
                hora2.s = GetGadgetText(h2)
                min2.s = GetGadgetText(m2)
                If hora1 = "" Or min1 = "" Or hora2 = "" Or min2 = ""
                  SetGadgetText(info_horario, "Hay campos vacíos")
                Else
                  If OpenDatabase(1, DatabaseFile$, "", "")
                    If DatabaseQuery(1, "SELECT hora_inicial, hora_final FROM horario;")
                      exist = NextDatabaseRow(1) ;Se comprueba si hay un horario o NO
                      If exist = 0               ;No hay horario previo, por tanto, insertamos
                        sql.s = "INSERT INTO horario (hora_inicial, hora_final) VALUES ('"+ hora1 + ":" + min1 +"','"+ hora2 + ":" + min2 +"')"
                        DatabaseUpdate(1, sql)
                      Else ;Hay un horario existente, por tanto, borramos el viejo y insertamos el nuevo
                        delete.s = "DELETE FROM horario"
                        DatabaseUpdate(1, delete)
                        insert.s = "INSERT INTO horario (hora_inicial, hora_final) VALUES ('"+ hora1 + ":" + min1 +"','"+ hora2 + ":" + min2 +"')"
                        DatabaseUpdate(1, insert)
                      EndIf
                      SetGadgetColor(info_horario, #PB_Gadget_FrontColor, $15c206)
                      SetGadgetText(info_horario, "Nuevo horario añadido")
                      FinishDatabaseQuery(1)
                    EndIf
                    CloseDatabase(1)
                  Else
                    MessageRequester("Information","Error to Open DB")
                  EndIf
                EndIf
            EndSelect
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
          Case #prog, #prog_dom, #prog_shop: Debug "voy pa programar"
          Case #msg, #msg_dom, #msg_shop
            CloseWindow(EventWindow())
            Openpanel_mensajes()
            NewList msgfiles.s()
            obtainMsgFiles(DirectoryMsg$, msgfiles())
            ForEach msgfiles()
              AddGadgetItem(show_msg, 0, msgfiles())
            Next
          Case #dom, #dom_dom, #dom_shop ;(Boton de Menu: Dominios)
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
          Case #shop, #shop_dom, #shop_shop ;(Boton de Menu: Tienda)
            CloseWindow(EventWindow())
            Openpanel_tienda()
            NewList dats.s()
            ;Leemos el fichero de settings
            load_Domains(settings_file$, dats())
            ForEach dats()
              If CountString(dats(), "http:") >= 1
                ;Envíamos la IP de salida
                SetGadgetText(dir_ip, dats())
              EndIf
            Next
            formar_horas(h1)
            formar_minutos(m1)
            formar_horas(h2)
            formar_minutos(m2)
            If OpenDatabase(0, DatabaseFile$, "", "")
              If DatabaseQuery(0, "SELECT hora_inicial, hora_final FROM horario;")
                exist = NextDatabaseRow(0) ;Se comprueba si hay un horario o NO
                If exist = 1 
                  SetGadgetText(h1, StringField(GetDatabaseString(0, 0), 1, ":"))
                  SetGadgetText(m1, StringField(GetDatabaseString(0, 0), 2, ":"))
                  SetGadgetText(h2, StringField(GetDatabaseString(0, 1), 1, ":"))
                  SetGadgetText(m2, StringField(GetDatabaseString(0, 1), 2, ":"))
                EndIf
                FinishDatabaseQuery(0)
              EndIf
              CloseDatabase(0)
            Else
              MessageRequester("Information","Error to Open DB")
            EndIf
        EndSelect
    EndSelect
    Delay(1)
  Until Event = #PB_Event_CloseWindow
EndProcedure
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 267
; FirstLine = 213
; Folding = -
; EnableXP