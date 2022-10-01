# Windows API using golang
GO is awesome general purpose programming language.

UI Automation is designed for experienced C/C++ developers. In general, developers need a moderate level of understanding about Component Object Model (COM) objects and interfaces, Unicode, and Windows API programming.

The UI Automation element objects are provided to clients in a hierarchical tree structure called the UI Automation tree.

* uiautomationclient.h
  * GetCurrentPropertyValue

#### https://docs.microsoft.com/en-us/windows/win32/winauto/uiauto-clientsoverview
* To gain access to the core UI Automation service, 
* a client must create an instance of the CUIAutomation object and 
* retrieve an IUIAutomation interface pointer on the object
* https://www.codemag.com/article/0810122/Creating-UI-Automation-Client-Applications
* https://github.com/hnakamur/w32uiautomation ***

github.com/go-ole/go-ole

## Required DLL
* UIAutomationClient.dll
* UiAutomationCore.dll
* Client API (UIAutomationClient.dll and UIAutomationTypes.dll)

* https://docs.microsoft.com/en-us/windows/win32/winauto/uiauto-clientsoverview
* https://docs.microsoft.com/en-us/windows/win32/winauto/uiauto-creatingcuiautomation

#include <uiautomation.h>

# Resources
* https://developer.microsoft.com/en-us/windows/downloads/windows-10-sdk/
* https://docs.microsoft.com/en-us/windows/win32/winauto/entry-uiautocore-overview
* https://docs.microsoft.com/en-us/dotnet/framework/ui-automation/ui-automation-overview
* https://docs.microsoft.com/en-us/dotnet/api/system.windows.automation.automationelement.automationelementinformation.name?view=netcore-3.1
* https://docs.microsoft.com/en-us/windows/win32/api/_controls/
* https://docs.microsoft.com/en-us/windows/win32/api/_automat/