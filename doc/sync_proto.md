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
