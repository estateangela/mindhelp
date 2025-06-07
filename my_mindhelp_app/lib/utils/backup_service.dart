// lib/utils/backup_service.dart

import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:firebase_auth/firebase_auth.dart';

// 只从 models 里引入 ChatMessage
import 'package:my_mindhelp_app/models/chat_message.dart';

// 只从 utils/db_helper.dart 引入 DBHelper，屏蔽其中的 ChatMessage
import 'package:my_mindhelp_app/utils/db_helper.dart' show DBHelper;

class BackupService {
  final _fire = FirebaseFirestore.instance;
  final _auth = FirebaseAuth.instance;

  Future<void> uploadAll() async {
    final user = _auth.currentUser;
    if (user == null) throw '尚未登入';
    final uid = user.uid;

    final all = await DBHelper().getAllMessages(); // List<ChatMessage>
    final jsonList = all.map((m) => m.toMap()).toList();

    await _fire.collection('conversations').doc(uid).set({
      'updatedAt': FieldValue.serverTimestamp(),
      'messages': jsonList,
    });
  }

  Future<void> restoreAll() async {
    final user = _auth.currentUser;
    if (user == null) throw '尚未登入';
    final uid = user.uid;

    final doc = await _fire.collection('conversations').doc(uid).get();
    if (!doc.exists) return;

    final rawList = (doc.data()!['messages'] as List<dynamic>);
    final msgs = rawList.map((e) {
      return ChatMessage.fromMap(Map<String, dynamic>.from(e as Map));
    }).toList();

    final db = DBHelper();
    await db.deleteAll(); // 清掉本機所有舊對話
    for (var m in msgs) {
      await db.insertMessage(m);
    }
  }
}
