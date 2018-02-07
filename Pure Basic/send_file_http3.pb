fichero$ = "C:/Users/Isaac/Desktop/pajaro.txt"

Procedure.i NetworkSendFile(Connection.i, Filename$)
  
  Protected Result.i, File.i, Size.q, *Buffer, Offset.i
  
  File = ReadFile(#PB_Any, Filename$)
  If File
    Size = Lof(File)
    Filename$ = "file:" + GetFilePart(Filename$) + ":" + Str(Size)
    *Buffer =  AllocateMemory(Size + Len(Filename$) + 1)
    If *Buffer
      PokeS(*Buffer, Filename$)
      If ReadData(File, *Buffer + Len(Filename$) + 1, Size) = Size
        If SendNetworkData(Connection, *Buffer, MemorySize(*Buffer)) = MemorySize(*Buffer)
          Result = #True
        EndIf
      EndIf
    EndIf
    CloseFile(File)
  EndIf
  
  ProcedureReturn Result
  
EndProcedure




If InitNetwork()
  
  Filename$ = OpenFileRequester("Chose a file to send", "", "*.*|*.*", 0)
  If Filename$ <> ""
    Connection = OpenNetworkConnection("192.168.4.22", 9999)
    If Connection
      If NetworkSendFile(Connection, Filename$)
        MessageRequester("Info", Filename$ + " transmitted.")
      Else
        MessageRequester("Error", "Failed to send " + Filename$)
      EndIf
      CloseNetworkConnection(Connection)
    EndIf
  EndIf
  
EndIf
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 31
; Folding = -
; EnableXP