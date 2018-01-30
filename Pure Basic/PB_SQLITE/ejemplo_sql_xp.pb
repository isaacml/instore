UseSQLiteDatabase()

DatabaseFile$ = "C:\Documents and Settings\Antonio\Escritorio\sql_prueba\servext_data.db"

If OpenConsole()
  
If OpenDatabase(0, DatabaseFile$, "", "")
  
  If DatabaseQuery(0, "select id, tienda from tiendas;")
    
    While NextDatabaseRow(0)
      Debug GetDatabaseString(0, 0) + " " + GetDatabaseString(0, 1)
    Wend
    
    FinishDatabaseQuery(0)
  Else
    Debug DatabaseError()

  EndIf
  
  CloseDatabase(0)
  
Else
  Debug "Can't open database !"
  Debug DatabaseError()
EndIf


  Print(#CRLF$ + #CRLF$ + "Press ENTER to exit")
  Input()
  CloseConsole()
EndIf
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 4
; EnableXP