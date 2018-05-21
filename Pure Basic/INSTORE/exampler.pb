;- Constants
  #DateFormat = "%dd/%mm/%yyyy"
  
  ;- Enumerations
  Enumeration 
    #WindowClassButton = 1
    #WindowClassDate
    #WindowClassTrack
  EndEnumeration
  
  ; The menu commands will be the same on all the windows.
  Enumeration
    #MenuNewButton
    #MenuNewDate
    #MenuNewTrack
    #MenuClose
  EndEnumeration
  
  ;- Structures
  Structure BASEWINDOW
    WindowClass.i 
    Menu.i  
  EndStructure
  
  Structure BUTTONWINDOW Extends BASEWINDOW
    ListView.i
    AddButton.i
    RemoveButton.i  
  EndStructure
  
  Structure DATEWINDOW Extends BASEWINDOW
    Calendar.i
    AddButton.i
    RemoveButton.i
    ListView.i
  EndStructure
  
  Structure TRACKWINDOW Extends BASEWINDOW
    TrackBar1.i
    TrackBar2.i
    Label.i
    Difference.i
  EndStructure
  
  ;- Variables
  ; This map will track all the active windows as before, however as the structure for each window class
  ; differs we will be storing pointers to structures in the map not the gadget references directly.
  NewMap ActiveWindows.i()
  
  ; These values will be used to give new windows a unique label.
  Define.i Buttons, Dates, Tracks
  
  ; Event variables.
  ; Notice the type of *EventGadgets.
  Define.i Event, EventWindow, EventGadget, EventType, EventMenu, EventQuit
  Define.s EventWindowKey
  Define *EventGadgets.BASEWINDOW  
  
  ;- Button Window
  Procedure.i CreateButtonWindow()
    ; Creates a new Button window and adds it to the tracking map, 
    ; allocates memory for gadget references, creates the gadgets 
    ; and stores these references in the memory structure.
    Shared Buttons, ActiveWindows()
    Protected *ThisData.BUTTONWINDOW
    Protected.i ThisWindow
    Protected.s ThisKey, ThisTitle
    
    ; Set the window caption.
    Buttons + 1
    ThisTitle = "Button Window " + StrU(Buttons)
    
    ; Open the window.
    ThisWindow = OpenWindow(#PB_Any, 30, 30, 225, 300, ThisTitle, #PB_Window_SystemMenu | #PB_Window_SizeGadget | #PB_Window_MinimizeGadget | #PB_Window_TitleBar)
    
    ; Check that the OpenWindow command worked.
    If ThisWindow
      ; Set minimum window dimensions.
      WindowBounds(ThisWindow, 220, 100, #PB_Ignore, #PB_Ignore)
      
      ; Convert the window reference to a string to use as the map key value.
      ThisKey = StrU(ThisWindow)
      
      ; Allocate memory to store the gadget references in.
      *ThisData = AllocateMemory(SizeOf(BUTTONWINDOW))
    EndIf
    
    ; Check that the memory allocation worked.
    If *ThisData
      ; Store the window reference and memory pointer values in the map.
      ActiveWindows(ThisKey) = *ThisData
      
      ; Set the window class.
      *ThisData\WindowClass = #WindowClassButton
      
      ; Create the menu bar.
      *ThisData\Menu = CreateMenu(#PB_Any, WindowID(ThisWindow))
      
      ; If the menu creation worked, create the menu items.
      If *ThisData\Menu
        MenuTitle("Window")
        MenuItem(#MenuNewButton, "New Button Window")
        MenuItem(#MenuNewDate, "New Date Window")
        MenuItem(#MenuNewTrack, "New Track Window")
        MenuItem(#MenuClose, "Close Window")
      EndIf
      
      ; Create the window gadgets.
      *ThisData\ListView = ListViewGadget(#PB_Any, 10, 10, 200, 225)
      *ThisData\AddButton = ButtonGadget(#PB_Any, 15, 245, 90, 30, "Add")
      *ThisData\RemoveButton = ButtonGadget(#PB_Any, 115, 245, 90, 30, "Remove")
    Else
      ; Memory allocation failed.
      CloseWindow(ThisWindow)
    EndIf
    
    ; Set the return value.
    If ThisWindow > 0 And *ThisData > 0
      ; Return the reference to the new window.
      ProcedureReturn ThisWindow
    Else
      ; Return 0
      ProcedureReturn 0
    EndIf
  EndProcedure
  
  Procedure.i DestroyButtonWindow(Window.i)
    ; Remove Window from the ActiveWindows map, release the allocated memory,
    ; close the window and set the quit flag, if appropriate.
    Shared EventQuit, ActiveWindows()
    Protected *ThisData.BUTTONWINDOW
    Protected.s ThisKey
    
    ; Convert the integer Window to a string.
    ThisKey = StrU(Window)
    
    ; Obtain the reference structure pointer.
    *ThisData = ActiveWindows(ThisKey)
    
    ; Check that a valid pointer was obtained, if not stop.
    If *ThisData = 0
      ProcedureReturn #False
    EndIf
    
    ; Check that it is the correct window type, if not stop.
    If *ThisData\WindowClass <> #WindowClassButton
      ProcedureReturn #False
    EndIf
    
    ; Release the memory allocation.
    FreeMemory(*ThisData)
    
    ; Delete the map entry.
    DeleteMapElement(ActiveWindows(), ThisKey)
    
    ; Close the window.
    CloseWindow(Window)
    
    ; Check if there are still open windows.
    If MapSize(ActiveWindows()) = 0
      EventQuit = #True
    EndIf
    
    ; Set the successful return value.
    ProcedureReturn #True
  EndProcedure
  
  Procedure.i ResizeButtonWindow(Window.i)
    ; Resize the child gadgets on Window.
    Shared ActiveWindows()
    Protected *ThisData.BUTTONWINDOW
    Protected.i X, Y, W, H
    Protected.s ThisKey
    
    ; Obtain the reference structure pointer.
    ThisKey = StrU(Window)
    *ThisData = ActiveWindows(ThisKey)
    
    ; Check that a valid pointer was obtained, if not stop.
    If *ThisData = 0
      ProcedureReturn #False
    EndIf
    
    ; Check that it is the correct window type, if not stop.
    If *ThisData\WindowClass <> #WindowClassButton
      ProcedureReturn #False
    EndIf
    
    ; Resize list view.
    W = WindowWidth(Window) - 25
    H = WindowHeight(Window) - 85
    ResizeGadget(*ThisData\ListView, #PB_Ignore, #PB_Ignore, W, H)
    
    ; Recenter buttons.
    X = WindowWidth(Window)/2 - 95
    Y = WindowHeight(Window) - 65
    ResizeGadget(*ThisData\AddButton, X, Y, #PB_Ignore, #PB_Ignore)
    
    X = WindowWidth(Window)/2 + 5
    ResizeGadget(*ThisData\RemoveButton, X, Y, #PB_Ignore, #PB_Ignore)
    
    ProcedureReturn #True
  EndProcedure
  
  Procedure.i EventsButtonWindow(Window, Gadget, Type)
    ; Handle events for a button window.
    Shared Buttons, ActiveWindows()
    Protected *ThisData.BUTTONWINDOW
    Protected.i NewValue, Index
    Protected.s ThisKey
    
    ; Convert the integer Window to a string.
    ThisKey = StrU(Window)
    
    ; Obtain the reference structure pointer.
    *ThisData = ActiveWindows(ThisKey)
    
    ; Check that a valid pointer was obtained, if not stop.
    If *ThisData = 0
      ProcedureReturn #False
    EndIf
    
    ; Check that it is the correct window type, if not stop.
    If *ThisData\WindowClass <> #WindowClassButton
      ProcedureReturn #False
    EndIf
    
    Select Gadget
      Case *ThisData\AddButton
        NewValue = Random(2147483647)
        AddGadgetItem(*ThisData\ListView, -1, StrU(NewValue))
        
      Case *ThisData\RemoveButton
        Index = GetGadgetState(*ThisData\ListView)
        If Index >= 0 And Index <= CountGadgetItems(*ThisData\ListView)
          RemoveGadgetItem(*ThisData\ListView, Index)
        EndIf
        
      Case *ThisData\ListView
        ; Do nothing.
    EndSelect
  EndProcedure
  
  ;- Date Window
  Procedure.i CreateDateWindow()
    ; Creates a new Date window and adds it to the tracking map, 
    ; allocates memory for gadget references, creates the gadgets 
    ; and stores these references in the memory Structure.
    Shared Dates, ActiveWindows()
    Protected *ThisData.DATEWINDOW
    Protected.i ThisWindow
    Protected.s ThisKey, ThisTitle
    
    Dates + 1
    ThisTitle = "Date Window " + StrU(Dates)
    ThisWindow = OpenWindow(#PB_Any, 30, 30, 310, 420, ThisTitle , #PB_Window_SystemMenu | #PB_Window_SizeGadget | #PB_Window_MinimizeGadget | #PB_Window_TitleBar)
    
    ; Check that the OpenWindow command worked.
    If ThisWindow
      ; Set minimum window dimensions.
      WindowBounds(ThisWindow, 310, 245, #PB_Ignore, #PB_Ignore)
      
      ; Convert the window reference to a string to use as the map key value.
      ThisKey = StrU(ThisWindow)
      
      ; Allocate memory to store the gadget references in.
      *ThisData = AllocateMemory(SizeOf(DATEWINDOW))
    EndIf
    
    ; Check that the memory allocation worked.
    If *ThisData
      ; Store the window reference and memory pointer values in the map.
      ActiveWindows(ThisKey) = *ThisData
      
      ; Set the window class.
      *ThisData\WindowClass = #WindowClassDate
      
      ; Create the menu bar.
      *ThisData\Menu = CreateMenu(#PB_Any, WindowID(ThisWindow))
      
      ; If the menu creation worked, create the menu items.
      If *ThisData\Menu
        MenuTitle("Window")
        MenuItem(#MenuNewButton, "New Button Window")
        MenuItem(#MenuNewDate, "New Date Window")
        MenuItem(#MenuNewTrack, "New Track Window")
        MenuItem(#MenuClose, "Close Window")
      EndIf
      
      ; Create the window gadgets.
      *ThisData\Calendar = CalendarGadget(#PB_Any, 10, 10, 182, 162)
      *ThisData\AddButton = ButtonGadget(#PB_Any, 210, 10, 90, 30, "Add")
      *ThisData\RemoveButton = ButtonGadget(#PB_Any, 210, 45, 90, 30, "Remove")
      *ThisData\ListView = ListViewGadget(#PB_Any, 10, 187, 290, 200)
    Else
      ; Memory allocation failed.
      CloseWindow(ThisWindow)
    EndIf
    
    ; Set the return value.
    If ThisWindow > 0 And *ThisData > 0 
      ; Return the reference to the new window.
      ProcedureReturn ThisWindow
    Else
      ; Return 0
      ProcedureReturn 0
    EndIf
    
  EndProcedure
  
  Procedure.i DestroyDateWindow(Window.i)
    ; Remove Window from the ActiveWindows map, release the allocated memory,
    ; close the window and set the quit flag, if appropriate.
    Shared EventQuit, ActiveWindows()
    Protected *ThisData.DATEWINDOW
    Protected.s ThisKey
    
    ; Convert the integer Window to a string.
    ThisKey = StrU(Window)
    
    ; Obtain the reference structure pointer.
    *ThisData = ActiveWindows(ThisKey)
    
    ; Check that a valid pointer was obtained.
    If *ThisData = 0
      ProcedureReturn #False
    EndIf
    
    ; Check that it is the correct window type, if not stop as this procedure can't destroy this window.
    If *ThisData\WindowClass <> #WindowClassDate
      ProcedureReturn #False
    EndIf
    
    ; Release the memory allocation.
    FreeMemory(*ThisData)
    
    ; Delete the map entry.
    DeleteMapElement(ActiveWindows(), ThisKey)
    
    ; Close the window.
    CloseWindow(Window)
    
    ; Check if there are still open windows.
    If MapSize(ActiveWindows()) = 0
      EventQuit = #True
    EndIf
    
    ; Set the successful return value.
    ProcedureReturn #True
  EndProcedure
  
  Procedure.i ResizeDateWindow(Window.i)
    ; Resize the child gadgets on Window.
    Shared ActiveWindows()
    Protected *ThisData.DATEWINDOW
    Protected.i X, Y, W, H
    Protected.s ThisKey
    
    ; Obtain the reference structure pointer.
    ThisKey = StrU(Window)
    *ThisData = ActiveWindows(ThisKey)
    
    ; Check that a valid pointer was obtained, if not stop.
    If *ThisData = 0
      ProcedureReturn #False
    EndIf
    
    ; Check that it is the correct window type, if not stop.
    If *ThisData\WindowClass <> #WindowClassDate
      ProcedureReturn #False
    EndIf
    
    ; Resize list view.
    W = WindowWidth(Window) - 20
    H = WindowHeight(Window) - 220
    ResizeGadget(*ThisData\ListView, #PB_Ignore, #PB_Ignore, W, H)
    
    ProcedureReturn #True
  EndProcedure
  
  Procedure.i EventsDateWindow(Window, Gadget, Type)
    ; Handle events for a Date window.
    Shared Buttons, ActiveWindows()
    Protected *ThisData.DATEWINDOW
    Protected.i NewValue, Index
    Protected.s ThisKey
    
    ; Convert the integer Window to a string.
    ThisKey = StrU(Window)
    
    ; Obtain the reference structure pointer.
    *ThisData = ActiveWindows(ThisKey)
    
    ; Check that a valid pointer was obtained, if not stop.
    If *ThisData = 0
      ProcedureReturn #False
    EndIf
    
    ; Check that it is the correct window type, if not stop.
    If *ThisData\WindowClass <> #WindowClassDate
      ProcedureReturn #False
    EndIf
    
    Select Gadget
      Case *ThisData\AddButton
        NewValue = GetGadgetState(*ThisData\Calendar)
        AddGadgetItem(*ThisData\ListView, -1, FormatDate(#DateFormat, NewValue))
        
      Case *ThisData\RemoveButton
        Index = GetGadgetState(*ThisData\ListView)
        If Index >= 0 And Index <= CountGadgetItems(*ThisData\ListView)
          RemoveGadgetItem(*ThisData\ListView, Index)
        EndIf
        
      Case *ThisData\Calendar, *ThisData\ListView
        ; Do nothing.
    EndSelect
  EndProcedure
  
  ;- Track Window
  Procedure.i CreateTrackWindow()
    ; Creates a new Track window and adds it to the tracking map, 
    ; allocates memory for gadget references, creates the gadgets 
    ; and stores these references in the memory Structure.
    Shared Tracks, ActiveWindows()
    Protected *ThisData.TRACKWINDOW
    Protected.i ThisWindow, ThisSum
    Protected.s ThisKey, ThisTitle
    
    Tracks + 1
    ThisTitle = "Track Bar Window " + StrU(Tracks)
    ThisWindow = OpenWindow(#PB_Any, 30, 30, 398, 130, ThisTitle, #PB_Window_SystemMenu | #PB_Window_SizeGadget | #PB_Window_MinimizeGadget | #PB_Window_TitleBar)
    
    ; Check that the OpenWindow command worked.
    If ThisWindow
      ; Set minimum window dimensions.
      WindowBounds(ThisWindow, 135, 130, #PB_Ignore, 130)
      
      ; Convert the window reference to a string to use as the map key value.
      ThisKey = StrU(ThisWindow)
      
      ; Allocate memory to store the gadget references in.
      *ThisData = AllocateMemory(SizeOf(TRACKWINDOW))
    EndIf
    
    ; Check that the memory allocation worked.
    If *ThisData
      ; Store the window reference and memory pointer values in the map.
      ActiveWindows(ThisKey) = *ThisData
      
      ; Set the window class.
      *ThisData\WindowClass = #WindowClassTrack
      
      ; Create the menu bar.
      *ThisData\Menu = CreateMenu(#PB_Any, WindowID(ThisWindow))
      
      ; If the menu creation worked, create the menu items.
      If *ThisData\Menu
        MenuTitle("Window")
        MenuItem(#MenuNewButton, "New Button Window")
        MenuItem(#MenuNewDate, "New Date Window")
        MenuItem(#MenuNewTrack, "New Track Window")
        MenuItem(#MenuClose, "Close Window")
      EndIf
      
      ; Create the window gadgets.
      *ThisData\TrackBar1 = TrackBarGadget(#PB_Any, 10, 10, 375, 25, 0, 100, #PB_TrackBar_Ticks)
      *ThisData\TrackBar2 = TrackBarGadget(#PB_Any, 10, 40, 375, 25, 0, 100, #PB_TrackBar_Ticks)
      *ThisData\Label = TextGadget(#PB_Any, 10, 75, 80, 25, "Difference:")
      *ThisData\Difference = StringGadget(#PB_Any, 90, 75, 290, 25, "0", #PB_String_ReadOnly)
    Else
      ; Memory allocation failed.
      CloseWindow(ThisWindow)
    EndIf
    
    ; Set the return value.
    If ThisWindow > 0 And *ThisData > 0
      ; Return the reference to the new window.
      ProcedureReturn ThisWindow
    Else
      ; Return 0
      ProcedureReturn 0
    EndIf
  EndProcedure
  
  Procedure.i DestroyTrackWindow(Window.i)
    ; Remove Window from the ActiveWindows map, release the allocated memory,
    ; close the window and set the quit flag, if appropriate.
    Shared EventQuit, ActiveWindows()
    Protected *ThisData.DATEWINDOW
    Protected.s ThisKey
    
    ; Convert the integer Window to a string.
    ThisKey = StrU(Window)
    
    ; Obtain the reference structure pointer.
    *ThisData = ActiveWindows(ThisKey)
    
    ; Check that a valid pointer was obtained.
    If *ThisData = 0
      ProcedureReturn #False
    EndIf
    
    ; Check that it is the correct window type, if not stop as this procedure can't destroy this window.
    If *ThisData\WindowClass <> #WindowClassTrack
      ProcedureReturn #False
    EndIf
    
    ; Release the memory allocation.
    FreeMemory(*ThisData)
    
    ; Delete the map entry.
    DeleteMapElement(ActiveWindows(), ThisKey)
    
    ; Close the window.
    CloseWindow(Window)
    
    ; Check if there are still open windows.
    If MapSize(ActiveWindows()) = 0
      EventQuit = #True
    EndIf
    
    ; Set the successful return value.
    ProcedureReturn #True
  EndProcedure
  
  Procedure.i ResizeTrackWindow(Window.i)
    ; Resize the child gadgets on Window.
    Shared ActiveWindows()
    Protected *ThisData.TRACKWINDOW
    Protected.i X, Y, W, H
    Protected.s ThisKey
    
    ; Obtain the reference structure pointer.
    ThisKey = StrU(Window)
    *ThisData = ActiveWindows(ThisKey)
    
    ; Check that a valid pointer was obtained, if not stop.
    If *ThisData = 0
      ProcedureReturn #False
    EndIf
    
    ; Check that it is the correct window type, if not stop.
    If *ThisData\WindowClass <> #WindowClassTrack
      ProcedureReturn #False
    EndIf
    
    ; Resize track bars.
    W = WindowWidth(Window) - 20
    ResizeGadget(*ThisData\TrackBar1, #PB_Ignore, #PB_Ignore, W, #PB_Ignore)
    ResizeGadget(*ThisData\TrackBar2, #PB_Ignore, #PB_Ignore, W, #PB_Ignore)
    
    ; Resize string.
    W = WindowWidth(Window) - 110
    ResizeGadget(*ThisData\Difference, #PB_Ignore, #PB_Ignore, W, #PB_Ignore)
    
    ProcedureReturn #True
  EndProcedure
  
  Procedure.i EventsTrackWindow(Window, Gadget, Type)
    ; Handle events for a Track window.
    Shared Buttons, ActiveWindows()
    Protected *ThisData.TRACKWINDOW
    Protected.i NewValue
    Protected.s ThisKey
    
    ; Convert the integer Window to a string.
    ThisKey = StrU(Window)
    
    ; Obtain the reference structure pointer.
    *ThisData = ActiveWindows(ThisKey)
    
    ; Check that a valid pointer was obtained, if not stop.
    If *ThisData = 0
      ProcedureReturn #False
    EndIf
    
    ; Check that it is the correct window type, if not stop.
    If *ThisData\WindowClass <> #WindowClassTrack
      ProcedureReturn #False
    EndIf
    
    Select Gadget
      Case *ThisData\TrackBar1, *ThisData\TrackBar2
        NewValue = GetGadgetState(*ThisData\TrackBar1) - GetGadgetState(*ThisData\TrackBar2)
        SetGadgetText(*ThisData\Difference, Str(NewValue))
        
      Case *ThisData\Label, *ThisData\Difference
        ; Do nothing.
    EndSelect
  EndProcedure
  
  ;- Main
  
  ; Create the first window.
  EventWindow = CreateButtonWindow()
  ResizeButtonWindow(EventWindow)
  
  ;- Event loop
  Repeat
    Event = WaitWindowEvent()
    EventWindow = EventWindow()
    EventWindowKey = StrU(EventWindow)
    EventGadget = EventGadget()
    EventType = EventType()
    EventMenu = EventMenu()
    *EventGadgets = ActiveWindows(EventWindowKey)
    
    Select Event
      Case #PB_Event_Gadget
        ; Check that a valid pointer was obtained.
        If *EventGadgets > 0
          ; Use *EventGadgets\WindowClass to despatch events to the correct event procedure.
          Select *EventGadgets\WindowClass
            Case #WindowClassButton
              EventsButtonWindow(EventWindow, EventGadget, EventType)
              
            Case #WindowClassDate
              EventsDateWindow(EventWindow, EventGadget, EventType)
              
            Case #WindowClassTrack
              EventsTrackWindow(EventWindow, EventGadget, EventType)
              
            Default
              ; Do nothing
          EndSelect
        EndIf
        
      Case #PB_Event_Menu
        Select EventMenu
          Case #MenuNewButton
            EventWindow = CreateButtonWindow()
            If EventWindow > 0
              ResizeButtonWindow(EventWindow)
            EndIf
            
          Case #MenuNewDate
            EventWindow = CreateDateWindow()
            If EventWindow > 0
              ResizeDateWindow(EventWindow)
            EndIf
            
          Case #MenuNewTrack
            EventWindow = CreateTrackWindow()
            If EventWindow > 0
              ResizeTrackWindow(EventWindow)
            EndIf
            
          Case #MenuClose
            ; Check that a valid pointer was obtained.
            If *EventGadgets > 0
              ; Use *EventGadgets\WindowClass to call the correct destroy window procedure.    
              Select *EventGadgets\WindowClass
                Case #WindowClassButton
                  DestroyButtonWindow(EventWindow)
                  
                Case #WindowClassDate
                  DestroyDateWindow(EventWindow)
                  
                Case #WindowClassTrack
                  DestroyTrackWindow(EventWindow)
                  
                Default
                  ; Do nothing
              EndSelect
            EndIf
        EndSelect
        
      Case #PB_Event_CloseWindow
        ; Check that a valid pointer was obtained.
        If *EventGadgets > 0
          ; Use *EventGadgets\WindowClass to call the correct destroy window procedure.
          Select *EventGadgets\WindowClass
            Case #WindowClassButton
              DestroyButtonWindow(EventWindow)
              
            Case #WindowClassDate
              DestroyDateWindow(EventWindow)
              
            Case #WindowClassTrack
              DestroyTrackWindow(EventWindow)
              
            Default
              ; Do nothing
          EndSelect
        EndIf
        
      Case #PB_Event_SizeWindow     
        If *EventGadgets > 0
          ; Use *EventGadgets\WindowClass to call the correct resize window procedure.
          Select *EventGadgets\WindowClass
            Case #WindowClassButton
              ResizeButtonWindow(EventWindow)
              
            Case #WindowClassDate
              ResizeDateWindow(EventWindow)
              
            Case #WindowClassTrack
              ResizeTrackWindow(EventWindow)
              
            Default
              ; Do nothing
          EndSelect
        EndIf
        
    EndSelect
    
  Until EventQuit = #True
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 707
; FirstLine = 651
; Folding = ---
; EnableXP