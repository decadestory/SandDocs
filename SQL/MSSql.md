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

10.获取数据库死锁(放在master数据库中)
##

```
USE [master]
GO
/****** Object:  StoredProcedure [dbo].[sp_who_lock]    Script Date: 2017/1/11 15:50:41 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE procedure [dbo].[sp_who_lock]
as
begin
declare @spid int,@bl int,
 @intTransactionCountOnEntry  int,
        @intRowcount    int,
        @intCountProperties   int,
        @intCounter    int

 create table #tmp_lock_who (
 id int identity(1,1),
 spid smallint,
 bl smallint)
 
 IF @@ERROR<>0 RETURN @@ERROR
 
 insert into #tmp_lock_who(spid,bl) select  0 ,blocked
   from (select * from sysprocesses where  blocked>0 ) a 
   where not exists(select * from (select * from sysprocesses where  blocked>0 ) b 
   where a.blocked=spid)
   union select spid,blocked from sysprocesses where  blocked>0

 IF @@ERROR<>0 RETURN @@ERROR 
  
-- 找到临时表的记录数
 select  @intCountProperties = Count(*),@intCounter = 1
 from #tmp_lock_who
 
 IF @@ERROR<>0 RETURN @@ERROR 
 
 if @intCountProperties=0
  select '现在没有阻塞和死锁信息' as message

-- 循环开始
while @intCounter <= @intCountProperties
begin
-- 取第一条记录
  select  @spid = spid,@bl = bl
  from #tmp_lock_who where Id = @intCounter 
 begin
  if @spid =0 
            select '引起数据库死锁的是: '+ CAST(@bl AS VARCHAR(10)) + '进程号,其执行的SQL语法如下'
 else
            select '进程号SPID：'+ CAST(@spid AS VARCHAR(10))+ '被' + '进程号SPID：'+ CAST(@bl AS VARCHAR(10)) +'阻塞,其当前进程执行的SQL语法如下'
 DBCC INPUTBUFFER (@bl )
 end 

-- 循环指针下移
 set @intCounter = @intCounter + 1
end

drop table #tmp_lock_who

return 0
end
```
```
用法：exec sp_who_lock
```

11.杀掉死锁进程（放在master数据库中）
##

```
USE [master]
GO
/****** Object:  StoredProcedure [dbo].[p_killspid]    Script Date: 2017/1/11 15:59:36 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

ALTER proc [dbo].[p_killspid]
@dbname varchar(200)    --要关闭进程的数据库名
as  
	begin 

    declare @sql  nvarchar(500)  
    declare @spid nvarchar(20)

    declare #tb cursor for
        select spid=cast(spid as varchar(20)) from master..sysprocesses where dbid=db_id(@dbname)
    open #tb
    fetch next from #tb into @spid
    while @@fetch_status=0
    begin  
        exec('kill '+@spid)
        fetch next from #tb into @spid
    end  
    close #tb
    deallocate #tb

	end
```
```
用法：exec p_killspid  '数据库名'
```

12.查询死锁表信息（放在master数据库中）
##

```
USE [master]
GO
/****** Object:  StoredProcedure [dbo].[sp_lock_table]    Script Date: 2017/1/11 16:02:01 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

ALTER PROCEDURE [dbo].[sp_lock_table]
AS
BEGIN
	
--查看锁信息
create table #t(req_spid int,obj_name sysname)

declare @s nvarchar(4000)
    ,@rid int,@dbname sysname,@id int,@objname sysname

declare tb cursor for 
    select distinct req_spid,dbname=db_name(rsc_dbid),rsc_objid
    from master..syslockinfo where rsc_type in(4,5)
open tb
fetch next from tb into @rid,@dbname,@id
while @@fetch_status=0
begin
    set @s='select @objname=name from ['+@dbname+']..sysobjects where id=@id'
    exec sp_executesql @s,N'@objname sysname out,@id int',@objname out,@id
    insert into #t values(@rid,@objname)
    fetch next from tb into @rid,@dbname,@id
end
close tb
deallocate tb

select 进程id=a.req_spid
    ,数据库=db_name(rsc_dbid)
    ,类型=case rsc_type when 1 then 'NULL 资源（未使用）'
        when 2 then '数据库'
        when 3 then '文件'
        when 4 then '索引'
        when 5 then '表'
        when 6 then '页'
        when 7 then '键'
        when 8 then '扩展盘区'
        when 9 then 'RID（行 ID)'
        when 10 then '应用程序'
    end
    ,对象id=rsc_objid
    ,对象名=b.obj_name
    ,rsc_indid
 from master..syslockinfo a left join #t b on a.req_spid=b.req_spid

drop table #t

END
```
```
用法：exec sp_lock_table
```

13. 查询执行慢的SQL 
##

```
SELECT
(total_elapsed_time / execution_count)/1000 N'平均时间ms'
,total_elapsed_time/1000 N'总花费时间ms'
,total_worker_time/1000 N'所用的CPU总时间ms'
,total_physical_reads N'物理读取总次数'
,total_logical_reads/execution_count N'每次逻辑读次数'
,total_logical_reads N'逻辑读取总次数'
,total_logical_writes N'逻辑写入总次数'
,execution_count N'执行次数'
,SUBSTRING(st.text, (qs.statement_start_offset/2) + 1,
((CASE statement_end_offset
WHEN -1 THEN DATALENGTH(st.text)
ELSE qs.statement_end_offset END
- qs.statement_start_offset)/2) + 1) N'执行语句'
,creation_time N'语句编译时间'
,last_execution_time N'上次执行时间'
FROM
sys.dm_exec_query_stats AS qs CROSS APPLY sys.dm_exec_sql_text(qs.sql_handle) st
WHERE
SUBSTRING(st.text, (qs.statement_start_offset/2) + 1,
((CASE statement_end_offset
WHEN -1 THEN DATALENGTH(st.text)
ELSE qs.statement_end_offset END
- qs.statement_start_offset)/2) + 1) not like '%fetch%'
ORDER BY
total_elapsed_time / execution_count DESC;
```