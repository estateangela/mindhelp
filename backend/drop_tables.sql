-- SQL Server 資料表清理腳本
-- ⚠️  警告：此腳本會刪除所有資料！僅用於開發環境

USE [114-MindHelp];
GO

-- 刪除外鍵約束（如果存在）
IF OBJECT_ID('FK_chat_messages_user_id') IS NOT NULL
    ALTER TABLE chat_messages DROP CONSTRAINT FK_chat_messages_user_id;

IF OBJECT_ID('FK_locations_user_id') IS NOT NULL  
    ALTER TABLE locations DROP CONSTRAINT FK_locations_user_id;

-- 刪除資料表（依序刪除以避免外鍵約束問題）
IF OBJECT_ID('chat_messages') IS NOT NULL
    DROP TABLE chat_messages;

IF OBJECT_ID('locations') IS NOT NULL
    DROP TABLE locations;

IF OBJECT_ID('users') IS NOT NULL  
    DROP TABLE users;

-- 確認刪除完成
SELECT 'Tables dropped successfully' AS Result;
