fichero$ = "C:/Users/Isaac/Desktop/pajaro.txt"

Procedure.s ReadBlock64(File.i, Size.i)
  Protected BytesRead.i
  Protected Result.s = ""
 
  If (IsFile(File))
    *Block = AllocateMemory(Size)
    If (*Block)
      BytesRead = ReadData(File, *Block, Size)
      If (BytesRead > 0)
        *Encode = AllocateMemory(Int(BytesRead * 1.4))
        If (*Encode)
          Base64EncoderBuffer(*Block, BytesRead, *Encode, Int(BytesRead * 1.4))
          Result = PeekS(*Encode)
          FreeMemory(*Encode)
        EndIf
      Else
      EndIf
      FreeMemory(*Block)
    EndIf
  EndIf
 
  ProcedureReturn Result
EndProcedure

; Write 100 random bytes between 0 and 32
If CreateFile(0, "C:/Users/Isaac/Desktop/algo.data")
  For i = 1 To 100
    WriteByte(0, Random(32))
  Next i
  CloseFile(0)
EndIf

; Read the data back as a safe, alphanumeric string
If ReadFile(0, "C:/Users/Isaac/Desktop/algo.data")
  Encoded.s = ReadBlock64(0, 1000)
  Debug Encoded
  CloseFile(0)
EndIf

; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 40
; Folding = -
; EnableXP