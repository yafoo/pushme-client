This is a command line application - run from cmd/PowerShell/any shell of your choice.

Usage: makelnk.exe <aumid> <exe_path> <shortcut_name> [arguments]

Generates <shortcut_name> pointing to <exe_path> (with optional arguments) with the AppUserModeId set to <aumid>.
If <shortcut_name> is not a path, '%APPDATA%\Microsoft\Windows\Start Menu\Programs' is prepended.

Github: https://github.com/Robertof/make-shortcut-with-appusermodelid