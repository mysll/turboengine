# Service 执行顺序

## Start：
* prepare
* handler:OnPrepare
* module:Init
* init
* module:Start
* handler:OnStart

## Close：
* handler:OnShut
* Shut

## Shut：
* module:Close

