InitNetwork()
  Pattern$ = "All files (*.*)|*.*"
  Filename$ = SaveFileRequester("Guardamos", "C:\Users\Isaac\Documents\prueba.php", Pattern$, 1)

  If ReceiveHTTPFile("http://www.purebasic.com/index.php", Filename$)
    Debug "Success"
  Else
    Debug "Failed"
  EndIf
; IDE Options = PureBasic 5.61 (Windows - x86)
; CursorPosition = 2
; EnableXP