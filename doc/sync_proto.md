# SYNC协议

##  数据源数据结构： 
    Data:	DataVersion
            [数据]

    Patch:	[
                {
                PatchVersion
                [数据]
                }
            ]

    Node:	[
                {
                NodeVersion
                NodeId
                SyncTimeOut
                }
            ]
##  初始化： 
    Data:  
    DataVersion=0
    []
    Patch:  
    []
    Node:
    []

##  初始数据加载完成后
    Data:
        DataVersion++
        [数据]
等待Node加入

## Node加入处理：
向Node中插入数据，初始值为
#   
    {
        NodeVersion=0
        NodeId=id
        SyncTimeOut=0
    }
并向Node发送Data数据段,SyncTimeOut设置超时时间  
Node同步数据后，回发确认消息，并带回DataVersion，收到DataVersion后更新NodeVersion。并检查Patch中版本号，如果比Patch中小，则依次更新Patch中数据。每次更新后都要确认版本号，并记录在Node中。  
如果超时后仍未收到确认,则直接删除Node信息。如果是网络断开，超时后再联上，则按初次加入的步骤，全量更新。