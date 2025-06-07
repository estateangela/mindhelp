// lib/utils/db_helper.dart

import 'dart:async';
import 'dart:io';

import 'package:my_mindhelp_app/models/chat_message.dart';
import 'package:path/path.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';

class DBHelper {
  // 單例
  DBHelper._privateConstructor();
  static final DBHelper _instance = DBHelper._privateConstructor();
  factory DBHelper() => _instance;

  static Database? _db;

  /// 取得資料庫實例
  Future<Database> get database async {
    if (_db != null) return _db!;
    _db = await _initDB();
    return _db!;
  }

  /// 初始化並開啟資料庫
  Future<Database> _initDB() async {
    // 取得應用程式文件資料夾路徑
    Directory docsDir = await getApplicationDocumentsDirectory();
    String path = join(docsDir.path, 'mindhelp.db');

    // openDatabase 會自動建立檔案
    return await openDatabase(
      path,
      version: 1,
      onCreate: _onCreate,
    );
  }

  /// 建表
  FutureOr<void> _onCreate(Database db, int version) async {
    await db.execute('''
      CREATE TABLE messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        role TEXT NOT NULL,
        content TEXT NOT NULL,
        timestamp INTEGER NOT NULL
      );
    ''');
  }

  /// 插入一筆訊息
  Future<int> insertMessage(ChatMessage msg) async {
    final db = await database;
    return await db.insert('messages', msg.toMap());
  }

  /// 讀取所有訊息（依 timestamp 升序）
  Future<List<ChatMessage>> getAllMessages() async {
    final db = await database;
    final rows = await db.query(
      'messages',
      orderBy: 'timestamp ASC',
    );
    return rows.map((r) => ChatMessage.fromMap(r)).toList();
  }

  /// 刪除所有訊息
  Future<void> deleteAll() async {
    final db = await database;
    await db.delete('messages');
  }
}
