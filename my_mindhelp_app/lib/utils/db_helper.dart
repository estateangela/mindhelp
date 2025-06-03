// lib/utils/db_helper.dart

import 'dart:async';
import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:sqflite/sqflite.dart';
import 'package:path/path.dart';
import 'package:path_provider/path_provider.dart';

/// 聊天訊息的資料模型
class ChatMessage {
  final int? id;         // 自動增主鍵 (mobile 有用，web 返回 null)
  final String role;     // "user" 或 "assistant"
  final String content;  // 訊息內容
  final int timestamp;   // Unix 毫秒

  ChatMessage({
    this.id,
    required this.role,
    required this.content,
    required this.timestamp,
  });

  Map<String, dynamic> toMap() {
    return {
      'id': id,
      'role': role,
      'content': content,
      'timestamp': timestamp,
    };
  }

  factory ChatMessage.fromMap(Map<String, dynamic> map) {
    return ChatMessage(
      id: map['id'] as int?,
      role: map['role'] as String,
      content: map['content'] as String,
      timestamp: map['timestamp'] as int,
    );
  }
}

/// SQLite 資料庫助手 (Web 平台上不使用 SQLite，直接跳過)
class DBHelper {
  DBHelper._internal();
  static final DBHelper _instance = DBHelper._internal();
  factory DBHelper() => _instance;

  Database? _db;

  /// 取得 Database：若非 Web 平台，就呼叫 _initDB；否則直接拋 Exception（不會真的被呼叫）
  Future<Database> get database async {
    if (kIsWeb) {
      // 在 Web 上，根本不會用到此 database getter，因為 insert/getAll 都先檢查 kIsWeb
      throw Exception('Web 平台不使用 SQLite');
    }
    if (_db != null) return _db!;
    _db = await _initDB();
    return _db!;
  }

  Future<Database> _initDB() async {
    // 只會在非 kIsWeb 情況下被呼叫
    final docsDir = await getApplicationDocumentsDirectory();
    final path = join(docsDir.path, 'chat_history.db');
    return await openDatabase(
      path,
      version: 1,
      onCreate: (db, version) async {
        await db.execute('''
          CREATE TABLE chat_messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            role TEXT NOT NULL,
            content TEXT NOT NULL,
            timestamp INTEGER NOT NULL
          )
        ''');
      },
    );
  }

  /// 插入一筆訊息
  Future<int> insertMessage(ChatMessage msg) async {
    if (kIsWeb) {
      // Web 平台只做空操作，回傳 0
      return 0;
    }
    final dbClient = await database;
    return await dbClient.insert(
      'chat_messages',
      msg.toMap(),
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  /// 取得所有訊息 (依 timestamp 升序)
  Future<List<ChatMessage>> getAllMessages() async {
    if (kIsWeb) {
      // Web 平台回傳空列表
      return [];
    }
    final dbClient = await database;
    final List<Map<String, dynamic>> maps = await dbClient.query(
      'chat_messages',
      orderBy: 'timestamp ASC',
    );
    return maps.map((map) => ChatMessage.fromMap(map)).toList();
  }

  /// 刪除所有歷史（登出或清除用）
  Future<int> clearAll() async {
    if (kIsWeb) {
      return 0;
    }
    final dbClient = await database;
    return await dbClient.delete('chat_messages');
  }
}
