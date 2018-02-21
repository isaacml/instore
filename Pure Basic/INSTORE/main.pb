IncludeFile "panel_inicio.pbf"
IncludeFile "menu.pbf"
IncludeFile "config.pbi"
IncludeFile "config_shop.pbf"
IncludeFile  "../LIBS/libs.pb"
Define output.s

Procedure get_values(*cur_node, nodeName$, attribute$, *valueResults.String)
  ;If nodeName$ and attribute$ are matched then the value
  ;will be added to the string structure pointed to by *valueResults .
  Protected result$
 
  While *cur_node
    If XMLNodeType(*cur_node) = #PB_XML_Normal
 
      result$ = GetXMLNodeName(*cur_node)
      If result$ = nodeName$
        If ExamineXMLAttributes(*cur_node)
          While NextXMLAttribute(*cur_node)
            If XMLAttributeName(*cur_node) = attribute$
              If *valueResults <> #Null
                *valueResults\s + XMLAttributeValue(*cur_node) + Chr(13) ;value + carriage-return
              EndIf 
            EndIf
          Wend
        EndIf
      EndIf 
 
    EndIf 
 
    get_values(ChildXMLNode(*cur_node), nodeName$, attribute$, *valueResults)
    *cur_node = NextXMLNode(*cur_node)
  Wend 
EndProcedure

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
              ents.s = POST_PB(ConnectionID, server$, "/transf_orgs.cgi", "user=" + user + "&action=entidad")
              Debug ents
              If CatchXML(0, @ents, Len(ents))
                *MainNode = MainXMLNode(0)   
                *SubNode = ChildXMLNode(*MainNode)
              
                Debug XMLNodeType(*SubNode)
                Debug "Text: " + GetXMLNodeText(*SubNode)
                Debug "Name: " + GetXMLNodeName(*SubNode)
                Debug ""
                
                *SubNode = NextXMLNode(*SubNode)
                Debug XMLNodeType(*SubNode)
                Debug "Text: " + GetXMLNodeText(XMLNodeFromPath(*SubNode, "#cdata"))
                Debug "Name: " + GetXMLNodeName(*SubNode)
             EndIf
            EndIf
          Else
            info_login = TextGadget(#PB_Any, 220, 220, 180, 25, "Fallo de login", #PB_Text_Center)
            SetGadgetColor(info_login, #PB_Gadget_FrontColor,RGB(200, 1, 0))
          EndIf
      EndSelect
    Case #PB_Event_CloseWindow
        eventClose = #True
  EndSelect
Until eventClose = #True
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 75
; Folding = -
; EnableXP