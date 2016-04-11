常用Mssql收集
--------

1、获取表字段描述
##
```
select t.name, c.name, ep.value  from sys.tables t
INNER JOIN sys.columns c ON t.object_id = c.object_id
LEFT JOIN sys.extended_properties ep ON ep.major_id = c.object_id AND ep.minor_id = c.column_id 
WHERE ep.class =1 AND t.name='TABLENAME'
```

2、获取表字段信息
##
```
sp_columns TABLENAME
```

3、添加描述
##
```
EXEC sp_addextendedproperty 
'MS_Description','description_string',
'user',dbo,
'table',TABLENAME,
'column',COLUMNNAME
```

4.生成实体类
##
```

select 
'///<summary>'+char(13)+ 
'///'+cast(ep.value as varchar)+char(13)+
'///</summary>'+char(13),
'public',(
case 
when  ty.[name] in ('text','ntext' ,'char','nchar', 'varchar', 'nvarchar') then 'string'
when ty.[name] in ('date' , 'datetime' , 'datetime2') then 'DateTime'
when ty.[name] in ('bit') then 'bool'
when ty.[name] in ('smallint') then 'short'
when ty.[name] in ('bigint') then 'long'
when ty.[name] in ('real') then 'float'
when ty.[name] in ('float') then 'double'
when ty.[name] in ('money') then 'decimal'
when ty.[name] in ('uniqueidentifier') then 'Guid'
else ty.[name] end
) as typeName,
(case c.[is_nullable] when 1 then 
	case when ty.[name] not in('text','ntext' ,'char','nchar', 'varchar', 'nvarchar') then '?' else '' end
 else ''end) as isnulls
,
c.name+' {set;get;}' 
from sys.tables t
INNER JOIN sys.columns c ON t.object_id = c.object_id
LEFT JOIN sys.extended_properties ep ON ep.major_id = c.object_id AND ep.minor_id = c.column_id 
left JOIN sys.types ty on ty.[system_type_id]=c.[user_type_id] and ty.[name]!='sysname'
WHERE ep.class =1 AND t.name='TABLENAME'
```
5.服务器获取所有数据库名
##
```
SELECT * FROM Master..SysDatabases ORDER BY Name
```
6.根据数据库获取所有表
##
```
SELECT * FROM DatabaseName..SysObjects Where XType='U' ORDER BY Name
 
XType='U':表示所有用户表;
XType='S':表示所有系统表;
```
7.根据表获取所有字段
##
```
SELECT * FROM SysColumns WHERE id=Object_Id('TableName')
```
8.获取所有系统类型
##
```
select * from sys.types
```
9.跨数据库查询
##
```
第一种

--开启
exec sp_configure 'show advanced options',1
reconfigure
exec sp_configure 'Ad Hoc Distributed Queries',1
reconfigure
--执行
select top 100 * from OPENDATASOURCE('SQLOLEDB','Data Source=.;User ID=sa;Password=sa').DataBaseName.dbo.TableName
--关闭
exec sp_configure 'Ad Hoc Distributed Queries',0
reconfigure
exec sp_configure 'show advanced options',0
reconfigure 

第二种

-- 创建链接服务器
exec sp_addlinkedserver 'svr_link','','sqloledb','192.168.1.1'
-- 创建登录信息
exec sp_addlinkedsrvlogin 'svr_link','false',null,'sa','sa'
--查询 格式为：链接服务器.数据库名.架构名.表名
select top 100 * from svr_link.DatabaseName.dbo.TableName
-- 删除链接服务器
exec sp_dropserver 'svr_link','droplogins'
```
