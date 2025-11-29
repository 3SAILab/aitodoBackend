1. 创建任务
谁（比如你）可以新建一个任务，选择任务类型
干什么用的（描述，可不填）
谁负责（指定一个人）
什么时候必须做完（截止时间）
放到哪个看板里（比如“设计项目”看板）

2. 编辑任务
任务创建后，可以随时修改：任务类型、描述
负责人（比如换人做）
截止时间（比如延期了）
状态（比如从“没开始”变成“正在做”或“做完了”）

3. 自动记录所有改动
只要上面这些信息被改了，系统就自动记一笔账，比如：“2025年11月26日上午10点，张三把任务负责人从李四改成了王五”
“2025年11月26日下午2点，王五把截止时间从12月1日改成了12月5日”
这些记录谁也删不掉，就像聊天记录一样，随时可以翻出来看“这事儿是谁什么时候改的”。

4. 团队看板（Kanban）
所有任务按状态分成几列，比如：
待办｜进行中｜已完成
每个任务是一张卡片，你可以直接拖来拖去：拖到“进行中” → 状态自动更新，系统自动记一笔“谁把任务开始做了”
拖到“已完成” → 自动记录“谁完成的、什么时候完成的”
整个团队都能看到所有任务在哪个阶段，谁在负责，有没有快超时。

5. 任务类型修改

6.登录

7.后台管理用户密码和账号

backend/app/user/
├── api/
│   ├── etc/              # 配置文件 (已存在)
│   ├── internal/
│   │   ├── config/       # 配置结构体映射
│   │   ├── handler/      # HTTP 路由处理 (Controller)
│   │   ├── logic/        # 业务逻辑 (Service)
│   │   ├── svc/          # 依赖注入上下文 (DB, Redis等)
│   │   └── types/        # 请求/响应结构体 (DTO)
│   └── user.go           # main 入口文件
└── model/                # 数据库模型 (DAO)


====================================================================================================
PostgreSQL 数据库表结构汇总（public模式）
导出时间：2025-11-28 08:57:30
====================================================================================================

【表名】: sales_persons
------------------------------------------------------------------------------------------
列名                   列类型                       列描述                                     
------------------------------------------------------------------------------------------
id                   character varying(64)     主键 UUID                                 
name                 text                      销售人员姓名                                  
phone                text                      销售联系电话 (可选)                             
is_active            boolean                   在职状态：true=正常显示，false=离职/隐藏 (软删除)        
created_at           timestamp without time zone 创建时间                                    
updated_at           timestamp without time zone 最后更新时间                                  

--------------------------------------------------

【表名】: task_activity_logs
------------------------------------------------------------------------------------------
列名                   列类型                       列描述                                     
------------------------------------------------------------------------------------------
id                   uuid                      主键 UUID                                 
task_id              uuid                      被修改的任务 ID                               
actor_id             character varying(64)     执行操作的用户 ID                              
action_type          text                      操作类型：CREATE, UPDATE, DELETE, STATUS_CHANGE
field_changed        text                      被修改的字段名 (如 status, assignee_id)         
old_value            jsonb                     修改前的旧值 (JSONB 格式)                       
new_value            jsonb                     修改后的新值 (JSONB 格式)                       
description          text                      人类可读的操作描述 (如“张三修改了截止时间”)                
created_at           timestamp without time zone 操作发生的时间                                 

--------------------------------------------------

【表名】: task_types
------------------------------------------------------------------------------------------
列名                   列类型                       列描述                                     
------------------------------------------------------------------------------------------
id                   uuid                      主键 UUID                                 
name                 text                      类型名称                                    
color_code           text                      看板上显示的颜色代码 (如 #FF0000)                  
created_by           uuid                      创建该类型的管理员 ID (关联 users.id)              
created_at           timestamp without time zone 创建时间                                    
updated_at           timestamp without time zone 最后更新时间                                  

--------------------------------------------------

【表名】: tasks
------------------------------------------------------------------------------------------
列名                   列类型                       列描述                                     
------------------------------------------------------------------------------------------
id                   character varying(64)     主键 UUID                                 
type_id              uuid                      任务类型 ID (关联 task_types)                 
creator_id           character varying(64)     任务创建人 ID (关联 users)                     
assignee_id          character varying(64)     当前任务负责人 ID (关联 users，可为空)               
sales_person_id      character varying(64)     关联的销售人员 ID (关联 sales_persons，可为空)       
title                text                      任务标题                                    
description          text                      任务详细描述                                  
status               text                      当前状态：TODO(待办), DOING(进行中), DONE(已完成)    
sort_order           integer                   排序权重：用于看板内卡片的上下拖拽排序                     
due_date             timestamp without time zone 截止时间 (不带时区，仅年月日时分)                      
is_deleted           boolean                   软删除标记：true=已删除，false=正常                 
created_at           timestamp without time zone 任务创建时间                                  
updated_at           timestamp without time zone 任务最后修改时间                                
deleted_at           timestamp with time zone                                          
content              text                                                              

--------------------------------------------------

【表名】: user_behavior_logs
------------------------------------------------------------------------------------------
列名                   列类型                       列描述                                     
------------------------------------------------------------------------------------------
id                   uuid                      主键 UUID                                 
user_id              character varying(64)     用户 ID (游客可能为空)                          
event_type           text                      事件类型：CLICK(点击), VIEW(浏览), STAY(停留)      
event_target         text                      事件目标标识 (如 btn_submit, page_home)        
page_url             text                      事件发生的页面 URL                             
meta_data            jsonb                     元数据 (JSONB)：存储停留时长、错误信息等额外参数            
created_at           timestamp without time zone 事件发生时间                                  

--------------------------------------------------

【表名】: user_login_logs
------------------------------------------------------------------------------------------
列名                   列类型                       列描述                                     
------------------------------------------------------------------------------------------
id                   uuid                      主键 UUID                                 
user_id              character varying(64)     登录的用户 ID                                
ip_address           character varying(45)     登录时的 IP 地址                              
user_agent           text                      浏览器或设备信息 (User-Agent)                   
login_at             timestamp without time zone 登录时间                                    

--------------------------------------------------

【表名】: users
------------------------------------------------------------------------------------------
列名                   列类型                       列描述                                     
------------------------------------------------------------------------------------------
id                   character varying(64)     主键 UUID                                 
username             text                      用户显示名称                                  
email                text                      电子邮箱，作为唯一的登录账号                          
password_hash        text                      加密后的密码哈希值 (不要存明文)                       
role                 text                      用户角色：admin(管理员) 或 user(普通用户)            
avatar_url           text                      用户头像图片的 URL 地址                          
is_active            boolean                   账户状态：true=正常，false=冻结/禁用                
created_at           timestamp without time zone 注册时间 (服务器本地时间)                          
updated_at           timestamp without time zone 最后一次资料更新时间 (服务器本地时间)                    

--------------------------------------------------

